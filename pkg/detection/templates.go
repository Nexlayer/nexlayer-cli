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
	fmt.Fprintf(&builder, "  name: %s\n", info.Name)

	// Add pods section
	fmt.Fprintf(&builder, "  pods:\n")

	// Add default pod if no project type detected
	if info.Type == "" {
		fmt.Fprintf(&builder, "    - name: app\n")
		fmt.Fprintf(&builder, "      type: frontend\n")
		fmt.Fprintf(&builder, "      image: nginx:latest\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - 80\n")
		return builder.String(), nil
	}

	// Determine pod configuration based on project type
	switch info.Type {
	case TypeNextjs, TypeReact:
		// Frontend pod
		fmt.Fprintf(&builder, "    - name: web\n")
		fmt.Fprintf(&builder, "      type: frontend\n")
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s%s\n", info.Name, getImageTag(info))
		fmt.Fprintf(&builder, "      path: /\n")
		fmt.Fprintf(&builder, "      vars:\n")
		fmt.Fprintf(&builder, "        NODE_ENV: production\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	case TypeNode:
		// Backend API pod
		fmt.Fprintf(&builder, "    - name: api\n")
		fmt.Fprintf(&builder, "      type: backend\n")
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s%s\n", info.Name, getImageTag(info))
		fmt.Fprintf(&builder, "      path: /api\n")
		fmt.Fprintf(&builder, "      vars:\n")
		fmt.Fprintf(&builder, "        NODE_ENV: production\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	case TypePython:
		// Python backend pod
		fmt.Fprintf(&builder, "    - name: api\n")
		fmt.Fprintf(&builder, "      type: backend\n")
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s%s\n", info.Name, getImageTag(info))
		fmt.Fprintf(&builder, "      path: /api\n")
		fmt.Fprintf(&builder, "      vars:\n")
		fmt.Fprintf(&builder, "        PYTHON_ENV: production\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	case TypeGo:
		// Go backend pod
		fmt.Fprintf(&builder, "    - name: api\n")
		fmt.Fprintf(&builder, "      type: backend\n")
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s%s\n", info.Name, getImageTag(info))
		fmt.Fprintf(&builder, "      path: /api\n")
		fmt.Fprintf(&builder, "      vars:\n")
		fmt.Fprintf(&builder, "        GO_ENV: production\n")
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	case TypeDockerRaw:
		// Raw Docker pod
		fmt.Fprintf(&builder, "    - name: app\n")
		fmt.Fprintf(&builder, "      image: ghcr.io/nexlayer/%s%s\n", info.Name, getImageTag(info))
		fmt.Fprintf(&builder, "      servicePorts:\n")
		fmt.Fprintf(&builder, "        - %d\n", info.Port)

	default:
		return "", fmt.Errorf("unsupported project type: %s", info.Type)
	}

	return builder.String(), nil
}

func getImageTag(info *ProjectInfo) string {
	if info.ImageTag != "" {
		return ":" + info.ImageTag
	}
	return ":latest"
}
