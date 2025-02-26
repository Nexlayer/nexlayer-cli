// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package detect

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
	coretypes "github.com/Nexlayer/nexlayer-cli/pkg/core/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// AnalyzeDirectory analyzes a directory to detect application type and configuration
// Deprecated: Use pkg/detection.DetectorRegistry.DetectProject instead
func AnalyzeDirectory(dir string) (*types.AppConfig, error) {
	// Use the newer detection mechanism internally
	registry := detection.NewDetectorRegistry()
	projectInfo, err := registry.DetectProject(dir)
	if err != nil {
		// Fall back to the original implementation if the new detection fails
		return analyzeDirectoryLegacy(dir)
	}

	// Convert ProjectInfo to AppConfig
	return convertToAppConfig(projectInfo, dir)
}

// analyzeDirectoryLegacy contains the original implementation for backward compatibility
func analyzeDirectoryLegacy(dir string) (*types.AppConfig, error) {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize default configuration
	config := &types.AppConfig{
		Name: filepath.Base(dir),
		Type: "generic",
		Container: &types.Container{
			Ports: []int{8080},
		},
		Resources: &types.Resources{
			CPU:    "100m",
			Memory: "128Mi",
		},
	}

	// Check for Dockerfile
	if _, err := os.Stat(filepath.Join(dir, "Dockerfile")); err == nil {
		config.Container.UseDockerfile = true
	}

	// Check for package.json (Node.js)
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		config.Type = "nodejs"
		config.Container.Command = "npm start"
	}

	// Check for requirements.txt (Python)
	if _, err := os.Stat(filepath.Join(dir, "requirements.txt")); err == nil {
		config.Type = "python"
		config.Container.Command = "python app.py"
	}

	// Check for go.mod (Go)
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		config.Type = "go"
		config.Container.Command = "./app"
	}

	// Check for .env file
	if envFile, err := os.Open(filepath.Join(dir, ".env")); err == nil {
		defer envFile.Close()
		if content, err := ioutil.ReadAll(envFile); err == nil {
			config.Env = []string{string(content)}
		}
	}

	return config, nil
}

// convertToAppConfig converts a ProjectInfo to an AppConfig
func convertToAppConfig(info *coretypes.ProjectInfo, dir string) (*types.AppConfig, error) {
	config := &types.AppConfig{
		Name: info.Name,
		Type: string(info.Type),
		Container: &types.Container{
			Ports: []int{info.Port},
		},
		Resources: &types.Resources{
			CPU:    "100m",
			Memory: "128Mi",
		},
	}

	// Set HasExistingImage based on Docker detection
	config.HasExistingImage = info.HasDocker

	// Set container command based on project type
	switch info.Type {
	case "node":
		config.Type = "nodejs"
		config.Container.Command = "npm start"
	case "python":
		config.Container.Command = "python app.py"
	case "go":
		config.Container.Command = "./app"
	}

	// Check for Dockerfile
	if info.HasDocker {
		config.Container.UseDockerfile = true
	}

	// Read .env file for environment variables
	if envFile, err := os.Open(filepath.Join(dir, ".env")); err == nil {
		defer envFile.Close()
		if content, err := ioutil.ReadAll(envFile); err == nil {
			config.Env = []string{string(content)}
		}
	}

	return config, nil
}
