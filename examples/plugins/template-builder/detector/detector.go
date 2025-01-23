package detector

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/types"
)

// DetectStack analyzes a project directory to determine the technology stack
func DetectStack(projectDir string) (*types.ProjectStack, error) {
	// Check if directory is readable
	if _, err := os.ReadDir(projectDir); err != nil {
		return nil, fmt.Errorf("error accessing project directory: %v", err)
	}

	stack := &types.ProjectStack{}

	// Check for Node.js
	if packageJSON, err := os.ReadFile(filepath.Join(projectDir, "package.json")); err == nil {
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if err := json.Unmarshal(packageJSON, &pkg); err != nil {
			return nil, fmt.Errorf("error parsing package.json: %v", err)
		}

		stack.Language = "nodejs"

		// Check for frameworks
		if _, hasReact := pkg.Dependencies["react"]; hasReact {
			stack.Framework = "react"
		} else if _, hasExpress := pkg.Dependencies["express"]; hasExpress {
			stack.Framework = "express"
		}

		// Check for databases
		if _, hasMongo := pkg.Dependencies["mongodb"]; hasMongo {
			stack.Database = "mongodb"
		}
	}

	// Check for Python
	if requirementsTxt, err := os.ReadFile(filepath.Join(projectDir, "requirements.txt")); err == nil {
		stack.Language = "python"

		content := string(requirementsTxt)
		if strings.Contains(strings.ToLower(content), "django") {
			stack.Framework = "django"
		} else if strings.Contains(strings.ToLower(content), "flask") {
			stack.Framework = "flask"
		}

		if strings.Contains(strings.ToLower(content), "psycopg2") {
			stack.Database = "postgres"
		}
	}

	// Check for Go
	if goMod, err := os.ReadFile(filepath.Join(projectDir, "go.mod")); err == nil {
		stack.Language = "go"

		content := string(goMod)
		if strings.Contains(content, "gin-gonic/gin") {
			stack.Framework = "gin"
		}
	}

	// Check for Java
	if pomXML, err := os.ReadFile(filepath.Join(projectDir, "pom.xml")); err == nil {
		var pom struct {
			Dependencies struct {
				Dependency []struct {
					GroupID    string `xml:"groupId"`
					ArtifactID string `xml:"artifactId"`
				} `xml:"dependency"`
			} `xml:"dependencies"`
		}

		if err := xml.Unmarshal(pomXML, &pom); err != nil {
			return nil, fmt.Errorf("error parsing pom.xml: %v", err)
		}

		stack.Language = "java"

		for _, dep := range pom.Dependencies.Dependency {
			if dep.GroupID == "org.springframework.boot" {
				stack.Framework = "spring"
				break
			}
		}
	}

	return stack, nil
}

// contains checks if a string contains another string
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
