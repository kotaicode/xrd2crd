package converter

import (
	"encoding/json"
	"fmt"

	crossplanev1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConvertXRDToCRD converts a Crossplane CompositeResourceDefinition (XRD) to a Kubernetes CustomResourceDefinition (CRD).
// It specifically generates the claim CRD based on the XRD's claim names and specifications.
func ConvertXRDToCRD(xrd *crossplanev1.CompositeResourceDefinition) (*apiextensionsv1.CustomResourceDefinition, error) {
	// Validate input
	if xrd == nil {
		return nil, fmt.Errorf("input CompositeResourceDefinition is nil")
	}

	// Map the XRD spec to CRD spec
	crd := &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s.%s", xrd.Spec.ClaimNames.Plural, xrd.Spec.Group),
			Labels:      xrd.Labels,
			Annotations: xrd.Annotations,
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: xrd.Spec.Group,
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Categories: []string{"claim"},
				Plural:     xrd.Spec.ClaimNames.Plural,
				Kind:       xrd.Spec.ClaimNames.Kind,
				Singular:   xrd.Spec.ClaimNames.Singular,
			},
			Scope:    "Namespaced",
			Versions: convertXRDVersionsToCRDVersions(xrd.Spec.Versions),
		},
	}

	crd.Kind = "CustomResourceDefinition"
	crd.APIVersion = apiextensionsv1.SchemeGroupVersion.String()
	return crd, nil
}

// convertXRDVersionsToCRDVersions converts the version specifications from an XRD format to CRD format.
// It handles the conversion of schema definitions and version-specific attributes.
func convertXRDVersionsToCRDVersions(xrdVersions []crossplanev1.CompositeResourceDefinitionVersion) []apiextensionsv1.CustomResourceDefinitionVersion {
	var crdVersions []apiextensionsv1.CustomResourceDefinitionVersion
	for _, version := range xrdVersions {
		crdVersion := apiextensionsv1.CustomResourceDefinitionVersion{
			Name:    version.Name,
			Served:  version.Served,
			Storage: version.Referenceable,
		}

		if version.Schema != nil {
			var schema apiextensionsv1.JSONSchemaProps
			raw, err := version.Schema.OpenAPIV3Schema.MarshalJSON()
			if err != nil {
				fmt.Printf("error marshalling %v", err)
				continue
			}
			if err := json.Unmarshal(raw, &schema); err != nil {
				fmt.Printf("error unmarshalling %v", err)
				continue
			}
			crdVersion.Schema = &apiextensionsv1.CustomResourceValidation{
				OpenAPIV3Schema: &schema,
			}
		}
		crdVersion.Schema.OpenAPIV3Schema.Required = []string{"spec", "apiVersion", "kind", "metadata"}
		crdVersions = append(crdVersions, crdVersion)
	}
	return crdVersions
}
