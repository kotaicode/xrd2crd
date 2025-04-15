package xrd2crd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/kotaicode/xrd2crd/pkg/converter"
	"github.com/kotaicode/xrd2crd/pkg/fileio"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Version information (populated by GoReleaser)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// OutputFormat represents the output format and destination
type OutputFormat struct {
	Format string
	Path   string
}

// Decode implements kong.MapperValue interface for custom output format parsing
func (o *OutputFormat) Decode(ctx *kong.DecodeContext) error {
	var value string
	if err := ctx.Scan.PopValueInto("value", &value); err != nil {
		return err
	}

	// Check for path specification
	if strings.HasPrefix(value, "path=") {
		o.Path = strings.TrimPrefix(value, "path=")
		o.Format = "yaml" // default format for file output
		return nil
	}

	// Handle format specification
	switch value {
	case "yaml", "json":
		o.Format = value
	default:
		return fmt.Errorf("invalid output format: %s (must be yaml, json, or path=<filepath>)", value)
	}

	return nil
}

// CLI represents the command-line interface structure
var CLI struct {
	Pattern string       `arg:"" optional:"" help:"File path or glob pattern for XRD files (e.g., 'xrd.yaml' or 'xrds/*.yaml')" type:"path"`
	Output  OutputFormat `help:"Output format and destination. Can be 'yaml', 'json', or 'path=/path/to/file'" short:"o"`
	Stdout  bool         `help:"Force output to stdout even if output file is specified" short:"s"`
	Version bool         `help:"Show version information" short:"v"`
}

// Main is the entry point for the CLI tool
func Main(v, c, d string) {
	// Update version information
	version = v
	commit = c
	date = d

	ctx := kong.Parse(&CLI,
		kong.Name("xrd2crd"),
		kong.Description("Converts Crossplane XRDs to Kubernetes CRDs"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))

	if CLI.Version {
		fmt.Printf("xrd2crd %s (%s) - %s\n", version, commit, date)
		os.Exit(0)
	}

	if CLI.Pattern == "" {
		fmt.Println("Error: pattern argument is required when not using --version")
		os.Exit(1)
	}

	matches, err := filepath.Glob(CLI.Pattern)
	if err != nil {
		fmt.Printf("Error finding files: %v\n", err)
		os.Exit(1)
	}

	if len(matches) == 0 {
		fmt.Printf("No files found matching pattern: %s\n", CLI.Pattern)
		os.Exit(1)
	}

	var allCRDs []*apiextensionsv1.CustomResourceDefinition

	for fileIndex, filePath := range matches {
		if fileIndex > 0 && (CLI.Stdout || CLI.Output.Path == "") {
			// Print separator between multiple files when outputting to stdout
			fmt.Println("---")
		}

		// Load and convert XRDs
		xrds, err := fileio.LoadXRDFromFile(filePath)
		if err != nil {
			fmt.Printf("Error loading XRD from %s: %v\n", filePath, err)
			continue
		}

		for _, xrd := range xrds {
			crd, err := converter.ConvertXRDToCRD(xrd)
			if err != nil {
				fmt.Printf("Error converting XRD from %s: %v\n", filePath, err)
				continue
			}
			allCRDs = append(allCRDs, crd)
		}
	}

	if len(allCRDs) == 0 {
		fmt.Println("No valid CRDs were generated")
		os.Exit(1)
	}

	// Handle output
	if CLI.Output.Path != "" {
		// Generate output filename for multiple files
		outputPath := CLI.Output.Path
		if len(matches) > 1 {
			ext := filepath.Ext(CLI.Output.Path)
			base := CLI.Output.Path[:len(CLI.Output.Path)-len(ext)]
			outputPath = fmt.Sprintf("%s-%d%s", base, len(matches), ext)
		}

		if err := fileio.WriteToFile(allCRDs, outputPath, CLI.Output.Format == "json"); err != nil {
			fmt.Printf("Error writing to file %s: %v\n", outputPath, err)
			os.Exit(1)
		}

		// Print to stdout as well if requested
		if CLI.Stdout {
			if output, err := fileio.FormatOutput(allCRDs, CLI.Output.Format == "json"); err != nil {
				fmt.Printf("Error formatting output: %v\n", err)
			} else {
				fmt.Println(output)
			}
		}
	} else {
		// Print to stdout
		if output, err := fileio.FormatOutput(allCRDs, CLI.Output.Format == "json"); err != nil {
			fmt.Printf("Error formatting output: %v\n", err)
		} else {
			fmt.Println(output)
		}
	}

	ctx.Exit(0)
}
