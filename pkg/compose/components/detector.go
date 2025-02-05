package components

import (
	"fmt"
	"os"
)

// DefaultDetector implements ComponentDetector interface
type DefaultDetector struct{}

// NewComponentDetector creates a new component detector
func NewComponentDetector() ComponentDetector {
	return &DefaultDetector{}
}

// DetectAndConfigure detects and configures a component based on its type
func (d *DefaultDetector) DetectAndConfigure(pod Pod) (DetectedComponent, error) {
	detected := DetectedComponent{
		Type: pod.Type,
		Config: ComponentConfig{
			Image: d.getDefaultImage(pod.Type),
		},
	}

	// Add default configuration based on component type
	switch pod.Type {
	case "frontend":
		detected.Config.Ports = []Port{{Container: 3000, Host: 80, Protocol: "tcp", Name: "web"}}
	case "backend":
		detected.Config.Ports = []Port{{Container: 8000, Host: 8000, Protocol: "tcp", Name: "api"}}
	case "database":
		detected.Config.Ports = []Port{{Container: 5432, Host: 5432, Protocol: "tcp", Name: "db"}}
	}

	return detected, nil
}

// detectFromDirectory detects component type from directory contents
func (d *DefaultDetector) detectFromDirectory(_ string) string {
	// TODO: Implement directory-based detection logic that will:
	// 1. Scan for package.json, requirements.txt, go.mod, etc.
	// 2. Analyze file contents to determine component type
	// 3. Return appropriate component type string
	return ""
}

// DefaultRegistry is the Nexlayer registry URL
const DefaultRegistry = "us-east1-docker.pkg.dev/nexlayer/components"

// getDefaultImage returns default Docker image for component type
func (d *DefaultDetector) getDefaultImage(componentType string) string {
	// Allow override through environment variable
	if registry := os.Getenv("NEXLAYER_REGISTRY"); registry != "" {
		return fmt.Sprintf("%s/%s:latest", registry, componentType)
	}

	switch componentType {
	// Databases
	case "postgres":
		return fmt.Sprintf("%s/postgres:latest", DefaultRegistry)
	case "redis":
		return fmt.Sprintf("%s/redis:7", DefaultRegistry)
	case "mongodb":
		return fmt.Sprintf("%s/mongodb:latest", DefaultRegistry)
	case "mysql":
		return fmt.Sprintf("%s/mysql:8", DefaultRegistry)
	case "clickhouse":
		return fmt.Sprintf("%s/clickhouse:latest", DefaultRegistry)

	// Message Queues
	case "rabbitmq":
		return fmt.Sprintf("%s/rabbitmq:3-management", DefaultRegistry)
	case "kafka":
		return fmt.Sprintf("%s/kafka:latest", DefaultRegistry)

	// Storage
	case "minio":
		return fmt.Sprintf("%s/minio:latest", DefaultRegistry)
	case "elasticsearch":
		return fmt.Sprintf("%s/elasticsearch:8", DefaultRegistry)

	// Web Servers
	case "nginx":
		return fmt.Sprintf("%s/nginx:latest", DefaultRegistry)
	case "traefik":
		return fmt.Sprintf("%s/traefik:v2.10", DefaultRegistry)

	// Language Runtimes
	case "node":
		return fmt.Sprintf("%s/node:18-alpine", DefaultRegistry)
	case "python":
		return fmt.Sprintf("%s/python:3.11-slim", DefaultRegistry)
	case "golang":
		return fmt.Sprintf("%s/golang:1.21-alpine", DefaultRegistry)
	case "java":
		return fmt.Sprintf("%s/java:17-slim", DefaultRegistry)

	// AI/ML Services
	case "jupyter":
		return fmt.Sprintf("%s/jupyter:latest", DefaultRegistry)

	// Frontend Frameworks
	case "react", "angular", "vue":
		return fmt.Sprintf("%s/node:18-alpine", DefaultRegistry)

	// Backend Frameworks
	case "express":
		return fmt.Sprintf("%s/node:18-alpine", DefaultRegistry)
	case "django", "fastapi":
		return fmt.Sprintf("%s/python:3.11-slim", DefaultRegistry)

	default:
		return fmt.Sprintf("%s/%s:latest", DefaultRegistry, componentType)
	}
}
