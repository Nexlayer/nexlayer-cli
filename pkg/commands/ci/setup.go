package ci

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up CI/CD pipelines",
	Long: `Set up CI/CD pipelines for your project.
Currently supports:
- GitHub Actions workflow setup`,
}

// githubActionsSetupCmd represents the github-actions command
var githubActionsSetupCmd = &cobra.Command{
	Use:   "github-actions",
	Short: "Set up GitHub Actions workflow",
	Long: `Set up GitHub Actions workflow for your project.
This will create a basic workflow file in .github/workflows/build.yml`,
	RunE: runGithubActionsSetup,
}

func init() {
	setupCmd.AddCommand(githubActionsSetupCmd)
	
	// Add flags for registry configuration
	githubActionsSetupCmd.Flags().StringVar(&vars.RegistryType, "registry-type", "ghcr", "Container registry type (ghcr or dockerhub)")
	githubActionsSetupCmd.Flags().StringVar(&vars.Registry, "registry", "", "Container registry URL (optional, defaults based on type)")
	githubActionsSetupCmd.Flags().StringVar(&vars.RegistryUsername, "registry-username", "", "Registry username (optional, defaults to github.actor for GHCR)")
}

func runGithubActionsSetup(cmd *cobra.Command, args []string) error {
	// Validate and set registry defaults
	vars.RegistryType = strings.ToLower(vars.RegistryType)
	if vars.RegistryType != "ghcr" && vars.RegistryType != "dockerhub" {
		return fmt.Errorf("invalid registry type: %s. Must be 'ghcr' or 'dockerhub'", vars.RegistryType)
	}

	// Set registry defaults based on type
	if vars.Registry == "" {
		switch vars.RegistryType {
		case "ghcr":
			vars.Registry = "ghcr.io"
		case "dockerhub":
			vars.Registry = "docker.io"
		}
	}

	// Create .github/workflows directory if it doesn't exist
	workflowDir := ".github/workflows"
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflow directory: %w", err)
	}

	// Check if workflow file already exists
	workflowPath := filepath.Join(workflowDir, "build.yml")
	if _, err := os.Stat(workflowPath); err == nil {
		return fmt.Errorf("workflow file already exists: %s", workflowPath)
	}

	// Create workflow file with registry-specific configuration
	var workflow string
	switch vars.RegistryType {
	case "ghcr":
		workflow = fmt.Sprintf(`name: Build and Push Docker Image

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: %s
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
    - uses: actions/checkout@v4

    - name: Log in to GitHub Container Registry
      uses: docker/login-action@v3
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Build and push Docker image
      uses: docker/build-push-action@v5
      with:
        context: %s
        push: true
        tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:%s
`, vars.Registry, vars.BuildContext, vars.ImageTag)

	case "dockerhub":
		workflow = fmt.Sprintf(`name: Build and Push Docker Image

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: %s
  IMAGE_NAME: ${{ secrets.DOCKERHUB_USERNAME }}/${{ github.event.repository.name }}

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Log in to Docker Hub
      uses: docker/login-action@v3
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ secrets.DOCKERHUB_USERNAME }}
        password: ${{ secrets.DOCKERHUB_TOKEN }}

    - name: Build and push Docker image
      uses: docker/build-push-action@v5
      with:
        context: %s
        push: true
        tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:%s
`, vars.Registry, vars.BuildContext, vars.ImageTag)
	}

	if err := os.WriteFile(workflowPath, []byte(workflow), 0644); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}

	// Print registry-specific setup instructions
	fmt.Printf("✅ Successfully created GitHub Actions workflow in %s\n", workflowPath)
	if vars.RegistryType == "dockerhub" {
		fmt.Println("\n⚠️ Important: You need to configure the following secrets in your GitHub repository:")
		fmt.Println("  - DOCKERHUB_USERNAME: Your Docker Hub username")
		fmt.Println("  - DOCKERHUB_TOKEN: Your Docker Hub access token (create at https://hub.docker.com/settings/security)")
	}
	
	return nil
}
