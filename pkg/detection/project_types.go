// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package detection provides project type detection and configuration generation.
package detection

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
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

// ProjectAnalysis contains AI-generated analysis of a project
type ProjectAnalysis struct {
	Description     string   `json:"description"`     // Brief description of what the project does
	Technologies    []string `json:"technologies"`    // List of main technologies used
	Dependencies    []string `json:"dependencies"`    // Key dependencies and their purposes
	Architecture    string   `json:"architecture"`    // High-level architecture overview
	Recommendations []string `json:"recommendations"` // Suggested improvements or best practices
	SecurityRisks   []string `json:"security_risks"`  // Potential security concerns
	NextSteps       []string `json:"next_steps"`      // Recommended next steps for development
	Notes           []string `json:"notes,omitempty"` // Additional observations or notes
}

// String returns a string representation of the ProjectInfo
func (p *ProjectInfo) String() string {
	return fmt.Sprintf("%s project: %s", p.Type, p.Name)
}

// Detector defines the interface for project detection
type Detector interface {
	// Detect attempts to determine the project type and gather relevant info
	Detect(dir string) (*types.ProjectInfo, error)
	// Priority returns the detection priority (higher runs first)
	Priority() int
}

// DetectorOption represents an option for configuring a detector
type DetectorOption func(Detector) error

// WithTimeout sets a timeout for detection operations
func WithTimeout(timeout int) DetectorOption {
	return func(d Detector) error {
		// Implementation would depend on detector type
		return nil
	}
}

// WithRecursive enables recursive detection in subdirectories
func WithRecursive(enabled bool) DetectorOption {
	return func(d Detector) error {
		// Implementation would depend on detector type
		return nil
	}
}

// WithCache enables caching of detection results
func WithCache(enabled bool) DetectorOption {
	return func(d Detector) error {
		// Implementation would depend on detector type
		return nil
	}
}

// DetectorRegistry holds all registered project detectors
type DetectorRegistry struct {
	detectors []Detector
}

// NewDetectorRegistry creates a new registry with all available detectors
func NewDetectorRegistry() *DetectorRegistry {
	return &DetectorRegistry{
		detectors: []Detector{
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

// DetectProject attempts to detect project type using all registered detectors
func (r *DetectorRegistry) DetectProject(dir string) (*types.ProjectInfo, error) {
	// Sort detectors by priority
	detectors := append([]Detector{}, r.detectors...)
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
