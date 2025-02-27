// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detectors

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
)

// AIFrameworkDetector detects AI/ML frameworks
type AIFrameworkDetector struct{}

// Priority returns the detector's priority level
func (d *AIFrameworkDetector) Priority() int {
	return 150 // Higher priority than standard framework detectors
}

// Detect analyzes a directory to detect AI framework usage
func (d *AIFrameworkDetector) Detect(dir string) (*types.ProjectInfo, error) {
	projectInfo := &types.ProjectInfo{
		Type:         types.TypeUnknown,
		Dependencies: make(map[string]string),
		Scripts:      make(map[string]string),
	}

	// JS/TS project detection
	if isJSProject, jsInfo := d.detectJSProject(dir); isJSProject {
		return jsInfo, nil
	}

	// Python project detection
	if isPythonProject, pyInfo := d.detectPythonProject(dir); isPythonProject {
		return pyInfo, nil
	}

	return projectInfo, nil
}

// detectJSProject checks for JavaScript/TypeScript AI frameworks
func (d *AIFrameworkDetector) detectJSProject(dir string) (bool, *types.ProjectInfo) {
	projectInfo := &types.ProjectInfo{
		Type:         types.TypeUnknown,
		Dependencies: make(map[string]string),
		Scripts:      make(map[string]string),
	}

	// Check for package.json
	packageJSONPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		return false, projectInfo
	}

	// Read and parse package.json
	packageJSONBytes, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return false, projectInfo
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(packageJSONBytes, &packageJSON); err != nil {
		return false, projectInfo
	}

	// Extract project name and version
	if name, ok := packageJSON["name"].(string); ok {
		projectInfo.Name = name
	}
	if version, ok := packageJSON["version"].(string); ok {
		projectInfo.Version = version
	}

	// Check for dependencies
	dependencies := make(map[string]interface{})
	if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		dependencies = deps
	}
	if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
		for k, v := range devDeps {
			dependencies[k] = v
		}
	}

	// Extract all dependencies
	for name, version := range dependencies {
		if vStr, ok := version.(string); ok {
			projectInfo.Dependencies[name] = vStr
		}
	}

	// Check for specific AI frameworks
	hasLangChain := d.hasJSDependency(dependencies, "langchain")
	hasOpenAI := d.hasJSDependency(dependencies, "openai")
	hasNextJS := d.hasJSDependency(dependencies, "next")
	hasReact := d.hasJSDependency(dependencies, "react")
	hasVectorDB := d.hasVectorDBDependency(dependencies)
	hasRAG := false

	// Check for RAG implementation by scanning files
	if (hasLangChain || hasOpenAI) && hasVectorDB {
		hasRAG = d.hasRAGImplementation(dir)
	}

	// Determine the primary AI framework type
	if hasRAG {
		projectInfo.Type = types.TypeRAG
		projectInfo.LLMProvider = "Multiple"
	} else if hasLangChain && hasNextJS {
		projectInfo.Type = types.TypeLangchainNextjs
		projectInfo.LLMProvider = "LangChain"
	} else if hasLangChain {
		projectInfo.Type = types.TypeLangchainJS
		projectInfo.LLMProvider = "LangChain"
	} else if hasOpenAI {
		projectInfo.Type = types.TypeOpenAINode
		projectInfo.LLMProvider = "OpenAI"
	} else {
		// Not an AI project
		return false, projectInfo
	}

	// Check for scripts
	if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
		for name, cmd := range scripts {
			if cmdStr, ok := cmd.(string); ok {
				projectInfo.Scripts[name] = cmdStr
			}
		}
	}

	// Check for Docker configuration
	if _, err := os.Stat(filepath.Join(dir, "Dockerfile")); err == nil {
		projectInfo.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yml")); err == nil {
		projectInfo.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yaml")); err == nil {
		projectInfo.HasDocker = true
	}

	// Set default port
	if hasNextJS {
		projectInfo.Port = 3000
	} else if hasReact {
		projectInfo.Port = 5173
	} else {
		projectInfo.Port = 3000
	}

	// Try to identify specific LLM model from code
	if model := d.detectLLMModel(dir); model != "" {
		projectInfo.LLMModel = model
	}

	return true, projectInfo
}

