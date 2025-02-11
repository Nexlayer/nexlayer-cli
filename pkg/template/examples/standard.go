// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package examples

import "github.com/Nexlayer/nexlayer-cli/pkg/template"

// StandardTemplate returns a standard example template that follows
// the Nexlayer YAML schema v2 format
func StandardTemplate() *template.NexlayerYAML {
	return &template.NexlayerYAML{
		Application: template.Application{
			// REQUIRED: Unique deployment name
			Name: "my-app",
			// OPTIONAL: Permanent domain
			URL:  "my-app.nexlayer.dev",
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
					Name:  "frontend",
					// OPTIONAL: Route path for frontend
					Path:  "/",
					// REQUIRED: Pod type
					Type:  template.React,
					// REQUIRED: Fully qualified image path
					Image: "docker.io/my-org/frontend:latest",
					Vars: []template.EnvVar{
						{Key: "API_URL", Value: "http://backend.pod:8000"},
						{Key: "NODE_ENV", Value: "production"},
					},
					Ports: []template.Port{
						{
							ContainerPort: 3000,
							ServicePort:   3000,
							Name:          "http",
						},
					},
				},
				{
					// REQUIRED: Pod name (lowercase alphanumeric)
					Name:  "backend",
					// OPTIONAL: Route path for API
					Path:  "/api",
					// REQUIRED: Pod type
					Type:  template.FastAPI,
					// REQUIRED: Fully qualified image path
					Image: "docker.io/my-org/backend:latest",
					Vars: []template.EnvVar{
						{Key: "DATABASE_URL", Value: "postgres://user:pass@db.pod:5432/db"},
						{Key: "PORT", Value: "8000"},
					},
					Ports: []template.Port{
						{
							ContainerPort: 8000,
							ServicePort:   8000,
							Name:          "http",
						},
					},
				},
				{
					Name:  "db",
					Type:  template.Postgres,
					Image: "postgres:latest",
					Volumes: []template.Volume{
						{
							Name:      "pg-data-volume",
							Size:      "5Gi",
							MountPath: "/var/lib/postgresql/data",
						},
					},
					Ports: []template.Port{
						{
							ContainerPort: 5432,
							ServicePort:   5432,
							Name:          "postgresql",
						},
					},
				},
				{
					Name:  "llm",
					Type:  template.Ollama,
					Image: "ollama/ollama:latest",
					Ports: []template.Port{
						{
							ContainerPort: 11434,
							ServicePort:   11434,
							Name:          "ollama",
						},
					},
				},
			},
		},
	}
}
