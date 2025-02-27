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

// VueDetector detects Vue.js and Nuxt.js projects
type VueDetector struct{}

// Priority returns the detector's priority level
func (d *VueDetector) Priority() int {
	return 100
}

// Detect analyzes a directory to detect Vue projects
func (d *VueDetector) Detect(dir string) (*types.ProjectInfo, error) {
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

	// Check for Vue dependencies
	isVue := false
	isNuxt := false

	for dep, version := range dependencies {
		if dep == "vue" {
			isVue = true
			// Add to dependencies
			if vStr, ok := version.(string); ok {
				projectInfo.Dependencies["vue"] = vStr
			}
		}
		if dep == "nuxt" || dep == "nuxt3" || dep == "@nuxt/core" {
			isNuxt = true
			// Add to dependencies
			if vStr, ok := version.(string); ok {
				projectInfo.Dependencies[dep] = vStr
			}
		}
	}

	// Check for config files (additional confirmation)
	if _, err := os.Stat(filepath.Join(dir, "vue.config.js")); err == nil {
		isVue = true
	}
	if _, err := os.Stat(filepath.Join(dir, "nuxt.config.js")); err == nil {
		isNuxt = true
	}
	if _, err := os.Stat(filepath.Join(dir, "nuxt.config.ts")); err == nil {
		isNuxt = true
	}

	// Determine project type
	if isNuxt {
		projectInfo.Type = types.TypeNuxt
	} else if isVue {
		projectInfo.Type = types.TypeVue
	} else {
		// Not a Vue project
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
			projectInfo.Port = 3000 // Check for custom port
		} else {
			projectInfo.Port = 5173 // Vite default for Vue 3
		}
	} else if startCmd, ok := projectInfo.Scripts["serve"]; ok {
		if strings.Contains(startCmd, "--port") || strings.Contains(startCmd, "-p") {
			// Try to extract port
			projectInfo.Port = 8080 // Default Vue CLI port
		} else {
			projectInfo.Port = 8080 // Default Vue CLI port
		}
	} else {
		// Default ports
		if isNuxt {
			projectInfo.Port = 3000 // Default for Nuxt
		} else {
			projectInfo.Port = 8080 // Default for Vue
		}
	}

	return projectInfo, nil
}

// TODO: Integrate this detector into the detector registry in pkg/detection/detectors.go
// by adding it to the detectors slice in the NewDetectorRegistry function.
