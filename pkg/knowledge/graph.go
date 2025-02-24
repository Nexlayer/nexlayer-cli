// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
)

// NodeType represents the type of a node in the knowledge graph
type NodeType string

const (
	// Common node types across languages
	TypeFunction    NodeType = "function"
	TypeClass       NodeType = "class"
	TypeMethod      NodeType = "method"
	TypeModule      NodeType = "module"
	TypePackage     NodeType = "package"
	TypeDependency  NodeType = "dependency"
	TypeAPIEndpoint NodeType = "api_endpoint"
	TypeFile        NodeType = "file"
	TypeConfigFile  NodeType = "config"
	TypeVariable    NodeType = "variable"
	TypeType        NodeType = "type"
	TypeInterface   NodeType = "interface"
	TypeConstant    NodeType = "constant"
	TypeAnnotation  NodeType = "annotation"
	TypeDecorator   NodeType = "decorator"
)

// EdgeType represents the type of relationship between nodes
type EdgeType string

const (
	// Common edge types across languages
	EdgeCalls            EdgeType = "calls"
	EdgeImports          EdgeType = "imports"
	EdgeImplements       EdgeType = "implements"
	EdgeExtends          EdgeType = "extends"
	EdgeDependsOn        EdgeType = "depends_on"
	EdgeCommunicatesWith EdgeType = "communicates_with"
	EdgeDefines          EdgeType = "defines"
	EdgeUses             EdgeType = "uses"
	EdgeDecorates        EdgeType = "decorates"
)

// MetadataType represents the type of metadata stored in node properties
type MetadataType string

const (
	// Metadata types for deployment and configuration
	MetadataLocation   MetadataType = "location"   // File location info
	MetadataVisibility MetadataType = "visibility" // Public/private/exported
	MetadataDeployment MetadataType = "deployment" // Deployment-specific info
	MetadataResource   MetadataType = "resource"   // Resource requirements
	MetadataNetwork    MetadataType = "network"    // Network configuration
	MetadataAuth       MetadataType = "auth"       // Authentication requirements
	MetadataStorage    MetadataType = "storage"    // Storage requirements
)

// Node represents a code entity in the knowledge graph
type Node struct {
	ID          string                       `json:"id"`
	Type        NodeType                     `json:"type"`
	Name        string                       `json:"name"`
	Path        string                       `json:"path,omitempty"`
	Language    string                       `json:"language,omitempty"`
	Metadata    map[MetadataType]interface{} `json:"metadata"`
	Annotations map[string]string            `json:"annotations"`
}

// Edge represents a relationship between nodes
type Edge struct {
	Source      string                       `json:"source"`
	Target      string                       `json:"target"`
	Type        EdgeType                     `json:"type"`
	Metadata    map[MetadataType]interface{} `json:"metadata"`
	Annotations map[string]string            `json:"annotations"`
}

// Graph represents the project's knowledge graph
type Graph struct {
	nodes sync.Map // type: map[string]*Node
	edges sync.Map // type: map[string][]*Edge
}

// NewGraph creates a new knowledge graph
func NewGraph() *Graph {
	return &Graph{}
}

