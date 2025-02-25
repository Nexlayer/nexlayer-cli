// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package template

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
		Application: ApplicationYAML{
			Name: name,
			Pods: make([]PodYAML, 0),
		},
	}

	// Add pod based on type
	pod, err := g.createPodForType(name, podType, port)
	if err != nil {
		return nil, fmt.Errorf("failed to create pod: %w", err)
	}

	// Convert Pod to PodYAML
	podYAML := PodYAML{
		Name:         pod.Name,
		Type:         pod.Type,
		Path:         pod.Path,
		Image:        pod.Image,
		Command:      pod.Command,
		Entrypoint:   pod.Entrypoint,
		ServicePorts: pod.ServicePorts,
		Vars:         pod.Vars,
		Volumes:      pod.Volumes,
		Secrets:      pod.Secrets,
		Annotations:  pod.Annotations,
	}
	tmpl.Application.Pods = append(tmpl.Application.Pods, podYAML)

	return tmpl, nil
}

// createPodForType creates a pod configuration based on the project type
func (g *Generator) createPodForType(name, podType string, port int) (*Pod, error) {
	// Use default port if none specified
	if port == 0 {
		if defaultPort, ok := DefaultPorts[podType]; ok {
			port = defaultPort
		} else {
			port = 8080 // Fallback default
		}
	}

	// Create base pod
	pod := &Pod{
		Name: "web",
		Type: podType,
		// Use schema-compliant registry format for private images
		Image: fmt.Sprintf("%s/%s:%s", RegistryPlaceholder, name, g.Tag),
		ServicePorts: []ServicePort{
			{
				Name:       "http",
				Port:       port,
				TargetPort: port,
				Protocol:   ProtocolTCP,
			},
		},
	}

	// Add default environment variables
	if defaultVars, ok := DefaultEnvVars[podType]; ok {
		pod.Vars = defaultVars
	}

	// Configure pod based on type
	switch podType {
	case PodTypeNextJS, PodTypeReact, PodTypeVue:
		pod.Path = "/"
		pod.Type = podType // Ensure specific frontend type
		// Add URL reference for frontend
		pod.Vars = append(pod.Vars, EnvVar{
			Key:   "PUBLIC_URL",
			Value: URLPlaceholder,
		})

	case PodTypeNode, PodTypeExpress:
		pod.Name = "api"
		pod.Path = "/api"
		pod.Type = podType // Keep original type for Node.js
		// Add dynamic pod reference for database if needed
		pod.Vars = append(pod.Vars, EnvVar{
			Key:   "DATABASE_URL",
			Value: "postgresql://postgres:postgres@postgres.pod:5432/app",
		})

	case PodTypePython, PodTypeDjango, PodTypeFastAPI:
		pod.Name = "api"
		pod.Path = "/api"
		pod.Type = PodTypeBackend
		// Add dynamic pod reference for database if needed
		pod.Vars = append(pod.Vars, EnvVar{
			Key:   "DATABASE_URL",
			Value: "postgresql://postgres:postgres@postgres.pod:5432/app",
		})

	case PodTypeGolang:
		pod.Name = "api"
		pod.Path = "/api"
		pod.Type = PodTypeBackend
		// Add dynamic pod reference for database if needed
		pod.Vars = append(pod.Vars, EnvVar{
			Key:   "DATABASE_URL",
			Value: "postgresql://postgres:postgres@postgres.pod:5432/app",
		})

	case PodTypePostgres, PodTypeMongoDB, PodTypeRedis:
		pod.Name = strings.ToLower(strings.TrimPrefix(podType, "PodType"))
		pod.Type = podType
		// Use public image for databases
		pod.Image = fmt.Sprintf("%s:latest", pod.Name)
		// Add default volume for databases
		pod.Volumes = []Volume{
			{
				Name:     fmt.Sprintf("%s-data", pod.Name),
				Path:     fmt.Sprintf("/var/lib/%s/data", pod.Name),
				Size:     "1Gi",
				Type:     VolumeTypePersistent,
				ReadOnly: false,
			},
		}
		// Add default environment variables for databases
		switch podType {
		case PodTypePostgres:
			pod.Vars = append(pod.Vars, []EnvVar{
				{Key: "POSTGRES_USER", Value: "postgres"},
				{Key: "POSTGRES_PASSWORD", Value: "<% DB_PASSWORD %>"},
				{Key: "POSTGRES_DB", Value: "app"},
			}...)
		case PodTypeMongoDB:
			pod.Vars = append(pod.Vars, []EnvVar{
				{Key: "MONGO_INITDB_ROOT_USERNAME", Value: "root"},
				{Key: "MONGO_INITDB_ROOT_PASSWORD", Value: "<% MONGO_ROOT_PASSWORD %>"},
			}...)
		case PodTypeRedis:
			pod.Vars = append(pod.Vars, EnvVar{
				Key:   "REDIS_PASSWORD",
				Value: "<% REDIS_PASSWORD %>",
			})
		}
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

	// Convert Pod to PodYAML
	podYAML := PodYAML{
		Name:         pod.Name,
		Type:         pod.Type,
		Path:         pod.Path,
		Image:        pod.Image,
		Command:      pod.Command,
		Entrypoint:   pod.Entrypoint,
		ServicePorts: pod.ServicePorts,
		Vars:         pod.Vars,
		Volumes:      pod.Volumes,
		Secrets:      pod.Secrets,
		Annotations:  pod.Annotations,
	}
	tmpl.Application.Pods = append(tmpl.Application.Pods, podYAML)
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
