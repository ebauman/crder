package crder

import (
	"context"
	v13 "k8s.io/api/admissionregistration/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

	c.AddValidation("foo.example.org", func(vv *Validation) {
		vv.
			MatchPolicyExact().
			WithVersions("v1").
			WithService(v13.ServiceReference{
				Namespace: "default",
				Name:      "service",
				Path:      pointer.String("/foo.example.org"),
				Port:      pointer.Int32(6443),
			}).
			AddRules(v13.RuleWithOperations{
				Operations: []v13.OperationType{v13.Create},
				Rule: v13.Rule{
					APIGroups:   []string{"foo.example.org"},
					APIVersions: []string{"v1"},
					Resources:   []string{"foo"},
				},
			})
	})

	cfg, err := clientcmd.BuildConfigFromFlags("", "/tmp/k3s-default")
	if err != nil {
		t.Error(err)
	}

	objs, err := InstallUpdateCRDsWithRecordedObjects(cfg, *c)
	if err != nil {
		t.Error(err)
	}

	defer cleanup(t, cfg, objs)
}

func cleanup(t *testing.T, cfg *rest.Config, objs []client.Object) {
	cli, err := client.New(cfg, client.Options{Scheme: scheme}) // scheme defined in install.go
	if err != nil {
		t.Errorf("error during cleanup: %s", err.Error())
	}

	for _, o := range objs {
		err = cli.Delete(context.Background(), o)
		if err != nil {
			t.Errorf("error during cleanup: %s", err.Error())
		}
	}
}
