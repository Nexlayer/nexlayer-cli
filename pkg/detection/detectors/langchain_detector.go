// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package detectors provides project-specific detection implementations.

// DEPRECATED: This specialized detector is deprecated and will be removed in a future version.
// The functionality has been replaced by the unified StackDetector in pkg/detection/stack_detector.go,
// which uses pattern-based detection to identify technology components including LangChain.
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

// LangchainDetector detects the presence of LangChain in a project
type LangchainDetector struct {
	detection.BaseDetector
}

// NewLangchainDetector creates a new detector for LangChain
func NewLangchainDetector() *LangchainDetector {
	detector := &LangchainDetector{}
	detector.BaseDetector = *detection.NewBaseDetector("LangChain Detector", 0.8)
	return detector
}

// detectLangchainJS detects LangChain in JavaScript/TypeScript projects
func (d *LangchainDetector) detectLangchainJS(projectPath string) (bool, float64) {
	// Check package.json for LangChain dependencies
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		packageJSONContent, err := os.ReadFile(packageJSONPath)
		if err == nil {
			contentStr := string(packageJSONContent)
			if strings.Contains(contentStr, "\"langchain\"") ||
				strings.Contains(contentStr, "\"@langchain/core\"") ||
				strings.Contains(contentStr, "\"@langchain/openai\"") ||
				strings.Contains(contentStr, "\"langchain-js\"") {
				return true, 0.9
			}
		}
	}

	// Check for import statements in JS/TS files
	langchainImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]langchain`)
	langchainImportCorePattern := regexp.MustCompile(`import\s+.*?from\s+['"]@langchain/core`)
	langchainImportOpenAIPattern := regexp.MustCompile(`import\s+.*?from\s+['"]@langchain/openai`)

	// Check JavaScript/TypeScript files
	jsExtensions := []string{".js", ".jsx", ".ts", ".tsx"}
	for _, ext := range jsExtensions {
		matches, _ := filepath.Glob(filepath.Join(projectPath, "**/*"+ext))
		for _, file := range matches {
			content, err := os.ReadFile(file)
			if err == nil {
				contentStr := string(content)
				if langchainImportPattern.MatchString(contentStr) ||
					langchainImportCorePattern.MatchString(contentStr) ||
					langchainImportOpenAIPattern.MatchString(contentStr) {
					return true, 0.9
				}

				// Check for common LangChain usage patterns
				if strings.Contains(contentStr, "LLMChain") ||
					strings.Contains(contentStr, "PromptTemplate") ||
					strings.Contains(contentStr, "ConversationChain") ||
					strings.Contains(contentStr, "ChatMessageHistory") {
					return true, 0.8
				}
			}
		}
	}

	return false, 0.0
}

// detectLangchainPython detects LangChain in Python projects
func (d *LangchainDetector) detectLangchainPython(projectPath string) (bool, float64) {
	// Check requirements.txt for LangChain
	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		content, err := os.ReadFile(requirementsPath)
		if err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "langchain") ||
				strings.Contains(contentStr, "langchain_core") ||
				strings.Contains(contentStr, "langchain-openai") {
				return true, 0.9
			}
		}
	}

	// Check Poetry pyproject.toml
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	if _, err := os.Stat(pyprojectPath); err == nil {
		content, err := os.ReadFile(pyprojectPath)
		if err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "langchain") ||
				strings.Contains(contentStr, "langchain_core") ||
				strings.Contains(contentStr, "langchain-openai") {
				return true, 0.9
			}
		}
	}

	// Check for import statements in Python files
	importPattern := regexp.MustCompile(`(?m)^import\s+langchain`)
	fromImportPattern := regexp.MustCompile(`(?m)^from\s+langchain\s+import`)
	fromImportCorePattern := regexp.MustCompile(`(?m)^from\s+langchain_core\s+import`)
	fromImportOpenAIPattern := regexp.MustCompile(`(?m)^from\s+langchain_openai\s+import`)

	// Check Python files
	pyMatches, _ := filepath.Glob(filepath.Join(projectPath, "**/*.py"))
	for _, file := range pyMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			contentStr := string(content)
			if importPattern.MatchString(contentStr) ||
				fromImportPattern.MatchString(contentStr) ||
				fromImportCorePattern.MatchString(contentStr) ||
				fromImportOpenAIPattern.MatchString(contentStr) {
				return true, 0.9
			}

			// Check for common LangChain usage patterns
			if strings.Contains(contentStr, "LLMChain") ||
				strings.Contains(contentStr, "PromptTemplate") ||
				strings.Contains(contentStr, "ConversationChain") ||
				strings.Contains(contentStr, "ChatMessageHistory") {
				return true, 0.8
			}
		}
	}

	return false, 0.0
}

// Detect implementation for LangchainDetector
func (d *LangchainDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	// Emit deprecation warning
	detection.EmitDeprecationWarning("LangchainDetector")

	// Create a basic project info
	projectInfo := &detection.ProjectInfo{
		Type:       "unknown",
		Path:       dir,
		Confidence: 0.0,
		Metadata:   make(map[string]interface{}),
	}

	// Try to detect LangChain in JavaScript/TypeScript
	jsFound, jsConf := d.detectLangchainJS(dir)
	if jsFound {
		projectInfo.Type = "langchain"
		projectInfo.Confidence = jsConf
		projectInfo.Language = "JavaScript"
		projectInfo.LLMProvider = "LangChain"
		return projectInfo, nil
	}

	// Try to detect LangChain in Python
	pyFound, pyConf := d.detectLangchainPython(dir)
	if pyFound {
		projectInfo.Type = "langchain"
		projectInfo.Confidence = pyConf
		projectInfo.Language = "Python"
		projectInfo.LLMProvider = "LangChain"
		return projectInfo, nil
	}

	return projectInfo, nil
}
