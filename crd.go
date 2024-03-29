package crder

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
)

type CRD struct {
	// default false
	preserveUnknown bool

	versions []Version

	gvk schema.GroupVersionKind

	object interface{}

	conversion *Conversion
	validation []*Validation

	namespaced bool

	singularName string
	pluralName   string
	shortNames   []string

	categories []string
}

func NewCRD(obj interface{}, group string, customize func(c *CRD)) *CRD {
	c := &CRD{
		object: obj,
		gvk: schema.GroupVersionKind{
			Kind:  getType(obj).Name(),
			Group: group,
		},
		versions: []Version{},
	}

	if customize != nil {
		customize(c)
	}

	return c
}

// WithPreserveUnknown sets preserveUnknown to true
func (c *CRD) WithPreserveUnknown() {
	c.preserveUnknown = true
}

func (c *CRD) OverrideGVK(group string, version string, kind string) *CRD {
	c.gvk = schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}

	return c
}

func (c *CRD) AddVersion(version string, object interface{}, customize versionCustomizer) *CRD {
	v := Version{
		version: version,
		object:  object,
		served:  true,
		stored:  true,
	}

	if customize != nil {
		customize(&v)
	}

	c.versions = append(c.versions, v)

	return c
}

func (c *CRD) WithConversion(customizer conversionCustomizer) *CRD {
	conv := &Conversion{}

	customizer(conv)

	c.conversion = conv

	return c
}

func (c *CRD) AddValidation(name string, customizer validationCustomizer) *CRD {
	v := &Validation{
		name: name,
	}

	v.defaults()

	customizer(v)

	c.validation = append(c.validation, v)

	return c
}

func (c *CRD) IsNamespaced(namespaced bool) *CRD {
	c.namespaced = namespaced

	return c
}

func (c *CRD) WithNames(singular string, plural string) *CRD {
	c.singularName = singular
	c.pluralName = plural

	return c
}

func (c *CRD) WithShortNames(names ...string) *CRD {
	c.shortNames = names

	return c
}

func (c *CRD) WithCategories(categories ...string) *CRD {
	c.categories = categories

	return c
}

func getType(obj interface{}) reflect.Type {
	if t, ok := obj.(reflect.Type); ok {
		return t
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
