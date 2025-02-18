// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package ai

import (
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/analysis"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	"github.com/Nexlayer/nexlayer-cli/pkg/knowledge"
	tmpl "github.com/Nexlayer/nexlayer-cli/pkg/template"
)

// createTemplate creates a template based on project analysis
func createTemplate(info *detection.ProjectInfo, analysis *analysis.ProjectAnalysis, graph *knowledge.Graph, req TemplateRequest) tmpl.NexlayerYAML {
	var yamlTemplate tmpl.NexlayerYAML

	// Map detected stack to pod type
	podType := detectPodType(info, analysis)

	// Create template with detected pod type
	if podType != "" {
		pod := createPod(podType, analysis, info, req.ProjectName)
		yamlTemplate = tmpl.NexlayerYAML{
			Application: tmpl.Application{
				Name: req.ProjectName,
				Pods: []tmpl.Pod{pod},
			},
		}

		// Use knowledge graph for additional insights if available
		if graph != nil {
			enrichTemplateWithGraph(&yamlTemplate, graph)
		}
	}

	return yamlTemplate
}

// detectPodType determines the pod type based on project info and analysis
func detectPodType(info *detection.ProjectInfo, analysis *analysis.ProjectAnalysis) tmpl.PodType {
	switch info.Type {
	case detection.TypeReact:
		return tmpl.React
	case detection.TypePython:
		if analysis != nil && len(analysis.Frameworks) > 0 {
			if containsAny(analysis.Frameworks, "django") {
				return tmpl.Django
			} else if containsAny(analysis.Frameworks, "fastapi") {
				return tmpl.FastAPI
			}
			return tmpl.Backend
		}
		return tmpl.Backend
	default:
		return tmpl.Backend
	}
}

// createPod creates a pod configuration based on detected settings
func createPod(podType tmpl.PodType, analysis *analysis.ProjectAnalysis, info *detection.ProjectInfo, projectName string) tmpl.Pod {
	// Get default ports for the pod type
	ports := detectPorts(analysis, info)

	// Get default environment variables
	vars := tmpl.DefaultEnvVars[podType]

	// Add database environment variables if needed
	if analysis != nil && len(analysis.DatabaseTypes) > 0 {
		vars = append(vars, generateDatabaseEnvVars(analysis.DatabaseTypes)...)
	}

	return tmpl.Pod{
		Name:         "app",
		Type:         podType,
		Image:        fmt.Sprintf("ghcr.io/nexlayer/%s:latest", projectName),
		ServicePorts: ports,
		Vars:         vars,
	}
}

// detectPorts detects ports from analysis or project info
func detectPorts(analysis *analysis.ProjectAnalysis, info *detection.ProjectInfo) []int {
	var ports []int

	if len(analysis.APIEndpoints) > 0 {
		for _, endpoint := range analysis.APIEndpoints {
			if port := extractPortFromEndpoint(endpoint.Path); port > 0 {
				ports = append(ports, port)
			}
		}
	}

	// Fallback to info.Port if no ports detected
	if len(ports) == 0 && info.Port > 0 {
		ports = []int{info.Port}
	}

	return ports
}

// enrichTemplateWithGraph adds graph-based insights to the template
func enrichTemplateWithGraph(yamlTemplate *tmpl.NexlayerYAML, graph *knowledge.Graph) {
	for _, node := range graph.Nodes {
		if node.Type == "api_endpoint" {
			yamlTemplate.Application.Pods[0].Annotations = node.Annotations
			break
		}
	}
}

// generateDatabaseEnvVars generates environment variables for database configuration
func generateDatabaseEnvVars(dbTypes []string) []tmpl.EnvVar {
	var vars []tmpl.EnvVar
	for _, dbType := range dbTypes {
		switch strings.ToLower(dbType) {
		case "postgres", "postgresql":
			vars = append(vars, []tmpl.EnvVar{
				{Key: "DB_HOST", Value: "localhost"},
				{Key: "DB_PORT", Value: "5432"},
				{Key: "DB_NAME", Value: "app"},
				{Key: "DB_USER", Value: "postgres"},
				{Key: "DB_PASSWORD", Value: "<% DB_PASSWORD %>"},
			}...)
		case "mysql", "mariadb":
			vars = append(vars, []tmpl.EnvVar{
				{Key: "DB_HOST", Value: "localhost"},
				{Key: "DB_PORT", Value: "3306"},
				{Key: "DB_NAME", Value: "app"},
				{Key: "DB_USER", Value: "root"},
				{Key: "DB_PASSWORD", Value: "<% DB_PASSWORD %>"},
			}...)
		case "mongodb":
			vars = append(vars, []tmpl.EnvVar{
				{Key: "MONGODB_URI", Value: "mongodb://localhost:27017/app"},
			}...)
		}
	}
	return vars
}

// extractPortFromEndpoint extracts port number from endpoint path
func extractPortFromEndpoint(path string) int {
	// Simple port extraction from URL-like paths
	// e.g., "http://localhost:3000" -> 3000
	parts := strings.Split(path, ":")
	if len(parts) > 1 {
		lastPart := parts[len(parts)-1]
		// Remove any trailing path
		if idx := strings.Index(lastPart, "/"); idx != -1 {
			lastPart = lastPart[:idx]
		}
		var port int
		if _, err := fmt.Sscanf(lastPart, "%d", &port); err == nil {
			return port
		}
	}
	return 0
}

// containsAny checks if any of the items are in the slice
func containsAny(slice []string, items ...string) bool {
	for _, item := range items {
		for _, s := range slice {
			if strings.Contains(strings.ToLower(s), strings.ToLower(item)) {
				return true
			}
		}
	}
	return false
}
