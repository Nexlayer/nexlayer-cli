package ci

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	imageTag     string
	buildContext string
)

// customizeCmd represents the customize command
var customizeCmd = &cobra.Command{
	Use:   "customize",
	Short: "Customize CI/CD workflows",
	Long: `Customize existing CI/CD workflow files.
Currently supports:
- GitHub Actions workflow customization`,
}

// githubActionsCustomizeCmd represents the github-actions customize command
var githubActionsCustomizeCmd = &cobra.Command{
	Use:   "github-actions",
	Short: "Customize GitHub Actions workflow",
	Long:  `Customize an existing GitHub Actions workflow file with specific parameters`,
	RunE:  runGithubActionsCustomize,
}

func init() {
	customizeCmd.AddCommand(githubActionsCustomizeCmd)

	// Add flags
	githubActionsCustomizeCmd.Flags().StringVar(&imageTag, "image-tag", "latest", "Docker image tag")
	githubActionsCustomizeCmd.Flags().StringVar(&buildContext, "build-context", ".", "Docker build context path")
}

type WorkflowFile struct {
	Name string                 `yaml:"name"`
	On   map[string]interface{} `yaml:"on"`
	Env  map[string]string      `yaml:"env"`
	Jobs map[string]Job         `yaml:"jobs"`
}

type Job struct {
	RunsOn      string            `yaml:"runs-on"`
	Permissions map[string]string `yaml:"permissions"`
	Steps       []Step            `yaml:"steps"`
}

type Step struct {
	Name string                 `yaml:"name"`
	Uses string                 `yaml:"uses,omitempty"`
	With map[string]interface{} `yaml:"with,omitempty"`
}

func runGithubActionsCustomize(cmd *cobra.Command, args []string) error {
	workflowPath := ".github/workflows/docker-publish.yml"
	
	// Check if workflow file exists
	if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
		return fmt.Errorf("workflow file not found. Run 'nexlayer ci setup github-actions' first")
	}

	// Read existing workflow file
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		return fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow WorkflowFile
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return fmt.Errorf("failed to parse workflow file: %w", err)
	}

	// Update workflow configuration
	modified := false
	for jobName, job := range workflow.Jobs {
		for i, step := range job.Steps {
			if step.Uses == "docker/build-push-action@v5" {
				// Update build context if specified
				if buildContext != "." {
					if step.With == nil {
						step.With = make(map[string]interface{})
					}
					step.With["context"] = buildContext
					modified = true
				}

				// Update image tag if specified
				if imageTag != "latest" {
					if step.With == nil {
						step.With = make(map[string]interface{})
					}
					tags := fmt.Sprintf("${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:%s", imageTag)
					step.With["tags"] = tags
					modified = true
				}

				job.Steps[i] = step
				workflow.Jobs[jobName] = job
			}
		}
	}

	if !modified {
		fmt.Println("ℹ️  No changes were necessary")
		return nil
	}

	// Save updated workflow file
	updatedData, err := yaml.Marshal(workflow)
	if err != nil {
		return fmt.Errorf("failed to generate workflow file: %w", err)
	}

	if err := os.WriteFile(workflowPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to save workflow file: %w", err)
	}

	fmt.Printf("✅ Successfully updated GitHub Actions workflow in %s\n", workflowPath)
	if buildContext != "." {
		fmt.Printf("📁 Build context set to: %s\n", buildContext)
	}
	if imageTag != "latest" {
		fmt.Printf("🏷️  Image tag set to: %s\n", imageTag)
	}

	return nil
}
