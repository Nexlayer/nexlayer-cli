// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

package components

import (
	"fmt"
	"os"
	"strings"
)

// DefaultDetector implements the ComponentDetector interface.
// It is responsible for detecting a component’s type and applying default configurations.
type DefaultDetector struct{
	analyzer *ProjectAnalyzer
}

// NewComponentDetector creates a new component detector.
func NewComponentDetector() ComponentDetector {
	return &DefaultDetector{
		analyzer: NewProjectAnalyzer(),
	}
}

// DetectAndConfigure detects and configures a component based on its type.
func (d *DefaultDetector) DetectAndConfigure(pod Pod) (DetectedComponent, error) {
	if pod.Name == "" {
		return DetectedComponent{}, fmt.Errorf("pod name cannot be empty. Please provide a valid name for your component")
	}

	// 🔥 Detect component type using image or directory structure
	detectedType := d.detectFromImage(pod.Image)
	if detectedType == "" {
		detectedType = d.detectFromDirectory(pod.Name)
	}

	// 🚀 Build the default configuration based on detected type
	image := pod.Image
	if image == "" {
		// Default images for common components
		imageMap := map[string]string{
			"postgres":      "docker.io/library/postgres:latest",
			"redis":         "docker.io/library/redis:7",
			"mongodb":       "docker.io/library/mongo:latest",
			"mysql":         "docker.io/library/mysql:8",
			"clickhouse":    "docker.io/clickhouse/clickhouse-server:latest",
			"vector-db":     "docker.io/pinecone/pinecone:latest",
			"llm-inference": "docker.io/ollama/ollama:latest",
			"rabbitmq":      "docker.io/library/rabbitmq:3-management",
			"kafka":         "docker.io/confluentinc/cp-kafka:latest",
			"minio":         "docker.io/minio/minio:latest",
			"elasticsearch": "docker.io/elasticsearch:8",
			"nginx":         "docker.io/library/nginx:latest",
			"react":         "docker.io/library/node:18-alpine",
			"express":       "docker.io/library/node:18-alpine",
			"fastapi":       "docker.io/library/python:3.11-slim",
		}

		if defaultImage, exists := imageMap[detectedType]; exists {
			image = defaultImage
		} else {
			image = "docker.io/library/node:18-alpine"
		}
	}

	detected := DetectedComponent{
		Type: detectedType,
		Config: ComponentConfig{
			Image:       image,
			HealthCheck: d.getHealthCheck(detectedType),
		},
	}

	// 🎯 Configure ports, volumes, secrets, and paths based on component type
	switch detectedType {
	case "frontend":
		detected.Config.ServicePorts = []int{3000}
		detected.Config.Path = "/"
	case "backend":
		detected.Config.ServicePorts = []int{8000}
	case "postgres":
		detected.Config.ServicePorts = []int{5432}
		detected.Config.Volumes = []Volume{{
			Name:      fmt.Sprintf("%s-data", pod.Name),
			Size:      "1Gi",
			MountPath: "/var/lib/postgresql/data",
		}}
	case "mongodb":
		detected.Config.ServicePorts = []int{27017}
		detected.Config.Volumes = []Volume{{
			Name:      fmt.Sprintf("%s-data", pod.Name),
			Size:      "1Gi",
			MountPath: "/data/db",
		}}
	case "redis":
		detected.Config.ServicePorts = []int{6379}
	}

	return detected, nil
}

// detectFromImage attempts to detect the component type from the image name.
func (d *DefaultDetector) detectFromImage(image string) string {
	if image == "" {
		return ""
	}

	// Common AI, DB, and web service images
	imageMap := map[string]string{
		"postgres":      "postgres",
		"redis":         "redis",
		"mongo":         "mongodb",
		"mysql":         "mysql",
		"clickhouse":    "clickhouse",
		"pinecone":      "vector-db",
		"ollama":        "llm-inference",
		"fastapi":       "backend",
		"express":       "backend",
		"django":        "backend",
		"react":         "frontend",
		"nextjs":        "frontend",
		"vue":           "frontend",
		"node":          "backend",
		"python":        "backend",
		"openjdk":       "java",
		"jupyter":       "ml-notebook",
		"minio":         "storage",
		"elasticsearch": "search-engine",
		"nginx":         "web-server",
		"traefik":       "load-balancer",
	}

	for key, val := range imageMap {
		if strings.Contains(strings.ToLower(image), key) {
			return val
		}
	}

	return ""
}

// detectFromDirectory attempts to detect component type based on project files.
func (d *DefaultDetector) detectFromDirectory(name string) string {
	// Check name for common patterns first
	nameLower := strings.ToLower(name)
	if strings.Contains(nameLower, "api") || strings.Contains(nameLower, "server") {
		return "backend"
	}
	if strings.Contains(nameLower, "ui") || strings.Contains(nameLower, "web") || strings.Contains(nameLower, "app") {
		return "frontend"
	}

	// Check directory contents
	files, err := os.ReadDir(".")
	if err != nil {
		return ""
	}

	for _, file := range files {
		switch file.Name() {
		case "package.json":
			return "node"
		case "requirements.txt", "Pipfile", "pyproject.toml":
			return "python"
		case "go.mod":
			return "golang"
		case "pom.xml", "build.gradle":
			return "java"
		}
	}

	return ""
}

// getHealthCheck returns a health check command based on component type.
func (d *DefaultDetector) getHealthCheck(componentType string) *Healthcheck {
	healthChecks := map[string]*Healthcheck{
		"postgres": {
			Command:  []string{"pg_isready", "-U", "postgres"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  3,
		},
		"redis": {
			Command:  []string{"redis-cli", "ping"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  3,
		},
		"mongodb": {
			Command:  []string{"mongo", "--eval", "db.runCommand({ping:1})"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  3,
		},
		"mysql": {
			Command:  []string{"mysqladmin", "ping", "-h", "localhost"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  3,
		},
		"nginx": {
			Command:  []string{"service", "nginx", "status"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  3,
		},
		"vector-db": {
			Command:  []string{"curl", "-f", "http://localhost:8000/health"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  3,
		},
		"llm-inference": {
			Command:  []string{"curl", "-f", "http://localhost:11434/health"},
			Interval: "10s",
			Timeout:  "5s",
			Retries:  3,
		},
	}

	if hc, exists := healthChecks[componentType]; exists {
		return hc
	}
	return nil
}

