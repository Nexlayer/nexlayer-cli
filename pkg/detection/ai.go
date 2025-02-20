// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detection

// ProjectAnalysis represents the AI analysis of a project
type ProjectAnalysis struct {
	Functions    []Function    `json:"functions"`
	Endpoints    []APIEndpoint `json:"endpoints"`
	Dependencies []Dependency  `json:"dependencies"`
	LLMProvider  string        `json:"llm_provider,omitempty"`
	LLMModel     string        `json:"llm_model,omitempty"`
}

// Function represents a detected function in the codebase
type Function struct {
	Name       string   `json:"name"`
	FilePath   string   `json:"file_path"`
	LineNumber int      `json:"line_number"`
	Language   string   `json:"language"`
	Type       string   `json:"type"` // e.g., "http_handler", "llm_prompt", "utility"
	Tags       []string `json:"tags,omitempty"`
}

// APIEndpoint represents a detected API endpoint
type APIEndpoint struct {
	Path       string   `json:"path"`
	Method     string   `json:"method"`
	Handler    string   `json:"handler"`
	FilePath   string   `json:"file_path"`
	LineNumber int      `json:"line_number"`
	Tags       []string `json:"tags,omitempty"`
}

// Dependency represents a project dependency with AI-specific metadata
type Dependency struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Type        string   `json:"type"` // e.g., "llm", "framework", "utility"
	IsAIRelated bool     `json:"is_ai_related"`
	Features    []string `json:"features,omitempty"`
}
