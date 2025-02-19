// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package ai

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/template"
)

// Component represents a detected application component
type Component struct {
	Name    string
	Type    string
	Image   string
	Ports   []int
	EnvVars []template.EnvVar
	Volumes []template.Volume
	Secrets []template.Secret
}

// AnalysisResult represents the project analysis output
type AnalysisResult struct {
	Components []*Component
}

// GraphResult represents the knowledge graph output
type GraphResult struct {
	Nodes []*GraphNode
}

// GraphNode represents a node in the knowledge graph
type GraphNode struct {
	Name        string
	EnvVars     map[string]string
	Annotations map[string]string
}

// createTemplate creates a template based on project analysis
func createTemplate(projectInfo *ProjectInfo, analysis *AnalysisResult, graph *GraphResult) (*template.NexlayerYAML, error) {
	// Create a new template parser
	parser, err := template.NewParser("")
	if err != nil {
		return nil, fmt.Errorf("failed to create template parser: %w", err)
	}

	// Load the base template
	if err := parser.LoadTemplate(); err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	// Create detected settings
	detected := &template.NexlayerYAML{
		Application: template.Application{
			Name: projectInfo.Name,
			Pods: []template.Pod{},
		},
	}

	// Add detected pods
	for _, component := range analysis.Components {
		pod := template.Pod{
			Name: component.Name,
			Type: detectPodType(projectInfo, component),
		}

		// Set image if available
		if component.Image != "" {
			pod.Image = component.Image
		}

		// Add service ports
		for _, port := range component.Ports {
			pod.ServicePorts = append(pod.ServicePorts, template.ServicePort{
				Name:       fmt.Sprintf("%s-%d", component.Name, port),
				Port:       port,
				TargetPort: port,
			})
		}

		// Add environment variables
		for _, env := range component.EnvVars {
			pod.Vars = append(pod.Vars, env)
		}

		// Add volumes if needed
		if len(component.Volumes) > 0 {
			pod.Volumes = append(pod.Volumes, component.Volumes...)
		}

		// Add secrets if needed
		if len(component.Secrets) > 0 {
			pod.Secrets = append(pod.Secrets, component.Secrets...)
		}

		detected.Application.Pods = append(detected.Application.Pods, pod)
	}

	// Enrich with graph insights
	if graph != nil {
		enrichTemplateWithGraph(detected, graph)
	}

	// Merge with base template
	final, err := parser.MergeWithDetected(detected)
	if err != nil {
		return nil, fmt.Errorf("failed to merge template: %w", err)
	}

	return final, nil
}

// detectPodType determines the pod type based on project info and analysis
func detectPodType(info *ProjectInfo, component *Component) string {
	// First check if we have an explicit type from analysis
	if component.Type != "" {
		return component.Type
	}

	// Detect based on project info
	switch info.Type {
	case "nextjs":
		return "nextjs"
	case "react":
		return "react"
	case "node":
		return "node"
	case "python":
		return "python"
	case "go":
		return "golang"
	default:
		// Default to raw type if we can't determine
		return "raw"
	}
}

// enrichTemplateWithGraph adds graph-based insights to the template
func enrichTemplateWithGraph(tmpl *template.NexlayerYAML, graph *GraphResult) {
	// Add any insights from the graph analysis to enrich the template
	for _, node := range graph.Nodes {
		for i, pod := range tmpl.Application.Pods {
			if pod.Name == node.Name {
				// Add any additional environment variables
				for k, v := range node.EnvVars {
					tmpl.Application.Pods[i].Vars = append(tmpl.Application.Pods[i].Vars, template.EnvVar{
						Key:   k,
						Value: v,
					})
				}

				// Add any annotations
				if len(node.Annotations) > 0 {
					if tmpl.Application.Pods[i].Annotations == nil {
						tmpl.Application.Pods[i].Annotations = make(map[string]string)
					}
					for k, v := range node.Annotations {
						tmpl.Application.Pods[i].Annotations[k] = v
					}
				}
			}
		}
	}
}
