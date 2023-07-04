package bar

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/redhat-appstudio/operator-toolkit-example/api/v1alpha1"
	"github.com/redhat-appstudio/operator-toolkit-example/loader"
	"github.com/redhat-appstudio/operator-toolkit/controller"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Adapter holds the objects needed to reconcile a Bar resource.
type adapter struct {
	bar    *v1alpha1.Bar // this is the kind of resource this adapter reconciles
	client client.Client
	ctx    context.Context
	loader loader.ObjectLoader
	logger *logr.Logger
}

// NewAdapter creates and returns an Adapter instance.
func NewAdapter(ctx context.Context, client client.Client, bar *v1alpha1.Bar, loader loader.ObjectLoader, logger *logr.Logger) *adapter {
	return &adapter{
		bar:    bar,
		client: client,
		ctx:    ctx,
		loader: loader,
		logger: logger,
	}
}

// EnsureOwnerReferenceIsSet is an operation that will ensure that the owner reference is set.
func (a *adapter) EnsureOwnerReferenceIsSet() (controller.OperationResult, error) {
	foo, err := a.loader.GetFoo(a.ctx, a.client, a.bar.Spec.Foo, a.bar.Namespace)
	if err != nil {
		return controller.RequeueWithError(err)
	}

	patch := client.MergeFrom(a.bar.DeepCopy())
	err = ctrl.SetControllerReference(foo, a.bar, a.client.Scheme())
	if err != nil {
		return controller.RequeueWithError(err)
	}

	err = a.client.Patch(a.ctx, a.bar, patch)
	if err != nil && !errors.IsNotFound(err) {
		return controller.RequeueWithError(err)
	}

	return controller.ContinueProcessing()
}
