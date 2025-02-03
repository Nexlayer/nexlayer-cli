package compose

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type DockerCompose struct {
	Version  string                 `yaml:"version"`
	Services map[string]Service     `yaml:"services"`
	Networks map[string]interface{} `yaml:"networks"`
}

type Service struct {
	Image       string            `yaml:"image"`
	Build       string            `yaml:"build,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Environment []string          `yaml:"environment,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compose",
		Short: "Manage local development with Docker Compose",
		Long:  `Commands for managing local development environment using Docker Compose.`,
	}

	cmd.AddCommand(newGenerateCommand())
	cmd.AddCommand(newUpCommand())
	cmd.AddCommand(newDownCommand())
	cmd.AddCommand(newLogsCommand())

	return cmd
}

func getDefaultImage(podType string, tag string) string {
	// If tag is provided, use it
	if tag != "" {
		return tag
	}

	// Default images for each pod type
	switch podType {
	case "database":
		return "mongo:latest"  // Official MongoDB image
	case "backend":
		return "node:18"
	case "frontend":
		return "node:18"
	default:
		return "node:18"
	}
}

func newGenerateCommand() *cobra.Command {
	var configFile string
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate docker-compose.yml from nexlayer.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use specified file or find nexlayer.yaml in current directory
			if configFile == "" {
				files, err := filepath.Glob("*.yaml")
				if err != nil {
					return fmt.Errorf("failed to search for yaml files: %w", err)
				}

				for _, f := range files {
					if f != "docker-compose.yaml" && f != "docker-compose.yml" {
						configFile = f
						break
					}
				}

				if configFile == "" {
					return fmt.Errorf("no nexlayer.yaml file found in current directory")
				}
			}

			// Read nexlayer.yaml
			data, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to read config file %s: %w", configFile, err)
			}

			var config struct {
				Application struct {
					Template struct {
						Name           string `yaml:"name"`
						DeploymentName string `yaml:"deploymentname"`
						Pods          []struct {
							Type       string `yaml:"type"`
							Name       string `yaml:"name"`
							Tag        string `yaml:"tag"`
							ExposeHttp bool   `yaml:"exposehttp"`
							Vars       []struct {
								Key   string `yaml:"key"`
								Value string `yaml:"value"`
							} `yaml:"vars"`
						} `yaml:"pods"`
					} `yaml:"template"`
				} `yaml:"application"`
			}

			if err := yaml.Unmarshal(data, &config); err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}

			// Create docker-compose config
			compose := DockerCompose{
				Version:  "3.8",
				Services: make(map[string]Service),
				Networks: map[string]interface{}{
					"app-network": nil,
				},
			}

			// Add services for each pod
			for _, pod := range config.Application.Template.Pods {
				service := Service{
					Image:    getDefaultImage(pod.Type, pod.Tag),
					Networks: []string{"app-network"},
				}

				// Add environment variables
				for _, v := range pod.Vars {
					service.Environment = append(service.Environment, fmt.Sprintf("%s=%s", v.Key, v.Value))
				}

				// Add ports for HTTP-exposed services
				if pod.ExposeHttp {
					switch pod.Type {
					case "frontend":
						service.Ports = []string{"3000:3000"}
					case "backend":
						service.Ports = []string{"5000:5000"}
					case "nginx":
						service.Ports = []string{"80:80"}
					}
				}

				compose.Services[pod.Name] = service
			}

			// Write docker-compose.yml
			output, err := yaml.Marshal(compose)
			if err != nil {
				return fmt.Errorf("failed to generate docker-compose.yml: %w", err)
			}

			if err := os.WriteFile("docker-compose.yml", output, 0644); err != nil {
				return fmt.Errorf("failed to write docker-compose.yml: %w", err)
			}

			fmt.Println("Successfully generated docker-compose.yml")
			return nil
		},
	}

	cmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to nexlayer.yaml file")
	return cmd
}

func newUpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Start local development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDockerCompose("up", "-d")
		},
	}
}

func newDownCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Stop local development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDockerCompose("down")
		},
	}
}

func newLogsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "logs [service]",
		Short: "View service logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return runDockerCompose("logs", "-f", args[0])
			}
			return runDockerCompose("logs", "-f")
		},
	}
}

func runDockerCompose(args ...string) error {
	// Check if docker-compose.yml exists
	if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
		return fmt.Errorf("docker-compose.yml not found. Run 'nexlayer compose generate' first")
	}

	// Execute docker-compose command
	dockerCmd := exec.Command("docker-compose", args...)
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr
	return dockerCmd.Run()
}
