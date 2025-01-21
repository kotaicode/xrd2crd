package fileio

import (
	"encoding/json"
	"fmt"
	"os"

	crossplanev1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

// LoadXRDFromFile reads and parses a CompositeResourceDefinition from a YAML file.
func LoadXRDFromFile(filePath string) (*crossplanev1.CompositeResourceDefinition, error) {
	// Read the file content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal the YAML content into an XRD object
	var xrd crossplanev1.CompositeResourceDefinition
	if err := yaml.Unmarshal(fileContent, &xrd); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XRD: %w", err)
	}

	return &xrd, nil
}

// WriteToFile writes the CRD content to a file in either YAML or JSON format
func WriteToFile(crd *apiextensionsv1.CustomResourceDefinition, filePath string, asJSON bool) error {
	var content []byte
	var err error

	if asJSON {
		content, err = json.MarshalIndent(crd, "", "  ")
	} else {
		content, err = yaml.Marshal(crd)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal CRD: %w", err)
	}

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// FormatOutput formats the CRD as either YAML or JSON string
func FormatOutput(crd *apiextensionsv1.CustomResourceDefinition, asJSON bool) (string, error) {
	var content []byte
	var err error

	if asJSON {
		content, err = json.MarshalIndent(crd, "", "  ")
	} else {
		content, err = yaml.Marshal(crd)
	}
	if err != nil {
		return "", fmt.Errorf("failed to marshal CRD: %w", err)
	}

	return string(content), nil
}
