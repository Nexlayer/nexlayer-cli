// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
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
func enhancePromptWithAnalysis(basePrompt string, analysis *types.ProjectAnalysis) string {
	var sb strings.Builder
	sb.WriteString(basePrompt)
	sb.WriteString("\n\nProject Analysis:\n")

	// Add detected functions
	if len(analysis.Functions) > 0 {
		sb.WriteString("\nFunctions:\n")
		for file, functions := range analysis.Functions {
			sb.WriteString(fmt.Sprintf("\nFile: %s\n", file))
			for _, fn := range functions {
				sb.WriteString(fmt.Sprintf("- %s\n", fn.Name))
			}
		}
	}

	// Add detected API endpoints
	if len(analysis.APIEndpoints) > 0 {
		sb.WriteString("\nAPI Endpoints:\n")
		for _, ep := range analysis.APIEndpoints {
			sb.WriteString(fmt.Sprintf("- %s %s\n", ep.Method, ep.Path))
		}
	}

	// Add detected dependencies
	if len(analysis.Dependencies) > 0 {
		sb.WriteString("\nDependencies:\n")
		for pkg, deps := range analysis.Dependencies {
			sb.WriteString(fmt.Sprintf("\nPackage: %s\n", pkg))
			for _, dep := range deps {
				sb.WriteString(fmt.Sprintf("- %s@%s (%s)\n", dep.Name, dep.Version, dep.Type))
			}
		}
	}

	return sb.String()
}
