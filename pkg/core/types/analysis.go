// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package types

// CodeFunction represents a detected function in the codebase
type CodeFunction struct {
	Name       string `json:"name"`
	Signature  string `json:"signature"`
	StartLine  int    `json:"start_line"`
	EndLine    int    `json:"end_line"`
	IsExported bool   `json:"is_exported"`
}

// APIEndpoint represents a detected API endpoint
type APIEndpoint struct {
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	Handler    string            `json:"handler"`
	Parameters map[string]string `json:"parameters"`
}

// ProjectDependency represents a project dependency
type ProjectDependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"` // direct, indirect, dev
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
