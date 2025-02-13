// Package detection provides project type detection and configuration generation.
package detection

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ProjectType represents the detected type of project
type ProjectType string

const (
	// Base project types
	TypeUnknown   ProjectType = "unknown"
	TypeNextjs    ProjectType = "nextjs"
	TypeReact     ProjectType = "react"
	TypeNode      ProjectType = "node"
	TypePython    ProjectType = "python"
	TypeGo        ProjectType = "go"
	TypeDockerRaw ProjectType = "docker"

	// AI/LLM project types
	TypeLangchainNextjs ProjectType = "langchain-nextjs"
	TypeOpenAINode      ProjectType = "openai-node"
	TypeLlamaPython     ProjectType = "llama-py"

	// Full-stack project types
	TypeMERN ProjectType = "mern" // MongoDB + Express + React + Node.js
	TypePERN ProjectType = "pern" // PostgreSQL + Express + React + Node.js
	TypeMEAN ProjectType = "mean" // MongoDB + Express + Angular + Node.js
)

// ProjectInfo contains detected information about a project
type ProjectInfo struct {
	Type         ProjectType       `json:"type"`
	Name         string            `json:"name"`
	Version      string            `json:"version,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
	Scripts      map[string]string `json:"scripts,omitempty"`
	Port         int               `json:"port,omitempty"`
	HasDocker    bool              `json:"has_docker"`
	LLMProvider  string            `json:"llm_provider,omitempty"` // AI-powered IDE
	LLMModel     string            `json:"llm_model,omitempty"`    // LLM Model being used
	ImageTag     string            `json:"image_tag,omitempty"`    // The Docker image tag to use (optional)
}

// ProjectDetector defines the interface for project detection
type ProjectDetector interface {
	// Detect attempts to determine the project type and gather relevant info
	Detect(dir string) (*ProjectInfo, error)
	// Priority returns the detection priority (higher runs first)
	Priority() int
}

// DetectorRegistry holds all registered project detectors
type DetectorRegistry struct {
	detectors []ProjectDetector
}

// DetectProject attempts to detect project type using all registered detectors
func (r *DetectorRegistry) DetectProject(dir string) (*ProjectInfo, error) {
	// Sort detectors by priority
	detectors := append([]ProjectDetector{}, r.detectors...)
	for i := 0; i < len(detectors)-1; i++ {
		for j := i + 1; j < len(detectors); j++ {
			if detectors[i].Priority() < detectors[j].Priority() {
				detectors[i], detectors[j] = detectors[j], detectors[i]
			}
		}
	}

	// Try each detector in order
	for _, detector := range detectors {
		if info, err := detector.Detect(dir); err == nil && info != nil {
			return info, nil
		}
	}
	return nil, fmt.Errorf("project type could not be detected")
}

// NewDetectorRegistry creates a new registry with all available detectors
func NewDetectorRegistry() *DetectorRegistry {
	return &DetectorRegistry{
		detectors: []ProjectDetector{
			// LLM Detector (runs first)
			&LLMDetector{},

			// Full-stack Detectors
			&MERNDetector{},
			&PERNDetector{},
			&MEANDetector{},

			// Base Detectors
			&NextjsDetector{},
			&ReactDetector{},
			&NodeDetector{},
			&PythonDetector{},
			&GoDetector{},
			&DockerDetector{},
		},
	}
}

// NextjsDetector detects Next.js projects
type NextjsDetector struct{}

func (d *NextjsDetector) Priority() int { return 100 }

func (d *NextjsDetector) Detect(dir string) (*ProjectInfo, error) {
	// Read package.json
	data, err := readFileIfExists(filepath.Join(dir, "package.json"))
	if err != nil || data == nil {
		return nil, nil
	}

	name, version, deps, scripts, err := parsePackageJSON(data)
	if err != nil {
		return nil, err
	}

	// Check for Next.js dependency
	if _, hasNext := deps["next"]; !hasNext {
		return nil, nil
	}

	// Look for Dockerfile
	hasDocker := hasAnyFile(dir, []string{"Dockerfile", "dockerfile"})
	port := 3000 // Default Next.js port

	// Try to get port from Dockerfile
	if hasDocker {
		if dockerPort, err := findPortInDockerfile(filepath.Join(dir, "Dockerfile")); err == nil && dockerPort > 0 {
			port = dockerPort
		}
	}

	return &ProjectInfo{
		Type:         TypeNextjs,
		Name:         name,
		Version:      version,
		Dependencies: deps,
		Scripts:      scripts,
		Port:         port,
		HasDocker:    hasDocker,
	}, nil
}

// ReactDetector detects React projects
type ReactDetector struct{}

func (d *ReactDetector) Priority() int { return 90 }

func (d *ReactDetector) Detect(dir string) (*ProjectInfo, error) {
	// Read package.json
	data, err := readFileIfExists(filepath.Join(dir, "package.json"))
	if err != nil || data == nil {
		return nil, nil
	}

	name, version, deps, scripts, err := parsePackageJSON(data)
	if err != nil {
		return nil, err
	}

	// Check for React dependency but no Next.js
	if _, hasReact := deps["react"]; !hasReact {
		return nil, nil
	}
	if _, hasNext := deps["next"]; hasNext {
		return nil, nil // Let Next.js detector handle this
	}

	// Look for Dockerfile
	hasDocker := hasAnyFile(dir, []string{"Dockerfile", "dockerfile"})
	port := 3000 // Default React port

	// Try to get port from Dockerfile
	if hasDocker {
		if dockerPort, err := findPortInDockerfile(filepath.Join(dir, "Dockerfile")); err == nil && dockerPort > 0 {
			port = dockerPort
		}
	}

	return &ProjectInfo{
		Type:         TypeReact,
		Name:         name,
		Version:      version,
		Dependencies: deps,
		Scripts:      scripts,
		Port:         port,
		HasDocker:    hasDocker,
	}, nil
}

// NodeDetector detects Node.js projects
type NodeDetector struct{}

func (d *NodeDetector) Priority() int { return 80 }

func (d *NodeDetector) Detect(dir string) (*ProjectInfo, error) {
	// Read package.json
	data, err := readFileIfExists(filepath.Join(dir, "package.json"))
	if err != nil || data == nil {
		return nil, nil
	}

	name, version, deps, scripts, err := parsePackageJSON(data)
	if err != nil {
		return nil, err
	}

	// Look for Dockerfile
	hasDocker := hasAnyFile(dir, []string{"Dockerfile", "dockerfile"})
	port := 3000 // Default Node port

	// Try to get port from Dockerfile or env files
	if hasDocker {
		if dockerPort, err := findPortInDockerfile(filepath.Join(dir, "Dockerfile")); err == nil && dockerPort > 0 {
			port = dockerPort
		}
	}

	return &ProjectInfo{
		Type:         TypeNode,
		Name:         name,
		Version:      version,
		Dependencies: deps,
		Scripts:      scripts,
		Port:         port,
		HasDocker:    hasDocker,
	}, nil
}

// PythonDetector detects Python projects
type PythonDetector struct{}

func (d *PythonDetector) Priority() int { return 70 }

func (d *PythonDetector) Detect(dir string) (*ProjectInfo, error) {
	// Check for Python files
	pyFiles, err := findFiles(dir, []string{"*.py"})
	if err != nil || len(pyFiles) == 0 {
		return nil, nil
	}

	// Look for requirements.txt or setup.py
	hasReqs := hasAnyFile(dir, []string{"requirements.txt", "setup.py", "pyproject.toml"})
	if !hasReqs {
		return nil, nil
	}

	// Look for Dockerfile
	hasDocker := hasAnyFile(dir, []string{"Dockerfile", "dockerfile"})
	port := 8000 // Default Python web port

	// Try to get port from Dockerfile
	if hasDocker {
		if dockerPort, err := findPortInDockerfile(filepath.Join(dir, "Dockerfile")); err == nil && dockerPort > 0 {
			port = dockerPort
		}
	}

	// Try to get name from setup.py or pyproject.toml
	name := filepath.Base(dir)
	version := ""

	if setupData, err := readFileIfExists(filepath.Join(dir, "setup.py")); err == nil && setupData != nil {
		// Very basic parsing - could be improved
		for _, line := range strings.Split(string(setupData), "\n") {
			if strings.Contains(line, "name=") {
				parts := strings.Split(line, "=")
				if len(parts) > 1 {
					name = strings.Trim(strings.TrimSpace(parts[1]), "'\"")
				}
			}
		}
	}

	return &ProjectInfo{
		Type:      TypePython,
		Name:      name,
		Version:   version,
		Port:      port,
		HasDocker: hasDocker,
	}, nil
}

// GoDetector detects Go projects
type GoDetector struct{}

func (d *GoDetector) Priority() int { return 60 }

func (d *GoDetector) Detect(dir string) (*ProjectInfo, error) {
	// Check for go.mod
	data, err := readFileIfExists(filepath.Join(dir, "go.mod"))
	if err != nil || data == nil {
		return nil, nil
	}

	// Parse module name from go.mod
	name := filepath.Base(dir)
	version := ""
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			name = strings.TrimSpace(strings.TrimPrefix(line, "module "))
			break
		}
	}

	// Look for Dockerfile
	hasDocker := hasAnyFile(dir, []string{"Dockerfile", "dockerfile"})
	port := 8080 // Default Go web port

	// Try to get port from Dockerfile
	if hasDocker {
		if dockerPort, err := findPortInDockerfile(filepath.Join(dir, "Dockerfile")); err == nil && dockerPort > 0 {
			port = dockerPort
		}
	}

	return &ProjectInfo{
		Type:      TypeGo,
		Name:      name,
		Version:   version,
		Port:      port,
		HasDocker: hasDocker,
	}, nil
}

// DockerDetector detects raw Docker projects
type DockerDetector struct{}

func (d *DockerDetector) Priority() int { return 50 }

func (d *DockerDetector) Detect(dir string) (*ProjectInfo, error) {
	// Look for Dockerfile
	if !hasAnyFile(dir, []string{"Dockerfile", "dockerfile"}) {
		return nil, nil
	}

	// Try to get port from Dockerfile
	port := 80 // Default HTTP port
	if dockerPort, err := findPortInDockerfile(filepath.Join(dir, "Dockerfile")); err == nil && dockerPort > 0 {
		port = dockerPort
	}

	return &ProjectInfo{
		Type:      TypeDockerRaw,
		Name:      filepath.Base(dir),
		Port:      port,
		HasDocker: true,
	}, nil
}

// LLMDetector detects AI-powered IDEs or LLM-based coding assistants
type LLMDetector struct{}

func (d *LLMDetector) Priority() int { return 250 } // Highest priority

func (d *LLMDetector) Detect(dir string) (*ProjectInfo, error) {
	// Detect AI-powered IDE
	aiIDE := DetectAIIDE()

	// Detect LLM Model being used
	llmModel := DetectLLMModel()

	// If neither IDE nor model is found, return nil
	if aiIDE == "Unknown" && llmModel == "Unknown" {
		return nil, nil
	}

	return &ProjectInfo{
		Type:        TypeUnknown, // This detector does not detect project types
		LLMProvider: aiIDE,
		LLMModel:    llmModel,
	}, nil
}

// DetectAIIDE detects the AI-powered IDE in use
func DetectAIIDE() string {
	// 1️⃣ Check for environment variables
	if ai := os.Getenv("CURSOR_LLM"); ai != "" {
		return "Cursor"
	}
	if ai := os.Getenv("WINDSURF_LLM"); ai != "" {
		return "Windsurf"
	}
	if ai := os.Getenv("VSCODE_COPILOT_LLM"); ai != "" {
		return "VSCode"
	}
	if ai := os.Getenv("ZED_LLM"); ai != "" {
		return "Zed"
	}
	if ai := os.Getenv("AIDER_LLM"); ai != "" {
		return "Aider"
	}

	// 2️⃣ Check running processes (as fallback)
	processes := []string{"cursor", "windsurf", "vscode", "zed", "aider"}
	for _, proc := range processes {
		if _, err := os.Stat("/proc/" + proc); err == nil {
			return strings.Title(proc)
		}
	}

	// 3️⃣ Default / Fallback Option
	return "Unknown"
}

// DetectLLMModel detects the LLM Model being used
func DetectLLMModel() string {
	// Check known environment variables
	if model := os.Getenv("AI_MODEL"); model != "" {
		return model
	}
	if model := os.Getenv("CURSOR_LLM_MODEL"); model != "" {
		return model
	}
	if model := os.Getenv("VSCODE_LLM_MODEL"); model != "" {
		return model
	}
	if model := os.Getenv("WINDSURF_LLM_MODEL"); model != "" {
		return model
	}

	// List of common LLMs used in AI-powered IDEs
	llmModels := []string{
		"gpt-4o", "gpt-4o-mini", "o1-mini", "o1-preview", "o1",
		"o3-mini", "claude-3.5-sonnet", "deepseek-v3", "deepseek-r1",
		"gemini-2.0-flash",
	}

	// Check system for running LLM model processes
	for _, model := range llmModels {
		if _, err := os.Stat("/proc/" + model); err == nil {
			return model
		}
	}

	return "Unknown"
}

// Helper functions

// readFileIfExists reads a file's contents if it exists
func readFileIfExists(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}
	return os.ReadFile(path)
}

// parsePackageJSON parses a package.json file and extracts relevant info
func parsePackageJSON(data []byte) (name string, version string, deps map[string]string, scripts map[string]string, err error) {
	var pkg struct {
		Name         string            `json:"name"`
		Version      string            `json:"version"`
		Dependencies map[string]string `json:"dependencies"`
		Scripts      map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", "", nil, nil, err
	}
	return pkg.Name, pkg.Version, pkg.Dependencies, pkg.Scripts, nil
}

// findFiles looks for files matching any of the given patterns
func findFiles(dir string, patterns []string) ([]string, error) {
	var matches []string
	for _, pattern := range patterns {
		files, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return nil, err
		}
		matches = append(matches, files...)
	}
	return matches, nil
}

// hasAnyFile checks if any of the given files exist in the directory
func hasAnyFile(dir string, files []string) bool {
	for _, file := range files {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			return true
		}
	}
	return false
}

// findPortInDockerfile attempts to find the exposed port in a Dockerfile
func findPortInDockerfile(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "EXPOSE") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				return strconv.Atoi(parts[1])
			}
		}
	}
	return 0, fmt.Errorf("no port found in Dockerfile")
}

// MERNDetector detects MERN stack projects (MongoDB + Express + React + Node.js)
type MERNDetector struct{}

func (d *MERNDetector) Priority() int { return 150 }

func (d *MERNDetector) Detect(dir string) (*ProjectInfo, error) {
	// Read package.json
	data, err := readFileIfExists(filepath.Join(dir, "package.json"))
	if err != nil || data == nil {
		return nil, nil
	}

	name, version, deps, scripts, err := parsePackageJSON(data)
	if err != nil {
		return nil, err
	}

	// Check for MERN stack dependencies
	if _, hasMongoDB := deps["mongodb"]; !hasMongoDB {
		return nil, nil
	}
	if _, hasExpress := deps["express"]; !hasExpress {
		return nil, nil
	}
	if _, hasReact := deps["react"]; !hasReact {
		return nil, nil
	}

	// Look for Dockerfile
	hasDocker := hasAnyFile(dir, []string{"Dockerfile", "dockerfile"})
	port := 3000 // Default port

	return &ProjectInfo{
		Type:         TypeMERN,
		Name:         name,
		Version:      version,
		Dependencies: deps,
		Scripts:      scripts,
		Port:         port,
		HasDocker:    hasDocker,
	}, nil
}

// PERNDetector detects PERN stack projects (PostgreSQL + Express + React + Node.js)
type PERNDetector struct{}

func (d *PERNDetector) Priority() int { return 140 }

func (d *PERNDetector) Detect(dir string) (*ProjectInfo, error) {
	// Read package.json
	data, err := readFileIfExists(filepath.Join(dir, "package.json"))
	if err != nil || data == nil {
		return nil, nil
	}

	name, version, deps, scripts, err := parsePackageJSON(data)
	if err != nil {
		return nil, err
	}

	// Check for PERN stack dependencies
	if _, hasPostgres := deps["pg"]; !hasPostgres {
		return nil, nil
	}
	if _, hasExpress := deps["express"]; !hasExpress {
		return nil, nil
	}
	if _, hasReact := deps["react"]; !hasReact {
		return nil, nil
	}

	// Look for Dockerfile
	hasDocker := hasAnyFile(dir, []string{"Dockerfile", "dockerfile"})
	port := 3000 // Default port

	return &ProjectInfo{
		Type:         TypePERN,
		Name:         name,
		Version:      version,
		Dependencies: deps,
		Scripts:      scripts,
		Port:         port,
		HasDocker:    hasDocker,
	}, nil
}

// MEANDetector detects MEAN stack projects (MongoDB + Express + Angular + Node.js)
type MEANDetector struct{}

func (d *MEANDetector) Priority() int { return 130 }

func (d *MEANDetector) Detect(dir string) (*ProjectInfo, error) {
	// Read package.json
	data, err := readFileIfExists(filepath.Join(dir, "package.json"))
	if err != nil || data == nil {
		return nil, nil
	}

	name, version, deps, scripts, err := parsePackageJSON(data)
	if err != nil {
		return nil, err
	}

	// Check for MEAN stack dependencies
	if _, hasMongoDB := deps["mongodb"]; !hasMongoDB {
		return nil, nil
	}
	if _, hasExpress := deps["express"]; !hasExpress {
		return nil, nil
	}
	if _, hasAngular := deps["@angular/core"]; !hasAngular {
		return nil, nil
	}

	// Look for Dockerfile
	hasDocker := hasAnyFile(dir, []string{"Dockerfile", "dockerfile"})
	port := 4200 // Default Angular port

	return &ProjectInfo{
		Type:         TypeMEAN,
		Name:         name,
		Version:      version,
		Dependencies: deps,
		Scripts:      scripts,
		Port:         port,
		HasDocker:    hasDocker,
	}, nil
}