// BuildFromAnalysis constructs the graph from project analysis
func (g *Graph) BuildFromAnalysis(ctx context.Context, projectAnalysis *types.ProjectAnalysis) error {
	// Clear existing graph for full rebuild
	g.nodes = sync.Map{}
	g.edges = sync.Map{}

	// Add nodes for functions (metadata only)
	for file, functions := range projectAnalysis.Functions {
		for _, fn := range functions {
			nodeID := fmt.Sprintf("%s:%s:%s", TypeFunction, file, fn.Name)
			g.nodes.Store(nodeID, &Node{
				ID:   nodeID,
				Type: TypeFunction,
				Name: fn.Name,
				Path: file,
				Metadata: map[MetadataType]interface{}{
					MetadataLocation: map[string]int{
						"startLine": fn.StartLine,
						"endLine":   fn.EndLine,
					},
					MetadataVisibility: fn.IsExported,
				},
				Annotations: make(map[string]string),
			})
		}
	}

	// Add nodes for API endpoints (deployment-focused)
	for _, endpoint := range projectAnalysis.APIEndpoints {
		nodeID := fmt.Sprintf("%s:%s:%s", TypeAPIEndpoint, endpoint.Method, endpoint.Path)
		g.nodes.Store(nodeID, &Node{
			ID:   nodeID,
			Type: TypeAPIEndpoint,
			Name: endpoint.Path,
			Metadata: map[MetadataType]interface{}{
				MetadataNetwork: map[string]interface{}{
					"method": endpoint.Method,
					"path":   endpoint.Path,
				},
				MetadataAuth: map[string]interface{}{
					"handler": endpoint.Handler,
				},
			},
			Annotations: make(map[string]string),
		})
	}

	// Add nodes for dependencies (deployment-focused)
	for depName, depVersion := range projectAnalysis.Dependencies {
		nodeID := fmt.Sprintf("%s:%s", TypeDependency, depName)
		g.nodes.Store(nodeID, &Node{
			ID:   nodeID,
			Type: TypeDependency,
			Name: depName,
			Metadata: map[MetadataType]interface{}{
				MetadataDeployment: map[string]interface{}{
					"version": depVersion,
					"type":    getDeploymentType(depName),
				},
			},
			Annotations: make(map[string]string),
		})
	}

	// Add edges for imports and dependencies (deployment-focused)
	for file, imports := range projectAnalysis.Imports {
		sourceID := fmt.Sprintf("%s:%s", TypeFile, file)
		if _, exists := g.nodes.Load(sourceID); !exists {
			g.nodes.Store(sourceID, &Node{
				ID:   sourceID,
				Type: TypeFile,
				Name: file,
				Metadata: map[MetadataType]interface{}{
					MetadataLocation: map[string]string{
						"path": file,
					},
				},
				Annotations: make(map[string]string),
			})
		}

		for _, imp := range imports {
			targetID := fmt.Sprintf("%s:%s", TypeModule, imp)
			g.nodes.Store(targetID, &Node{
				ID:   targetID,
				Type: TypeModule,
				Name: imp,
				Metadata: map[MetadataType]interface{}{
					MetadataDeployment: map[string]interface{}{
						"path": imp,
					},
				},
				Annotations: make(map[string]string),
			})

			edgeID := fmt.Sprintf("%s-%s", sourceID, targetID)
			g.edges.Store(edgeID, &Edge{
				Source: sourceID,
				Target: targetID,
				Type:   EdgeImports,
				Metadata: map[MetadataType]interface{}{
					MetadataDeployment: map[string]interface{}{
						"importPath": imp,
					},
				},
				Annotations: make(map[string]string),
			})
		}
	}

	return nil
}

// getDeploymentType determines the deployment type based on dependency name
func getDeploymentType(depName string) string {
	switch {
	case strings.Contains(strings.ToLower(depName), "database"),
		strings.Contains(strings.ToLower(depName), "db"),
		strings.Contains(strings.ToLower(depName), "sql"):
		return "database"
	case strings.Contains(strings.ToLower(depName), "cache"),
		strings.Contains(strings.ToLower(depName), "redis"):
		return "cache"
	case strings.Contains(strings.ToLower(depName), "queue"),
		strings.Contains(strings.ToLower(depName), "mq"):
		return "queue"
	case strings.Contains(strings.ToLower(depName), "storage"),
		strings.Contains(strings.ToLower(depName), "s3"):
		return "storage"
	default:
		return "library"
	}
}

// UpdateFromAnalysis incrementally updates the graph with new analysis data
func (g *Graph) UpdateFromAnalysis(ctx context.Context, projectAnalysis *types.ProjectAnalysis) error {
	// Update functions
	for file, functions := range projectAnalysis.Functions {
		for _, fn := range functions {
			nodeID := fmt.Sprintf("%s:%s:%s", TypeFunction, file, fn.Name)
			g.nodes.Store(nodeID, &Node{
				ID:   nodeID,
				Type: TypeFunction,
				Name: fn.Name,
				Path: file,
				Metadata: map[MetadataType]interface{}{
					MetadataLocation: map[string]int{
						"startLine": fn.StartLine,
						"endLine":   fn.EndLine,
					},
					MetadataVisibility: fn.IsExported,
				},
				Annotations: make(map[string]string),
			})
		}
	}

	// Update API endpoints
	for _, endpoint := range projectAnalysis.APIEndpoints {
		nodeID := fmt.Sprintf("%s:%s:%s", TypeAPIEndpoint, endpoint.Method, endpoint.Path)
		g.nodes.Store(nodeID, &Node{
			ID:   nodeID,
			Type: TypeAPIEndpoint,
			Name: endpoint.Path,
			Metadata: map[MetadataType]interface{}{
				MetadataNetwork: map[string]interface{}{
					"method": endpoint.Method,
					"path":   endpoint.Path,
				},
				MetadataAuth: map[string]interface{}{
					"handler": endpoint.Handler,
				},
			},
			Annotations: make(map[string]string),
		})
	}

	return nil
}

