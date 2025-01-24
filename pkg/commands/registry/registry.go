package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/registry"
	"github.com/spf13/cobra"
)

// NewRegistryCmd creates a new registry command group
func NewRegistryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Manage container registry operations",
		Long:  `Build, tag, and push Docker images to GitHub Container Registry (GHCR)`,
	}

	cmd.AddCommand(
		newLoginCmd(),
		newBuildCmd(),
		newPushCmd(),
	)

	return cmd
}

func newLoginCmd() *cobra.Command {
	var (
		username string
		token    string
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to GitHub Container Registry",
		Long:  `Authenticate with GHCR using GitHub username and Personal Access Token (PAT)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If token not provided via flag, try environment variable
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
				if token == "" {
					return fmt.Errorf("GitHub token not provided. Use --token flag or set GITHUB_TOKEN environment variable")
				}
			}

			client := registry.NewClient(&registry.RegistryConfig{
				Username: username,
				Token:    token,
				Registry: "ghcr.io",
			})

			if err := client.Login(); err != nil {
				return fmt.Errorf("failed to login to GHCR: %w", err)
			}

			fmt.Println("Successfully logged in to GHCR")
			return nil
		},
	}

	cmd.Flags().StringVar(&username, "username", "", "GitHub username")
	cmd.Flags().StringVar(&token, "token", "", "GitHub Personal Access Token (PAT)")
	cmd.MarkFlagRequired("username")

	return cmd
}

func newBuildCmd() *cobra.Command {
	var (
		namespace string
		tags      []string
		services  []string
	)

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build Docker images",
		Long:  `Build and tag Docker images for multiple services`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !strings.HasPrefix(namespace, "ghcr.io/") {
				namespace = "ghcr.io/" + namespace
			}

			var images []registry.ImageConfig
			for _, svc := range services {
				images = append(images, registry.ImageConfig{
					ServiceName: filepath.Base(svc),
					Path:       svc,
					Tags:      tags,
					Namespace: namespace,
				})
			}

			client := registry.NewClient(&registry.RegistryConfig{
				Registry: "ghcr.io",
			})

			buildCfg := registry.BuildConfig{
				Images:    images,
				Namespace: namespace,
				Tags:     tags,
			}

			if err := client.BuildAndPushImages(buildCfg); err != nil {
				return fmt.Errorf("failed to build images: %w", err)
			}

			fmt.Println("Successfully built all images")
			return nil
		},
	}

	cmd.Flags().StringVar(&namespace, "namespace", "", "Registry namespace (e.g., my-org)")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{"latest"}, "Image tags (comma-separated)")
	cmd.Flags().StringSliceVar(&services, "services", []string{}, "Paths to service directories (comma-separated)")
	cmd.MarkFlagRequired("namespace")
	cmd.MarkFlagRequired("services")

	return cmd
}

func newPushCmd() *cobra.Command {
	var (
		namespace string
		tags      []string
		services  []string
	)

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push Docker images to GHCR",
		Long:  `Push built Docker images to GitHub Container Registry`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !strings.HasPrefix(namespace, "ghcr.io/") {
				namespace = "ghcr.io/" + namespace
			}

			var images []registry.ImageConfig
			for _, svc := range services {
				images = append(images, registry.ImageConfig{
					ServiceName: filepath.Base(svc),
					Path:       svc,
					Tags:      tags,
					Namespace: namespace,
				})
			}

			client := registry.NewClient(&registry.RegistryConfig{
				Registry: "ghcr.io",
			})

			for _, img := range images {
				if err := client.PushImage(img); err != nil {
					return fmt.Errorf("failed to push image %s: %w", img.ServiceName, err)
				}
			}

			fmt.Println("Successfully pushed all images to GHCR")
			return nil
		},
	}

	cmd.Flags().StringVar(&namespace, "namespace", "", "Registry namespace (e.g., my-org)")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{"latest"}, "Image tags (comma-separated)")
	cmd.Flags().StringSliceVar(&services, "services", []string{}, "Service names to push (comma-separated)")
	cmd.MarkFlagRequired("namespace")
	cmd.MarkFlagRequired("services")

	return cmd
}
