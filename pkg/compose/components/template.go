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

// Template represents the structure of the Nexlayer deployment template.
// It includes the template name, deployment name, and a list of pods.
type Template struct {
	Name           string `yaml:"name"`
	DeploymentName string `yaml:"deploymentName"`
	Pods           []Pod  `yaml:"pods"`
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
	template := Template{
		Name:           projectName,
		DeploymentName: projectName,
		Pods:           []Pod{},
	}

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

			pod := Pod{
				Name:  filepath.Base(file),
				Type:  detected.Type,
				Image: detected.Config.Image,
			}

			if pod.Type != "" {
				template.Pods = append(template.Pods, pod)
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
