package provisioning

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/grafana/grafana/pkg/apimachinery/identity"
	provisioning "github.com/grafana/grafana/pkg/apis/provisioning/v0alpha1"
	client "github.com/grafana/grafana/pkg/generated/clientset/versioned/typed/provisioning/v0alpha1"
	informer "github.com/grafana/grafana/pkg/generated/informers/externalversions/provisioning/v0alpha1"
	listers "github.com/grafana/grafana/pkg/generated/listers/provisioning/v0alpha1"
	"github.com/grafana/grafana/pkg/registry/apis/provisioning/auth"
	"github.com/grafana/grafana/pkg/registry/apis/provisioning/jobs"
)

type operation int

const (
	operationCreate operation = iota
	operationUpdate
	operationDelete
)

const maxAttempts = 3

type queueItem struct {
	key      string
	op       operation
	obj      interface{}
	attempts int
}

// RepositoryController controls how and when CRD is established.
type RepositoryController struct {
	client     client.ProvisioningV0alpha1Interface
	repoLister listers.RepositoryLister
	repoSynced cache.InformerSynced

	jobs jobs.JobQueue

	// Converts config to instance
	repoGetter RepoGetter
	identities auth.BackgroundIdentityService
	tester     *RepositoryTester

	// To allow injection for testing.
	processFn         func(item *queueItem) error
	enqueueRepository func(op operation, obj any)
	keyFunc           func(obj any) (string, error)

	queue  workqueue.TypedRateLimitingInterface[*queueItem]
	logger *slog.Logger
}

// NewRepositoryController creates new RepositoryController.
func NewRepositoryController(
	provisioningClient client.ProvisioningV0alpha1Interface,
	repoInformer informer.RepositoryInformer,
	repoGetter RepoGetter,
	identities auth.BackgroundIdentityService,
	tester *RepositoryTester,
	jobs jobs.JobQueue,
) (*RepositoryController, error) {
	rc := &RepositoryController{
		client:     provisioningClient,
		repoLister: repoInformer.Lister(),
		repoSynced: repoInformer.Informer().HasSynced,
		queue: workqueue.NewTypedRateLimitingQueueWithConfig(
			workqueue.DefaultTypedControllerRateLimiter[*queueItem](),
			workqueue.TypedRateLimitingQueueConfig[*queueItem]{
				Name: "provisioningRepositoryController",
			},
		),
		repoGetter: repoGetter,
		identities: identities,
		tester:     tester,
		logger:     slog.Default().With("logger", "provisioning-repository-controller"),
		jobs:       jobs,
	}

	_, err := repoInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    rc.addRepository,
		UpdateFunc: rc.updateRepository,
		DeleteFunc: rc.deleteRepository,
	})
	if err != nil {
		return nil, err
	}

	rc.processFn = rc.process
	rc.enqueueRepository = rc.enqueue
	rc.keyFunc = repoKeyFunc

	return rc, nil
}

func repoKeyFunc(obj any) (string, error) {
	repo, ok := obj.(*provisioning.Repository)
	if !ok {
		return "", fmt.Errorf("expected a Repository but got %T", obj)
	}
	return cache.DeletionHandlingMetaNamespaceKeyFunc(repo)
}

// Run starts the RepositoryController.
func (rc *RepositoryController) Run(ctx context.Context, workerCount int) {
	defer utilruntime.HandleCrash()
	defer rc.queue.ShutDown()

	rc.logger.Info("Starting RepositoryController")
	defer rc.logger.Info("Shutting down RepositoryController")

	if !cache.WaitForCacheSync(ctx.Done(), rc.repoSynced) {
		return
	}

	rc.logger.Info("Starting workers", "count", workerCount)
	for i := 0; i < workerCount; i++ {
		go wait.UntilWithContext(ctx, rc.runWorker, time.Second)
	}

	rc.logger.Info("Started workers")
	<-ctx.Done()
	rc.logger.Info("Shutting down workers")
}

func (rc *RepositoryController) runWorker(ctx context.Context) {
	for rc.processNextWorkItem(ctx) {
	}
}

func (rc *RepositoryController) enqueue(op operation, obj interface{}) {
	key, err := rc.keyFunc(obj)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("couldn't get key for object: %v", err))
		return
	}

	item := queueItem{key: key, obj: obj, op: op}
	rc.queue.Add(&item)
}

func (rc *RepositoryController) addRepository(obj interface{}) {
	rc.enqueueRepository(operationCreate, obj)
}

func (rc *RepositoryController) updateRepository(oldObj, newObj interface{}) {
	rc.enqueueRepository(operationUpdate, newObj)
}

func (rc *RepositoryController) deleteRepository(obj interface{}) {
	rc.enqueueRepository(operationDelete, obj)
}

