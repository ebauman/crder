package crder

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type HasGVK interface {
	GroupVersionKind() schema.GroupVersionKind
}

type CRD struct {
	// default false
	preserveUnknown bool

	versions []CRDVersion

	// determined by object used when creating CRD
	gvk schema.GroupVersionKind

	object any

	conversion *CRDConversion

	namespaced bool

	singularName string
	pluralName   string
	shortNames   []string

	categories []string
}

func NewCRD(obj HasGVK, customize func(c *CRD) *CRD) {
	c := &CRD{
		object:   obj,
		gvk:      obj.GroupVersionKind(),
		versions: []CRDVersion{},
	}

	if customize != nil {
		customize(c)
	}
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

func (c *CRD) AddVersion(version string, object HasGVK, customize versionCustomizer) *CRD {
	v := CRDVersion{
		version: version,
		object:  object,
	}

	if customize != nil {
		customize(&v)
	}

	c.versions = append(c.versions, v)

	return c
}

func (c *CRD) WithConversion(customizer conversionCustomizer) *CRD {
	conv := &CRDConversion{}

	customizer(conv)

	c.conversion = conv

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
