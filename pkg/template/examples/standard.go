// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package examples

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/template"
)

// StandardTemplate returns a standard example template that follows
// the Nexlayer YAML schema v1.0 format
func StandardTemplate() *template.NexlayerYAML {
	return &template.NexlayerYAML{
		Application: template.Application{
			// REQUIRED: Unique deployment name
			Name: "my-app",
			// OPTIONAL: Permanent domain
			URL: "my-app.nexlayer.dev",
			// REQUIRED for private images
			RegistryLogin: &template.RegistryLogin{
				Registry:            "docker.io/my-org",
				Username:            "myuser",
				PersonalAccessToken: "mytoken",
			},
			// REQUIRED: List of pod configurations
			Pods: []template.Pod{
				{
					// REQUIRED: Pod name (lowercase alphanumeric)
					Name: "frontend",
					// OPTIONAL: Route path for frontend
					Path: "/",
					// REQUIRED: Pod type
					Type: "react",
					// REQUIRED: Fully qualified image path
					Image: fmt.Sprintf("%s/frontend:latest", template.RegistryPlaceholder),
					Vars: []template.EnvVar{
						{Key: "API_URL", Value: "http://backend.pod:8000"},
						{Key: "NODE_ENV", Value: "production"},
					},
					ServicePorts: []template.ServicePort{
						{Name: "http", Port: 3000, TargetPort: 3000},
					},
				},
				{
					// REQUIRED: Pod name (lowercase alphanumeric)
					Name: "backend",
					// OPTIONAL: Route path for API
					Path: "/api",
					// REQUIRED: Pod type
					Type: "fastapi",
					// REQUIRED: Fully qualified image path
					Image: fmt.Sprintf("%s/backend:latest", template.RegistryPlaceholder),
					Vars: []template.EnvVar{
						{Key: "DATABASE_URL", Value: "postgresql://user:pass@db.pod:5432/db"},
						{Key: "PORT", Value: "8000"},
					},
					ServicePorts: []template.ServicePort{
						{Name: "http", Port: 8000, TargetPort: 8000},
					},
				},
				{
					Name:  "db",
					Type:  "postgres",
					Image: "postgres:latest",
					Volumes: []template.Volume{
						{
							Name: "pg-data-volume",
							Path: "/var/lib/postgresql/data",
							Size: "5Gi",
						},
					},
					ServicePorts: []template.ServicePort{
						{Name: "postgres", Port: 5432, TargetPort: 5432},
					},
				},
				{
					Name:  "llm",
					Type:  "ollama",
					Image: "ollama/ollama:latest",
					ServicePorts: []template.ServicePort{
						{Name: "api", Port: 11434, TargetPort: 11434},
					},
				},
			},
		},
	}
}
