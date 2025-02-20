// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// aiDetectionTimeout is the maximum time to wait for AI-based detection
const aiDetectionTimeout = 10 * time.Second

// DetectStack attempts to detect the project stack using both static analysis and AI.
// It follows a multi-step detection process:
// 1. Static detection using the DetectorRegistry
// 2. AI-based detection if static detection is inconclusive
// 3. Fallback to unknown type if all detection methods fail
func DetectStack(dir string) (*detection.ProjectInfo, error) {
	// Try static detection first using the registry
	registry := detection.NewDetectorRegistry()
	info, err := registry.DetectProject(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to detect project: %w", err)
	}
	if info != nil && info.Type != detection.TypeUnknown {
		info.Name = sanitizeName(info.Name) // Ensure name is sanitized

		// For Node.js projects, try to detect custom port from package.json
		if info.Type == detection.TypeNode || info.Type == detection.TypeNextjs {
			if pkgData, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
				if customPort := parsePortFromPackageJSON(pkgData); customPort > 0 {
					info.Port = customPort
				}
			}
		}

		return info, nil
	}

	// Try AI-based detection with timeout if static detection is inconclusive
	ctx, cancel := context.WithTimeout(context.Background(), aiDetectionTimeout)
	defer cancel()

	projectType, err := DetectProjectTypeWithAI(ctx, dir)
	if err == nil && projectType != "" {
		// Get default port for the detected type
		port := getDefaultPort(projectType)

		// For Node.js projects, check for custom port
		if strings.Contains(strings.ToLower(projectType), "node") || strings.Contains(strings.ToLower(projectType), "next") {
			if pkgData, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
				if customPort := parsePortFromPackageJSON(pkgData); customPort > 0 {
					port = customPort
				}
			}
		}

		return &detection.ProjectInfo{
			Type: detection.ProjectType(projectType),
			Name: sanitizeName(filepath.Base(dir)),
			Port: port,
		}, nil
	}

	// Return unknown type if all detection fails
	return &detection.ProjectInfo{
		Type: detection.TypeUnknown,
		Name: sanitizeName(filepath.Base(dir)),
		Port: 8080, // Default fallback port
	}, nil
}

// getDefaultPort returns the default port for a given project type.
// This function provides sensible defaults for common project types
// while allowing for future expansion with more project types.
func getDefaultPort(projectType string) int {
	switch strings.ToLower(projectType) {
	// Node.js ecosystem
	case "node", "nodejs", "express":
		return 3000
	case "nextjs", "react", "vue":
		return 3000
	case "nuxt", "angular":
		return 4200

	// Python ecosystem
	case "python", "flask":
		return 5000
	case "django", "fastapi":
		return 8000

	// Go ecosystem
	case "go", "golang", "gin", "echo":
		return 8080
	case "fiber":
		return 3000

	// Java ecosystem
	case "java", "spring", "springboot":
		return 8080
	case "quarkus":
		return 8080
	case "micronaut":
		return 8080

	// Ruby ecosystem
	case "ruby", "rails", "sinatra":
		return 3000

	// PHP ecosystem
	case "php", "laravel", "symfony":
		return 8000

	// Rust ecosystem
	case "rust", "actix", "rocket":
		return 8000

	// Database ports (for reference)
	case "postgres", "postgresql":
		return 5432
	case "mysql", "mariadb":
		return 3306
	case "mongodb":
		return 27017
	case "redis":
		return 6379

	// Default fallback
	default:
		return 8080
	}
}

// sanitizeName ensures the name follows Nexlayer naming conventions.
// It applies the following rules:
// - Convert to lowercase
// - Replace invalid characters with hyphens
// - Ensure name starts with a letter
// - Provide default name if empty
func sanitizeName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace invalid characters with hyphens
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return '-'
	}, name)

	// Clean up multiple consecutive hyphens
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}

	// Trim hyphens from start and end
	name = strings.Trim(name, "-")

	// Ensure starts with a letter
	if name == "" || (name[0] < 'a' || name[0] > 'z') {
		name = "app-" + name
	}

	// If empty after sanitization, use default
	if name == "" {
		name = "app"
	}

	return name
}

// parsePortFromPackageJSON attempts to extract a port number from a package.json start script.
// This is useful for Node.js projects that specify a custom port.
// Returns 0 if no port is found or if there's an error parsing the file.
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
