package detector

import (
	"os"
	"path/filepath"
	"strings"
)

// StackInfo represents detected project stack information
type StackInfo struct {
	Language    string
	Framework   string
	Database    string
	HasDocker   bool
	HasAI       bool
	Template    string
}

// DetectStack analyzes the project directory to determine the tech stack
func DetectStack(projectPath string) (*StackInfo, error) {
	info := &StackInfo{}

	// Walk through the project directory
	err := filepath.Walk(projectPath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip vendor and node_modules
		if f.IsDir() && (f.Name() == "vendor" || f.Name() == "node_modules" || f.Name() == ".git") {
			return filepath.SkipDir
		}

		switch f.Name() {
		// JavaScript/Node.js detection
		case "package.json":
			info.Language = "JavaScript"
			content, err := os.ReadFile(path)
			if err == nil {
				if strings.Contains(string(content), "next") {
					info.Framework = "Next.js"
				} else if strings.Contains(string(content), "react") {
					info.Framework = "React"
				} else if strings.Contains(string(content), "vue") {
					info.Framework = "Vue"
				}
				// AI framework detection
				if strings.Contains(string(content), "langchain") {
					info.HasAI = true
					info.Template = "langchain-nextjs"
				} else if strings.Contains(string(content), "openai") {
					info.HasAI = true
					info.Template = "openai-node"
				}
			}

		// Python detection
		case "requirements.txt", "pyproject.toml":
			info.Language = "Python"
			content, err := os.ReadFile(path)
			if err == nil {
				if strings.Contains(string(content), "fastapi") {
					info.Framework = "FastAPI"
				} else if strings.Contains(string(content), "flask") {
					info.Framework = "Flask"
				}
				// AI framework detection
				if strings.Contains(string(content), "langchain") {
					info.HasAI = true
					info.Template = "langchain-fastapi"
				} else if strings.Contains(string(content), "openai") {
					info.HasAI = true
					info.Template = "openai-py"
				} else if strings.Contains(string(content), "transformers") {
					info.HasAI = true
					info.Template = "huggingface"
				}
			}


		// Docker detection
		case "Dockerfile":
			info.HasDocker = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Set default template based on stack if no AI template was detected
	if info.Template == "" {
		switch {
		case info.Language == "JavaScript" && info.Database == "MongoDB":
			info.Template = "mern-stack"
		case info.Language == "JavaScript" && info.Database == "PostgreSQL":
			info.Template = "pern-stack"
		case info.Framework == "FastAPI":
			info.Template = "fastapi-starter"
		default:
			info.Template = "blank"
		}
	}

	return info, nil
}
