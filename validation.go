package crder

import (
	v1 "k8s.io/api/admissionregistration/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Validation struct {
	name              string
	service           v1.ServiceReference
	caBundle          string
	url               string
	rules             []v1.RuleWithOperations
	namespaceSelector v12.LabelSelector
	objectSelector    v12.LabelSelector
	matchPolicy       v1.MatchPolicyType
	versions          []string
	sideEffect        v1.SideEffectClass
}

func (v *Validation) defaults() {
	v.SideEffectNone()
}

type validationCustomizer func(vv *Validation)

func (vv *Validation) MatchPolicyExact() *Validation {
	vv.matchPolicy = v1.Exact

	return vv
}

func (vv *Validation) MatchPolicyEquivalent() *Validation {
	vv.matchPolicy = v1.Equivalent

	return vv
}

func (vv *Validation) SideEffectNone() *Validation {
	vv.sideEffect = v1.SideEffectClassNone

	return vv
}

func (vv *Validation) SideEffectNoneOnDryRun() *Validation {
	vv.sideEffect = v1.SideEffectClassNoneOnDryRun

	return vv
}

func (vv *Validation) WithVersions(versions ...string) *Validation {
	vv.versions = versions

	return vv
}

func (vv *Validation) WithCABundle(bundle string) *Validation {
	vv.caBundle = bundle

	return vv
}

func (vv *Validation) WithService(service v1.ServiceReference) *Validation {
	vv.service = service

	return vv
}

func (vv *Validation) WithURL(url string) *Validation {
	vv.url = url

	return vv
}

func (vv *Validation) AddRules(rules ...v1.RuleWithOperations) *Validation {
	if len(vv.rules) == 0 {
		vv.rules = []v1.RuleWithOperations{}
	}

	vv.rules = append(vv.rules, rules...)

	return vv
}

func (vv *Validation) SetNamespaceSelector(selector v12.LabelSelector) *Validation {
	vv.namespaceSelector = selector

	return vv
}

func (vv *Validation) SetObjectSelector(selector v12.LabelSelector) *Validation {
	vv.objectSelector = selector

	return vv
}
