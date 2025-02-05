// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package types

// Stack represents a project stack configuration
type Stack struct {
	Language  string          `json:"language"`
	Framework string          `json:"framework"`
	Database  string          `json:"database"`
	Resources *ResourceConfig `json:"resources,omitempty"`
	Network   *NetworkConfig  `json:"network,omitempty"`
	Security  *SecurityConfig `json:"security,omitempty"`
}

// ResourceConfig defines resource requirements
type ResourceConfig struct {
	CPU    string   `json:"cpu,omitempty"`
	Memory string   `json:"memory,omitempty"`
	GPU    []string `json:"gpu,omitempty"`
}

// NetworkConfig defines network settings
type NetworkConfig struct {
	Ingress []string `json:"ingress,omitempty"`
	Egress  []string `json:"egress,omitempty"`
}

// SecurityConfig defines security settings
type SecurityConfig struct {
	Issues []SecurityIssue `json:"issues,omitempty"`
}

// SecurityIssue represents a security concern
type SecurityIssue struct {
	Severity    string `json:"severity"`
	Description string `json:"description"`
	Context     string `json:"context,omitempty"`
}
