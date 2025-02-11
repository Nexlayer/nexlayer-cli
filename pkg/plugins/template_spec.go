// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

// TemplateSpec represents the Nexlayer template specification
type TemplateSpec struct {
	Version  string                 `yaml:"version"`
	Schema   map[string]interface{} `yaml:"schema"`
	Components map[string]interface{} `yaml:"componentTypes"`
	Examples map[string]interface{} `yaml:"examples"`
	BestPractices map[string]interface{} `yaml:"bestPractices"`
	Validation map[string]interface{} `yaml:"validation"`
}

// LoadTemplateSpec loads the template specification from the YAML file
func LoadTemplateSpec() (*TemplateSpec, error) {
	// Find the spec file relative to the executable
	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	specPath := filepath.Join(filepath.Dir(execPath), "..", "docs", "nexlayer_template_reference.yaml")
	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template spec file: %w", err)
	}

	var spec TemplateSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse template spec: %w", err)
	}

	return &spec, nil
}

// GetComponentDefaults returns the default configuration for a given component type
func (s *TemplateSpec) GetComponentDefaults(componentType, subType string) (map[string]interface{}, error) {
	components, ok := s.Components[componentType]
	if !ok {
		return nil, fmt.Errorf("unsupported component type: %s", componentType)
	}

	componentsMap, ok := components.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid component type configuration")
	}

	supported, ok := componentsMap["supported"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid supported components configuration")
	}

	defaults, ok := supported[subType].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unsupported component sub-type: %s", subType)
	}

	return defaults, nil
}

// ValidateTemplate validates a template against the specification
func (s *TemplateSpec) ValidateTemplate(template map[string]interface{}) []string {
	var violations []string

	// Check required fields
	if app, ok := template["application"].(map[string]interface{}); ok {
		// Check application.name
		if _, ok := app["name"]; !ok {
			violations = append(violations, "missing required field: application.name")
		}

		// Check application.pods
		if pods, ok := app["pods"].([]interface{}); ok {
			for i, pod := range pods {
				podMap := pod.(map[string]interface{})
				// Check required pod fields
				if _, ok := podMap["name"]; !ok {
					violations = append(violations, fmt.Sprintf("missing required field: application.pods[%d].name", i))
				}
				if _, ok := podMap["image"]; !ok {
					violations = append(violations, fmt.Sprintf("missing required field: application.pods[%d].image", i))
				}
				if _, ok := podMap["servicePorts"]; !ok {
					violations = append(violations, fmt.Sprintf("missing required field: application.pods[%d].servicePorts", i))
				}
			}
		} else {
			violations = append(violations, "missing required field: application.pods")
		}
	} else {
		violations = append(violations, "missing required field: application")
	}

	// Add more validation logic based on the specification...

	return violations
}
