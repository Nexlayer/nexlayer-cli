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
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/schema"
)

// LLMContext represents the enriched context for LLM interactions
type LLMContext struct {
	ProjectStructure map[string]interface{} `json:"project_structure"`
	Dependencies     map[string]string      `json:"dependencies"`
	APIEndpoints     []interface{}          `json:"api_endpoints"`
	PodFlows         []interface{}          `json:"pod_flows"`
	Patterns         []interface{}          `json:"patterns"`
	Languages        map[string]interface{} `json:"languages"`
	Frameworks       map[string]interface{} `json:"frameworks"`
	Resources        map[string]interface{} `json:"resources"`
	Network          map[string]interface{} `json:"network"`
	Storage          map[string]interface{} `json:"storage"`
}

// LLMEnricher enriches the knowledge graph with LLM metadata
type LLMEnricher struct {
	graph       *Graph
	metadata    map[string]interface{}
	metadataMu  sync.RWMutex
	metadataDir string
}

// NewLLMEnricher creates a new LLM metadata enricher
func NewLLMEnricher(graph *Graph, metadataDir string) *LLMEnricher {
	return &LLMEnricher{
		graph:       graph,
		metadata:    make(map[string]interface{}),
		metadataDir: metadataDir,
	}
}

// LoadMetadata loads LLM metadata from the tools directory with caching
func (e *LLMEnricher) LoadMetadata() error {
	e.metadataMu.Lock()
	defer e.metadataMu.Unlock()

	metadataPath := filepath.Join(e.metadataDir, "llm", "metadata.json")
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
func (e *LLMEnricher) EnrichContext(ctx context.Context, yamlConfig *template.NexlayerYAML) (*LLMContext, error) {
	enriched := &LLMContext{
		ProjectStructure: make(map[string]interface{}),
		Dependencies:     make(map[string]string),
		APIEndpoints:     make([]interface{}, 0),
		PodFlows:         make([]interface{}, 0),
		Patterns:         make([]interface{}, 0),
		Languages:        make(map[string]interface{}),
		Frameworks:       make(map[string]interface{}),
		Resources:        make(map[string]interface{}),
		Network:          make(map[string]interface{}),
		Storage:          make(map[string]interface{}),
	}

	// Extract project structure (only deployment-relevant files)
	e.graph.nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(*Node); ok {
			if node.Type == TypeFile {
				ext := strings.ToLower(filepath.Ext(node.Path))
				switch ext {
				case ".yaml", ".yml", ".json", ".toml":
					// Include configuration files
					parts := strings.Split(node.Path, string(os.PathSeparator))
					current := enriched.ProjectStructure
					for i, part := range parts {
						if i == len(parts)-1 {
							if metadata, ok := node.Metadata[MetadataDeployment]; ok {
								current[part] = metadata
							} else {
								current[part] = nil
							}
						} else {
							if _, exists := current[part]; !exists {
								current[part] = make(map[string]interface{})
							}
							if nextLevel, ok := current[part].(map[string]interface{}); ok {
								current = nextLevel
							} else {
								// Handle type mismatch gracefully
								current[part] = make(map[string]interface{})
								current = current[part].(map[string]interface{})
							}
						}
					}
				}
			}
		}
		return true
	})

	// Extract deployment-relevant information
	e.graph.nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(*Node); ok {
			switch node.Type {
			case TypeAPIEndpoint:
				if network, ok := node.Metadata[MetadataNetwork].(map[string]interface{}); ok {
					endpoint := map[string]interface{}{
						"path":   network["path"],
						"method": network["method"],
					}
					if auth, ok := node.Metadata[MetadataAuth].(map[string]interface{}); ok {
						endpoint["auth"] = auth
					}
					enriched.APIEndpoints = append(enriched.APIEndpoints, endpoint)
				}
			case TypeDependency:
				if deployment, ok := node.Metadata[MetadataDeployment].(map[string]interface{}); ok {
					depType := deployment["type"].(string)
					switch depType {
					case "database", "cache", "queue", "storage":
						enriched.Resources[node.Name] = deployment
					}
					enriched.Dependencies[node.Name] = deployment["version"].(string)
				}
			}
		}
		return true
	})

	// Extract pod communication flows from nexlayer.yaml
	if yamlConfig != nil {
		podMap := make(map[string]bool)
		podNetworking := make(map[string]interface{})
		podStorage := make(map[string]interface{})

		for _, pod := range yamlConfig.Application.Pods {
			podMap[pod.Name] = true

			// Collect networking config
			if len(pod.ServicePorts) > 0 {
				podNetworking[pod.Name] = map[string]interface{}{
					"ports": pod.ServicePorts,
					"path":  pod.Path,
				}
			}

			// Collect storage requirements
			if len(pod.Volumes) > 0 {
				podStorage[pod.Name] = pod.Volumes
			}
		}

		enriched.Network = podNetworking
		enriched.Storage = podStorage

		// Extract pod communication flows
		for _, pod := range yamlConfig.Application.Pods {
			for _, v := range pod.Vars {
				if strings.Contains(v.Value, ".pod") {
					targetPod := strings.Split(v.Value, ".pod")[0]
					if podMap[targetPod] {
						flow := map[string]interface{}{
							"source": pod.Name,
							"target": targetPod,
							"var":    v.Key,
							"value":  v.Value,
						}
						enriched.PodFlows = append(enriched.PodFlows, flow)
					}
				}
			}
		}
	}

	// Add patterns from cached metadata
	e.metadataMu.RLock()
	if patterns, ok := e.metadata["deployment_patterns"].([]interface{}); ok {
		enriched.Patterns = patterns
	}
	e.metadataMu.RUnlock()

	return enriched, nil
}

// GeneratePrompt generates an LLM prompt with enriched context
func (e *LLMEnricher) GeneratePrompt(ctx context.Context, basePrompt string, yamlConfig *template.NexlayerYAML) (string, error) {
	enriched, err := e.EnrichContext(ctx, yamlConfig)
	if err != nil {
		return "", fmt.Errorf("failed to enrich context: %w", err)
	}

	// Convert enriched context to JSON
	contextJSON, err := json.MarshalIndent(enriched, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal context: %w", err)
	}

	// Combine base prompt with enriched context
	prompt := fmt.Sprintf("%s\n\nContext:\n%s", basePrompt, string(contextJSON))
	return prompt, nil
}
