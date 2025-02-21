// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package types

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
	ImageTag     string            `json:"image_tag,omitempty"`    // Docker image tag
}

// ProjectAnalysis contains AI-generated analysis of a project
type ProjectAnalysis struct {
	Functions     map[string][]CodeFunction      `json:"functions"`       // Functions by file
	APIEndpoints  []APIEndpoint                  `json:"api_endpoints"`   // API endpoints
	Imports       map[string][]string            `json:"imports"`         // Imports by file
	Dependencies  map[string][]ProjectDependency `json:"dependencies"`    // Dependencies by type
	Description   string                         `json:"description"`     // Brief description
	Technologies  []string                       `json:"technologies"`    // Main technologies used
	Architecture  string                         `json:"architecture"`    // Architecture overview
	SecurityRisks []string                       `json:"security_risks"`  // Security concerns
	NextSteps     []string                       `json:"next_steps"`      // Next steps
	Notes         []string                       `json:"notes,omitempty"` // Additional notes
}

// CodeFunction represents a detected function in the codebase
type CodeFunction struct {
	Name       string   `json:"name"`
	Signature  string   `json:"signature"`
	StartLine  int      `json:"start_line"`
	EndLine    int      `json:"end_line"`
	IsExported bool     `json:"is_exported"`
	FilePath   string   `json:"file_path"`
	Language   string   `json:"language"`
	Type       string   `json:"type"` // e.g., "http_handler", "llm_prompt", "utility"
	Tags       []string `json:"tags,omitempty"`
}

// APIEndpoint represents a detected API endpoint
type APIEndpoint struct {
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	Handler    string            `json:"handler"`
	Parameters map[string]string `json:"parameters"`
	FilePath   string            `json:"file_path"`
	LineNumber int               `json:"line_number"`
	Tags       []string          `json:"tags,omitempty"`
}

// ProjectDependency represents a project dependency
type ProjectDependency struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Type        string   `json:"type"` // direct, indirect, dev
	IsAIRelated bool     `json:"is_ai_related"`
	Features    []string `json:"features,omitempty"`
}
