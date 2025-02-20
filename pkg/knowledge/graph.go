// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
)

// Node represents a code entity in the knowledge graph
type Node struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // function, class, interface, etc.
	Name        string                 `json:"name"`
	Path        string                 `json:"path"`
	Properties  map[string]interface{} `json:"properties"`
	Annotations map[string]string      `json:"annotations"`
}

// Edge represents a relationship between nodes
type Edge struct {
	Source      string                 `json:"source"`
	Target      string                 `json:"target"`
	Type        string                 `json:"type"` // calls, imports, implements, etc.
	Properties  map[string]interface{} `json:"properties"`
	Annotations map[string]string      `json:"annotations"`
}

// Graph represents the project's knowledge graph
type Graph struct {
	mu    sync.RWMutex
	Nodes map[string]*Node `json:"nodes"`
	Edges []*Edge          `json:"edges"`
}

// NewGraph creates a new knowledge graph
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make([]*Edge, 0),
	}
}

// BuildFromAnalysis constructs the graph from project analysis
func (g *Graph) BuildFromAnalysis(ctx context.Context, projectAnalysis *types.ProjectAnalysis) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Add nodes for functions
	for file, functions := range projectAnalysis.Functions {
		for _, fn := range functions {
			nodeID := fmt.Sprintf("func:%s:%s", file, fn.Name)
			g.Nodes[nodeID] = &Node{
				ID:   nodeID,
				Type: "function",
				Name: fn.Name,
				Path: file,
				Properties: map[string]interface{}{
					"signature":  fn.Signature,
					"startLine":  fn.StartLine,
					"endLine":    fn.EndLine,
					"isExported": fn.IsExported,
				},
				Annotations: make(map[string]string),
			}
		}
	}

	// Add nodes for API endpoints
	for _, endpoint := range projectAnalysis.APIEndpoints {
		nodeID := fmt.Sprintf("api:%s:%s", endpoint.Method, endpoint.Path)
		g.Nodes[nodeID] = &Node{
			ID:   nodeID,
			Type: "api_endpoint",
			Name: endpoint.Path,
			Properties: map[string]interface{}{
				"method":     endpoint.Method,
				"handler":    endpoint.Handler,
				"parameters": endpoint.Parameters,
			},
			Annotations: make(map[string]string),
		}
	}

	// Add edges for imports
	for file, imports := range projectAnalysis.Imports {
		for _, imp := range imports {
			sourceID := fmt.Sprintf("file:%s", file)
			targetID := fmt.Sprintf("package:%s", imp)

			// Add file node if it doesn't exist
			if _, exists := g.Nodes[sourceID]; !exists {
				g.Nodes[sourceID] = &Node{
					ID:   sourceID,
					Type: "file",
					Name: file,
					Properties: map[string]interface{}{
						"type": "source_file",
					},
					Annotations: make(map[string]string),
				}
			}

			g.Edges = append(g.Edges, &Edge{
				Source: sourceID,
				Target: targetID,
				Type:   "imports",
				Properties: map[string]interface{}{
					"importPath": imp,
				},
				Annotations: make(map[string]string),
			})
		}
	}

	return nil
}

// AddCallGraphData integrates call graph data into the knowledge graph
func (g *Graph) AddCallGraphData(callGraphData []byte) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	var callGraph struct {
		Nodes []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"nodes"`
		Edges []struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"edges"`
	}

	if err := json.Unmarshal(callGraphData, &callGraph); err != nil {
		return fmt.Errorf("failed to parse call graph data: %w", err)
	}

	// Add call relationships
	for _, edge := range callGraph.Edges {
		g.Edges = append(g.Edges, &Edge{
			Source: edge.Source,
			Target: edge.Target,
			Type:   "calls",
			Properties: map[string]interface{}{
				"type": "function_call",
			},
			Annotations: make(map[string]string),
		})
	}

	return nil
}

// GetNodeNeighbors returns all nodes connected to the given node
func (g *Graph) GetNodeNeighbors(nodeID string) ([]*Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	var neighbors []*Node
	for _, edge := range g.Edges {
		if edge.Source == nodeID {
			if node, exists := g.Nodes[edge.Target]; exists {
				neighbors = append(neighbors, node)
			}
		}
		if edge.Target == nodeID {
			if node, exists := g.Nodes[edge.Source]; exists {
				neighbors = append(neighbors, node)
			}
		}
	}

	return neighbors, nil
}

// ToJSON serializes the graph to JSON
func (g *Graph) ToJSON() ([]byte, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return json.Marshal(g)
}
