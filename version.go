package crder

import (
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/utils/pointer"
)

type Version struct {
	columns            []apiextv1.CustomResourceColumnDefinition
	served             bool
	stored             bool
	parent             *CRD
	object             interface{}
	version            string
	deprecated         bool
	deprecationMessage string
	scale              *apiextv1.CustomResourceSubresourceScale
	status             bool
	preserveUnknown    bool
}

type versionCustomizer func(cv *Version)

func (cv *Version) WithColumn(name string, jsonPath string) *Version {
	col := apiextv1.CustomResourceColumnDefinition{
		Name:     name,
		JSONPath: jsonPath,
		Type:     "string",
		Priority: 0,
	}

	cv.columns = append(cv.columns, col)
	return cv
}

func (cv *Version) WithCRDColumns(cols ...apiextv1.CustomResourceColumnDefinition) *Version {
	cv.columns = append(cv.columns, cols...)
	return cv
}

func (cv *Version) IsServed(served bool) *Version {
	cv.served = served
	return cv
}

func (cv *Version) IsStored(stored bool) *Version {
	cv.stored = stored

	return cv
}

func (cv *Version) IsDeprecated(deprecationWarning string) *Version {
	cv.deprecated = true
	cv.deprecationMessage = deprecationWarning

	return cv
}

func (cv *Version) WithObject(obj interface{}) *Version {
	cv.object = obj

	return cv
}

func (cv *Version) WithScale(labelSelectorPath string, specReplicasPath string, statusReplicaPath string) *Version {
	cv.scale = &apiextv1.CustomResourceSubresourceScale{
		SpecReplicasPath:   specReplicasPath,
		StatusReplicasPath: statusReplicaPath,
		LabelSelectorPath:  pointer.String(labelSelectorPath),
	}

	return cv
}

func (cv *Version) WithStatus() *Version {
	cv.status = true
	return cv
}

func (cv *Version) WithPreserveUnknown() *Version {
	cv.preserveUnknown = true
	return cv
}
