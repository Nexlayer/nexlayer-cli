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

	// Include AI-powered IDE & LLM model
	fmt.Fprintf(&builder, "  environment:\n")
	fmt.Fprintf(&builder, "    ai_ide: %s\n", info.LLMProvider)
	fmt.Fprintf(&builder, "    llm_model: %s\n", info.LLMModel)

	// Write pods section
	fmt.Fprintf(&builder, "  pods:\n")
	fmt.Fprintf(&builder, "    - name: app\n")

	// Determine image based on project type
	switch info.Type {
	case TypeNextjs:
		fmt.Fprintf(&builder, "      image: <%% REGISTRY %%>/%s\n", info.Name)
		fmt.Fprintf(&builder, "      path: /\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)
		if info.HasDocker {
			fmt.Fprintf(&builder, "      # Using custom Dockerfile\n")
		} else {
			fmt.Fprintf(&builder, "      # Using Next.js default configuration\n")
		}

	case TypeReact:
		fmt.Fprintf(&builder, "      image: <%% REGISTRY %%>/%s\n", info.Name)
		fmt.Fprintf(&builder, "      path: /\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)
		fmt.Fprintf(&builder, "      # Static file serving for React app\n")

	case TypeNode:
		fmt.Fprintf(&builder, "      image: <%% REGISTRY %%>/%s\n", info.Name)
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)
		if len(info.Scripts) > 0 {
			fmt.Fprintf(&builder, "      # Available npm scripts:\n")
			for script := range info.Scripts {
				fmt.Fprintf(&builder, "      #   - %s\n", script)
			}
		}

	case TypePython:
		fmt.Fprintf(&builder, "      image: <%% REGISTRY %%>/%s\n", info.Name)
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)
		fmt.Fprintf(&builder, "      # Python web application\n")

	case TypeGo:
		fmt.Fprintf(&builder, "      image: <%% REGISTRY %%>/%s\n", info.Name)
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)
		fmt.Fprintf(&builder, "      # Go web application\n")

	case TypeDockerRaw:
		fmt.Fprintf(&builder, "      image: <%% REGISTRY %%>/%s\n", info.Name)
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)
		fmt.Fprintf(&builder, "      # Raw Docker application\n")

	default:
		return "", fmt.Errorf("unsupported project type: %s", info.Type)
	}

	// Add registry login if using private registry
	fmt.Fprintf(&builder, "\n  # Registry authentication (required for private images)\n")
	fmt.Fprintf(&builder, "  registryLogin:\n")
	fmt.Fprintf(&builder, "    registry: <%% REGISTRY %%>\n")
	fmt.Fprintf(&builder, "    username: \"\"\n")
	fmt.Fprintf(&builder, "    personalAccessToken: \"\"\n")

	return builder.String(), nil
}
