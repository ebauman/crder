package crder

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
)

type Foo struct {
	v1.TypeMeta
	Metadata v1.ObjectMeta

	Spec FooSpec
}

type FooSpec struct {
	Bar bool
}

type Foo2 struct {
	v1.TypeMeta
	Metadata v1.ObjectMeta

	Spec Foo2Spec
}

type Foo2Spec struct {
	Bar bool
	Baz bool
}

func (f *Foo) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "example.org",
		Version: "v1",
		Kind:    "Foo",
	}
}

func (f *Foo2) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "example.org",
		Version: "v2",
		Kind:    "Foo",
	}
}

func Test_CRD(t *testing.T) {
	c := NewCRD(Foo{}, "example.org", nil)

	c.AddVersion("v1", Foo{}, func(cv *Version) {
		cv.
			IsServed(true).
			IsStored(true)
	})

	_, err := c.ToV1CustomResourceDefinition()
	if err != nil {
		t.Error(err)
	}
}

func Test_Install(t *testing.T) {
	c := NewCRD(Foo{}, "example.org", nil)

	c.AddVersion("v1", Foo{}, func(cv *Version) {
		cv.
			IsServed(true).
			IsStored(false)
	})

	c.AddVersion("v2", Foo2{}, func(cv *Version) {
		cv.
			IsServed(true).
			IsStored(true)

		cv.
			WithColumn("baz", ".spec.baz")
	})

	cfg, err := clientcmd.BuildConfigFromFlags("", "/tmp/k3s-default")
	if err != nil {
		t.Error(err)
	}

	err = InstallUpdateCRDs(cfg, *c)
	if err != nil {
		t.Error(err)
	}
}
