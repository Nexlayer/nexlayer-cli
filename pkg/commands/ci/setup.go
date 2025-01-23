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
	githubActionsSetupCmd.Flags().StringVar(&vars.RegistryType, "registry-type", "ghcr", 
		"Container registry type (ghcr, dockerhub, gcr, ecr, artifactory, gitlab)")
	githubActionsSetupCmd.Flags().StringVar(&vars.Registry, "registry", "", 
		"Container registry URL (optional, defaults based on type)")
	githubActionsSetupCmd.Flags().StringVar(&vars.RegistryUsername, "registry-username", "", 
		"Registry username (optional, defaults to github.actor for GHCR)")
	githubActionsSetupCmd.Flags().StringVar(&vars.RegistryRegion, "registry-region", "", 
		"Registry region (required for ECR)")
	githubActionsSetupCmd.Flags().StringVar(&vars.RegistryProject, "registry-project", "", 
		"Registry project ID (required for GCR)")
}

func runGithubActionsSetup(cmd *cobra.Command, args []string) error {
	// Validate and set registry defaults
	vars.RegistryType = strings.ToLower(vars.RegistryType)
	validTypes := map[string]bool{
		"ghcr":       true,
		"dockerhub":  true,
		"gcr":        true,
		"ecr":        true,
		"artifactory": true,
		"gitlab":     true,
	}
	if !validTypes[vars.RegistryType] {
		return fmt.Errorf("invalid registry type: %s. Must be one of: ghcr, dockerhub, gcr, ecr, artifactory, gitlab", vars.RegistryType)
	}

	// Set registry defaults based on type
	if vars.Registry == "" {
		switch vars.RegistryType {
		case "ghcr":
			vars.Registry = "ghcr.io"
		case "dockerhub":
			vars.Registry = "docker.io"
		case "gcr":
			if vars.RegistryProject == "" {
				return fmt.Errorf("registry-project is required for GCR")
			}
			vars.Registry = fmt.Sprintf("%s-docker.pkg.dev/%s", vars.RegistryRegion, vars.RegistryProject)
		case "ecr":
			if vars.RegistryRegion == "" {
				return fmt.Errorf("registry-region is required for ECR")
			}
			vars.Registry = fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", vars.RegistryProject, vars.RegistryRegion)
		case "artifactory":
			if vars.Registry == "" {
				return fmt.Errorf("registry URL is required for Artifactory")
			}
		case "gitlab":
			if vars.Registry == "" {
				vars.Registry = "registry.gitlab.com"
			}
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
		workflow = createGHCRWorkflow()
	case "dockerhub":
		workflow = createDockerHubWorkflow()
	case "gcr":
		workflow = createGCRWorkflow()
	case "ecr":
		workflow = createECRWorkflow()
	case "artifactory":
		workflow = createArtifactoryWorkflow()
	case "gitlab":
		workflow = createGitLabWorkflow()
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

func createGHCRWorkflow() string {
	return fmt.Sprintf(`name: Build and Push Docker Image

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
}

func createDockerHubWorkflow() string {
	return fmt.Sprintf(`name: Build and Push Docker Image

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

func createGCRWorkflow() string {
	return fmt.Sprintf(`name: Build and Push Docker Image

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: %s
  IMAGE_NAME: ${{ github.event.repository.name }}

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Log in to Google Container Registry
      uses: google-github-actions/login@v1
      with:
        credentials: ${{ secrets.GOOGLE_CREDENTIALS }}

    - name: Build and push Docker image
      uses: docker/build-push-action@v5
      with:
        context: %s
        push: true
        tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:%s
`, vars.Registry, vars.BuildContext, vars.ImageTag)
}

func createECRWorkflow() string {
	return fmt.Sprintf(`name: Build and Push Docker Image

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: %s
  IMAGE_NAME: ${{ github.event.repository.name }}

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Log in to Amazon ECR
      uses: aws-actions/login@v2
      with:
        aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
        aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
        aws-region: %s

    - name: Build and push Docker image
      uses: docker/build-push-action@v5
      with:
        context: %s
        push: true
        tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:%s
`, vars.Registry, vars.RegistryRegion, vars.BuildContext, vars.ImageTag)
}

func createArtifactoryWorkflow() string {
	return fmt.Sprintf(`name: Build and Push Docker Image

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: %s
  IMAGE_NAME: ${{ github.event.repository.name }}

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Log in to Artifactory
      uses: jfrog/jfrog-github-actions@v2
      with:
        server-id: ${{ secrets.ARTIFACTORY_SERVER_ID }}
        username: ${{ secrets.ARTIFACTORY_USERNAME }}
        password: ${{ secrets.ARTIFACTORY_PASSWORD }}

    - name: Build and push Docker image
      uses: docker/build-push-action@v5
      with:
        context: %s
        push: true
        tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:%s
`, vars.Registry, vars.BuildContext, vars.ImageTag)
}

func createGitLabWorkflow() string {
	return fmt.Sprintf(`name: Build and Push Docker Image

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

env:
  REGISTRY: %s
  IMAGE_NAME: ${{ github.event.repository.name }}

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Log in to GitLab Container Registry
      uses: docker/login-action@v3
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ secrets.GITLAB_USERNAME }}
        password: ${{ secrets.GITLAB_PASSWORD }}

    - name: Build and push Docker image
      uses: docker/build-push-action@v5
      with:
        context: %s
        push: true
        tags: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:%s
`, vars.Registry, vars.BuildContext, vars.ImageTag)
}
