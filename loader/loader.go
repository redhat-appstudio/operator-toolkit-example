package loader

import (
	"context"

	"github.com/konflux-ci/operator-toolkit-example/api/v1alpha1"
	toolkit "github.com/konflux-ci/operator-toolkit/loader"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ObjectLoader interface {
	GetBars(ctx context.Context, cli client.Client, foo *v1alpha1.Foo) ([]v1alpha1.Bar, error)
	GetFoo(ctx context.Context, cli client.Client, name, namespace string) (*v1alpha1.Foo, error)
}

type loader struct{}

func NewLoader() ObjectLoader {
	return &loader{}
}

// GetBars loads the list of Bar resources associated with the Foo resource passed as a parameter.
func (l *loader) GetBars(ctx context.Context, cli client.Client, foo *v1alpha1.Foo) ([]v1alpha1.Bar, error) {
	bars := &v1alpha1.BarList{}

	err := cli.List(ctx, bars,
		client.InNamespace(foo.Namespace),
		client.MatchingFields{"spec.foo": foo.Name})
	if err != nil {
		return nil, err
	}

	return bars.Items, nil
}

// GetFoo returns the Foo resource with the given name and namespace.
func (l *loader) GetFoo(ctx context.Context, cli client.Client, name, namespace string) (*v1alpha1.Foo, error) {
	foo := &v1alpha1.Foo{}
	return foo, toolkit.GetObject(name, namespace, cli, ctx, foo)
}
