// Package examples provides example Nexlayer YAML configurations
package examples

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/schema"
)

// StandardTemplate returns a standard example template that follows
// the Nexlayer YAML schema format
func StandardTemplate() *schema.NexlayerYAML {
	return &schema.NexlayerYAML{
		Application: schema.Application{
			// REQUIRED: Unique deployment name
			Name: "my-app",
			// OPTIONAL: Permanent domain
			URL: "my-app.nexlayer.dev",
			// REQUIRED for private images
			RegistryLogin: &schema.RegistryLogin{
				Registry:            "docker.io/my-org",
				Username:            "myuser",
				PersonalAccessToken: "mytoken",
			},
			// REQUIRED: List of pod configurations
			Pods: []schema.Pod{
				{
					// REQUIRED: Pod name (lowercase alphanumeric)
					Name: "frontend",
					// OPTIONAL: Route path for frontend
					Path: "/",
					// REQUIRED: Pod type
					Type: "react",
					// REQUIRED: Fully qualified image path
					Image: "<% REGISTRY %>/frontend:latest",
					Vars: []schema.EnvVar{
						{Key: "API_URL", Value: "http://backend.pod:8000"},
						{Key: "NODE_ENV", Value: "production"},
					},
					ServicePorts: []schema.ServicePort{
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
					Image: "<% REGISTRY %>/backend:latest",
					Vars: []schema.EnvVar{
						{Key: "DATABASE_URL", Value: "postgresql://user:pass@db.pod:5432/db"},
						{Key: "PORT", Value: "8000"},
					},
					ServicePorts: []schema.ServicePort{
						{Name: "http", Port: 8000, TargetPort: 8000},
					},
				},
				{
					Name:  "db",
					Type:  "postgres",
					Image: "postgres:latest",
					Volumes: []schema.Volume{
						{
							Name: "pg-data-volume",
							Path: "/var/lib/postgresql/data",
							Size: "5Gi",
						},
					},
					ServicePorts: []schema.ServicePort{
						{Name: "postgres", Port: 5432, TargetPort: 5432},
					},
				},
				{
					Name:  "llm",
					Type:  "ollama",
					Image: "ollama/ollama:latest",
					ServicePorts: []schema.ServicePort{
						{Name: "api", Port: 11434, TargetPort: 11434},
					},
				},
			},
		},
	}
}

// AITemplate returns an example template for an AI-powered application
func AITemplate() *schema.NexlayerYAML {
	return &schema.NexlayerYAML{
		Application: schema.Application{
			Name: "ai-app",
			URL:  "ai-app.nexlayer.dev",
			// REQUIRED for private images
			RegistryLogin: &schema.RegistryLogin{
				Registry:            "docker.io/my-ai-org",
				Username:            "aiuser",
				PersonalAccessToken: "aitoken",
			},
			Pods: []schema.Pod{
				{
					Name:  "web",
					Path:  "/",
					Type:  "langchain-nextjs",
					Image: "<% REGISTRY %>/web:latest",
					Vars: []schema.EnvVar{
						{Key: "API_URL", Value: "http://api.pod:8000"},
						{Key: "NODE_ENV", Value: "production"},
						{Key: "OPENAI_API_KEY", Value: "<% OPENAI_API_KEY %>"},
					},
					ServicePorts: []schema.ServicePort{
						{Name: "http", Port: 3000, TargetPort: 3000},
					},
					Annotations: map[string]string{
						"ai.nexlayer.io/provider": "openai",
						"ai.nexlayer.io/enabled":  "true",
					},
				},
				{
					Name:  "api",
					Path:  "/api",
					Type:  "openai-node",
					Image: "<% REGISTRY %>/api:latest",
					Vars: []schema.EnvVar{
						{Key: "PORT", Value: "8000"},
						{Key: "OPENAI_API_KEY", Value: "<% OPENAI_API_KEY %>"},
						{Key: "VECTOR_DB_URL", Value: "http://vector-db.pod:6333"},
					},
					ServicePorts: []schema.ServicePort{
						{Name: "http", Port: 8000, TargetPort: 8000},
					},
					Annotations: map[string]string{
						"ai.nexlayer.io/provider": "openai",
						"ai.nexlayer.io/enabled":  "true",
					},
				},
				{
					Name:  "vector-db",
					Type:  "qdrant",
					Image: "qdrant/qdrant:latest",
					Volumes: []schema.Volume{
						{
							Name: "vector-data",
							Path: "/qdrant/storage",
							Size: "10Gi",
						},
					},
					ServicePorts: []schema.ServicePort{
						{Name: "http", Port: 6333, TargetPort: 6333},
						{Name: "grpc", Port: 6334, TargetPort: 6334},
					},
				},
			},
		},
	}
}

// MicroservicesTemplate returns an example template for a microservices application
func MicroservicesTemplate() *schema.NexlayerYAML {
	return &schema.NexlayerYAML{
		Application: schema.Application{
			Name: "micro-app",
			URL:  "micro-app.nexlayer.dev",
			// REQUIRED for private images
			RegistryLogin: &schema.RegistryLogin{
				Registry:            "docker.io/my-micro-org",
				Username:            "microuser",
				PersonalAccessToken: "microtoken",
			},
			Pods: []schema.Pod{
				{
					Name:  "gateway",
					Path:  "/",
					Type:  "nginx",
					Image: "nginx:alpine",
					ServicePorts: []schema.ServicePort{
						{Name: "http", Port: 80, TargetPort: 80},
					},
				},
				{
					Name:  "auth",
					Path:  "/auth",
					Type:  "node",
					Image: "<% REGISTRY %>/auth:latest",
					Vars: []schema.EnvVar{
						{Key: "PORT", Value: "3000"},
						{Key: "JWT_SECRET", Value: "<% JWT_SECRET %>"},
						{Key: "REDIS_URL", Value: "redis://cache.pod:6379"},
					},
					ServicePorts: []schema.ServicePort{
						{Name: "http", Port: 3000, TargetPort: 3000},
					},
				},
				{
					Name:  "users",
					Path:  "/users",
					Type:  "go",
					Image: "<% REGISTRY %>/users:latest",
					Vars: []schema.EnvVar{
						{Key: "PORT", Value: "8000"},
						{Key: "DB_URL", Value: "postgresql://user:pass@users-db.pod:5432/users"},
					},
					ServicePorts: []schema.ServicePort{
						{Name: "http", Port: 8000, TargetPort: 8000},
					},
				},
				{
					Name:  "users-db",
					Type:  "postgres",
					Image: "postgres:latest",
					Volumes: []schema.Volume{
						{
							Name: "users-data",
							Path: "/var/lib/postgresql/data",
							Size: "5Gi",
						},
					},
					ServicePorts: []schema.ServicePort{
						{Name: "postgres", Port: 5432, TargetPort: 5432},
					},
				},
				{
					Name:  "cache",
					Type:  "redis",
					Image: "redis:alpine",
					ServicePorts: []schema.ServicePort{
						{Name: "redis", Port: 6379, TargetPort: 6379},
					},
				},
			},
		},
	}
}
