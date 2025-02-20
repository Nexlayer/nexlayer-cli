package detection

import (
	"fmt"
	"strings"
)

// GenerateYAMLFromTemplate generates a nexlayer.yaml based on detected project info using templates
func GenerateYAMLFromTemplate(info *ProjectInfo) (string, error) {
	var builder strings.Builder

	// Write application section
	fmt.Fprintf(&builder, "application:\n")
	fmt.Fprintf(&builder, "  name: %s\n", sanitizeName(info.Name))

	// Add pods section
	fmt.Fprintf(&builder, "  pods:\n")

	// Add default pod if no project type detected
	if info.Type == TypeUnknown {
		fmt.Fprintf(&builder, "    - name: app\n")
		fmt.Fprintf(&builder, "      type: frontend\n")
		fmt.Fprintf(&builder, "      image: docker.io/library/nginx:latest\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - name: http\n")
		fmt.Fprintf(&builder, "          port: 80\n")
		fmt.Fprintf(&builder, "          targetPort: 80\n")
		fmt.Fprintf(&builder, "      vars:\n")
		fmt.Fprintf(&builder, "        - key: NODE_ENV\n")
		fmt.Fprintf(&builder, "          value: production\n")
		return builder.String(), nil
	}

	// Determine pod configuration based on project type
	switch info.Type {
	case TypeNextjs, TypeReact:
		// Frontend pod
		fmt.Fprintf(&builder, "    - name: web\n")
		fmt.Fprintf(&builder, "      type: nextjs\n")
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s%s\n", sanitizeName(info.Name), getImageTag(info))
		fmt.Fprintf(&builder, "      path: /\n")
		fmt.Fprintf(&builder, "      vars:\n")
		fmt.Fprintf(&builder, "        - key: NODE_ENV\n")
		fmt.Fprintf(&builder, "          value: production\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - name: http\n")
		fmt.Fprintf(&builder, "          port: %d\n", info.Port)
		fmt.Fprintf(&builder, "          targetPort: %d\n", info.Port)
		fmt.Fprintf(&builder, "          protocol: TCP\n")

	case TypeNode:
		// Backend API pod
		fmt.Fprintf(&builder, "    - name: api\n")
		fmt.Fprintf(&builder, "      type: backend\n")
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s%s\n", sanitizeName(info.Name), getImageTag(info))
		fmt.Fprintf(&builder, "      path: /api\n")
		fmt.Fprintf(&builder, "      vars:\n")
		fmt.Fprintf(&builder, "        - key: NODE_ENV\n")
		fmt.Fprintf(&builder, "          value: production\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - name: http\n")
		fmt.Fprintf(&builder, "          port: %d\n", info.Port)
		fmt.Fprintf(&builder, "          targetPort: %d\n", info.Port)

	case TypePython:
		// Python backend pod
		fmt.Fprintf(&builder, "    - name: api\n")
		fmt.Fprintf(&builder, "      type: backend\n")
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s%s\n", sanitizeName(info.Name), getImageTag(info))
		fmt.Fprintf(&builder, "      path: /api\n")
		fmt.Fprintf(&builder, "      vars:\n")
		fmt.Fprintf(&builder, "        - key: PYTHON_ENV\n")
		fmt.Fprintf(&builder, "          value: production\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - name: http\n")
		fmt.Fprintf(&builder, "          port: %d\n", info.Port)
		fmt.Fprintf(&builder, "          targetPort: %d\n", info.Port)

	case TypeGo:
		// Go backend pod
		fmt.Fprintf(&builder, "    - name: api\n")
		fmt.Fprintf(&builder, "      type: backend\n")
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s%s\n", sanitizeName(info.Name), getImageTag(info))
		fmt.Fprintf(&builder, "      path: /api\n")
		fmt.Fprintf(&builder, "      vars:\n")
		fmt.Fprintf(&builder, "        - key: GO_ENV\n")
		fmt.Fprintf(&builder, "          value: production\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - name: http\n")
		fmt.Fprintf(&builder, "          port: %d\n", info.Port)
		fmt.Fprintf(&builder, "          targetPort: %d\n", info.Port)

	case TypeDockerRaw:
		// Raw Docker pod
		fmt.Fprintf(&builder, "    - name: app\n")
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s%s\n", sanitizeName(info.Name), getImageTag(info))
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - name: http\n")
		fmt.Fprintf(&builder, "          port: %d\n", info.Port)
		fmt.Fprintf(&builder, "          targetPort: %d\n", info.Port)

	default:
		return "", fmt.Errorf("unsupported project type: %s", info.Type)
	}

	// Print the generated YAML for debugging
	fmt.Printf("Generated YAML:\n%s\n", builder.String())

	return builder.String(), nil
}

func getImageTag(info *ProjectInfo) string {
	if info.ImageTag != "" {
		return ":" + info.ImageTag
	}
	return ":latest"
}

// sanitizeName converts a string to a valid Nexlayer name
// - must start with a lowercase letter
// - can include only alphanumeric characters, '-', '.'
func sanitizeName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace invalid characters with hyphens
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '.' {
			return r
		}
		return '-'
	}, name)

	// Ensure starts with a lowercase letter
	if len(name) > 0 && (name[0] < 'a' || name[0] > 'z') {
		name = "app-" + name
	}

	// If empty after sanitization, use default
	if name == "" {
		name = "app"
	}

	return name
}