// processNextWorkItem deals with one key off the queue.
// It returns false when it's time to quit.
func (rc *RepositoryController) processNextWorkItem(ctx context.Context) bool {
	item, quit := rc.queue.Get()
	if quit {
		return false
	}
	defer rc.queue.Done(item)

	logger := rc.logger.With("key", item.key)
	logger.InfoContext(ctx, "RepositoryController processing key")

	err := rc.processFn(item)
	if err == nil {
		rc.queue.Forget(item)
		return true
	}

	item.attempts++
	logger = logger.With("error", err, "attempts", item.attempts)
	if item.attempts >= maxAttempts {
		logger.ErrorContext(ctx, "RepositoryController failed too many times")
		rc.queue.Forget(item)
		return true
	}

	if !apierrors.IsServiceUnavailable(err) {
		logger.InfoContext(ctx, "RepositoryController will not retry")
		rc.queue.Forget(item)
		return true
	} else {
		logger.InfoContext(ctx, "RepositoryController will retry as service is unavailable")
	}

	utilruntime.HandleError(fmt.Errorf("%v failed with: %v", item, err))
	rc.queue.AddRateLimited(item)

	return true
}

// process is the business logic of the controller.
func (rc *RepositoryController) process(item *queueItem) error {
	logger := rc.logger.With("key", item.key)
	namespace, name, err := cache.SplitMetaNamespaceKey(item.key)
	if err != nil {
		return err
	}

	ctx := context.Background()
	if item.op == operationDelete {
		logger.InfoContext(ctx, "handle repository deletion")
		cfg, ok := item.obj.(*provisioning.Repository)
		if !ok {
			return errors.New("object is not a repository")
		}

		repo, err := rc.repoGetter.AsRepository(ctx, cfg)
		if err != nil {
			return fmt.Errorf("unable to create repository from object: %w", err)
		}

		return repo.OnDelete(ctx, logger)
	}

	cachedRepo, err := rc.repoLister.Repositories(namespace).Get(name)
	switch {
	case apierrors.IsNotFound(err):
		return errors.New("repository not found in cache")
	case err != nil:
		return err
	}

	id, err := rc.identities.WorkerIdentity(ctx, cachedRepo.Namespace)
	if err != nil {
		return err
	}
	ctx = identity.WithRequester(ctx, id)
	logger = logger.With("repository", cachedRepo.Name, "namespace", cachedRepo.Namespace)

	repo, err := rc.repoGetter.AsRepository(ctx, cachedRepo)
	if err != nil {
		return fmt.Errorf("unable to create repository from configuration: %w", err)
	}

	hasSpecChanged := cachedRepo.Generation != cachedRepo.Status.ObservedGeneration
	if !hasSpecChanged {
		logger.InfoContext(ctx, "repository spec unchanged")
		return nil
	}

	logger.InfoContext(ctx, "repository spec changed", "previous", cachedRepo.Status.ObservedGeneration, "current", cachedRepo.Generation)

	var status *provisioning.RepositoryStatus
	if cachedRepo.Status.ObservedGeneration > 0 {
		logger.InfoContext(ctx, "handle repository update")
		status, err = repo.OnUpdate(ctx, logger)
		if err != nil {
			return fmt.Errorf("handle repository update: %w", err)
		}
	} else {
		logger.InfoContext(ctx, "handle repository init")
		status, err = repo.OnCreate(ctx, logger)
		if err != nil {
			return fmt.Errorf("handle repository create: %w", err)
		}
	}

	if status == nil {
		status = cachedRepo.Status.DeepCopy()
	}
	status.ObservedGeneration = cachedRepo.Generation

	if time.Since(time.UnixMilli(cachedRepo.Status.Health.Checked)) < 200*time.Millisecond {
		logger.InfoContext(ctx, "skipping health check as it was recently checked")
	} else {
		res, err := rc.tester.TestRepository(ctx, repo)
		if err != nil {
			res = &provisioning.TestResults{
				Success: false,
				Errors: []string{
					"error running test repository",
					err.Error(),
				},
			}
		}

		if !res.Success {
			logger.ErrorContext(ctx, "repository is unhealthy", "errors", res.Errors)
		}

		status.Health = provisioning.HealthStatus{
			Checked: status.Health.Checked,
			Healthy: res.Success,
			Message: res.Errors,
		}
	}

	cfg := cachedRepo.DeepCopy()
	cfg.Status = *status
	if _, err := rc.client.Repositories(cachedRepo.GetNamespace()).
		UpdateStatus(ctx, cfg, v1.UpdateOptions{}); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	job, err := rc.jobs.Add(ctx, &provisioning.Job{
		ObjectMeta: v1.ObjectMeta{
			Namespace: cachedRepo.Namespace,
			Labels: map[string]string{
				"repository": cachedRepo.Name,
			},
		},
		Spec: provisioning.JobSpec{
			Action: provisioning.JobActionSync,
		},
	})
	if err != nil {
		return fmt.Errorf("trigger sync job: %w", err)
	}
	logger.InfoContext(ctx, "sync job triggered", "job", job.Name)

	return nil
}
