// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// aiDetectionTimeout is the maximum time to wait for AI-based detection
const aiDetectionTimeout = 10 * time.Second

// DetectProjectTypeWithAI uses AI to detect the project type
func DetectProjectTypeWithAI(ctx context.Context, dir string) (string, error) {
	// Get the AI provider
	provider := NewDefaultProvider()

	// Create detection prompt
	prompt := fmt.Sprintf("Analyze the project in directory %s and determine its type (e.g., Node.js, Go, etc.). Respond with 'Project type: <type>'", dir)

	// Call AI provider with timeout
	ctx, cancel := context.WithTimeout(ctx, aiDetectionTimeout)
	defer cancel()

	response, err := provider.GenerateText(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("AI detection failed: %w", err)
	}

	// Parse response
	if strings.Contains(strings.ToLower(response), "project type:") {
		parts := strings.Split(response, ":")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", fmt.Errorf("could not determine project type from AI response")
}

// getDefaultPort returns the default port for a given project type
func getDefaultPort(projectType string) int {
	switch strings.ToLower(projectType) {
	case "node", "nodejs", "express":
		return 3000
	case "nextjs", "react", "vue":
		return 3000
	case "nuxt", "angular":
		return 4200
	case "python", "flask":
		return 5000
	case "django", "fastapi":
		return 8000
	case "go", "golang", "gin", "echo":
		return 8080
	case "fiber":
		return 3000
	case "java", "spring", "springboot":
		return 8080
	case "quarkus":
		return 8080
	case "micronaut":
		return 8080
	case "ruby", "rails", "sinatra":
		return 3000
	case "php", "laravel", "symfony":
		return 8000
	case "rust", "actix", "rocket":
		return 8000
	default:
		return 8080
	}
}

// parsePortFromPackageJSON attempts to extract a port number from a package.json start script
func parsePortFromPackageJSON(data []byte) int {
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return 0
	}

	// Look for port in start script
	if startScript, ok := pkg.Scripts["start"]; ok {
		// Common patterns:
		// - "PORT=3000 node server.js"
		// - "node server.js --port 3000"
		// - "next start -p 3000"
		for _, pattern := range []string{"PORT=", "--port", "-p"} {
			if idx := strings.Index(startScript, pattern); idx != -1 {
				fields := strings.Fields(startScript[idx:])
				if len(fields) > 1 {
					if port, err := strconv.Atoi(fields[1]); err == nil {
						return port
					}
				}
			}
		}
	}

	return 0
}
