package loader

import (
	"context"
	"github.com/redhat-appstudio/operator-toolkit-example/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
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

// getObject loads an object from the cluster. This is a generic function that requires the object to be passed as an
// argument. The object is modified during the invocation.
func getObject(name, namespace string, cli client.Client, ctx context.Context, object client.Object) error {
	return cli.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, object)
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
	return foo, getObject(name, namespace, cli, ctx, foo)
}
