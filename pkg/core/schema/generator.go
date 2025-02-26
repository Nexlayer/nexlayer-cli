// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

import (
	"fmt"
	"strings"
)

// Generator handles the generation of Nexlayer YAML templates
type Generator struct {
	Registry string // Container registry (e.g., "ghcr.io/nexlayer")
	Tag      string // Default image tag
}

// NewGenerator creates a new template generator with default settings
func NewGenerator() *Generator {
	return &Generator{
		Registry: DefaultRegistry,
		Tag:      DefaultTag,
	}
}

// GenerateFromProjectInfo generates a template based on project information
func (g *Generator) GenerateFromProjectInfo(name, podType string, port int) (*NexlayerYAML, error) {
	// Clean the project name
	name = sanitizeName(name)

	// Create base template
	tmpl := &NexlayerYAML{
		Application: Application{
			Name: name,
			Pods: make([]Pod, 0),
		},
	}

	// Add pod based on type
	pod, err := g.createPodForType(name, podType, port)
	if err != nil {
		return nil, fmt.Errorf("failed to create pod: %w", err)
	}

	tmpl.Application.Pods = append(tmpl.Application.Pods, pod)
	return tmpl, nil
}

// createPodForType creates a pod configuration based on the project type
func (g *Generator) createPodForType(name, podType string, port int) (Pod, error) {
	pod := Pod{
		Name:  name,
		Type:  podType,
		Image: g.getImageForType(podType),
		ServicePorts: []ServicePort{
			{
				Name:       "http",
				Port:       port,
				TargetPort: port,
				Protocol:   "TCP",
			},
		},
	}

	// Add type-specific configuration
	switch podType {
	case "nextjs", "react":
		pod.Path = "/"
	case "express", "fastapi":
		pod.Path = "/api"
	}

	return pod, nil
}

// sanitizeName ensures the name follows Nexlayer naming conventions
func sanitizeName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace invalid characters with hyphens
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return '-'
	}, name)

	// Ensure starts with a letter
	if len(name) > 0 && (name[0] < 'a' || name[0] > 'z') {
		name = "app-" + name
	}

	// If empty after sanitization, use default
	if name == "" {
		name = "app"
	}

	return name
}

// AddPod adds a new pod to an existing template
func (g *Generator) AddPod(tmpl *NexlayerYAML, podType string, port int) error {
	pod, err := g.createPodForType(tmpl.Application.Name, podType, port)
	if err != nil {
		return fmt.Errorf("failed to create pod: %w", err)
	}

	// Ensure unique pod name
	existingNames := make(map[string]bool)
	for _, p := range tmpl.Application.Pods {
		existingNames[p.Name] = true
	}

	// If pod name exists, append a number
	baseName := pod.Name
	counter := 1
	for existingNames[pod.Name] {
		pod.Name = fmt.Sprintf("%s-%d", baseName, counter)
		counter++
	}

	tmpl.Application.Pods = append(tmpl.Application.Pods, pod)
	return nil
}

// SetRegistry updates the registry for all pods in a template
func (g *Generator) SetRegistry(tmpl *NexlayerYAML, registry string) {
	for i, pod := range tmpl.Application.Pods {
		// Only update images that use our registry
		if strings.HasPrefix(pod.Image, DefaultRegistry) {
			tmpl.Application.Pods[i].Image = strings.Replace(pod.Image, DefaultRegistry, registry, 1)
		}
	}
}

// SetTag updates the tag for all pods in a template
func (g *Generator) SetTag(tmpl *NexlayerYAML, tag string) {
	for i, pod := range tmpl.Application.Pods {
		// Extract current tag
		parts := strings.Split(pod.Image, ":")
		if len(parts) > 1 {
			// Replace tag
			tmpl.Application.Pods[i].Image = fmt.Sprintf("%s:%s", parts[0], tag)
		} else {
			// Add tag if none exists
			tmpl.Application.Pods[i].Image = fmt.Sprintf("%s:%s", pod.Image, tag)
		}
	}
}

// getImageForType returns the appropriate image for a given pod type
func (g *Generator) getImageForType(podType string) string {
	switch podType {
	case "nextjs":
		return fmt.Sprintf("%v/nextjs:%v", "<% REGISTRY %>", g.Tag)
	case "react":
		return fmt.Sprintf("%v/react:%v", "<% REGISTRY %>", g.Tag)
	case "express":
		return fmt.Sprintf("%v/express:%v", "<% REGISTRY %>", g.Tag)
	case "fastapi":
		return fmt.Sprintf("%v/fastapi:%v", "<% REGISTRY %>", g.Tag)
	default:
		return fmt.Sprintf("%v/%v:%v", "<% REGISTRY %>", podType, g.Tag)
	}
}

// GenerateFromTemplate generates a template from an existing template
func (g *Generator) GenerateFromTemplate(source *NexlayerYAML) (*NexlayerYAML, error) {
	if source == nil {
		return nil, fmt.Errorf("source configuration is nil")
	}

	// Create a copy of the source configuration
	copy := &NexlayerYAML{
		Application: Application{
			Name: source.Application.Name,
			URL:  source.Application.URL,
			Pods: make([]Pod, len(source.Application.Pods)),
		},
	}

	// Copy pods
	for i, pod := range source.Application.Pods {
		copy.Application.Pods[i] = pod
	}

	return copy, nil
}

// AddAIConfigurations adds AI-specific configurations to the template
func (g *Generator) AddAIConfigurations(tmpl *NexlayerYAML, provider string) {
	// Add environment variables for AI configuration
	for _, pod := range tmpl.Application.Pods {
		// Only add AI configs to application pods
		if pod.Type == "nextjs" || pod.Type == "react" || pod.Type == "express" || pod.Type == "fastapi" {
			pod.Vars = append(pod.Vars, EnvVar{
				Key:   "LLM_PROVIDER",
				Value: provider,
			})

			// Add provider-specific configurations
			switch provider {
			case "openai":
				pod.Vars = append(pod.Vars, EnvVar{
					Key:   "OPENAI_API_KEY",
					Value: "<% OPENAI_API_KEY %>",
				})
			case "anthropic":
				pod.Vars = append(pod.Vars, EnvVar{
					Key:   "ANTHROPIC_API_KEY",
					Value: "<% ANTHROPIC_API_KEY %>",
				})
			case "cohere":
				pod.Vars = append(pod.Vars, EnvVar{
					Key:   "COHERE_API_KEY",
					Value: "<% COHERE_API_KEY %>",
				})
			}
		}
	}
}
