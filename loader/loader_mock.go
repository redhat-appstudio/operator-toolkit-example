package loader

import (
	"context"
	"github.com/redhat-appstudio/operator-toolkit-example/api/v1alpha1"
	toolkit "github.com/redhat-appstudio/operator-toolkit/loader"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	FooContextKey  toolkit.ContextKey = iota
	BarsContextKey toolkit.ContextKey = iota
)

type mockLoader struct {
	loader ObjectLoader
}

func NewMockLoader() ObjectLoader {
	return &mockLoader{
		loader: NewLoader(),
	}
}

// GetBars returns the resource and error passed as values of the context.
func (l *mockLoader) GetBars(ctx context.Context, cli client.Client, foo *v1alpha1.Foo) ([]v1alpha1.Bar, error) {
	if ctx.Value(BarsContextKey) == nil {
		return l.loader.GetBars(ctx, cli, foo)
	}
	return toolkit.GetMockedResourceAndErrorFromContext(ctx, BarsContextKey, []v1alpha1.Bar{})
}

// GetFoo returns the resource and error passed as values of the context.
func (l *mockLoader) GetFoo(ctx context.Context, cli client.Client, name, namespace string) (*v1alpha1.Foo, error) {
	if ctx.Value(FooContextKey) == nil {
		return l.loader.GetFoo(ctx, cli, name, namespace)
	}
	return toolkit.GetMockedResourceAndErrorFromContext(ctx, FooContextKey, &v1alpha1.Foo{})
}
