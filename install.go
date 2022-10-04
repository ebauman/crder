package crder

import (
	"context"
	"fmt"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func InstallUpdateCRDs(config *rest.Config, crds ...CRD) error {
	cli, err := clientset.NewForConfig(config)
	if err != nil {
		return err
	}

	var convertedCrds = make([]apiextv1.CustomResourceDefinition, len(crds))
	for _, c := range crds {
		converted, err := c.ToV1CustomResourceDefinition()
		if err != nil {
			return err
		}

		convertedCrds = append(convertedCrds, *converted)
	}

	for _, c := range convertedCrds {
		uc, err := cli.ApiextensionsV1().CustomResourceDefinitions().Get(context.Background(), c.Name, metav1.GetOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}

			// not found, thus requires installation
			_, err := cli.ApiextensionsV1().CustomResourceDefinitions().Create(context.Background(), &c, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("error installing crd: %s", err.Error())
			}
		} else {
			// object found, update it in case new version
			updateable := uc.DeepCopy()
			updateable.Spec = c.Spec

			_, err = cli.ApiextensionsV1().CustomResourceDefinitions().Update(context.Background(), updateable, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("error updating crd: %s", err.Error())
			}
		}
	}

	return nil
}
