// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package detectors provides project-specific detection implementations.

// DEPRECATED: This specialized detector is deprecated and will be removed in a future version.
// The functionality has been replaced by the unified StackDetector in pkg/detection/stack_detector.go,
// which uses pattern-based detection to identify technology components including OpenAI integrations.
// Please migrate any direct references to this detector to use StackDetector instead.
// Planned removal: v1.x.0 (next major/minor release)

package detectors

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// OpenAIDetector detects the presence of OpenAI/ChatGPT API in a project
type OpenAIDetector struct {
	detection.BaseDetector
}

// NewOpenAIDetector creates a new detector for OpenAI integration
func NewOpenAIDetector() *OpenAIDetector {
	detector := &OpenAIDetector{}
	detector.BaseDetector = *detection.NewBaseDetector("OpenAI Detector", 0.8)
	return detector
}

// detectOpenAIJS detects OpenAI API usage in JavaScript/TypeScript projects
func (d *OpenAIDetector) detectOpenAIJS(projectPath string) (bool, float64) {
	// Check package.json for OpenAI dependencies
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		packageJSONContent, err := os.ReadFile(packageJSONPath)
		if err == nil {
			if strings.Contains(string(packageJSONContent), "\"openai\"") ||
				strings.Contains(string(packageJSONContent), "\"@langchain/openai\"") {
				return true, 0.9
			}
		}
	}

	// Check for import statements in JS/TS files
	openaiImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]openai['"]`)
	langchainOpenaiImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]@langchain/openai['"]`)

	// Check for OpenAI API key patterns
	apiKeyPattern := regexp.MustCompile(`sk-[a-zA-Z0-9]{32,}`)

	// Check JavaScript files
	jsExtensions := []string{".js", ".jsx", ".ts", ".tsx"}
	for _, ext := range jsExtensions {
		matches, _ := filepath.Glob(filepath.Join(projectPath, "**/*"+ext))
		for _, file := range matches {
			content, err := os.ReadFile(file)
			if err == nil {
				contentStr := string(content)
				if openaiImportPattern.MatchString(contentStr) || langchainOpenaiImportPattern.MatchString(contentStr) {
					return true, 0.9
				}
				// Check for API keys
				if apiKeyPattern.MatchString(contentStr) {
					return true, 0.8
				}
				// Check for common OpenAI API usage
				if strings.Contains(contentStr, "OpenAI") ||
					strings.Contains(contentStr, "ChatGPT") ||
					(strings.Contains(contentStr, "gpt-") &&
						(strings.Contains(contentStr, "gpt-3.5") ||
							strings.Contains(contentStr, "gpt-4"))) {
					return true, 0.7
				}
			}
		}
	}

	// Check for environment variables in .env or similar files
	envMatches, _ := filepath.Glob(filepath.Join(projectPath, ".env*"))
	for _, file := range envMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "OPENAI_API_KEY") ||
				strings.Contains(contentStr, "OPENAI_KEY") ||
				apiKeyPattern.MatchString(contentStr) {
				return true, 0.9
			}
		}
	}

	return false, 0.0
}

// detectOpenAIPython detects OpenAI API usage in Python projects
func (d *OpenAIDetector) detectOpenAIPython(projectPath string) (bool, float64) {
	// Check requirements.txt for OpenAI
	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		content, err := os.ReadFile(requirementsPath)
		if err == nil {
			if strings.Contains(string(content), "openai") ||
				strings.Contains(string(content), "langchain-openai") {
				return true, 0.9
			}
		}
	}

	// Check Poetry pyproject.toml
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	if _, err := os.Stat(pyprojectPath); err == nil {
		content, err := os.ReadFile(pyprojectPath)
		if err == nil {
			if strings.Contains(string(content), "openai") ||
				strings.Contains(string(content), "langchain-openai") {
				return true, 0.9
			}
		}
	}

	// Check for import statements in Python files
	importPattern := regexp.MustCompile(`(?m)^import\s+openai`)
	fromImportPattern := regexp.MustCompile(`(?m)^from\s+openai\s+import`)
	langchainOpenaiPattern := regexp.MustCompile(`(?m)^from\s+langchain_openai\s+import`)

	// Check for OpenAI API key patterns
	apiKeyPattern := regexp.MustCompile(`sk-[a-zA-Z0-9]{32,}`)

	// Check Python files
	pyMatches, _ := filepath.Glob(filepath.Join(projectPath, "**/*.py"))
	for _, file := range pyMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			contentStr := string(content)
			if importPattern.MatchString(contentStr) ||
				fromImportPattern.MatchString(contentStr) ||
				langchainOpenaiPattern.MatchString(contentStr) {
				return true, 0.9
			}
			// Check for API keys
			if apiKeyPattern.MatchString(contentStr) {
				return true, 0.8
			}
			// Check for common OpenAI API usage
			if strings.Contains(contentStr, "ChatCompletion") ||
				strings.Contains(contentStr, "ChatGPT") ||
				(strings.Contains(contentStr, "gpt-") &&
					(strings.Contains(contentStr, "gpt-3.5") ||
						strings.Contains(contentStr, "gpt-4"))) {
				return true, 0.7
			}
		}
	}

	// Check for environment variables
	envMatches, _ := filepath.Glob(filepath.Join(projectPath, ".env*"))
	for _, file := range envMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "OPENAI_API_KEY") ||
				strings.Contains(contentStr, "OPENAI_KEY") ||
				apiKeyPattern.MatchString(contentStr) {
				return true, 0.9
			}
		}
	}

	return false, 0.0
}

// Detect implementation for OpenAIDetector
func (d *OpenAIDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	// Emit deprecation warning
	detection.EmitDeprecationWarning("OpenAIDetector")

	// Create a basic project info
	projectInfo := &detection.ProjectInfo{
		Type:       "unknown",
		Path:       dir,
		Confidence: 0.0,
		Metadata:   make(map[string]interface{}),
	}

	// Try to detect OpenAI in JavaScript/TypeScript
	jsFound, jsConf := d.detectOpenAIJS(dir)
	if jsFound {
		projectInfo.Type = "openai"
		projectInfo.Confidence = jsConf
		projectInfo.Language = "JavaScript"
		projectInfo.LLMProvider = "OpenAI"
		return projectInfo, nil
	}

	// Try to detect OpenAI in Python
	pyFound, pyConf := d.detectOpenAIPython(dir)
	if pyFound {
		projectInfo.Type = "openai"
		projectInfo.Confidence = pyConf
		projectInfo.Language = "Python"
		projectInfo.LLMProvider = "OpenAI"
		return projectInfo, nil
	}

	return projectInfo, nil
}
