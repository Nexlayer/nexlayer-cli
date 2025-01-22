package ci

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	stack    string
	registry string
	token    string
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up CI/CD workflows",
	Long: `Generate CI/CD workflow files for different platforms.
Currently supports:
- GitHub Actions`,
}

// githubActionsCmd represents the github-actions command
var githubActionsCmd = &cobra.Command{
	Use:   "github-actions",
	Short: "Generate GitHub Actions workflow",
	Long:  `Generate a GitHub Actions workflow file for building and publishing Docker images`,
	RunE:  runGithubActionsSetup,
}

func init() {
	setupCmd.AddCommand(githubActionsCmd)

	// Add flags
	githubActionsCmd.Flags().StringVar(&stack, "stack", "", "Application stack (e.g., mern, next, django)")
	githubActionsCmd.Flags().StringVar(&registry, "registry", "ghcr.io", "Container registry")
	githubActionsCmd.Flags().StringVar(&token, "token", "", "GitHub token")

	// Mark required flags
	githubActionsCmd.MarkFlagRequired("stack")
}

func runGithubActionsSetup(cmd *cobra.Command, args []string) error {
	// Create .github/workflows directory if it doesn't exist
	workflowDir := ".github/workflows"
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflow directory: %w", err)
	}

	// Get workflow template based on stack
	tmpl, err := getWorkflowTemplate(stack)
	if err != nil {
		return err
	}

	// Create workflow file
	workflowFile := filepath.Join(workflowDir, "docker-publish.yml")
	f, err := os.Create(workflowFile)
	if err != nil {
		return fmt.Errorf("failed to create workflow file: %w", err)
	}
	defer f.Close()

	// Execute template
	data := struct {
		Registry string
		Stack    string
	}{
		Registry: registry,
		Stack:    stack,
	}

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to generate workflow file: %w", err)
	}

	fmt.Printf("✅ Successfully created GitHub Actions workflow in %s\n", workflowFile)
	return nil
}

func getWorkflowTemplate(stack string) (*template.Template, error) {
	const workflowTemplate = `name: Docker Image CI

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

env:
  REGISTRY: {{ .Registry }}
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
        registry: {{ .Registry }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    
    - name: Extract metadata for Docker
      id: meta
      uses: docker/metadata-action@v5
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
    
    - name: Build and push Docker image
      uses: docker/build-push-action@v5
      with:
        context: .
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
{{- if eq .Stack "mern" }}
        # MERN stack specific settings
        build-args: |
          NODE_ENV=production
          REACT_APP_API_URL=${{ secrets.API_URL }}
{{- else if eq .Stack "next" }}
        # Next.js specific settings
        build-args: |
          NEXT_PUBLIC_API_URL=${{ secrets.API_URL }}
{{- end }}`

	return template.New("workflow").Parse(workflowTemplate)
}
