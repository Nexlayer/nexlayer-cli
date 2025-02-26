// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultTemplatePath is the path to the reference template
	DefaultTemplatePath = "docs/reference/schemas/yaml/nexlayer-template.v1.yaml"
)

// Parser handles loading and merging of Nexlayer YAML templates
type Parser struct {
	templatePath string
	baseTemplate *NexlayerYAML
}

// NewParser creates a new template parser
func NewParser(templatePath string) (*Parser, error) {
	if templatePath == "" {
		templatePath = DefaultTemplatePath
	}

	// Ensure template exists
	if _, err := os.Stat(templatePath); err != nil {
		return nil, fmt.Errorf("template not found at %s: %w", templatePath, err)
	}

	return &Parser{
		templatePath: templatePath,
	}, nil
}

// LoadTemplate loads the base template
func (p *Parser) LoadTemplate() error {
	data, err := os.ReadFile(p.templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	var template NexlayerYAML
	if err := yaml.Unmarshal(data, &template); err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	p.baseTemplate = &template
	return nil
}

// MergeWithDetected merges the base template with detected settings
func (p *Parser) MergeWithDetected(detected *NexlayerYAML) (*NexlayerYAML, error) {
	if p.baseTemplate == nil {
		if err := p.LoadTemplate(); err != nil {
			return nil, err
		}
	}

	// Create a copy of the base template
	merged := *p.baseTemplate

	// Merge application settings
	if detected.Application.Name != "" {
		merged.Application.Name = detected.Application.Name
	}
	if detected.Application.URL != "" {
		merged.Application.URL = detected.Application.URL
	}
	if detected.Application.RegistryLogin != nil {
		merged.Application.RegistryLogin = detected.Application.RegistryLogin
	}

	// Merge pods
	if len(detected.Application.Pods) > 0 {
		merged.Application.Pods = mergePods(p.baseTemplate.Application.Pods, detected.Application.Pods)
	}

	return &merged, nil
}

// mergePods combines pod configurations from base and detected settings
func mergePods(basePods, detectedPods []Pod) []Pod {
	// Create a map of base pods by name for quick lookup
	basePodMap := make(map[string]Pod)
	for _, pod := range basePods {
		basePodMap[pod.Name] = pod
	}

	// Create merged pod list
	var mergedPods []Pod
	for _, pod := range detectedPods {
		if basePod, exists := basePodMap[pod.Name]; exists {
			mergedPods = append(mergedPods, mergePod(basePod, pod))
		} else {
			mergedPods = append(mergedPods, pod)
		}
	}

	return mergedPods
}

// mergePod combines settings from two pod configurations
func mergePod(base, detected Pod) Pod {
	merged := base

	// Merge fields if they are set in detected
	if detected.Type != "" {
		merged.Type = detected.Type
	}
	if detected.Path != "" {
		merged.Path = detected.Path
	}
	if detected.Image != "" {
		merged.Image = detected.Image
	}
	if detected.Command != "" {
		merged.Command = detected.Command
	}
	if detected.Entrypoint != "" {
		merged.Entrypoint = detected.Entrypoint
	}

	// Merge service ports
	if len(detected.ServicePorts) > 0 {
		merged.ServicePorts = detected.ServicePorts
	}

	// Merge environment variables
	if len(detected.Vars) > 0 {
		merged.Vars = detected.Vars
	}

	// Merge volumes
	if len(detected.Volumes) > 0 {
		merged.Volumes = detected.Volumes
	}

	// Merge secrets
	if len(detected.Secrets) > 0 {
		merged.Secrets = detected.Secrets
	}

	// Merge annotations
	if len(detected.Annotations) > 0 {
		if merged.Annotations == nil {
			merged.Annotations = make(map[string]string)
		}
		for k, v := range detected.Annotations {
			merged.Annotations[k] = v
		}
	}

	return merged
}

// mergeEnvVars combines environment variables from two configurations
func mergeEnvVars(base, detected []EnvVar) []EnvVar {
	envMap := make(map[string]string)

	// Add base vars to map
	for _, env := range base {
		envMap[env.Key] = env.Value
	}

	// Merge or add detected vars
	for _, env := range detected {
		envMap[env.Key] = env.Value
	}

	// Convert map back to slice
	var merged []EnvVar
	for k, v := range envMap {
		merged = append(merged, EnvVar{
			Key:   k,
			Value: v,
		})
	}

	return merged
}
