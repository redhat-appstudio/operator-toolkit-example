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
	"fmt"
	"github.com/go-logr/logr"
	"github.com/redhat-appstudio/operator-toolkit-example/api/v1alpha1"
	"github.com/redhat-appstudio/operator-toolkit-example/loader"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Webhook describes the data structure for the bar webhook
type Webhook struct {
	client client.Client
	loader loader.ObjectLoader
	log    logr.Logger
}

// Register registers the webhook with the passed manager and log.
func (w *Webhook) Register(mgr ctrl.Manager, log *logr.Logger) error {
	w.client = mgr.GetClient()
	w.loader = loader.NewLoader()
	w.log = log.WithName("bar")

	return ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.Bar{}).
		WithDefaulter(w).
		WithValidator(w).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-appstudio-redhat-com-v1alpha1-bar,mutating=true,failurePolicy=fail,sideEffects=None,groups=appstudio.redhat.com,resources=bars,verbs=create;update,versions=v1alpha1,name=mbar.kb.io,admissionReviewVersions=v1

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (w *Webhook) Default(ctx context.Context, obj runtime.Object) error {
	return nil
}

// +kubebuilder:webhook:path=/validate-appstudio-redhat-com-v1alpha1-bar,mutating=false,failurePolicy=fail,sideEffects=None,groups=appstudio.redhat.com,resources=bars,verbs=create;update,versions=v1alpha1,name=vbar.kb.io,admissionReviewVersions=v1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (w *Webhook) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	bar := obj.(*v1alpha1.Bar)

	_, err := w.loader.GetFoo(ctx, w.client, bar.Spec.Foo, bar.Namespace)
	if err != nil {
		return fmt.Errorf("resource references an unexistent Foo resource (%s/%s)", bar.Namespace, bar.Spec.Foo)
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (w *Webhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (w *Webhook) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	return nil
}
