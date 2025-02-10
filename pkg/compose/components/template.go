// template.go
// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

package components

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// RegistryLogin represents private registry authentication details
type RegistryLogin struct {
	Registry           string `yaml:"registry"`
	Username           string `yaml:"username"`
	PersonalAccessToken string `yaml:"personalAccessToken"`
}



// Template represents the structure of the Nexlayer deployment template.
// It follows the v2.0 schema specification.
type Template struct {
	Application struct {
		Name          string        `yaml:"name"`
		URL           string        `yaml:"url,omitempty"`
		RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty"`
		Pods          []Pod         `yaml:"pods"`
	} `yaml:"application"`
}

// GenerateTemplate creates a nexlayer.yaml template for the given project.
// It uses the provided ComponentDetector to scan the current directory for components,
// constructs a Template struct, and marshals it into YAML format.
// Returns the YAML string or an error if template generation fails.
func GenerateTemplate(projectName string, detector ComponentDetector) (string, error) {
	// Validate input parameters.
	if projectName == "" {
		return "", fmt.Errorf("project name cannot be empty")
	}
	if detector == nil {
		return "", fmt.Errorf("component detector cannot be nil")
	}

	// Create a basic template with the project name.
	template := Template{}
	template.Application.Name = projectName

	// Scan the current directory for components.
	files, err := filepath.Glob("*")
	if err != nil {
		return "", fmt.Errorf("failed to scan directory: %w", err)
	}

	// For each directory, attempt to detect a component.
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		if info.IsDir() {
			// Detect component type based on the directory name.
			detected, err := detector.DetectAndConfigure(Pod{
				Name: filepath.Base(file),
			})
			if err != nil {
				// Skip directories that do not yield a valid component.
				continue
			}

			// Create pod configuration with detected settings
			pod := Pod{
				Name:         filepath.Base(file),
				Image:        detected.Config.Image,
				Path:         detected.Config.Path,
				ServicePorts: detected.Config.ServicePorts,
				Volumes:      detected.Config.Volumes,
				Secrets:      detected.Config.Secrets,
				Vars:         detected.Config.Environment,
			}

			// Add pod to template if it has valid configuration
			if pod.Image != "" {
				template.Application.Pods = append(template.Application.Pods, pod)
			}
		}
	}

	// Marshal the Template struct into YAML.
	yamlData, err := yaml.Marshal(template)
	if err != nil {
		return "", fmt.Errorf("failed to generate YAML: %w", err)
	}

	return string(yamlData), nil
}
