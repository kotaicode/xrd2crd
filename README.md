# XRD to CRD Converter

A Go tool that converts Crossplane CompositeResourceDefinitions (XRDs) into Kubernetes CustomResourceDefinitions (CRDs). This tool specifically generates the claim CRDs from a given XRD.

## Overview

When working with Crossplane, CompositeResourceDefinitions (XRDs) are used to define composite resources. Each XRD can have associated claims that users can create to request instances of the composite resource. This tool helps by converting the XRD into its corresponding claim CRD format that Kubernetes can understand.

## Features

- Converts XRD to claim CRD format
- Preserves metadata (labels and annotations)
- Handles version-specific configurations
- Maintains OpenAPI v3 schema validation
- Supports multiple output formats (YAML/JSON)
- Supports file output and stdout
- Supports processing multiple files using wildcards
- Separates multiple CRDs with YAML document separators
- Command-line interface with built-in help

## Prerequisites

- Go 1.22 or higher
- Task (task runner) - [Installation guide](https://taskfile.dev/installation/)
- Dependencies are managed via Go modules

## Development Setup

This project uses [Task](https://taskfile.dev) for development workflows. To get started:

1. Install Task if you haven't already:
   ```bash
   # On macOS
   brew install go-task

   # On Windows (with scoop)
   scoop install task

   # On Linux
   sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
   ```

2. Set up the development environment:
   ```bash
   task setup-dev
   ```

This will install required development tools like `golangci-lint` and ensure dependencies are up to date.

## Development Commands

The project includes several Task commands to help with development:

```bash
# List all available tasks
task

# Build the binary
task build

# Run the application (for development)
task run -- <args>  # e.g., task run -- "xrds/*.yaml"

# Run tests
task test

# Run tests with coverage
task test-coverage

# Format code
task fmt

# Run linters
task lint

# Update dependencies
task update-deps

# Run all checks (lint, test, build)
task check

# Try the example
task example
```

### Docker Development

The project includes Docker support for building and running the tool:

```bash
# Build the Docker image locally
task docker-build

# Run the tool using Docker
task docker-run -- --help
task docker-run -- "xrds/*.yaml"

# For development, you can mount the current directory
docker run --rm -v $(pwd):/work -w /work ghcr.io/kotaicode/xrd2crd:dev "xrds/*.yaml"
```

The Docker image is based on Google's distroless base image, resulting in a minimal and secure container that only includes the necessary components to run the tool.

## Release Workflow

This project uses GoReleaser for building and publishing releases. You can test releases locally or publish them through GitHub Actions.

### Local Release Testing

1. Test the release configuration:
   ```bash
   task release-check
   ```

2. Build binaries locally without publishing:
   ```bash
   task release-local
   ```
   This will create binaries in the `dist/` directory without publishing them.

3. Create a full snapshot release:
   ```bash
   task release-snapshot
   ```
   This simulates a complete release process locally, including archive creation.

### Publishing Official Releases

1. Configure GitHub Repository:
   - The workflow uses the default `GITHUB_TOKEN` which is automatically provided by GitHub Actions
   - No additional configuration is needed for basic releases

2. Create and Push a Release:
   ```bash
   # Ensure all changes are committed
   git add .
   git commit -m "Prepare for release"

   # Create an annotated tag
   git tag -a v1.0.0 -m "First release"

   # Push the tag
   git push origin v1.0.0
   ```

3. The GitHub Actions workflow will automatically:
   - Build the release binaries for all platforms
   - Create a GitHub release with the binaries
   - Generate checksums
   - Generate a changelog
   - Build and push Docker images to GitHub Container Registry:
     - `ghcr.io/kotaicode/xrd2crd:latest`
     - `ghcr.io/kotaicode/xrd2crd:v1.0.0` (version tag)

4. Verify the Release:
   - Check the GitHub Actions workflow status
   - Verify the release on the GitHub releases page
   - Download and test the released binaries

### Installing Released Versions

After release, users can install the tool in several ways:

1. Using `go install`:
   ```bash
   go install github.com/kotaicode/xrd2crd@latest
   ```

2. Using Docker:
   ```bash
   # Pull the latest version
   docker pull ghcr.io/kotaicode/xrd2crd:latest

   # Run the tool
   docker run --rm ghcr.io/kotaicode/xrd2crd:latest --help

   # Process files (mount a directory)
   docker run --rm -v $(pwd):/work -w /work ghcr.io/kotaicode/xrd2crd:latest "xrds/*.yaml"
   ```

3. Downloading directly from GitHub Releases:
   - Visit the releases page on GitHub
   - Download the appropriate binary for your platform:
     - Linux: `xrd2crd_Linux_x86_64.tar.gz` or `xrd2crd_Linux_arm64.tar.gz`
     - macOS: `xrd2crd_Darwin_x86_64.tar.gz` or `xrd2crd_Darwin_arm64.tar.gz`
     - Windows: `xrd2crd_Windows_x86_64.zip`
   - Extract the archive and move the binary to a location in your PATH

## Usage

The tool provides a flexible command-line interface with various output options:

```bash
# Show help and version
xrd2crd --help
xrd2crd -v

# Basic usage - outputs YAML to stdout (default)
xrd2crd xrd.yaml

# Output as YAML explicitly
xrd2crd xrd.yaml -o yaml

# Output as JSON
xrd2crd xrd.yaml -o json

# Output to a file (defaults to YAML)
xrd2crd xrd.yaml -o path=/tmp/output.yaml

# Output to file and stdout
xrd2crd xrd.yaml -o path=/tmp/output.yaml -s

# Process multiple files
xrd2crd "xrds/*.yaml" -o yaml

# Process multiple files to separate output files
xrd2crd "xrds/*.yaml" -o path=output.yaml  # Creates output-1.yaml, output-2.yaml, etc.
```

### Command-line Options

- `<pattern>`: File path or glob pattern for input XRD files (required)
- `-o, --output`: Output format and destination. Can be:
  - `yaml`: Output as YAML to stdout (default)
  - `json`: Output as JSON to stdout
  - `path=/path/to/file`: Write output to specified file (in YAML format)
- `-s, --stdout`: Force output to stdout even if output file is specified
- `-v, --version`: Show version information
- `--help`: Show help message

When processing multiple files:
- With stdout: CRDs are separated by `---` document separators
- With file output: Each CRD is written to a separate file with a numeric suffix

## Project Structure

```
xrd2crd/
├── cmd/
│   └── xrd2crd/       # CLI application
│       └── main.go
├── pkg/
│   ├── converter/     # Core conversion logic
│   │   └── converter.go
│   └── fileio/        # File operations
│       └── fileio.go
├── examples/          # Example XRDs
│   └── example.yaml
├── .goreleaser.yaml  # Release configuration
├── Taskfile.yml      # Development tasks
├── go.mod
└── README.md
```

## Error Handling

The tool includes error handling for common scenarios:
- Missing or invalid input files
- Invalid XRD format
- Schema conversion errors
- Invalid file patterns
- No matching files found
- Output file writing errors
- Invalid output format combinations

## Limitations

- Currently only processes claim CRDs (not composite resource CRDs)
- All versions of the CRD are set to "Namespaced" scope 