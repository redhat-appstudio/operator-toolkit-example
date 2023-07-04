package foo

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/redhat-appstudio/operator-toolkit-example/api/v1alpha1"
	"github.com/redhat-appstudio/operator-toolkit-example/loader"
	"github.com/redhat-appstudio/operator-toolkit/controller"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// Adapter holds the objects needed to reconcile a Foo resource.
type adapter struct {
	client client.Client
	ctx    context.Context
	foo    *v1alpha1.Foo
	loader loader.ObjectLoader
	logger *logr.Logger
}

// NewAdapter creates and returns an Adapter instance.
func NewAdapter(ctx context.Context, client client.Client, foo *v1alpha1.Foo, loader loader.ObjectLoader, logger *logr.Logger) *adapter {
	return &adapter{
		client: client,
		ctx:    ctx,
		foo:    foo,
		loader: loader,
		logger: logger,
	}
}

// finalizerName is the finalizer name to be added to the Foo resource
const finalizerName string = "appstudio.redhat.com/finalizer"

// EnsureFinalizersAreCalled is an operation that will ensure that finalizers are called whenever the Foo resource being
// processed is marked for deletion. Once finalizers get called, the finalizer will be removed and the Foo resource will go
// back to the queue, so it gets deleted. If a finalizer function fails its execution or a finalizer fails to be removed,
// the Foo resource will be requeued with the error attached.
func (a *adapter) EnsureFinalizersAreCalled() (controller.OperationResult, error) {
	// Check if the Foo resource is marked for deletion and continue processing other operations otherwise
	if a.foo.GetDeletionTimestamp() == nil {
		return controller.ContinueProcessing()
	}

	if controllerutil.ContainsFinalizer(a.foo, finalizerName) {
		if err := a.finalizeResource(); err != nil {
			return controller.RequeueWithError(err)
		}

		patch := client.MergeFrom(a.foo.DeepCopy())
		controllerutil.RemoveFinalizer(a.foo, finalizerName)
		err := a.client.Patch(a.ctx, a.foo, patch)
		if err != nil {
			return controller.RequeueWithError(err)
		}
	}

	// Requeue the Foo resource again so it gets deleted and other operations are not executed
	return controller.Requeue()
}

// EnsureFinalizerIsAdded is an operation that will ensure that the Foo resource being processed contains a finalizer.
func (a *adapter) EnsureFinalizerIsAdded() (controller.OperationResult, error) {
	var finalizerFound bool
	for _, finalizer := range a.foo.GetFinalizers() {
		if finalizer == finalizerName {
			finalizerFound = true
		}
	}

	if !finalizerFound {
		a.logger.Info("Adding Finalizer to the Foo resource")
		patch := client.MergeFrom(a.foo.DeepCopy())
		controllerutil.AddFinalizer(a.foo, finalizerName)
		err := a.client.Patch(a.ctx, a.foo, patch)

		return controller.RequeueOnErrorOrContinue(err)
	}

	return controller.ContinueProcessing()
}

// EnsureMaximumReplicas is an operation that will ensure that the number of replicas for this resource doesn't go beyond
// the desired number of replicas, deleting Bar resources if needed.
func (a *adapter) EnsureMaximumReplicas() (controller.OperationResult, error) {
	replicas, err := a.loader.GetBars(a.ctx, a.client, a.foo)
	if err != nil {
		return controller.RequeueWithError(err)
	}

	replicasDelta := a.foo.Spec.DesiredReplicas - len(replicas)

	if replicasDelta >= 0 {
		return controller.ContinueProcessing()
	}

	for _, replica := range replicas[a.foo.Spec.DesiredReplicas:] {
		err := a.client.Delete(a.ctx, &replica)
		if err != nil && errors.IsNotFound(err) {
			return controller.RequeueWithError(err)
		}
		a.logger.Info("Bar deleted", "Bar.Name", replica.Name, "Bar.Namespace", replica.Namespace)
	}

	return controller.ContinueProcessing()
}

// EnsureMinimumReplicas is an operation that will ensure that the number of replicas for this resource doesn't go below
// the desired number of replicas, creating Bar resources if needed.
func (a *adapter) EnsureMinimumReplicas() (controller.OperationResult, error) {
	replicas, err := a.loader.GetBars(a.ctx, a.client, a.foo)
	if err != nil {
		return controller.RequeueWithError(err)
	}

	replicasDelta := a.foo.Spec.DesiredReplicas - len(replicas)

	if replicasDelta < 0 {
		return controller.ContinueProcessing()
	}

	for i := 0; i < replicasDelta; i++ {
		replica := &v1alpha1.Bar{
			ObjectMeta: v1.ObjectMeta{
				GenerateName: a.foo.Name + "-",
				Namespace:    a.foo.Namespace,
			},
			Spec: v1alpha1.BarSpec{
				Foo: a.foo.Name,
			},
		}

		err := a.client.Create(a.ctx, replica)
		if err != nil {
			return controller.RequeueWithError(err)
		}
		a.logger.Info("Bar created", "Bar.Name", replica.Name, "Bar.Namespace", replica.Namespace)
	}

	return controller.ContinueProcessing()
}

// EnsureReplicaDataConsistency is an operation that will ensure that the list of replicas in the Foo resource's status
// is kept up to date. It will also update the health condition type when needed.
func (a *adapter) EnsureReplicaDataConsistency() (controller.OperationResult, error) {
	replicas, err := a.loader.GetBars(a.ctx, a.client, a.foo)
	if err != nil {
		return controller.RequeueWithError(err)
	}

	patch := client.MergeFrom(a.foo.DeepCopy())
	a.foo.Status.Replicas = []string{}

	for _, replica := range replicas {
		a.foo.Status.Replicas = append(a.foo.Status.Replicas, replica.Name)
	}

	replicasDelta := len(replicas) - a.foo.Spec.DesiredReplicas
	if replicasDelta == 0 {
		a.foo.MarkHealthy(v1alpha1.HealthyReason)
	} else if replicasDelta > 0 {
		a.foo.MarkHealthy(v1alpha1.TooManyReplicasReason)
	} else {
		a.foo.MarkUnhealthy()
	}

	return controller.RequeueOnErrorOrContinue(a.client.Status().Patch(a.ctx, a.foo, patch))
}

// finalizeResource deletes all the Bar resources associated with this resource.
func (a *adapter) finalizeResource() error {
	bars, err := a.loader.GetBars(a.ctx, a.client, a.foo)
	if err != nil {
		return err
	}

	for _, bar := range bars {
		err = a.client.Delete(a.ctx, &bar)
		if err != nil && !errors.IsNotFound(err) {
			return err
		}

		a.logger.Info("Deleted bar", "bar.Name", bar.Name, "bar.Namespace", bar.Namespace)
	}

	a.logger.Info("Successfully finalized Foo")

	return nil
}