// AddCallGraphData integrates call graph data into the knowledge graph
func (g *Graph) AddCallGraphData(callGraphData []byte) error {
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

	// Validate nodes exist before adding edges
	for _, edge := range callGraph.Edges {
		if _, exists := g.nodes.Load(edge.Source); !exists {
			return fmt.Errorf("source node %s not found in graph", edge.Source)
		}
		if _, exists := g.nodes.Load(edge.Target); !exists {
			return fmt.Errorf("target node %s not found in graph", edge.Target)
		}
		g.edges.Store(fmt.Sprintf("%s-%s", edge.Source, edge.Target), &Edge{
			Source: edge.Source,
			Target: edge.Target,
			Type:   EdgeCalls,
			Metadata: map[MetadataType]interface{}{
				MetadataDeployment: map[string]interface{}{
					"type": "function_call",
				},
			},
			Annotations: make(map[string]string),
		})
	}

	return nil
}

// GetNodeNeighbors returns all nodes connected to a given node
func (g *Graph) GetNodeNeighbors(nodeID string) ([]*Node, error) {
	var neighbors []*Node
	g.edges.Range(func(key, value interface{}) bool {
		edge := value.(*Edge)
		if edge.Source == nodeID {
			if nodeVal, exists := g.nodes.Load(edge.Target); exists {
				if node, ok := nodeVal.(*Node); ok {
					neighbors = append(neighbors, node)
				}
			}
		}
		if edge.Target == nodeID {
			if nodeVal, exists := g.nodes.Load(edge.Source); exists {
				if node, ok := nodeVal.(*Node); ok {
					neighbors = append(neighbors, node)
				}
			}
		}
		return true
	})
	return neighbors, nil
}

// GetAPIs returns all API endpoint nodes
func (g *Graph) GetAPIs() []*Node {
	var apis []*Node
	g.nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(*Node); ok {
			if node.Type == "api_endpoint" {
				apis = append(apis, node)
			}
		}
		return true
	})
	return apis
}

// ToJSON serializes the graph to JSON, excluding sensitive data
func (g *Graph) ToJSON() ([]byte, error) {
	graph := struct {
		Nodes map[string]*Node `json:"nodes"`
		Edges []*Edge          `json:"edges"`
	}{
		Nodes: make(map[string]*Node),
		Edges: make([]*Edge, 0),
	}

	// Sanitize and collect nodes
	g.nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(*Node); ok {
			// Create a sanitized copy
			sanitizedNode := &Node{
				ID:          node.ID,
				Type:        node.Type,
				Name:        node.Name,
				Path:        node.Path,
				Language:    node.Language,
				Metadata:    make(map[MetadataType]interface{}),
				Annotations: node.Annotations,
			}

			// Only include deployment-relevant metadata
			for mType, mValue := range node.Metadata {
				switch mType {
				case MetadataDeployment, MetadataNetwork, MetadataResource, MetadataAuth, MetadataStorage:
					sanitizedNode.Metadata[mType] = mValue
				case MetadataLocation:
					// Include only file location, not line numbers
					if loc, ok := mValue.(map[string]interface{}); ok {
						sanitizedNode.Metadata[mType] = map[string]interface{}{
							"path": loc["path"],
						}
					}
				}
			}

			graph.Nodes[key.(string)] = sanitizedNode
		}
		return true
	})

	// Sanitize and collect edges
	g.edges.Range(func(key, value interface{}) bool {
		if edge, ok := value.(*Edge); ok {
			// Create a sanitized copy
			sanitizedEdge := &Edge{
				Source:      edge.Source,
				Target:      edge.Target,
				Type:        edge.Type,
				Metadata:    make(map[MetadataType]interface{}),
				Annotations: edge.Annotations,
			}

			// Only include deployment-relevant metadata
			for mType, mValue := range edge.Metadata {
				switch mType {
				case MetadataDeployment, MetadataNetwork, MetadataResource:
					sanitizedEdge.Metadata[mType] = mValue
				}
			}

			graph.Edges = append(graph.Edges, sanitizedEdge)
		}
		return true
	})

	return json.Marshal(graph)
}
