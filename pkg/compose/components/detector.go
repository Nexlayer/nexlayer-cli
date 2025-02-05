// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package components

import (
	"fmt"
	"os"
)

// Package components provides structures and functions for handling Nexlayer components.

// DefaultDetector implements the ComponentDetector interface.
// It is responsible for detecting a component’s type and applying default configurations.
type DefaultDetector struct{}

// NewComponentDetector creates a new component detector and returns a pointer to DefaultDetector.
func NewComponentDetector() ComponentDetector {
	return &DefaultDetector{}
}

// DetectAndConfigure detects and configures a component based on its type.
// It sets default Docker images, ports, and health check configurations based on the component type.
// Returns the detected component configuration or an error.
func (d *DefaultDetector) DetectAndConfigure(pod Pod) (DetectedComponent, error) {
	// Validate input: pod name must not be empty.
	if pod.Name == "" {
		return DetectedComponent{}, fmt.Errorf("pod name cannot be empty")
	}

	// If the component type is not specified, attempt to detect it from directory contents.
	if pod.Type == "" {
		detectedType := d.detectFromDirectory(pod.Name)
		if detectedType == "" {
			return DetectedComponent{}, fmt.Errorf("could not detect component type for %s", pod.Name)
		}
		pod.Type = detectedType
	}

	// Build the default configuration based on the component type.
	detected := DetectedComponent{
		Type: pod.Type,
		Config: ComponentConfig{
			// Get the default image based on component type.
			Image:       d.getDefaultImage(pod.Type),
			HealthCheck: d.getHealthCheck(pod.Type),
		},
	}

	// Set default port configurations based on well-known component categories.
	// Frontend frameworks: react, angular, vue, nextjs, svelte, reactjs.
	// Backend frameworks: express, django, fastapi.
	// Database systems: mongodb, postgres, redis, pinecone, neo4j.
	switch pod.Type {
	case "react", "angular", "vue", "nextjs", "svelte", "reactjs", "frontend":
		detected.Config.Ports = []Port{{Container: 3000, Host: 80, Protocol: "tcp", Name: "web"}}
	case "express", "django", "fastapi", "backend":
		detected.Config.Ports = []Port{{Container: 8000, Host: 8000, Protocol: "tcp", Name: "api"}}
	case "mongodb", "postgres", "redis", "pinecone", "neo4j", "database":
		// For databases, use common default ports (modify as needed for Pinecone).
		// For example, PostgreSQL and MongoDB typically use port 5432/27017 respectively.
		// For Pinecone, we assume a default port (example: 8100) if needed.
		switch pod.Type {
		case "postgres":
			detected.Config.Ports = []Port{{Container: 5432, Host: 5432, Protocol: "tcp", Name: "db"}}
		case "mongodb":
			detected.Config.Ports = []Port{{Container: 27017, Host: 27017, Protocol: "tcp", Name: "db"}}
		case "redis":
			detected.Config.Ports = []Port{{Container: 6379, Host: 6379, Protocol: "tcp", Name: "db"}}
		case "pinecone":
			detected.Config.Ports = []Port{{Container: 8100, Host: 8100, Protocol: "tcp", Name: "db"}} // Example port for Pinecone.
		default:
			detected.Config.Ports = []Port{{Container: 5432, Host: 5432, Protocol: "tcp", Name: "db"}}
		}
	}

	return detected, nil
}

// detectFromDirectory attempts to detect the component type from the directory contents.
// TODO: Implement directory-based detection logic that scans for common files such as package.json,
// requirements.txt, go.mod, etc., and returns the appropriate component type string.
func (d *DefaultDetector) detectFromDirectory(_ string) string {
	// TODO: Analyze file names and contents to detect component type.
	return ""
}

