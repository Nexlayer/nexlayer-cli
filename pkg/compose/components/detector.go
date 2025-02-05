// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package components

import (
	"fmt"
	"os"
)

// Package components provides structures and functions for handling Nexlayer components.

// DefaultDetector implements the ComponentDetector interface
// It is responsible for detecting component types and configuring them.
type DefaultDetector struct{}

// NewComponentDetector creates a new component detector
// Returns a pointer to DefaultDetector.
func NewComponentDetector() ComponentDetector {
	return &DefaultDetector{}
}

// DetectAndConfigure detects and configures a component based on its type.
// It sets default ports and images based on the component type.
// Returns the detected component configuration or an error.
func (d *DefaultDetector) DetectAndConfigure(pod Pod) (DetectedComponent, error) {
	// Validate input
	if pod.Name == "" {
		return DetectedComponent{}, fmt.Errorf("pod name cannot be empty")
	}

	// If type is not specified, try to detect it from directory contents
	if pod.Type == "" {
		detectedType := d.detectFromDirectory(pod.Name)
		if detectedType == "" {
			return DetectedComponent{}, fmt.Errorf("could not detect component type for %s", pod.Name)
		}
		pod.Type = detectedType
	}

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
// TODO: Implement directory-based detection logic that will:
// 1. Scan for package.json, requirements.txt, go.mod, etc.
// 2. Analyze file contents to determine component type
// 3. Return appropriate component type string
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
// It checks for an environment variable override, then selects an image
// based on predefined mappings for known component types.
func (d *DefaultDetector) getDefaultImage(componentType string) string {
	// Allow override through environment variable
	if registry := os.Getenv("NEXLAYER_REGISTRY"); registry != "" {
		return fmt.Sprintf("%s/%s:latest", registry, componentType)
	}

	// Return official Docker images based on component type
	switch componentType {
	// Databases
	case "postgres":
		return "docker.io/library/postgres:latest"
	case "redis":
		return "docker.io/library/redis:7"
	case "mongodb":
		return "docker.io/library/mongo:latest"
	case "mysql":
		return "docker.io/library/mysql:8"
	case "clickhouse":
		return "docker.io/clickhouse/clickhouse-server:latest"

	// Message Queues
	case "rabbitmq":
		return "docker.io/library/rabbitmq:3-management"
	case "kafka":
		return "docker.io/confluentinc/cp-kafka:latest"

	// Storage
	case "minio":
		return "docker.io/minio/minio:latest"
	case "elasticsearch":
		return "docker.io/elasticsearch:8"

	// Web Servers
	case "nginx":
		return "docker.io/library/nginx:latest"
	case "traefik":
		return "docker.io/library/traefik:v2.10"

	// Language Runtimes
	case "node":
		return "docker.io/library/node:18-alpine"
	case "python":
		return "docker.io/library/python:3.11-slim"
	case "golang":
		return "docker.io/library/golang:1.21-alpine"
	case "java":
		return "docker.io/library/openjdk:17-slim"

	// AI/ML Services
	case "jupyter":
		return "docker.io/jupyter/minimal-notebook:latest"

	// Frontend Frameworks (all use Node.js)
	case "react", "angular", "vue":
		return "docker.io/library/node:18-alpine"

	// Backend Frameworks
	case "express":
		return "docker.io/library/node:18-alpine"
	case "django", "fastapi":
		return "docker.io/library/python:3.11-slim"

	// Default case - use Nexlayer's registry
	default:
		return fmt.Sprintf("%s/%s:latest", DefaultRegistry, componentType)
	}

}
