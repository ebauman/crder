package crder

import apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

type Conversion struct {
	Webhook  bool
	Service  apiextv1.ServiceReference
	CABundle string
	URL      string
	Versions []string
}

type conversionCustomizer func(cc *Conversion)

func (cc *Conversion) StrategyNone() *Conversion {
	cc.Webhook = false

	return cc
}

func (cc *Conversion) StrategyWebhook() *Conversion {
	cc.Webhook = true

	return cc
}

func (cc *Conversion) WithCABundle(bundle string) *Conversion {
	cc.CABundle = bundle

	return cc
}

func (cc *Conversion) WithService(service apiextv1.ServiceReference, path string) *Conversion {
	cc.Service = service

	if(path != ""){
		cc.Service.Path = path
	}

	return cc
}

func (cc *Conversion) WithURL(url string) *Conversion {
	cc.URL = url

	return cc
}

func (cc *Conversion) WithVersions(versions ...string) *Conversion {
	cc.Versions = versions

	return cc
}
