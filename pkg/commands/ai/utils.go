// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/analysis"
)

// CallGraph represents the structure of the call graph JSON output
type CallGraph struct {
	Nodes []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"nodes"`
	Edges []struct {
		Source string `json:"source"`
		Target string `json:"target"`
	} `json:"edges"`
}

// parseCallGraph parses the JSON call graph data into a CallGraph struct
func parseCallGraph(data string) (*CallGraph, error) {
	var graph CallGraph
	err := json.Unmarshal([]byte(data), &graph)
	if err != nil {
		return nil, fmt.Errorf("failed to parse call graph JSON: %v", err)
	}
	return &graph, nil
}

// enhancePromptWithAnalysis enhances the base prompt with project analysis information
func enhancePromptWithAnalysis(basePrompt string, analysis *analysis.ProjectAnalysis) string {
	var sb strings.Builder
	sb.WriteString(basePrompt)
	sb.WriteString("\n\nProject Analysis:\n")

	// Add detected frameworks
	if len(analysis.Frameworks) > 0 {
		sb.WriteString("\nFrameworks:\n")
		for _, fw := range analysis.Frameworks {
			sb.WriteString(fmt.Sprintf("- %s\n", fw))
		}
	}

	// Add detected API endpoints
	if len(analysis.APIEndpoints) > 0 {
		sb.WriteString("\nAPI Endpoints:\n")
		for _, ep := range analysis.APIEndpoints {
			sb.WriteString(fmt.Sprintf("- %s %s\n", ep.Method, ep.Path))
		}
	}

	// Add detected database types
	if len(analysis.DatabaseTypes) > 0 {
		sb.WriteString("\nDatabases:\n")
		for _, db := range analysis.DatabaseTypes {
			sb.WriteString(fmt.Sprintf("- %s\n", db))
		}
	}

	return sb.String()
}
