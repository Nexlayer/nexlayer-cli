package ci

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
	"github.com/spf13/cobra"
)

// imagesCmd represents the images command
var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Manage Docker images",
	Long: `Manage Docker images for your Nexlayer deployment.
Examples:
  nexlayer ci images build --tag v1.0.0
  nexlayer ci images push --tag v1.0.0`,
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Docker image",
	Long: `Build a Docker image for your application.
Example:
  nexlayer ci images build --tag v1.0.0`,
	RunE: runBuild,
}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push Docker image",
	Long: `Push a Docker image to the registry.
Example:
  nexlayer ci images push --tag v1.0.0`,
	RunE: runPush,
}

func init() {
	imagesCmd.AddCommand(buildCmd)
	imagesCmd.AddCommand(pushCmd)

	// Add flags
	buildCmd.Flags().StringVar(&vars.ImageTag, "tag", "latest", "Image tag")
	pushCmd.Flags().StringVar(&vars.ImageTag, "tag", "latest", "Image tag")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Build Docker image
	buildCmd := exec.Command("docker", "build", "-t", fmt.Sprintf("ghcr.io/%s:%s", vars.AppName, vars.ImageTag), ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	fmt.Printf("🔨 Building Docker image ghcr.io/%s:%s...\n", vars.AppName, vars.ImageTag)
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}

	fmt.Printf("✅ Successfully built image ghcr.io/%s:%s\n", vars.AppName, vars.ImageTag)
	return nil
}

func runPush(cmd *cobra.Command, args []string) error {
	// Push Docker image
	pushCmd := exec.Command("docker", "push", fmt.Sprintf("ghcr.io/%s:%s", vars.AppName, vars.ImageTag))
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr

	fmt.Printf("📤 Pushing Docker image ghcr.io/%s:%s...\n", vars.AppName, vars.ImageTag)
	if err := pushCmd.Run(); err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}

	fmt.Printf("✅ Successfully pushed image ghcr.io/%s:%s\n", vars.AppName, vars.ImageTag)
	return nil
}
