// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detection

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
)

// NextjsDetector detects Next.js projects
type NextjsDetector struct{}

func (d *NextjsDetector) Priority() int { return 100 }

func (d *NextjsDetector) Detect(dir string) (*types.ProjectInfo, error) {
	// Check for next.config.js/ts
	nextConfigPath := filepath.Join(dir, "next.config.js")
	if _, err := os.Stat(nextConfigPath); err != nil {
		nextConfigPath = filepath.Join(dir, "next.config.ts")
		if _, err := os.Stat(nextConfigPath); err != nil {
			return nil, nil
		}
	}

	// Check package.json for next.js dependency
	pkgJSON, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil, nil
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(pkgJSON, &pkg); err != nil {
		return nil, nil
	}

	// Check if next.js is a dependency
	if _, hasNext := pkg.Dependencies["next"]; !hasNext {
		if _, hasNext = pkg.DevDependencies["next"]; !hasNext {
			return nil, nil
		}
	}

	// Check for AI/LLM dependencies
	hasLangchain := false
	for dep := range pkg.Dependencies {
		if strings.Contains(dep, "langchain") {
			hasLangchain = true
			break
		}
	}
	if !hasLangchain {
		for dep := range pkg.DevDependencies {
			if strings.Contains(dep, "langchain") {
				hasLangchain = true
				break
			}
		}
	}

	projectType := types.TypeNextjs
	if hasLangchain {
		projectType = types.TypeLangchainNextjs
	}

	return &types.ProjectInfo{
		Type:    projectType,
		Port:    3000, // Default Next.js port
		Name:    filepath.Base(dir),
		Version: pkg.Dependencies["next"],
	}, nil
}

// ReactDetector detects React projects
type ReactDetector struct{}

func (d *ReactDetector) Priority() int { return 90 }

func (d *ReactDetector) Detect(dir string) (*types.ProjectInfo, error) {
	pkgJSON, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil, nil
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		Scripts         map[string]string `json:"scripts"`
	}

	if err := json.Unmarshal(pkgJSON, &pkg); err != nil {
		return nil, nil
	}

	// Check if react is a dependency
	_, hasReact := pkg.Dependencies["react"]
	if !hasReact {
		_, hasReact = pkg.DevDependencies["react"]
		if !hasReact {
			return nil, nil
		}
	}

	// Check if it's not a Next.js project
	if _, hasNext := pkg.Dependencies["next"]; hasNext {
		return nil, nil
	}
	if _, hasNext := pkg.DevDependencies["next"]; hasNext {
		return nil, nil
	}

	// Determine port from scripts
	port := 3000 // Default port
	for _, script := range pkg.Scripts {
		if strings.Contains(script, "--port") {
			parts := strings.Split(script, "--port")
			if len(parts) > 1 {
				portStr := strings.TrimSpace(strings.Split(parts[1], " ")[0])
				if portStr != "" {
					port = parsePort(portStr)
				}
			}
		}
	}

	return &types.ProjectInfo{
		Type:    types.TypeReact,
		Port:    port,
		Name:    filepath.Base(dir),
		Version: pkg.Dependencies["react"],
	}, nil
}

// NodeDetector detects Node.js projects
type NodeDetector struct{}

func (d *NodeDetector) Priority() int { return 80 }

func (d *NodeDetector) Detect(dir string) (*types.ProjectInfo, error) {
	pkgJSON, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil, nil
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		Scripts         map[string]string `json:"scripts"`
		Name            string            `json:"name"`
		Version         string            `json:"version"`
	}

	if err := json.Unmarshal(pkgJSON, &pkg); err != nil {
		return nil, nil
	}

	// Check if it's not a React or Next.js project
	if _, hasReact := pkg.Dependencies["react"]; hasReact {
		return nil, nil
	}
	if _, hasNext := pkg.Dependencies["next"]; hasNext {
		return nil, nil
	}

	// Check for OpenAI dependencies
	hasOpenAI := false
	for dep := range pkg.Dependencies {
		if strings.Contains(dep, "openai") {
			hasOpenAI = true
			break
		}
	}
	if !hasOpenAI {
		for dep := range pkg.DevDependencies {
			if strings.Contains(dep, "openai") {
				hasOpenAI = true
				break
			}
		}
	}

	projectType := types.TypeNode
	if hasOpenAI {
		projectType = types.TypeOpenAINode
	}

	// Determine port from scripts or environment
	port := 3000 // Default port
	for _, script := range pkg.Scripts {
		if strings.Contains(script, "--port") || strings.Contains(script, "-p") {
			parts := strings.Split(script, "--port")
			if len(parts) == 1 {
				parts = strings.Split(script, "-p")
			}
			if len(parts) > 1 {
				portStr := strings.TrimSpace(strings.Split(parts[1], " ")[0])
				if portStr != "" {
					port = parsePort(portStr)
				}
			}
		}
	}

	name := pkg.Name
	if name == "" {
		name = filepath.Base(dir)
	}

	return &types.ProjectInfo{
		Type:    projectType,
		Port:    port,
		Name:    name,
		Version: pkg.Version,
	}, nil
}

