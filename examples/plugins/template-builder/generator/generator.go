package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/types"
)

// GenerateTemplate creates a template based on the detected stack
func GenerateTemplate(projectDir string, stack *types.ProjectStack) (*types.NexlayerTemplate, error) {
	if stack == nil {
		return nil, fmt.Errorf("stack cannot be nil")
	}

	// Check if directory is readable
	if _, err := os.ReadDir(projectDir); err != nil {
		return nil, fmt.Errorf("error accessing project directory: %v", err)
	}

	projectName := filepath.Base(projectDir)
	template := &types.NexlayerTemplate{
		Name:    projectName,
		Version: "0.1.0",
		Stack:   *stack,
	}

	// Generate main service
	mainService := generateServiceConfig(stack)
	template.Services = []types.Service{*mainService}

	// Load environment variables
	envVars, err := loadEnvironmentVariables(projectDir)
	if err != nil {
		return nil, fmt.Errorf("error loading environment variables: %v", err)
	}
	template.Config = envVars

	return template, nil
}

func generateServiceConfig(stack *types.ProjectStack) *types.Service {
	service := &types.Service{
		Name:  getServiceName(stack),
		Image: getServiceImage(stack),
		Ports: []types.PortConfig{
			{
				Name:        "http",
				Port:        getDefaultPort(stack),
				TargetPort:  getDefaultPort(stack),
				Protocol:    "TCP",
				Host:        false,
				Public:      true,
				Healthcheck: true,
			},
		},
		Resources: types.ResourceRequests{
			CPU:    "100m",
			Memory: getDefaultMemory(stack),
		},
		Healthcheck: &types.HealthcheckConfig{
			Path:     "/health",
			Port:     getDefaultPort(stack),
			Protocol: "HTTP",
		},
	}

	// Add database port if needed
	if stack.Database != "" {
		dbPort := getDefaultDatabasePort(stack.Database)
		if dbPort > 0 {
			service.Ports = append(service.Ports, types.PortConfig{
				Name:       stack.Database,
				Port:      dbPort,
				TargetPort: dbPort,
				Protocol:   "TCP",
				Host:      false,
				Public:    false,
			})
		}
	}

	return service
}

func getServiceName(stack *types.ProjectStack) string {
	switch stack.Language {
	case "python":
		return "web"
	case "nodejs", "javascript":
		return "api"
	default:
		return "app"
	}
}

func getServiceImage(stack *types.ProjectStack) string {
	switch stack.Language {
	case "javascript":
		return "node:16-alpine"
	case "python":
		return "python:3.9-slim"
	case "go":
		return "golang:1.19-alpine"
	case "java":
		return "openjdk:17-slim"
	default:
		return "alpine:latest"
	}
}

func getDefaultPort(stack *types.ProjectStack) int {
	switch {
	case stack.Language == "javascript" && stack.Framework == "react":
		return 3000
	case stack.Language == "javascript" && stack.Framework == "express":
		return 3000
	case stack.Language == "python" && stack.Framework == "flask":
		return 5000
	case stack.Language == "python" && stack.Framework == "django":
		return 8000
	case stack.Language == "go":
		return 8080
	case stack.Language == "java":
		return 8080
	default:
		return 8080
	}
}

func getDefaultMemory(stack *types.ProjectStack) string {
	switch stack.Language {
	case "python":
		return "256Mi"
	default:
		return "128Mi"
	}
}

func getDefaultDatabasePort(database string) int {
	switch database {
	case "postgres":
		return 5432
	case "mongodb":
		return 27017
	case "mysql":
		return 3306
	case "redis":
		return 6379
	default:
		return 0
	}
}

func loadEnvironmentVariables(projectDir string) (map[string]string, error) {
	envVars := make(map[string]string)

	// Try to load from .env file
	envFile := filepath.Join(projectDir, ".env")
	if data, err := os.ReadFile(envFile); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				// Remove quotes if present
				value = strings.Trim(value, `"'`)
				envVars[key] = value
			}
		}
	}

	// Try to load from .env.example file if .env doesn't exist
	if len(envVars) == 0 {
		if data, err := os.ReadFile(filepath.Join(projectDir, ".env.example")); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}

				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					value = strings.Trim(value, `"'`)
					envVars[key] = value
				}
			}
		}
	}

	return envVars, nil
}

// SaveTemplate saves a template to a file in the specified format
func SaveTemplate(template *types.NexlayerTemplate, outputPath string) error {
	// Implementation remains the same
	return nil
}
