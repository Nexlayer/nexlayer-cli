// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package detectors provides project-specific detection implementations.

// DEPRECATED: This specialized detector is deprecated and will be removed in a future version.
// The functionality has been replaced by the unified StackDetector in pkg/detection/stack_detector.go,
// which uses pattern-based detection to identify technology components including Gemini AI integrations.
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

// GeminiDetector detects the presence of Google Gemini/Vertex AI in a project
type GeminiDetector struct {
	detection.BaseDetector
}

// NewGeminiDetector creates a new detector for Gemini/Vertex AI integration
func NewGeminiDetector() *GeminiDetector {
	detector := &GeminiDetector{}
	detector.BaseDetector = *detection.NewBaseDetector("Gemini/Vertex AI Detector", 0.8)
	return detector
}

// detectGeminiJS detects Google AI SDK for JavaScript
func (d *GeminiDetector) detectGeminiJS(projectPath string) (bool, float64) {
	// Check package.json for @google/generative-ai
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		packageJSONContent, err := os.ReadFile(packageJSONPath)
		if err == nil {
			if strings.Contains(string(packageJSONContent), "\"@google/generative-ai\"") {
				return true, 0.9
			}
		}
	}

	// Check for import statements in JS/TS files
	jsImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]@google/generative-ai['"]`)
	tsImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]@google/generative-ai['"]`)

	// Check JavaScript files
	jsExtensions := []string{".js", ".jsx", ".ts", ".tsx"}
	for _, ext := range jsExtensions {
		matches, _ := filepath.Glob(filepath.Join(projectPath, "**/*"+ext))
		for _, file := range matches {
			content, err := os.ReadFile(file)
			if err == nil {
				if jsImportPattern.MatchString(string(content)) || tsImportPattern.MatchString(string(content)) {
					return true, 0.9
				}
				// Check for GoogleGenerativeAI or Gemini API usage
				if strings.Contains(string(content), "GoogleGenerativeAI") ||
					strings.Contains(string(content), "genAI.") ||
					strings.Contains(string(content), "generationConfig") ||
					strings.Contains(string(content), "gemini-pro") {
					return true, 0.8
				}
			}
		}
	}

	return false, 0.0
}

// detectGeminiPython detects Google AI SDK for Python
func (d *GeminiDetector) detectGeminiPython(projectPath string) (bool, float64) {
	// Check requirements.txt for google-generativeai
	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		content, err := os.ReadFile(requirementsPath)
		if err == nil {
			if strings.Contains(string(content), "google-generativeai") ||
				strings.Contains(string(content), "vertexai") {
				return true, 0.9
			}
		}
	}

	// Check Poetry pyproject.toml
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	if _, err := os.Stat(pyprojectPath); err == nil {
		content, err := os.ReadFile(pyprojectPath)
		if err == nil {
			if strings.Contains(string(content), "google-generativeai") ||
				strings.Contains(string(content), "vertexai") {
				return true, 0.9
			}
		}
	}

	// Check for import statements in Python files
	importPattern1 := regexp.MustCompile(`(?m)^import\s+google\.generativeai`)
	importPattern2 := regexp.MustCompile(`(?m)^from\s+google\s+import\s+generativeai`)
	importPattern3 := regexp.MustCompile(`(?m)^import\s+vertexai`)
	importPattern4 := regexp.MustCompile(`(?m)^from\s+vertexai\s+import`)

	// Check Python files
	pyMatches, _ := filepath.Glob(filepath.Join(projectPath, "**/*.py"))
	for _, file := range pyMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			if importPattern1.MatchString(string(content)) ||
				importPattern2.MatchString(string(content)) ||
				importPattern3.MatchString(string(content)) ||
				importPattern4.MatchString(string(content)) {
				return true, 0.9
			}
			// Check for Gemini API usage
			if strings.Contains(string(content), "genai.") ||
				strings.Contains(string(content), "vertexai.") ||
				strings.Contains(string(content), "GenerativeModel") ||
				strings.Contains(string(content), "gemini-pro") {
				return true, 0.8
			}
		}
	}

	return false, 0.0
}

// detectGeminiGo detects Google AI SDK for Go
func (d *GeminiDetector) detectGeminiGo(projectPath string) (bool, float64) {
	// Check go.mod for google.golang.org/genai
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		content, err := os.ReadFile(goModPath)
		if err == nil {
			if strings.Contains(string(content), "google.golang.org/genai") ||
				strings.Contains(string(content), "cloud.google.com/go/vertexai") {
				return true, 0.9
			}
		}
	}

	// Check for import statements in Go files
	importPattern1 := regexp.MustCompile(`(?m)import\s+\(.*?["']google\.golang\.org/genai["'].*?\)`)
	importPattern2 := regexp.MustCompile(`(?m)import\s+["']google\.golang\.org/genai["']`)
	importPattern3 := regexp.MustCompile(`(?m)import\s+\(.*?["']cloud\.google\.com/go/vertexai["'].*?\)`)
	importPattern4 := regexp.MustCompile(`(?m)import\s+["']cloud\.google\.com/go/vertexai["']`)

	// Check Go files
	goMatches, _ := filepath.Glob(filepath.Join(projectPath, "**/*.go"))
	for _, file := range goMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			if importPattern1.MatchString(string(content)) ||
				importPattern2.MatchString(string(content)) ||
				importPattern3.MatchString(string(content)) ||
				importPattern4.MatchString(string(content)) {
				return true, 0.9
			}
			// Check for Gemini API usage
			if strings.Contains(string(content), "genai.") ||
				strings.Contains(string(content), "vertexai.") ||
				strings.Contains(string(content), "GenerativeModel") {
				return true, 0.8
			}
		}
	}

	return false, 0.0
}

// Detect implementation for GeminiDetector
func (d *GeminiDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	// Emit deprecation warning
	detection.EmitDeprecationWarning("GeminiDetector")

	// Create a basic project info
	projectInfo := &detection.ProjectInfo{
		Type:       "unknown",
		Path:       dir,
		Confidence: 0.0,
		Metadata:   make(map[string]interface{}),
	}

	// Try to detect Gemini in JavaScript/TypeScript
	jsFound, jsConf := d.detectGeminiJS(dir)
	if jsFound {
		projectInfo.Type = "gemini-node"
		projectInfo.Confidence = jsConf
		projectInfo.Language = "JavaScript"
		projectInfo.LLMProvider = "Google"
		return projectInfo, nil
	}

	// Try to detect Gemini in Python
	pyFound, pyConf := d.detectGeminiPython(dir)
	if pyFound {
		projectInfo.Type = "gemini-python"
		projectInfo.Confidence = pyConf
		projectInfo.Language = "Python"
		projectInfo.LLMProvider = "Google"
		return projectInfo, nil
	}

	// Try to detect Gemini in Go
	goFound, goConf := d.detectGeminiGo(dir)
	if goFound {
		projectInfo.Type = "gemini-go"
		projectInfo.Confidence = goConf
		projectInfo.Language = "Go"
		projectInfo.LLMProvider = "Google"
		return projectInfo, nil
	}

	return projectInfo, nil
}
