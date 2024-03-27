/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bar

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/konflux-ci/operator-toolkit-example/loader"
	"github.com/konflux-ci/operator-toolkit/controller"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/konflux-ci/operator-toolkit-example/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BarReconciler reconciles a Bar object
type Controller struct {
	client client.Client
	log    logr.Logger
}

//+kubebuilder:rbac:groups=appstudio.redhat.com,resources=bars,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=appstudio.redhat.com,resources=bars/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=appstudio.redhat.com,resources=bars/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (c *Controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := c.log.WithValues("Bar", req.NamespacedName)

	bar := &v1alpha1.Bar{}
	err := c.client.Get(ctx, req.NamespacedName, bar)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	adapter := NewAdapter(ctx, c.client, bar, loader.NewLoader(), &logger)

	return controller.ReconcileHandler([]controller.Operation{
		adapter.EnsureOwnerReferenceIsSet,
	})
}

// Register registers the controller with the passed manager and log.
func (c *Controller) Register(mgr ctrl.Manager, log *logr.Logger, _ cluster.Cluster) error {
	c.client = mgr.GetClient()
	c.log = log.WithName("bar")

	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Bar{}).
		Complete(c)
}
