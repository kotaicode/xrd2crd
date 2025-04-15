package fileio

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	crossplanev1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

// LoadXRDFromFile reads and parses one or more CompositeResourceDefinitions from a YAML file.
// The file can contain multiple YAML documents separated by '---'.
func LoadXRDFromFile(filePath string) ([]*crossplanev1.CompositeResourceDefinition, error) {
	// Read the file content
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Split the content by '---' to handle multiple documents
	documents := strings.Split(string(fileContent), "---")
	var xrds []*crossplanev1.CompositeResourceDefinition

	for i, doc := range documents {
		// Skip empty documents
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}

		// Unmarshal each document into an XRD object
		var xrd crossplanev1.CompositeResourceDefinition
		if err := yaml.Unmarshal([]byte(doc), &xrd); err != nil {
			return nil, fmt.Errorf("failed to unmarshal XRD at document %d: %w", i+1, err)
		}

		// Validate that it's actually an XRD
		if xrd.Kind != "CompositeResourceDefinition" {
			return nil, fmt.Errorf("document %d is not a CompositeResourceDefinition", i+1)
		}

		xrds = append(xrds, &xrd)
	}

	if len(xrds) == 0 {
		return nil, fmt.Errorf("no valid CompositeResourceDefinitions found in file")
	}

	return xrds, nil
}

// WriteToFile writes one or more CRDs to a file in either YAML or JSON format
func WriteToFile(crds []*apiextensionsv1.CustomResourceDefinition, filePath string, asJSON bool) error {
	var content []byte

	if asJSON {
		var err error
		content, err = json.MarshalIndent(crds, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal CRD: %w", err)
		}
	} else {
		// For YAML, we need to handle multiple documents
		var yamlDocs []string
		for _, crd := range crds {
			doc, err := yaml.Marshal(crd)
			if err != nil {
				return fmt.Errorf("failed to marshal CRD: %w", err)
			}
			yamlDocs = append(yamlDocs, string(doc))
		}
		// Join documents with proper YAML document separator
		content = []byte(strings.Join(yamlDocs, "\n---\n"))
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(filePath)
	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
		}
	}

	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// FormatOutput formats one or more CRDs as either YAML or JSON string
func FormatOutput(crds []*apiextensionsv1.CustomResourceDefinition, asJSON bool) (string, error) {
	var content []byte

	if asJSON {
		var err error
		content, err = json.MarshalIndent(crds, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal CRD: %w", err)
		}
	} else {
		// For YAML, we need to handle multiple documents
		var yamlDocs []string
		for _, crd := range crds {
			doc, err := yaml.Marshal(crd)
			if err != nil {
				return "", fmt.Errorf("failed to marshal CRD: %w", err)
			}
			yamlDocs = append(yamlDocs, string(doc))
		}
		// Join documents with proper YAML document separator
		content = []byte(strings.Join(yamlDocs, "\n---\n"))
	}

	return string(content), nil
}
