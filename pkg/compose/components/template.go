package components

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

// Template represents the nexlayer.yaml template
type Template struct {
	Name           string `yaml:"name"`
	DeploymentName string `yaml:"deploymentName"`
	Pods           []Pod  `yaml:"pods"`
}

// GenerateTemplate creates a nexlayer.yaml template for the given project
func GenerateTemplate(projectName string, detector ComponentDetector) (string, error) {
	// Create basic template structure
	template := Template{
		Name:           projectName,
		DeploymentName: projectName,
		Pods:           []Pod{},
	}

	// Analyze current directory for components
	files, err := filepath.Glob("*")
	if err != nil {
		return "", fmt.Errorf("failed to scan directory: %w", err)
	}

	// Detect components based on files
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		if info.IsDir() {
			// Try to detect component type from directory
			// Try to detect component type
			detected, err := detector.DetectAndConfigure(Pod{
				Name: filepath.Base(file),
			})
			if err != nil {
				continue
			}

			pod := Pod{
				Name: filepath.Base(file),
				Type: detected.Type,
				Image: detected.Config.Image,
			}
			
			if pod.Type != "" {
				template.Pods = append(template.Pods, pod)
			}
		}
	}

	// Convert template to YAML
	yamlData, err := yaml.Marshal(template)
	if err != nil {
		return "", fmt.Errorf("failed to generate YAML: %w", err)
	}

	return string(yamlData), nil
}
