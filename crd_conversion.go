package crder

import apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

type CRDConversion struct {
	Webhook  bool
	Service  apiextv1.ServiceReference
	CABundle string
	URL      string
	Versions []string
}

type conversionCustomizer func(cc *CRDConversion)

func (cc *CRDConversion) StrategyNone() *CRDConversion {
	cc.Webhook = false

	return cc
}

func (cc *CRDConversion) StrategyWebhook() *CRDConversion {
	cc.Webhook = true

	return cc
}

func (cc *CRDConversion) WithCABundle(bundle string) *CRDConversion {
	cc.CABundle = bundle

	return cc
}

func (cc *CRDConversion) WithService(service apiextv1.ServiceReference) *CRDConversion {
	cc.Service = service

	return cc
}

func (cc *CRDConversion) WithURL(url string) *CRDConversion {
	cc.URL = url

	return cc
}

func (cc *CRDConversion) WithVersions(versions ...string) *CRDConversion {
	cc.Versions = versions

	return cc
}
