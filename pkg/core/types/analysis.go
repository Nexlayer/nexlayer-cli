// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package types

// ProjectAnalysis represents the analysis results of a project
type ProjectAnalysis struct {
	Functions    map[string][]Function   `json:"functions"`
	APIEndpoints []APIEndpoint           `json:"api_endpoints"`
	Imports      map[string][]string     `json:"imports"`
	Dependencies map[string][]Dependency `json:"dependencies"`
}

// Function represents a detected function in the codebase
type Function struct {
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

// Dependency represents a project dependency
type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"` // direct, indirect, dev
}
