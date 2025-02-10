// compose.go
// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

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

// 🔥 Ensure the schema matches `nexlayer.yaml`
type NexlayerConfig struct {
	Application struct {
		Name          string `yaml:"name"`
		URL           string `yaml:"url,omitempty"`
		RegistryLogin struct {
			Registry            string `yaml:"registry"`
			Username            string `yaml:"username"`
			PersonalAccessToken string `yaml:"personalAccessToken"`
		} `yaml:"registryLogin,omitempty"`
		Pods []components.Pod `yaml:"pods"`
	} `yaml:"application"`
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

// 🛠 Configure a Docker Compose service from a Nexlayer Pod
func configureService(pod components.Pod) (Service, error) {
	service := Service{
		Image:    pod.Image,
		Networks: []string{"app-network"},
	}

	// 🌟 Convert environment variables properly
	if len(pod.Vars) > 0 {
		service.Environment = convertVars(pod.Vars)
	}

	// 🔥 Convert ports properly
	if len(pod.ServicePorts) > 0 {
		for _, port := range pod.ServicePorts {
			service.Ports = append(service.Ports,
				fmt.Sprintf("%d:%d", port, port))
		}
	}

	// 📂 Handle volumes correctly
	if len(pod.Volumes) > 0 {
		for _, vol := range pod.Volumes {
			service.Volumes = append(service.Volumes,
				fmt.Sprintf("%s:%s", vol.Name, vol.MountPath))
		}
	}

	// 🔑 Handle secrets properly
	if len(pod.Secrets) > 0 {
		for _, secret := range pod.Secrets {
			service.Volumes = append(service.Volumes,
				fmt.Sprintf("%s:%s/%s:ro", secret.Name, secret.MountPath, secret.FileName))
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

// 📌 Generate `docker-compose.yml` from `nexlayer.yaml`
func newGenerateCommand() *cobra.Command {
	var configFile string
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate docker-compose.yml from nexlayer.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if configFile == "" {
				files, err := filepath.Glob("*.yaml")
				if err != nil {
					return fmt.Errorf("failed to search for YAML files: %w", err)
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

			var config NexlayerConfig

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

			for _, pod := range config.Application.Pods {
				service, err := configureService(pod)
				if err != nil {
					return fmt.Errorf("failed to configure service %s: %w", pod.Name, err)
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

			fmt.Println("✅ Successfully generated docker-compose.yml")
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
		return fmt.Errorf("❌ docker-compose.yml not found. Run 'nexlayer compose generate' first")
	}

	dockerCmd := exec.Command("docker-compose", args...)
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr
	return dockerCmd.Run()
}
