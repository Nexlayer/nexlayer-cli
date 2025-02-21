// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package ai

import (
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// enhancePromptWithAnalysis enhances the base prompt with project analysis information.
// It adds structured information about functions, API endpoints, and dependencies to help
// the AI provider generate more accurate and context-aware responses.
func enhancePromptWithAnalysis(basePrompt string, analysis *detection.ProjectAnalysis) string {
	if analysis == nil {
		return basePrompt
	}

	var sb strings.Builder
	sb.WriteString(basePrompt)
	sb.WriteString("\n\nProject Analysis:\n")

	// Add detected functions with file context
	if len(analysis.Functions) > 0 {
		sb.WriteString("\nFunctions:\n")
		for _, fn := range analysis.Functions {
			// Include additional function details if available
			details := []string{fn.Name}
			if fn.FilePath != "" {
				details = append(details, fmt.Sprintf("file: %s", fn.FilePath))
			}
			if fn.Type != "" {
				details = append(details, fmt.Sprintf("type: %s", fn.Type))
			}
			if len(fn.Tags) > 0 {
				details = append(details, fmt.Sprintf("tags: %s", strings.Join(fn.Tags, ", ")))
			}
			sb.WriteString(fmt.Sprintf("- %s\n", strings.Join(details, ", ")))
		}
	}

	// Add detected API endpoints with method and path
	if len(analysis.Endpoints) > 0 {
		sb.WriteString("\nAPI Endpoints:\n")
		for _, ep := range analysis.Endpoints {
			// Include endpoint details
			details := []string{
				fmt.Sprintf("%s %s", ep.Method, ep.Path),
				fmt.Sprintf("handler: %s", ep.Handler),
			}
			if len(ep.Tags) > 0 {
				details = append(details, fmt.Sprintf("tags: %s", strings.Join(ep.Tags, ", ")))
			}
			sb.WriteString(fmt.Sprintf("- %s\n", strings.Join(details, ", ")))
		}
	}

	// Add detected dependencies with version information
	if len(analysis.Dependencies) > 0 {
		sb.WriteString("\nDependencies:\n")
		for _, dep := range analysis.Dependencies {
			// Include dependency details
			details := []string{
				fmt.Sprintf("%s@%s", dep.Name, dep.Version),
				fmt.Sprintf("type: %s", dep.Type),
			}
			if dep.IsAIRelated {
				details = append(details, "AI-related")
			}
			if len(dep.Features) > 0 {
				details = append(details, fmt.Sprintf("features: %s", strings.Join(dep.Features, ", ")))
			}
			sb.WriteString(fmt.Sprintf("- %s\n", strings.Join(details, ", ")))
		}
	}

	return sb.String()
}
