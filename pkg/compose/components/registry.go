// registry.go
// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

package components

import "fmt"

// NexlayerRegistry defines the default registry URL for Nexlayer components.
const NexlayerRegistry = "us-east1-docker.pkg.dev/nexlayer/components"

// ComponentRegistry maps known component types to their default configurations.
var ComponentRegistry = map[string]ComponentConfig{
	"postgres": {
		Image:        fmt.Sprintf("%s/pern-postgres-todo:latest", NexlayerRegistry),
		ServicePorts: []int{5432},
		Environment: []EnvVar{
			{Key: "POSTGRES_USER", Value: "postgres"},
			{Key: "POSTGRES_PASSWORD", Value: "db_password"},
			{Key: "POSTGRES_DB", Value: "electric"},
		},
		Secrets: []Secret{
			{
				Name:      "my-secret",
				MountPath: "/var/secrets/my-secret-volume",
				FileName:  "tldr-56b79-firebase-adminsdk-jnzk4-a1f2fa6ef4.json",
			},
		},
		Volumes: []Volume{
			{Name: "pg-data-volume", Size: "1Gi", MountPath: "/var/lib/postgresql"},
		},
		HealthCheck: &Healthcheck{
			Command:  []string{"CMD-SHELL", "pg_isready -U postgres"},
			Interval: "5s",
			Timeout:  "5s",
			Retries:  5,
		},
	},
	"api": {
		Image:        fmt.Sprintf("%s/pern-express-todo:latest", NexlayerRegistry),
		ServicePorts: []int{3000},
		Environment: []EnvVar{
			{Key: "DATABASE_URL", Value: "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres.pod:5432/${POSTGRES_DB}"},
		},
		HealthCheck: &Healthcheck{
			Command:  []string{"CMD", "curl", "-f", "http://localhost:3000/health"},
			Interval: "5s",
			Timeout:  "5s",
			Retries:  3,
		},
	},
	"react": {
		Image:        fmt.Sprintf("%s/pern-react-todo:latest", NexlayerRegistry),
		ServicePorts: []int{80},
		Environment: []EnvVar{
			{Key: "REACT_APP_API_URL", Value: "<% URL %>/api"},
		},
	},
	"redis": {
		Image:        fmt.Sprintf("%s/redis:7", NexlayerRegistry),
		ServicePorts: []int{6379},
		Command:      []string{"redis-server", "--requirepass", "${REDIS_PASSWORD:-redis}"},
		Environment: []EnvVar{
			{Key: "REDIS_PASSWORD", Value: "redis"},
		},
		HealthCheck: &Healthcheck{
			Command:  []string{"CMD", "redis-cli", "ping"},
			Interval: "5s",
			Timeout:  "5s",
			Retries:  5,
		},
	},
}

// GetComponentConfig returns the default configuration for the given component type.
func GetComponentConfig(componentType string) (ComponentConfig, error) {
	config, exists := ComponentRegistry[componentType]
	if !exists {
		return ComponentConfig{}, fmt.Errorf("unknown component type: %s", componentType)
	}
	return config, nil
}

// DetectComponentType analyzes a pod configuration and returns the component type.
func DetectComponentType(pod interface{}) string {
	// TODO: Implement enhanced component detection based on:
	// 1. Explicit type declarations
	// 2. Image name patterns
	// 3. Environment variables
	// 4. Port configurations
	// 5. Volume mounts
	return ""
}