// detectPythonProject checks for Python AI frameworks
func (d *AIFrameworkDetector) detectPythonProject(dir string) (bool, *types.ProjectInfo) {
	projectInfo := &types.ProjectInfo{
		Type:         types.TypeUnknown,
		Dependencies: make(map[string]string),
		Scripts:      make(map[string]string),
	}

	// Check for requirements.txt or setup.py
	requirementsPath := filepath.Join(dir, "requirements.txt")
	setupPath := filepath.Join(dir, "setup.py")
	pyprojectPath := filepath.Join(dir, "pyproject.toml")

	requirementsExists := false
	setupExists := false
	pyprojectExists := false

	if _, err := os.Stat(requirementsPath); err == nil {
		requirementsExists = true
	}
	if _, err := os.Stat(setupPath); err == nil {
		setupExists = true
	}
	if _, err := os.Stat(pyprojectPath); err == nil {
		pyprojectExists = true
	}

	if !requirementsExists && !setupExists && !pyprojectExists {
		return false, projectInfo
	}

	// Set default name from directory
	projectInfo.Name = filepath.Base(dir)

	// Parse dependencies
	dependencies := make(map[string]string)

	if requirementsExists {
		if deps, err := d.readRequirementsTxt(requirementsPath); err == nil {
			for k, v := range deps {
				dependencies[k] = v
			}
		}
	}

	if setupExists {
		if deps, err := d.readSetupPy(setupPath); err == nil {
			for k, v := range deps {
				dependencies[k] = v
			}
		}
	}

	if pyprojectExists {
		if deps, err := d.readPyprojectToml(pyprojectPath); err == nil {
			for k, v := range deps {
				dependencies[k] = v
			}
		}
	}

	// Set dependencies in project info
	projectInfo.Dependencies = dependencies

	// Check for specific AI frameworks
	hasLangChain := d.hasPythonDependency(dependencies, "langchain")
	hasLlamaIndex := d.hasPythonDependency(dependencies, "llama-index") || d.hasPythonDependency(dependencies, "llama_index")
	hasOpenAI := d.hasPythonDependency(dependencies, "openai")
	hasHuggingFace := d.hasPythonDependency(dependencies, "transformers") || d.hasPythonDependency(dependencies, "huggingface_hub")
	hasVectorDB := d.hasPythonVectorDBDependency(dependencies)
	hasRAG := false

	// Check for RAG implementation by scanning files
	if (hasLangChain || hasLlamaIndex || hasOpenAI) && hasVectorDB {
		hasRAG = d.hasRAGImplementation(dir)
	}

	// Determine the primary AI framework type
	if hasRAG {
		projectInfo.Type = types.TypeRAG
		projectInfo.LLMProvider = "Multiple"
	} else if hasHuggingFace {
		projectInfo.Type = types.TypeHuggingFace
		projectInfo.LLMProvider = "Hugging Face"
	} else if hasLlamaIndex {
		projectInfo.Type = types.TypeLlamaIndex
		projectInfo.LLMProvider = "LlamaIndex"
	} else if hasLangChain {
		projectInfo.Type = types.TypeLlamaPython
		projectInfo.LLMProvider = "LangChain"
	} else if hasOpenAI {
		projectInfo.Type = types.TypeOpenAIPython
		projectInfo.LLMProvider = "OpenAI"
	} else {
		// Not an AI project
		return false, projectInfo
	}

	// Check for Docker configuration
	if _, err := os.Stat(filepath.Join(dir, "Dockerfile")); err == nil {
		projectInfo.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yml")); err == nil {
		projectInfo.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yaml")); err == nil {
		projectInfo.HasDocker = true
	}

	// Set default port (common for Python web apps)
	projectInfo.Port = 8000

	// Try to identify specific LLM model from code
	if model := d.detectLLMModel(dir); model != "" {
		projectInfo.LLMModel = model
	}

	return true, projectInfo
}

// Helper functions

func (d *AIFrameworkDetector) hasJSDependency(deps map[string]interface{}, name string) bool {
	for dep := range deps {
		if strings.ToLower(dep) == name || strings.HasPrefix(strings.ToLower(dep), name+"-") {
			return true
		}
	}
	return false
}

func (d *AIFrameworkDetector) hasPythonDependency(deps map[string]string, name string) bool {
	for dep := range deps {
		if strings.ToLower(dep) == name || strings.HasPrefix(strings.ToLower(dep), name) {
			return true
		}
	}
	return false
}

func (d *AIFrameworkDetector) hasVectorDBDependency(deps map[string]interface{}) bool {
	vectorDBs := []string{"pinecone", "weaviate", "chromadb", "faiss", "qdrant", "redis", "mongodb", "pgvector"}
	for db := range vectorDBs {
		if d.hasJSDependency(deps, vectorDBs[db]) {
			return true
		}
	}
	return false
}

func (d *AIFrameworkDetector) hasPythonVectorDBDependency(deps map[string]string) bool {
	vectorDBs := []string{"pinecone", "weaviate", "chromadb", "faiss", "qdrant", "redis", "pymongo", "pgvector"}
	for db := range vectorDBs {
		if d.hasPythonDependency(deps, vectorDBs[db]) {
			return true
		}
	}
	return false
}

func (d *AIFrameworkDetector) hasRAGImplementation(dir string) bool {
	// This is a simplified implementation - in practice, you'd want to scan
	// for patterns that suggest RAG (Retrieval Augmented Generation) architecture
	return false
}

func (d *AIFrameworkDetector) detectLLMModel(dir string) string {
	// In a real implementation, this would scan code to find model references
	// like "gpt-4", "llama-2-70b", etc.
	return ""
}

func (d *AIFrameworkDetector) readRequirementsTxt(path string) (map[string]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	deps := make(map[string]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "==")
		if len(parts) == 2 {
			deps[parts[0]] = parts[1]
		} else {
			parts = strings.Split(line, ">=")
			if len(parts) == 2 {
				deps[parts[0]] = ">=" + parts[1]
			} else {
				deps[line] = ""
			}
		}
	}

	return deps, nil
}

func (d *AIFrameworkDetector) readSetupPy(path string) (map[string]string, error) {
	// This is a simplified implementation - parsing setup.py properly would require
	// more sophisticated parsing logic
	return make(map[string]string), nil
}

func (d *AIFrameworkDetector) readPyprojectToml(path string) (map[string]string, error) {
	// This is a simplified implementation - parsing pyproject.toml properly would require
	// a TOML parser
	return make(map[string]string), nil
}

// TODO: Integrate this detector into the detector registry in pkg/detection/detectors.go
// by adding it to the detectors slice in the NewDetectorRegistry function.