// PythonDetector detects Python projects
type PythonDetector struct{}

func (d *PythonDetector) Priority() int { return 70 }

func (d *PythonDetector) Detect(dir string) (*types.ProjectInfo, error) {
	// Check for requirements.txt or setup.py
	reqPath := filepath.Join(dir, "requirements.txt")
	setupPath := filepath.Join(dir, "setup.py")
	if _, err := os.Stat(reqPath); err != nil {
		if _, err := os.Stat(setupPath); err != nil {
			return nil, nil
		}
	}

	// Check for main.py or app.py
	mainPath := filepath.Join(dir, "main.py")
	appPath := filepath.Join(dir, "app.py")
	if _, err := os.Stat(mainPath); err != nil {
		if _, err := os.Stat(appPath); err != nil {
			return nil, nil
		}
	}

	// Try to determine port from common Python web frameworks
	port := 8000 // Default Python port
	files, err := filepath.Glob(filepath.Join(dir, "*.py"))
	if err == nil {
		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			contentStr := string(content)
			// Check for Flask port
			if strings.Contains(contentStr, "app.run") && strings.Contains(contentStr, "port") {
				port = parsePort(contentStr)
			}
			// Check for FastAPI port
			if strings.Contains(contentStr, "uvicorn.run") && strings.Contains(contentStr, "port") {
				port = parsePort(contentStr)
			}
		}
	}

	return &types.ProjectInfo{
		Type:    types.TypePython,
		Port:    port,
		Name:    filepath.Base(dir),
		Version: "", // Version could be extracted from requirements.txt if needed
	}, nil
}

// GoDetector detects Go projects
type GoDetector struct{}

func (d *GoDetector) Priority() int { return 60 }

func (d *GoDetector) Detect(dir string) (*types.ProjectInfo, error) {
	// Check for go.mod
	goModPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		return nil, nil
	}

	// Read go.mod to get module name and Go version
	modContent, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, nil
	}

	modLines := strings.Split(string(modContent), "\n")
	var moduleName, goVersion string
	for _, line := range modLines {
		if strings.HasPrefix(line, "module ") {
			moduleName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
		if strings.HasPrefix(line, "go ") {
			goVersion = strings.TrimSpace(strings.TrimPrefix(line, "go "))
		}
	}

	// Try to determine port from main.go or server.go
	port := 8080 // Default Go port
	files := []string{
		filepath.Join(dir, "main.go"),
		filepath.Join(dir, "server.go"),
		filepath.Join(dir, "cmd", "main.go"),
		filepath.Join(dir, "cmd", "server.go"),
	}

	for _, file := range files {
		if content, err := os.ReadFile(file); err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "ListenAndServe") {
				port = parsePort(contentStr)
				break
			}
		}
	}

	name := moduleName
	if name == "" {
		name = filepath.Base(dir)
	}

	return &types.ProjectInfo{
		Type:    types.TypeGo,
		Port:    port,
		Name:    name,
		Version: goVersion,
	}, nil
}

// DockerDetector detects Docker projects
type DockerDetector struct{}

func (d *DockerDetector) Priority() int { return 50 }

func (d *DockerDetector) Detect(dir string) (*types.ProjectInfo, error) {
	// Check for Dockerfile or docker-compose.yml
	dockerfilePath := filepath.Join(dir, "Dockerfile")
	composePath := filepath.Join(dir, "docker-compose.yml")
	if _, err := os.Stat(dockerfilePath); err != nil {
		if _, err := os.Stat(composePath); err != nil {
			return nil, nil
		}
	}

	// Try to determine port from Dockerfile or docker-compose.yml
	port := 80 // Default Docker port
	if content, err := os.ReadFile(dockerfilePath); err == nil {
		contentStr := string(content)
		if strings.Contains(contentStr, "EXPOSE") {
			port = parsePort(contentStr)
		}
	}

	return &types.ProjectInfo{
		Type:    types.TypeDockerRaw,
		Port:    port,
		Name:    filepath.Base(dir),
		Version: "", // Version could be extracted from Dockerfile if needed
	}, nil
}