// getHealthCheck returns the health check configuration for a given component type.
func (d *DefaultDetector) getHealthCheck(componentType string) *Healthcheck {
	switch componentType {
	case "postgres":
		return &Healthcheck{
			Command:             []string{"pg_isready", "-U", "postgres"},
			InitialDelaySeconds: 5,
			PeriodSeconds:       10,
			Interval:            "10s",
			Timeout:             "5s",
			Retries:             3,
		}
	case "redis":
		return &Healthcheck{
			Command:             []string{"redis-cli", "ping"},
			InitialDelaySeconds: 5,
			PeriodSeconds:       10,
			Interval:            "10s",
			Timeout:             "5s",
			Retries:             3,
		}
	case "mongodb":
		return &Healthcheck{
			Command:             []string{"mongo", "--eval", "db.adminCommand('ping')"},
			InitialDelaySeconds: 5,
			PeriodSeconds:       10,
			Interval:            "10s",
			Timeout:             "5s",
			Retries:             3,
		}
	case "clickhouse":
		return &Healthcheck{
			Command:             []string{"wget", "-q", "--spider", "http://localhost:8123/ping"},
			InitialDelaySeconds: 10,
			PeriodSeconds:       15,
			Interval:            "15s",
			Timeout:             "10s",
			Retries:             5,
		}
	case "minio":
		return &Healthcheck{
			Command:             []string{"curl", "-f", "http://localhost:9000/minio/health/live"},
			InitialDelaySeconds: 10,
			PeriodSeconds:       15,
			Interval:            "15s",
			Timeout:             "10s",
			Retries:             5,
		}
	}
	// Return nil if no specific health check is defined.
	return nil
}

// DefaultRegistry is the default registry URL for Nexlayer components.
const DefaultRegistry = "us-east1-docker.pkg.dev/nexlayer/components"

// getDefaultImage returns the default Docker image for a given component type.
// It first checks for an environment variable override (NEXLAYER_REGISTRY), then
// selects an image based on predefined mappings. This mapping aligns with Nexlayer Cloud’s expectations.
func (d *DefaultDetector) getDefaultImage(componentType string) string {
	// Map of component types to their default images.
	imageMap := map[string]string{
		// Databases
		"postgres":   "docker.io/library/postgres:latest",
		"redis":      "docker.io/library/redis:7",
		"mongodb":    "docker.io/library/mongo:latest",
		"mysql":      "docker.io/library/mysql:8",
		"clickhouse": "docker.io/clickhouse/clickhouse-server:latest",
		"pinecone":   "docker.io/pinecone/pinecone:latest", // Support for Pinecone vector DB

		// Message Queues
		"rabbitmq": "docker.io/library/rabbitmq:3-management",
		"kafka":    "docker.io/confluentinc/cp-kafka:latest",

		// Storage
		"minio":         "docker.io/minio/minio:latest",
		"elasticsearch": "docker.io/elasticsearch:8",

		// Web Servers
		"nginx":   "docker.io/library/nginx:latest",
		"traefik": "docker.io/library/traefik:v2.10",

		// Language Runtimes
		"node":   "docker.io/library/node:18-alpine",
		"python": "docker.io/library/python:3.11-slim",
		"golang": "docker.io/library/golang:1.21-alpine",
		"java":   "docker.io/library/openjdk:17-slim",

		// Frontend frameworks.
		"react":   "docker.io/library/node:18-alpine",
		"angular": "docker.io/library/node:18-alpine",
		"vue":     "docker.io/library/node:18-alpine",
		"nextjs":  "docker.io/library/node:18-alpine",
		"svelte":  "docker.io/library/node:18-alpine",
		"reactjs": "docker.io/library/node:18-alpine",

		// Backend frameworks.
		"express": "docker.io/library/node:18-alpine",
		"django":  "docker.io/library/python:3.11-slim",
		"fastapi": "docker.io/library/python:3.11-slim",
	}

	// Allow override via the NEXLAYER_REGISTRY environment variable.
	if registry := os.Getenv("NEXLAYER_REGISTRY"); registry != "" {
		return fmt.Sprintf("%s/%s:latest", registry, componentType)
	}

	// If the component type exists in the map, return its image.
	if image, ok := imageMap[componentType]; ok {
		return image
	}

	// Fallback default image.
	return "docker.io/library/node:18-alpine"
}
