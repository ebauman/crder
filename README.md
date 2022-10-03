# CRDer ("see-are-dee-er')

Helps you create CustomResourceDefinitions and install them into a cluster.

## Example

```go
package main

import (
	"github.com/ebauman/crder"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// usually present in pkg/apis/v1/example.org/types.g

type Foo struct {
	v1.TypeMeta
	Metadata v1.ObjectMeta

	Spec FooSpec
}

type FooSpec struct {
	Bar bool
}

func main() {
	c := NewCRD(Foo{})
	
	c.AddVersion("v1", Foo{}, func(cv *CRDVersion) {
		cv.
			IsServed(true).
			IsStored(true)
	})

	cfg, err := clientcmd.BuildConfigFromFlags("", "/path/to/kube/config/file")
	if err != nil {
		log.Fatal(err)
	}

	installable, err := c.ToV1CustomResourceDefinition()
	if err != nil {
		log.Fatal(err)
	}

	err = InstallUpdateCRDs(cfg, *installable)
	if err != nil {
		log.Fatal(err)
	}
}
```