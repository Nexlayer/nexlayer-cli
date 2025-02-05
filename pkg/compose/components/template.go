// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package components

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

// Package components provides structures and functions for handling Nexlayer components.

// Template represents the nexlayer.yaml template
// It contains the name, deployment name, and a list of pods.
// Name is the name of the template.
// DeploymentName is the name used for deployment.
// Pods is a list of Pod structures representing application components.
type Template struct {
	Name           string `yaml:"name"`
	DeploymentName string `yaml:"deploymentName"`
	Pods           []Pod  `yaml:"pods"`
}

// GenerateTemplate creates a nexlayer.yaml template for the given project.
// It takes the project name and a ComponentDetector as parameters.
// The function scans the current directory for components, detects their types,
// and constructs a Template struct. It then marshals the struct to YAML format.
// Returns the YAML string or an error if generation fails.
func GenerateTemplate(projectName string, detector ComponentDetector) (string, error) {
	// Validate input parameters
	if projectName == "" {
		return "", fmt.Errorf("project name cannot be empty")
	}
	if detector == nil {
		return "", fmt.Errorf("component detector cannot be nil")
	}

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
				Name:  filepath.Base(file),
				Type:  detected.Type,
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