// MERNDetector detects MERN stack projects (MongoDB + Express + React + Node.js)
type MERNDetector struct{}

func (d *MERNDetector) Priority() int { return 150 }

func (d *MERNDetector) Detect(dir string) (*types.ProjectInfo, error) {
	pkgJSON, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil, nil
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		Name            string            `json:"name"`
		Version         string            `json:"version"`
	}

	if err := json.Unmarshal(pkgJSON, &pkg); err != nil {
		return nil, nil
	}

	// Check for MERN stack dependencies
	hasMongoDB := false
	hasExpress := false
	hasReact := false

	for dep := range pkg.Dependencies {
		if strings.Contains(dep, "mongodb") || strings.Contains(dep, "mongoose") {
			hasMongoDB = true
		}
		if dep == "express" {
			hasExpress = true
		}
		if dep == "react" {
			hasReact = true
		}
	}

	if !hasMongoDB || !hasExpress || !hasReact {
		return nil, nil
	}

	name := pkg.Name
	if name == "" {
		name = filepath.Base(dir)
	}

	return &types.ProjectInfo{
		Type:    types.TypeMERN,
		Port:    3000, // Default MERN stack port
		Name:    name,
		Version: pkg.Version,
	}, nil
}

// PERNDetector detects PERN stack projects (PostgreSQL + Express + React + Node.js)
type PERNDetector struct{}

func (d *PERNDetector) Priority() int { return 140 }

func (d *PERNDetector) Detect(dir string) (*types.ProjectInfo, error) {
	pkgJSON, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil, nil
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		Name            string            `json:"name"`
		Version         string            `json:"version"`
	}

	if err := json.Unmarshal(pkgJSON, &pkg); err != nil {
		return nil, nil
	}

	// Check for PERN stack dependencies
	hasPostgres := false
	hasExpress := false
	hasReact := false

	for dep := range pkg.Dependencies {
		if strings.Contains(dep, "pg") || strings.Contains(dep, "postgres") {
			hasPostgres = true
		}
		if dep == "express" {
			hasExpress = true
		}
		if dep == "react" {
			hasReact = true
		}
	}

	if !hasPostgres || !hasExpress || !hasReact {
		return nil, nil
	}

	name := pkg.Name
	if name == "" {
		name = filepath.Base(dir)
	}

	return &types.ProjectInfo{
		Type:    types.TypePERN,
		Port:    3000, // Default PERN stack port
		Name:    name,
		Version: pkg.Version,
	}, nil
}

// MEANDetector detects MEAN stack projects (MongoDB + Express + Angular + Node.js)
type MEANDetector struct{}

func (d *MEANDetector) Priority() int { return 130 }

func (d *MEANDetector) Detect(dir string) (*types.ProjectInfo, error) {
	pkgJSON, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil, nil
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		Name            string            `json:"name"`
		Version         string            `json:"version"`
	}

	if err := json.Unmarshal(pkgJSON, &pkg); err != nil {
		return nil, nil
	}

	// Check for MEAN stack dependencies
	hasMongoDB := false
	hasExpress := false
	hasAngular := false

	for dep := range pkg.Dependencies {
		if strings.Contains(dep, "mongodb") || strings.Contains(dep, "mongoose") {
			hasMongoDB = true
		}
		if dep == "express" {
			hasExpress = true
		}
		if strings.Contains(dep, "@angular/core") {
			hasAngular = true
		}
	}

	if !hasMongoDB || !hasExpress || !hasAngular {
		return nil, nil
	}

	name := pkg.Name
	if name == "" {
		name = filepath.Base(dir)
	}

	return &types.ProjectInfo{
		Type:    types.TypeMEAN,
		Port:    4200, // Default Angular port
		Name:    name,
		Version: pkg.Version,
	}, nil
}

// Helper function to parse port numbers from strings
func parsePort(s string) int {
	// Default port if parsing fails
	defaultPort := 3000

	// Remove common prefixes and suffixes
	s = strings.TrimPrefix(s, ":")
	s = strings.TrimPrefix(s, "=")
	s = strings.TrimSpace(s)

	// Split by common delimiters and take first number-like part
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !('0' <= r && r <= '9')
	})

	if len(parts) == 0 {
		return defaultPort
	}

	// Parse the first number we find
	var port int
	for _, c := range parts[0] {
		if '0' <= c && c <= '9' {
			port = port*10 + int(c-'0')
		}
	}

	if port == 0 {
		return defaultPort
	}
	return port
}
