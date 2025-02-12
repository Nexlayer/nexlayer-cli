package detection

import (
	"fmt"
	"strings"
)

// GenerateYAML generates a nexlayer.yaml based on detected project info
func GenerateYAML(info *ProjectInfo) (string, error) {
	var builder strings.Builder

	// Write application section
	fmt.Fprintf(&builder, "application:\n")
	fmt.Fprintf(&builder, "  name: %s\n", info.Name)

	// Include AI-powered IDE & LLM model if available
	if info.LLMProvider != "" || info.LLMModel != "" {
		fmt.Fprintf(&builder, "  environment:\n")
		if info.LLMProvider != "" {
			fmt.Fprintf(&builder, "    ai_ide: %s\n", info.LLMProvider)
		}
		if info.LLMModel != "" {
			fmt.Fprintf(&builder, "    llm_model: %s\n", info.LLMModel)
		}
	}

	// Write pods section
	fmt.Fprintf(&builder, "  pods:\n")

	// Add default pod if no project type detected
	if info.Type == "" {
		fmt.Fprintf(&builder, "    - name: app\n")
		fmt.Fprintf(&builder, "      image: nginx:latest\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - 80\n")
		return builder.String(), nil
	}

	// Determine image based on project type
	fmt.Fprintf(&builder, "    - name: app\n")
	switch info.Type {
	case TypeNextjs:
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s\n", info.Name)
		fmt.Fprintf(&builder, "      path: /\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	case TypeReact:
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s\n", info.Name)
		fmt.Fprintf(&builder, "      path: /\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	case TypeNode:
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s\n", info.Name)
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	case TypePython:
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s\n", info.Name)
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	case TypeGo:
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s\n", info.Name)
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	case TypeDockerRaw:
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s\n", info.Name)
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	default:
		return "", fmt.Errorf("unsupported project type: %s", info.Type)
	}

	return builder.String(), nil
}
