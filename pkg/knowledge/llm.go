// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LLMContext represents the enriched context for LLM interactions
type LLMContext struct {
	ProjectStructure map[string]interface{} `json:"project_structure"`
	CodeEntities     map[string]interface{} `json:"code_entities"`
	Dependencies     map[string]string      `json:"dependencies"`
	APIEndpoints     []interface{}          `json:"api_endpoints"`
	Patterns         []interface{}          `json:"patterns"`
}

// LLMEnricher enriches the knowledge graph with LLM metadata
type LLMEnricher struct {
	graph    *Graph
	metadata map[string]interface{}
}

// NewLLMEnricher creates a new LLM metadata enricher
func NewLLMEnricher(graph *Graph) *LLMEnricher {
	return &LLMEnricher{
		graph:    graph,
		metadata: make(map[string]interface{}),
	}
}

// LoadMetadata loads LLM metadata from the tools directory
func (e *LLMEnricher) LoadMetadata(toolsDir string) error {
	metadataPath := filepath.Join(toolsDir, "llm", "metadata.json")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	if err := json.Unmarshal(data, &e.metadata); err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	return nil
}

// EnrichContext creates an enriched context for LLM interactions
func (e *LLMEnricher) EnrichContext(ctx context.Context) (*LLMContext, error) {
	enriched := &LLMContext{
		ProjectStructure: make(map[string]interface{}),
		CodeEntities:     make(map[string]interface{}),
		Dependencies:     make(map[string]string),
		APIEndpoints:     make([]interface{}, 0),
		Patterns:         make([]interface{}, 0),
	}

	// Extract project structure
	for _, node := range e.graph.Nodes {
		if node.Type == "file" {
			parts := strings.Split(node.Path, string(os.PathSeparator))
			current := enriched.ProjectStructure
			for i, part := range parts {
				if i == len(parts)-1 {
					current[part] = node.Properties
				} else {
					if _, exists := current[part]; !exists {
						current[part] = make(map[string]interface{})
					}
					current = current[part].(map[string]interface{})
				}
			}
		}
	}

	// Extract code entities
	for _, node := range e.graph.Nodes {
		switch node.Type {
		case "function":
			if enriched.CodeEntities["functions"] == nil {
				enriched.CodeEntities["functions"] = make([]interface{}, 0)
			}
			enriched.CodeEntities["functions"] = append(
				enriched.CodeEntities["functions"].([]interface{}),
				map[string]interface{}{
					"name":       node.Name,
					"path":       node.Path,
					"properties": node.Properties,
				},
			)
		case "api_endpoint":
			endpoint := map[string]interface{}{
				"path":       node.Name,
				"method":     node.Properties["method"],
				"handler":    node.Properties["handler"],
				"parameters": node.Properties["parameters"],
			}
			enriched.APIEndpoints = append(enriched.APIEndpoints, endpoint)
		}
	}

	// Add metadata-based patterns
	if patterns, ok := e.metadata["deployment_patterns"].([]interface{}); ok {
		enriched.Patterns = patterns
	}

	return enriched, nil
}

// GeneratePrompt creates an enhanced prompt for template generation
func (e *LLMEnricher) GeneratePrompt(ctx context.Context, basePrompt string) (string, error) {
	enrichedCtx, err := e.EnrichContext(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to enrich context: %w", err)
	}

	// Convert the enriched context to a natural language description
	var sb strings.Builder
	sb.WriteString(basePrompt)
	sb.WriteString("\n\nProject Analysis:\n")

	// Add API endpoints
	if len(enrichedCtx.APIEndpoints) > 0 {
		sb.WriteString("\nAPI Endpoints:\n")
		for _, endpoint := range enrichedCtx.APIEndpoints {
			ep := endpoint.(map[string]interface{})
			sb.WriteString(fmt.Sprintf("- %s %s (Handler: %s)\n",
				ep["method"], ep["path"], ep["handler"]))
		}
	}

	// Add code structure insights
	if functions, ok := enrichedCtx.CodeEntities["functions"].([]interface{}); ok {
		sb.WriteString("\nKey Functions:\n")
		for _, fn := range functions {
			f := fn.(map[string]interface{})
			if f["properties"].(map[string]interface{})["isExported"].(bool) {
				sb.WriteString(fmt.Sprintf("- %s (in %s)\n", f["name"], f["path"]))
			}
		}
	}

	// Add relevant patterns
	if len(enrichedCtx.Patterns) > 0 {
		sb.WriteString("\nRelevant Deployment Patterns:\n")
		for _, pattern := range enrichedCtx.Patterns {
			p := pattern.(map[string]interface{})
			sb.WriteString(fmt.Sprintf("- %s: %s\n", p["name"], p["description"]))
		}
	}

	return sb.String(), nil
}
