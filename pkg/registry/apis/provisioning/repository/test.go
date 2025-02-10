package repository

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	provisioning "github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Tester is a struct that implements the Tester interface
// it's temporary
// FIXME: remove as soon as controller and jobs refactoring PRs are merged
type Tester struct{}

func (t *Tester) TestRepository(ctx context.Context, repo Repository) (*provisioning.TestResults, error) {
	return TestRepository(ctx, repo)
}

func TestRepository(ctx context.Context, repo Repository) (*provisioning.TestResults, error) {
	errors := ValidateRepository(repo)
	if len(errors) > 0 {
		rsp := &provisioning.TestResults{
			Code:    http.StatusUnprocessableEntity, // Invalid
			Success: false,
			Errors:  make([]string, len(errors)),
		}
		for i, v := range errors {
			rsp.Errors[i] = v.Error()
		}
		return rsp, nil
	}

	return repo.Test(ctx)
}

func ValidateRepository(repo Repository) field.ErrorList {
	list := repo.Validate()
	cfg := repo.Config()

	if cfg.Spec.Title == "" {
		list = append(list, field.Required(field.NewPath("spec", "title"), "a repository title must be given"))
	}

	if cfg.Spec.Sync.Enabled && cfg.Spec.Sync.Target == "" {
		list = append(list, field.Required(field.NewPath("spec", "sync", "target"),
			"The target type is required when sync is enabled"))
	}

	if cfg.Spec.Sync.Enabled && cfg.Spec.Sync.IntervalSeconds < 10 {
		list = append(list, field.Invalid(field.NewPath("spec", "sync", "intervalSeconds"),
			cfg.Spec.Sync.IntervalSeconds, fmt.Sprintf("Interval must be at least %d seconds", 10)))
	}

	// Reserved names (for now)
	reserved := []string{"classic", "sql", "SQL", "plugins", "legacy", "new", "job", "github", "s3", "gcs", "file", "new", "create", "update", "delete"}
	if slices.Contains(reserved, cfg.Name) {
		list = append(list, field.Invalid(field.NewPath("metadata", "name"), cfg.Name, "Name is reserved, choose a different identifier"))
	}

	if cfg.Spec.Type != provisioning.LocalRepositoryType && cfg.Spec.Local != nil {
		list = append(list, field.Invalid(field.NewPath("spec", "local"),
			cfg.Spec.GitHub, "Local config only valid when type is local"))
	}

	if cfg.Spec.Type != provisioning.GitHubRepositoryType && cfg.Spec.GitHub != nil {
		list = append(list, field.Invalid(field.NewPath("spec", "github"),
			cfg.Spec.GitHub, "Github config only valid when type is github"))
	}

	if cfg.Spec.Type != provisioning.S3RepositoryType && cfg.Spec.S3 != nil {
		list = append(list, field.Invalid(field.NewPath("spec", "s3"),
			cfg.Spec.GitHub, "S3 config only valid when type is s3"))
	}
	return list
}
