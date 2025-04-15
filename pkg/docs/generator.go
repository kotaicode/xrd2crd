package docs

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

//go:embed template.html
var templateFS embed.FS

// PropertiesData holds the data needed for rendering properties
type PropertiesData struct {
	Properties map[string]apiextensionsv1.JSONSchemaProps
	Required   []string
	Level      int
}

// newPropertiesData creates a new PropertiesData instance
func newPropertiesData(properties map[string]apiextensionsv1.JSONSchemaProps, required []string, level int) PropertiesData {
	return PropertiesData{
		Properties: properties,
		Required:   required,
		Level:      level,
	}
}

// cleanJSON removes quotes from a JSON string value
func cleanJSON(val apiextensionsv1.JSON) string {
	var result interface{}
	if err := json.Unmarshal(val.Raw, &result); err == nil {
		if str, ok := result.(string); ok {
			return str
		}
		data, err := json.Marshal(result)
		if err == nil {
			str := string(data)
			str = strings.Trim(str, "\"")
			return str
		}
	}
	str := string(val.Raw)
	str = strings.Trim(str, "\"")
	return str
}

// toString converts a value to a string, handling byte arrays and JSON
func toString(v interface{}) string {
	switch val := v.(type) {
	case []byte:
		return string(val)
	case string:
		return val
	case apiextensionsv1.JSON:
		return cleanJSON(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// generateExampleYAML generates a YAML example for a CRD
func generateExampleYAML(crd *apiextensionsv1.CustomResourceDefinition) string {
	var version string
	for _, v := range crd.Spec.Versions {
		if v.Storage {
			version = v.Name
			break
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("apiVersion: %s/%s\n", crd.Spec.Group, version))
	sb.WriteString(fmt.Sprintf("kind: %s\n", crd.Spec.Names.Kind))
	sb.WriteString("metadata:\n")
	sb.WriteString("  name: example\n")
	sb.WriteString("spec:\n")

	for _, v := range crd.Spec.Versions {
		if v.Storage {
			if v.Schema != nil && v.Schema.OpenAPIV3Schema != nil {
				generatePropertiesYAML(&sb, v.Schema.OpenAPIV3Schema.Properties, v.Schema.OpenAPIV3Schema.Required, 1)
			}
		}
	}

	return sb.String()
}

// generatePropertiesYAML generates YAML for properties
func generatePropertiesYAML(sb *strings.Builder, properties map[string]apiextensionsv1.JSONSchemaProps, required []string, level int) {
	indent := strings.Repeat("  ", level)
	for name, property := range properties {
		sb.WriteString(fmt.Sprintf("%s%s:", indent, name))
		if property.Properties != nil {
			sb.WriteString("\n")
			generatePropertiesYAML(sb, property.Properties, property.Required, level+1)
		} else if property.Enum != nil && len(property.Enum) > 0 {
			sb.WriteString(fmt.Sprintf(" %s\n", cleanJSON(property.Enum[0])))
		} else {
			switch property.Type {
			case "string":
				sb.WriteString(" \"example\"\n")
			case "integer":
				sb.WriteString(" 1\n")
			case "boolean":
				sb.WriteString(" true\n")
			case "array":
				sb.WriteString(" []\n")
			default:
				sb.WriteString(" {}\n")
			}
		}
	}
}

type CRDDoc struct {
	CRDs []*apiextensionsv1.CustomResourceDefinition
}

// contains checks if a string is in a slice of strings
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// add adds two integers
func add(a, b int) int {
	return a + b
}

// GenerateHTML generates HTML documentation for the given CRDs
func GenerateHTML(crds []*apiextensionsv1.CustomResourceDefinition, outputPath string) error {
	templateContent, err := templateFS.ReadFile("template.html")
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	tmpl, err := template.New("crd-docs").Funcs(template.FuncMap{
		"contains":          contains,
		"add":               add,
		"newPropertiesData": newPropertiesData,
		"cleanJSON":         cleanJSON,
		"indent": func(level int) string {
			return strings.Repeat("  ", level)
		},
		"generateExampleYAML": generateExampleYAML,
	}).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	doc := CRDDoc{
		CRDs: crds,
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
		}
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Execute template
	if err := tmpl.Execute(file, doc); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
