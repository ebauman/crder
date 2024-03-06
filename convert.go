package crder

import (
	"errors"
	"fmt"
	"github.com/rancher/wrangler/pkg/schemas/openapi"
	v12 "k8s.io/api/admissionregistration/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"strings"
)

func (c CRD) ToV1CustomResourceDefinition() (*apiextv1.CustomResourceDefinition, error) {
	if len(c.versions) == 0 {
		return nil, errors.New("must define at least one version")
	}
	var scope = apiextv1.ClusterScoped

	if c.namespaced {
		scope = apiextv1.NamespaceScoped
	}

	var singular, plural string
	singular, plural = c.resolveNames()

	out := apiextv1.CustomResourceDefinition{
		TypeMeta: v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{
			Name: plural + "." + c.gvk.Group,
		},
		Spec: apiextv1.CustomResourceDefinitionSpec{
			Group: c.gvk.Group,
			Names: apiextv1.CustomResourceDefinitionNames{
				Plural:     plural,
				Singular:   singular,
				ShortNames: c.shortNames,
				Kind:       c.gvk.Kind,
				ListKind:   "",
				Categories: c.categories,
			},
			Scope:                 scope,
			PreserveUnknownFields: c.preserveUnknown,
		},
	}

	if c.conversion != nil {
		out.Spec.Conversion = &apiextv1.CustomResourceConversion{
			Strategy: func() apiextv1.ConversionStrategyType {
				if c.conversion.Webhook {
					return apiextv1.WebhookConverter
				}

				return apiextv1.NoneConverter
			}(),
			Webhook: func() *apiextv1.WebhookConversion {
				if c.conversion.Webhook {
					return &apiextv1.WebhookConversion{
						ClientConfig: &apiextv1.WebhookClientConfig{
							URL: func() *string {
								if c.conversion.URL == "" {
									return nil
								}
								return pointer.String(c.conversion.URL)
							}(),
							Service:  &c.conversion.Service,
							CABundle: []byte(c.conversion.CABundle),
						},
						ConversionReviewVersions: c.conversion.Versions,
					}
				}

				return nil
			}(),
		}
	}

	for _, cv := range c.versions {
		ver, err := cv.ToV1CustomResourceDefinitionVersion()
		if err != nil {
			return nil, err
		}

		out.Spec.Versions = append(out.Spec.Versions, *ver)
	}

	return &out, nil
}

func (cv Version) ToV1CustomResourceDefinitionVersion() (*apiextv1.CustomResourceDefinitionVersion, error) {
	schema, err := openapi.ToOpenAPIFromStruct(cv.object)
	if err != nil {
		return nil, err
	}

	if cv.preserveUnknown {
		schema.XPreserveUnknownFields = pointer.Bool(true)
	}

	out := apiextv1.CustomResourceDefinitionVersion{
		Name:    cv.version,
		Served:  cv.served,
		Storage: cv.stored,
		Schema: &apiextv1.CustomResourceValidation{
			OpenAPIV3Schema: schema,
		},
		Subresources: &apiextv1.CustomResourceSubresources{
			Scale: cv.scale,
		},
		AdditionalPrinterColumns: cv.columns,
	}

	if cv.status {
		out.Subresources.Status = &apiextv1.CustomResourceSubresourceStatus{}
	}

	if cv.deprecated {
		out.Deprecated = true
		out.DeprecationWarning = pointer.String(cv.deprecationMessage)
	}

	return &out, nil
}

func (c *CRD) GetValidatingWebhooks() (*[]v12.ValidatingWebhookConfiguration, error) {
	if len(c.validation) == 0 {
		return nil, fmt.Errorf("no validations in this crd")
	}

	vwcs := make([]v12.ValidatingWebhookConfiguration, len(c.validation))

	for i, v := range c.validation {
		out := v12.ValidatingWebhookConfiguration{
			TypeMeta: v1.TypeMeta{},
			ObjectMeta: v1.ObjectMeta{
				Name: v.name,
			},
			Webhooks: []v12.ValidatingWebhook{
				{
					Name: v.name,
					ClientConfig: v12.WebhookClientConfig{
						URL: func() *string {
							if v.url != "" {
								return pointer.String(v.url)
							}
							return nil
						}(),
						Service:  &v.service,
						CABundle: []byte(v.caBundle),
					},
					Rules:                   v.rules,
					MatchPolicy:             &v.matchPolicy,
					NamespaceSelector:       &v.namespaceSelector,
					ObjectSelector:          &v.objectSelector,
					AdmissionReviewVersions: v.versions,
					SideEffects:             &v.sideEffect,
				},
			},
		}
		vwcs[i] = out
	}

	return &vwcs, nil
}

func (c *CRD) resolveNames() (singular string, plural string) {
	if c.singularName != "" {
		singular = c.singularName
	}

	if c.pluralName != "" {
		plural = c.pluralName
	}

	if c.singularName == "" {
		singular = strings.ToLower(c.gvk.Kind)
	}

	if c.pluralName == "" {
		plural = strings.ToLower(c.gvk.Kind + "s")
	}

	return
}
