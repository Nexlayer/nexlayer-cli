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
		for file, functions := range analysis.Functions {
			if len(functions) == 0 {
				continue
			}
			sb.WriteString(fmt.Sprintf("\nFile: %s\n", file))
			for _, fn := range functions {
				// Include additional function details if available
				details := []string{fn.Name}
				if fn.IsExported {
					details = append(details, "exported")
				}
				if fn.Signature != "" {
					details = append(details, fmt.Sprintf("signature: %s", fn.Signature))
				}
				sb.WriteString(fmt.Sprintf("- %s\n", strings.Join(details, ", ")))
			}
		}
	}

	// Add detected API endpoints with method and path
	if len(analysis.APIEndpoints) > 0 {
		sb.WriteString("\nAPI Endpoints:\n")
		for _, ep := range analysis.APIEndpoints {
			// Include endpoint parameters if available
			if len(ep.Parameters) > 0 {
				var params []string
				for k, v := range ep.Parameters {
					params = append(params, fmt.Sprintf("%s: %s", k, v))
				}
				sb.WriteString(fmt.Sprintf("- %s %s [%s]\n", ep.Method, ep.Path, strings.Join(params, ", ")))
			} else {
				sb.WriteString(fmt.Sprintf("- %s %s\n", ep.Method, ep.Path))
			}
		}
	}

	// Add detected dependencies with version information
	if len(analysis.Dependencies) > 0 {
		sb.WriteString("\nDependencies:\n")
		for pkg, deps := range analysis.Dependencies {
			if len(deps) == 0 {
				continue
			}
			sb.WriteString(fmt.Sprintf("\nPackage: %s\n", pkg))
			for _, dep := range deps {
				// Include dependency type and version
				sb.WriteString(fmt.Sprintf("- %s@%s (%s)\n", dep.Name, dep.Version, dep.Type))
			}
		}
	}

	return sb.String()
}
