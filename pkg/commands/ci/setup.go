package ci

import (
	"fmt"
	"os"
	"path/filepath"

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
}

func runGithubActionsSetup(cmd *cobra.Command, args []string) error {
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

	// Create workflow file
	workflow := fmt.Sprintf(`name: Build and Push Docker Image

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
    - uses: actions/checkout@v4

    - name: Log in to the Container registry
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
`, vars.BuildContext, vars.ImageTag)

	if err := os.WriteFile(workflowPath, []byte(workflow), 0644); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}

	fmt.Printf("✅ Successfully created GitHub Actions workflow in %s\n", workflowPath)
	return nil
}
