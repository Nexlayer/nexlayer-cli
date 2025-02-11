// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package examples

import "github.com/Nexlayer/nexlayer-cli/pkg/template"

// StandardTemplate returns a standard example template
func StandardTemplate() *template.NexlayerYAML {
	return &template.NexlayerYAML{
		Application: template.Application{
			Name: "my-app",
			URL:  "https://myapp.example.com",
			RegistryLogin: &template.RegistryLogin{
				Registry:            "ghcr.io/my-org",
				Username:            "myuser",
				PersonalAccessToken: "mytoken",
			},
			Pods: []template.Pod{
				{
					Name:  "frontend",
					Path:  "/",
					Type:  template.React,
					Image: "<% REGISTRY %>/frontend:latest",
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
					Name:  "backend",
					Path:  "/api",
					Type:  template.FastAPI,
					Image: "<% REGISTRY %>/backend:latest",
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
