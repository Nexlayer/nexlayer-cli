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
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate docker-compose.yml from nexlayer.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find nexlayer.yaml in current directory
			files, err := filepath.Glob("*.yaml")
			if err != nil {
				return err
			}

			var nexlayerFile string
			for _, f := range files {
				if f != "docker-compose.yaml" && f != "docker-compose.yml" {
					nexlayerFile = f
					break
				}
			}

			if nexlayerFile == "" {
				return fmt.Errorf("no nexlayer.yaml file found in current directory")
			}

			// Read nexlayer.yaml
			data, err := os.ReadFile(nexlayerFile)
			if err != nil {
				return err
			}

			var config struct {
				Application struct {
					Template struct {
						Pods []struct {
							Type       string `yaml:"type"`
							Name       string `yaml:"name"`
							Tag        string `yaml:"tag"`
							ExposeHttp bool   `yaml:"exposeHttp"`
							Vars       []struct {
								Key   string `yaml:"key"`
								Value string `yaml:"value"`
							} `yaml:"vars"`
						} `yaml:"pods"`
					} `yaml:"template"`
				} `yaml:"application"`
			}

			if err := yaml.Unmarshal(data, &config); err != nil {
				return err
			}

			// Create docker-compose.yml
			compose := DockerCompose{
				Version:  "3.8",
				Services: make(map[string]Service),
				Networks: map[string]interface{}{
					"nexlayer": nil,
				},
			}

			// Convert pods to services
			for _, pod := range config.Application.Template.Pods {
				service := Service{
					Image:       getDefaultImage(pod.Type, pod.Tag),
					Networks:    []string{"nexlayer"},
					Environment: make([]string, 0),
				}

				// Add environment variables
				for _, v := range pod.Vars {
					service.Environment = append(service.Environment, fmt.Sprintf("%s=%s", v.Key, v.Value))
				}

				// Add ports if service exposes HTTP
				if pod.ExposeHttp {
					for _, v := range pod.Vars {
						if v.Key == "PORT" {
							service.Ports = append(service.Ports, fmt.Sprintf("%s:%s", v.Value, v.Value))
							break
						}
					}
				}

				compose.Services[pod.Name] = service
			}

			// Write docker-compose.yml
			output, err := yaml.Marshal(compose)
			if err != nil {
				return err
			}

			return os.WriteFile("docker-compose.yml", output, 0644)
		},
	}
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
