// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detectors

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
)

// SvelteDetector detects Svelte and SvelteKit projects
type SvelteDetector struct{}

// Priority returns the detector's priority level
func (d *SvelteDetector) Priority() int {
	return 100
}

// Detect analyzes a directory to detect Svelte projects
func (d *SvelteDetector) Detect(dir string) (*types.ProjectInfo, error) {
	projectInfo := &types.ProjectInfo{
		Type:         types.TypeUnknown,
		Dependencies: make(map[string]string),
		Scripts:      make(map[string]string),
	}

	// Check for package.json
	packageJSONPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return projectInfo, nil
	}

	// Read and parse package.json
	packageJSONBytes, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return projectInfo, err
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(packageJSONBytes, &packageJSON); err != nil {
		return projectInfo, err
	}

	// Extract project name and version
	if name, ok := packageJSON["name"].(string); ok {
		projectInfo.Name = name
	}
	if version, ok := packageJSON["version"].(string); ok {
		projectInfo.Version = version
	}

	// Check for dependencies
	dependencies := make(map[string]interface{})
	if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		dependencies = deps
	}
	if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
		for k, v := range devDeps {
			dependencies[k] = v
		}
	}

	// Check for Svelte dependencies
	isSvelte := false
	isSvelteKit := false

	for dep, version := range dependencies {
		if dep == "svelte" {
			isSvelte = true
			// Add to dependencies
			if vStr, ok := version.(string); ok {
				projectInfo.Dependencies["svelte"] = vStr
			}
		}
		if dep == "@sveltejs/kit" {
			isSvelteKit = true
			// Add to dependencies
			if vStr, ok := version.(string); ok {
				projectInfo.Dependencies["@sveltejs/kit"] = vStr
			}
		}
	}

	// Check for svelte.config.js file (additional confirmation)
	if _, err := os.Stat(filepath.Join(dir, "svelte.config.js")); err == nil {
		isSvelte = true
	}
	if _, err := os.Stat(filepath.Join(dir, "svelte.config.mjs")); err == nil {
		isSvelte = true
	}

	// Determine project type
	if isSvelteKit {
		projectInfo.Type = types.TypeSvelteKit
	} else if isSvelte {
		projectInfo.Type = types.TypeSvelte
	} else {
		// Not a Svelte project
		return projectInfo, nil
	}

	// Extract all dependencies
	for name, version := range dependencies {
		if vStr, ok := version.(string); ok {
			projectInfo.Dependencies[name] = vStr
		}
	}

	// Check for scripts
	if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
		for name, cmd := range scripts {
			if cmdStr, ok := cmd.(string); ok {
				projectInfo.Scripts[name] = cmdStr
			}
		}
	}

	// Check for Docker configuration
	if _, err := os.Stat(filepath.Join(dir, "Dockerfile")); err == nil {
		projectInfo.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yml")); err == nil {
		projectInfo.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yaml")); err == nil {
		projectInfo.HasDocker = true
	}

	// Look for default dev port
	if startCmd, ok := projectInfo.Scripts["dev"]; ok {
		if strings.Contains(startCmd, "--port") || strings.Contains(startCmd, "-p") {
			// Try to extract port
			projectInfo.Port = 3000 // Default for SvelteKit
		} else {
			projectInfo.Port = 5173 // Vite default for Svelte
		}
	} else {
		projectInfo.Port = 5173 // Default assumption
	}

	return projectInfo, nil
}

// TODO: Integrate this detector into the detector registry in pkg/detection/detectors.go
// by adding it to the detectors slice in the NewDetectorRegistry function.
