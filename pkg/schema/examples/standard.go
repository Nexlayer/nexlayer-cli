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
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "http",
							"port":       3000,
							"targetPort": 3000,
						},
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
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "http",
							"port":       8000,
							"targetPort": 8000,
						},
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
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "postgres",
							"port":       5432,
							"targetPort": 5432,
						},
					},
				},
				{
					Name:  "llm",
					Type:  "ollama",
					Image: "ollama/ollama:latest",
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "api",
							"port":       11434,
							"targetPort": 11434,
						},
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
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "http",
							"port":       3000,
							"targetPort": 3000,
						},
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
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "http",
							"port":       8000,
							"targetPort": 8000,
						},
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
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "http",
							"port":       6333,
							"targetPort": 6333,
						},
						map[string]interface{}{
							"name":       "grpc",
							"port":       6334,
							"targetPort": 6334,
						},
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
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "http",
							"port":       80,
							"targetPort": 80,
						},
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
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "http",
							"port":       3000,
							"targetPort": 3000,
						},
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
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "http",
							"port":       8000,
							"targetPort": 8000,
						},
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
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "postgres",
							"port":       5432,
							"targetPort": 5432,
						},
					},
				},
				{
					Name:  "cache",
					Type:  "redis",
					Image: "redis:alpine",
					ServicePorts: []interface{}{
						map[string]interface{}{
							"name":       "redis",
							"port":       6379,
							"targetPort": 6379,
						},
					},
				},
			},
		},
	}
}
