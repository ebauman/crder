package crder

import (
	"context"
	"fmt"
	v1 "k8s.io/api/admissionregistration/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = apiextv1.AddToScheme(scheme)
	_ = metav1.AddMetaToScheme(scheme)
	_ = v1.AddToScheme(scheme)
}

// InstallUpdateCRDs used to install and update CRDs (also validatingwebhookconfigurations)
func InstallUpdateCRDs(config *rest.Config, crds ...CRD) error {
	_, e := InstallUpdateCRDsWithRecordedObjects(config, crds...)

	return e
}

// InstallUpdateCRDsWithRecordedObjects same as InstallUpdateCRDs except this returns objects that were created or updated
// Mostly useful for tests where the test should clean up after itself
func InstallUpdateCRDsWithRecordedObjects(config *rest.Config, crds ...CRD) ([]client.Object, error) {
	cli, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, fmt.Errorf("error building kclient: %s", err.Error())
	}

	var created = []client.Object{}

	for _, c := range crds {
		converted, err := c.ToV1CustomResourceDefinition()
		if err != nil {
			return nil, fmt.Errorf("error converting CRD: %s", err.Error())
		}
		var uc = &apiextv1.CustomResourceDefinition{}
		err = cli.Get(context.Background(), client.ObjectKey{Name: converted.Name}, uc)
		if err != nil {
			if !errors.IsNotFound(err) {
				return nil, fmt.Errorf("error getting apiextv1.CRD from k: %s", err.Error())
			}

			// not found, thus requires installation
			err = cli.Create(context.Background(), converted)
			if err != nil {
				return nil, fmt.Errorf("error installing crd: %s", err.Error())
			}
			created = append(created, converted)
		} else {
			// object found, update it in case new version
			updateable := uc.DeepCopy()
			updateable.Spec = converted.Spec

			err = cli.Update(context.Background(), updateable)
			if err != nil {
				return nil, fmt.Errorf("error updating crd: %s", err.Error())
			}
			created = append(created, updateable)
		}

		err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (done bool, err error) {
			var crd = &apiextv1.CustomResourceDefinition{}
			err = cli.Get(context.Background(), client.ObjectKey{Name: converted.Name}, crd)
			if err != nil {
				return false, fmt.Errorf("errro getting apiextv1.CRD from k: %s", err.Error())
			}

			for _, cond := range crd.Status.Conditions {
				switch cond.Type {
				case apiextv1.Established:
					if cond.Status == apiextv1.ConditionTrue {
						return true, err
					}
				case apiextv1.NamesAccepted:
					if cond.Status == apiextv1.ConditionFalse {
						return true, fmt.Errorf("name conflict on %s: %v", converted.Name, cond.Reason)
					}
				}
			}

			return false, nil
		})
		if err != nil {
			return nil, fmt.Errorf("error waiting for crd readiness: %v", err.Error())
		}

		// install validations
		if len(c.validation) > 0 {
			validations, err := c.GetValidatingWebhooks()
			if err != nil {
				return nil, fmt.Errorf("error getting validating webhooks: %s", err.Error())
			}

			for _, v := range *validations {
				// check if this exists first
				var val = &v1.ValidatingWebhookConfiguration{}
				err = cli.Get(context.Background(), client.ObjectKey{Name: v.Name}, val)

				if errors.IsNotFound(err) {
					err = cli.Create(context.Background(), &v)
					if err != nil {
						return nil, fmt.Errorf("error creating validatingwebhookconfiguration: %s", err.Error())
					}
					created = append(created, &v)
				} else if err != nil {
					return nil, err
				} else {
					val = val.DeepCopy()
					val.Webhooks = v.Webhooks

					err = cli.Update(context.Background(), val)
					if err != nil {
						return nil, fmt.Errorf("error updating validatingwebhookconfiguration: %s", err.Error())
					}
					created = append(created, &v)
				}
			}
		}
	}

	return created, nil
}
