package compose

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/Nexlayer/nexlayer-cli/pkg/compose/components"
)

type DockerCompose struct {
	Version  string                 `yaml:"version"`
	Services map[string]Service     `yaml:"services"`
	Networks map[string]interface{} `yaml:"networks"`
}

type Port struct {
	ContainerPort int    `yaml:"containerPort"`
	ServicePort   int    `yaml:"servicePort"`
	Name          string `yaml:"name"`
}

type Pod = components.Pod

type Service struct {
	Image       string                 `yaml:"image"`
	Build       string                 `yaml:"build,omitempty"`
	Ports       []string               `yaml:"ports,omitempty"`
	Environment []string               `yaml:"environment,omitempty"`
	DependsOn   []string               `yaml:"depends_on,omitempty"`
	Networks    []string               `yaml:"networks,omitempty"`
	Command     []string               `yaml:"command,omitempty"`
	Volumes     []string               `yaml:"volumes,omitempty"`
	Healthcheck map[string]interface{} `yaml:"healthcheck,omitempty"`
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compose",
		Short: "Manage local development with Docker Compose",
		Long:  `Commands for managing the local development environment using Docker Compose.`,
	}

	cmd.AddCommand(newGenerateCommand())
	cmd.AddCommand(newUpCommand())
	cmd.AddCommand(newDownCommand())
	cmd.AddCommand(newLogsCommand())

	return cmd
}

func configureService(pod components.Pod, detector *components.ComponentDetector) (Service, error) {
	service := Service{
		Networks: []string{"app-network"},
	}

	if pod.Image != "" {
		service.Image = pod.Image
	} else if pod.Tag != "" {
		service.Image = pod.Tag
	} else {
		detected, err := detector.DetectAndConfigure(components.Pod{
			Type:       pod.Type,
			Name:       pod.Name,
			Tag:        pod.Tag,
			ExposeOn80: pod.ExposeOn80,
			Vars:       pod.Vars,
		})
		if err != nil {
			return Service{}, fmt.Errorf("failed to detect component type: %w", err)
		}
		service.Image = detected.Config.Image
	}

	if len(pod.Ports) > 0 {
		for _, port := range pod.Ports {
			protocol := port.Protocol
			if protocol == "" {
				protocol = "tcp"
			}
			service.Ports = append(service.Ports,
				fmt.Sprintf("%d:%d/%s", port.Host, port.Container, protocol))
		}
	}

	service.Environment = append(service.Environment, convertVars(pod.Vars)...)

	if len(pod.Command) > 0 {
		service.Command = pod.Command
	}

	detected, err := detector.DetectAndConfigure(pod)
	if err != nil {
		return Service{}, fmt.Errorf("failed to detect component: %w", err)
	}

	for _, vol := range detected.Config.Volumes {
		service.Volumes = append(service.Volumes,
			fmt.Sprintf("%s:%s:%s", vol.Source, vol.Target, vol.Type))
	}

	if detected.Config.Healthcheck != nil {
		service.Healthcheck = map[string]interface{}{
			"test":     detected.Config.Healthcheck.Test,
			"interval": detected.Config.Healthcheck.Interval,
			"timeout":  detected.Config.Healthcheck.Timeout,
			"retries":  detected.Config.Healthcheck.Retries,
		}
	}

	return service, nil
}

func convertVars(vars []components.EnvVar) []string {
	result := make([]string, len(vars))
	for i, v := range vars {
		result[i] = fmt.Sprintf("%s=%s", v.Key, v.Value)
	}
	return result
}

func newGenerateCommand() *cobra.Command {
	var configFile string
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate docker-compose.yml from nexlayer.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			data, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to read config file %s: %w", configFile, err)
			}

			var config struct {
				Application struct {
					Template struct {
						Name           string           `yaml:"name"`
						DeploymentName string           `yaml:"deploymentname"`
						Pods           []components.Pod `yaml:"pods"`
					} `yaml:"template"`
				} `yaml:"application"`
			}

			if err := yaml.Unmarshal(data, &config); err != nil {
				return fmt.Errorf("failed to parse config file: %w", err)
			}

			compose := DockerCompose{
				Version:  "3.8",
				Services: make(map[string]Service),
				Networks: map[string]interface{}{
					"app-network": nil,
				},
			}

			detector := components.NewComponentDetector()

			for _, pod := range config.Application.Template.Pods {
				service, err := configureService(pod, detector)
				if err != nil {
					return fmt.Errorf("failed to configure service %s: %w", pod.Name, err)
				}

				service.Environment = append(service.Environment, convertVars(pod.Vars)...)

				if pod.ExposeOn80 && len(service.Ports) == 0 {
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
	if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
		return fmt.Errorf("docker-compose.yml not found. Run 'nexlayer compose generate' first")
	}

	dockerCmd := exec.Command("docker-compose", args...)
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr
	return dockerCmd.Run()
}
