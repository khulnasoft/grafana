package export

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/grafana/grafana-app-sdk/logging"
	dashboards "github.com/grafana/grafana/pkg/apis/dashboard"
	provisioning "github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1"
	"github.com/grafana/grafana/pkg/registry/apis/dashboard/legacy"
	"github.com/grafana/grafana/pkg/storage/unified/parquet"
	"github.com/grafana/grafana/pkg/storage/unified/resource"
)

var (
	_ resource.BatchResourceWriter = (*resourceReader)(nil)
)

type resourceReader struct {
	job     *exportJob
	summary *provisioning.JobResourceSummary
	logger  logging.Logger
}

// Close implements resource.BatchResourceWriter.
func (f *resourceReader) Close() error {
	return nil
}

// CloseWithResults implements resource.BatchResourceWriter.
func (f *resourceReader) CloseWithResults() (*resource.BatchResponse, error) {
	return &resource.BatchResponse{}, nil
}

// Write implements resource.BatchResourceWriter.
func (f *resourceReader) Write(ctx context.Context, key *resource.ResourceKey, value []byte) error {
	item := &unstructured.Unstructured{}
	err := item.UnmarshalJSON(value)
	if err != nil {
		return err
	}
	err = f.job.add(ctx, f.summary, item)
	if err != nil {
		f.logger.Warn("error adding from legacy", "name", key.Name, "err", err)
		f.summary.Errors = append(f.summary.Errors, fmt.Sprintf("%s: %s", key.Name, err.Error()))
		if len(f.summary.Errors) > 50 {
			return err
		}
	}
	return nil
}

func (r *exportJob) loadResources(ctx context.Context) error {
	kinds := []schema.GroupVersionResource{{
		Group:    dashboards.GROUP,
		Resource: dashboards.DASHBOARD_RESOURCE,
		Version:  "v1alpha1",
	}}

	for _, kind := range kinds {
		r.jobStatus.Message = "Exporting " + kind.Resource + "..."
		if r.legacy != nil {
			gr := kind.GroupResource()
			reader := &resourceReader{
				summary: r.getSummary(gr),
				job:     r,
				logger:  r.logger,
			}
			opts := legacy.MigrateOptions{
				Namespace:   r.namespace,
				WithHistory: r.withHistory,
				Resources:   []schema.GroupResource{gr},
				Store:       parquet.NewBatchResourceWriterClient(reader),
				OnlyCount:   true, // first get the count
			}
			stats, err := r.legacy.Migrate(ctx, opts)
			if err != nil {
				return fmt.Errorf("unable to count legacy items %w", err)
			}
			if len(stats.Summary) > 0 {
				count := stats.Summary[0].Count
				history := stats.Summary[0].History
				if history > count {
					count = history // the number of items we will process
				}
				reader.summary.Total = count
			}

			opts.OnlyCount = false // this time actually write
			_, err = r.legacy.Migrate(ctx, opts)
			if err != nil {
				return fmt.Errorf("error running legacy migrate %s %w", kind.Resource, err)
			}
		}

		if err := r.loadResourcesFromAPIServer(ctx, kind); err != nil {
			return fmt.Errorf("error loading %s %w", kind.Resource, err)
		}
	}
	return nil
}

func (r *exportJob) loadResourcesFromAPIServer(ctx context.Context, kind schema.GroupVersionResource) error {
	r.maybeNotify(ctx)
	client := r.client.Resource(kind)
	summary := r.getSummary(kind.GroupResource())

	continueToken := ""
	for {
		list, err := client.List(ctx, metav1.ListOptions{Limit: 100, Continue: continueToken})
		if err != nil {
			return fmt.Errorf("error executing list: %w", err)
		}

		for _, item := range list.Items {
			if err = r.add(ctx, summary, &item); err != nil {
				return fmt.Errorf("error adding value: %w", err)
			}
		}

		continueToken = list.GetContinue()
		if continueToken == "" {
			break
		}
	}

	return nil
}
