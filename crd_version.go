package crder

import (
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/pointer"
)

type CRDVersion struct {
	columns            []apiextv1.CustomResourceColumnDefinition
	served             bool
	stored             bool
	parent             *CRD
	object             HasGVK
	version            string
	deprecated         bool
	deprecationMessage string
	scale              *apiextv1.CustomResourceSubresourceScale
	status             bool
}

type versionCustomizer func(cv *CRDVersion)

func (cv *CRDVersion) WithColumn(name string, jsonPath string) *CRDVersion {
	col := apiextv1.CustomResourceColumnDefinition{
		Name:     name,
		JSONPath: jsonPath,
		Type:     "string",
		Priority: 0,
	}

	cv.columns = append(cv.columns, col)
	return cv
}

func (cv *CRDVersion) WithCRDColumns(cols ...apiextv1.CustomResourceColumnDefinition) *CRDVersion {
	cv.columns = append(cv.columns, cols...)
	return cv
}

func (cv *CRDVersion) IsServed(served bool) *CRDVersion {
	cv.served = served
	return cv
}

func (cv *CRDVersion) IsStored(stored bool) *CRDVersion {
	cv.stored = stored

	return cv
}

func (cv *CRDVersion) IsDeprecated(deprecationWarning string) *CRDVersion {
	cv.deprecated = true
	cv.deprecationMessage = deprecationWarning

	return cv
}

func (cv *CRDVersion) WithObject(obj HasGVK) *CRDVersion {
	cv.object = obj

	return cv
}

func (cv *CRDVersion) WithScale(labelSelectorPath string, specReplicasPath string, statusReplicaPath string) *CRDVersion {
	cv.scale = &apiextv1.CustomResourceSubresourceScale{
		SpecReplicasPath:   specReplicasPath,
		StatusReplicasPath: statusReplicaPath,
		LabelSelectorPath:  pointer.String(labelSelectorPath),
	}

	return cv
}

func (cv *CRDVersion) WithStatus() *CRDVersion {
	cv.status = true
	return cv
}
