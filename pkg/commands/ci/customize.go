package ci

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
)

// customizeCmd represents the customize command
var customizeCmd = &cobra.Command{
	Use:   "customize",
	Short: "Customize CI/CD workflows",
	Long: `Customize existing CI/CD workflow files.
Currently supports:
- GitHub Actions workflow customization`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("please specify a valid subcommand")
	},
}

// githubActionsCustomizeCmd represents the github-actions command
var githubActionsCustomizeCmd = &cobra.Command{
	Use:   "github-actions",
	Short: "Customize GitHub Actions workflow",
	Long: `Customize the GitHub Actions workflow for your project.
Examples:
  nexlayer ci customize github-actions --build-context ./src --image-tag v1.0.0`,
	RunE: runGithubActionsCustomize,
}

func init() {
	customizeCmd.AddCommand(githubActionsCustomizeCmd)

	// Add flags
	githubActionsCustomizeCmd.Flags().StringVar(&vars.BuildContext, "build-context", ".", "Docker build context path")
	githubActionsCustomizeCmd.Flags().StringVar(&vars.ImageTag, "image-tag", "latest", "Docker image tag")
}

type WorkflowFile struct {
	Name string                 `yaml:"name"`
	On   map[string]interface{} `yaml:"on"`
	Jobs map[string]Job         `yaml:"jobs"`
}

type Job struct {
	RunsOn string `yaml:"runs-on"`
	Steps  []Step `yaml:"steps"`
}

type Step struct {
	Name string                 `yaml:"name"`
	Uses string                 `yaml:"uses,omitempty"`
	Run  string                 `yaml:"run,omitempty"`
	With map[string]interface{} `yaml:"with,omitempty"`
}

func runGithubActionsCustomize(cmd *cobra.Command, args []string) error {
	workflowPath := ".github/workflows/build.yml"
	if _, err := os.Stat(workflowPath); os.IsNotExist(err) {
		return fmt.Errorf("workflow file not found: %s", workflowPath)
	}

	data, err := os.ReadFile(workflowPath)
	if err != nil {
		return fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow WorkflowFile
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return fmt.Errorf("failed to parse workflow file: %w", err)
	}

	modified := false
	for jobName, job := range workflow.Jobs {
		for i, step := range job.Steps {
			if step.Uses == "docker/build-push-action@v5" {
				// Update build context if specified
				if vars.BuildContext != "." {
					if step.With == nil {
						step.With = make(map[string]interface{})
					}
					step.With["context"] = vars.BuildContext
					modified = true
				}

				// Update image tag if specified
				if vars.ImageTag != "latest" {
					if step.With == nil {
						step.With = make(map[string]interface{})
					}
					tags := fmt.Sprintf("${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:%s", vars.ImageTag)
					step.With["tags"] = tags
					modified = true
				}

				workflow.Jobs[jobName].Steps[i] = step
			}
		}
	}

	if !modified {
		return fmt.Errorf("no changes were made to the workflow file")
	}

	data, err = yaml.Marshal(&workflow)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow file: %w", err)
	}

	if err := os.WriteFile(workflowPath, data, 0o644); err != nil {
		return fmt.Errorf("failed to write workflow file: %w", err)
	}

	fmt.Printf("✅ Successfully updated GitHub Actions workflow in %s\n", workflowPath)
	if vars.BuildContext != "." {
		fmt.Printf("📁 Build context set to: %s\n", vars.BuildContext)
	}
	if vars.ImageTag != "latest" {
		fmt.Printf("🏷️  Image tag set to: %s\n", vars.ImageTag)
	}

	return nil
}
