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

// BunDetector detects Bun and Hono projects
type BunDetector struct{}

// Priority returns the detector's priority level
func (d *BunDetector) Priority() int {
	return 95
}

// Detect analyzes a directory to detect Bun projects
func (d *BunDetector) Detect(dir string) (*types.ProjectInfo, error) {
	projectInfo := &types.ProjectInfo{
		Type:         types.TypeUnknown,
		Dependencies: make(map[string]string),
		Scripts:      make(map[string]string),
	}

	// Check for typical Bun project files
	hasBunLock := false
	if _, err := os.Stat(filepath.Join(dir, "bun.lockb")); err == nil {
		hasBunLock = true
	}

	// Check for package.json
	packageJSONPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		// If no package.json but bun lock exists, it's likely a simple Bun project
		if hasBunLock {
			projectInfo.Type = types.TypeBun
			projectInfo.Port = 3000 // Default Bun port
		}
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

	// Check for Bun configuration
	hasBun := hasBunLock
	hasHono := false
	hasNextjs := false
	hasAstro := false
	hasExpo := false

	// Look for Bun scripts
	if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
		for name, cmd := range scripts {
			if cmdStr, ok := cmd.(string); ok {
				projectInfo.Scripts[name] = cmdStr

				// Check if using bun commands
				if strings.Contains(cmdStr, "bun run") || strings.Contains(cmdStr, "bunx") {
					hasBun = true
				}
			}
		}
	}

	// Extract all dependencies
	for name, version := range dependencies {
		if vStr, ok := version.(string); ok {
			projectInfo.Dependencies[name] = vStr

			// Check for key dependencies
			if name == "bun" || name == "bun-types" {
				hasBun = true
			} else if name == "hono" {
				hasHono = true
			} else if name == "next" {
				hasNextjs = true
			} else if name == "astro" {
				hasAstro = true
			} else if name == "expo" || name == "expo-cli" {
				hasExpo = true
			}
		}
	}

	// Check for Bun specific config
	if bunConfig, ok := packageJSON["bun"]; ok && bunConfig != nil {
		hasBun = true
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

	// Determine the project type based on detected features
	if hasBun {
		if hasHono {
			if hasNextjs && hasExpo {
				projectInfo.Type = types.TypeBunHonoNextjsExpo
				projectInfo.Port = 3000
			} else if hasAstro && hasExpo {
				projectInfo.Type = types.TypeBunHonoAstroExpo
				projectInfo.Port = 4321 // Astro default port
			} else if hasNextjs {
				projectInfo.Type = types.TypeBunHonoNextjs
				projectInfo.Port = 3000
			} else if hasAstro {
				projectInfo.Type = types.TypeBunHonoAstro
				projectInfo.Port = 4321
			} else if hasExpo {
				projectInfo.Type = types.TypeBunHonoExpo
				projectInfo.Port = 19000 // Expo default port
			} else {
				projectInfo.Type = types.TypeBunHono
				projectInfo.Port = 3000 // Default for Hono
			}
		} else {
			projectInfo.Type = types.TypeBun
			projectInfo.Port = 3000 // Default Bun port
		}
	}

	return projectInfo, nil
}

// TODO: Integrate this detector into the detector registry in pkg/detection/detectors.go
// by adding it to the detectors slice in the NewDetectorRegistry function.
