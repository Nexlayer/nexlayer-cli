// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detection

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
)

// TechStack represents a complete technology stack
type TechStack struct {
	Name          string
	Frontend      string
	Backend       string
	Database      string
	Deployment    string
	AIIntegration []string
	Confidence    float64
}

// PatternType defines the type of pattern to match
type PatternType string

const (
	PatternDependency  PatternType = "dependency"
	PatternFile        PatternType = "file"
	PatternImport      PatternType = "import"
	PatternContent     PatternType = "content"
	PatternEnvironment PatternType = "environment"
)

// DetectionPattern defines a pattern to match in a project
type DetectionPattern struct {
	Type       PatternType
	Pattern    string // Regex pattern or exact match
	Path       string // Where to look (package.json, requirements.txt, etc.)
	Confidence float64
}

// Components represents the components of a technology stack
type Components struct {
	Frontend   []string
	Backend    []string
	Database   []string
	AI         []string
	Deployment []string
}

// StackDefinition defines a technology stack and its detection patterns
type StackDefinition struct {
	Name               string
	Description        string
	Components         Components
	RequiredComponents []string
	OptionalComponents []string
	MainPatterns       []DetectionPattern
	ExtraPatterns      []DetectionPattern
}

// StackDetector is a unified detector for common technology stacks
type StackDetector struct {
	BaseDetector
	definitions map[string]StackDefinition
	fileCache   sync.Map // Cache for file existence and content
}

// fileReadCache stores the content of a file
type fileReadCache struct {
	exists  bool
	content string
}

// NewBaseDetector creates a new base detector with the given name and confidence
func NewBaseDetector(name string, confidence float64) *BaseDetector {
	return &BaseDetector{
		name:       name,
		confidence: confidence,
	}
}

// NewStackDetector creates a new detector for common technology stacks
func NewStackDetector() *StackDetector {
	detector := &StackDetector{
		definitions: TechStackDefinitions,
		fileCache:   sync.Map{},
	}
	detector.BaseDetector = BaseDetector{
		name:       "Stack Detector",
		confidence: 0.9,
	}
	return detector
}

// Priority returns the priority of this detector
func (d *StackDetector) Priority() int {
	return 150 // High priority, runs after LLM detector but before other detectors
}

// Detect analyzes a project to determine its technology stack
func (d *StackDetector) Detect(dir string) (*types.ProjectInfo, error) {
	// Create internal projectInfo for our detection logic
	internalInfo := &ProjectInfo{
		Type:         "unknown",
		Path:         dir,
		Confidence:   0.0,
		Dependencies: make(map[string]string),
		Metadata:     make(map[string]interface{}),
	}

	// Store detection metadata
	stacksMetadata := make(map[string]interface{})
	for stackID, stackDef := range d.definitions {
		stackMeta := map[string]interface{}{
			"name":        stackDef.Name,
			"description": stackDef.Description,
			"components":  stackDef.Components,
		}
		stacksMetadata[stackID] = stackMeta
	}
	internalInfo.Metadata["available_stacks"] = stacksMetadata

	// Check each stack definition in parallel
	var wg sync.WaitGroup
	resultChan := make(chan struct {
		stackID    string
		confidence float64
		components map[string]interface{}
	}, len(d.definitions))

	for stackID, stackDef := range d.definitions {
		wg.Add(1)
		go func(id string, def StackDefinition) {
			defer wg.Done()
			confidence, components := d.evaluateStack(dir, def)
			if confidence > 0.5 { // Only report stacks with decent confidence
				resultChan <- struct {
					stackID    string
					confidence float64
					components map[string]interface{}
				}{
					stackID:    id,
					confidence: confidence,
					components: components,
				}
			}
		}(stackID, stackDef)
	}

	// Wait for all goroutines to complete, then close the channel
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Process results
	bestConfidence := 0.0
	bestStackID := ""
	var bestComponents map[string]interface{}

	for result := range resultChan {
		if result.confidence > bestConfidence {
			bestConfidence = result.confidence
			bestStackID = result.stackID
			bestComponents = result.components
		}
	}

	// Create the external types.ProjectInfo that we'll return
	externalInfo := &types.ProjectInfo{
		Type:         types.ProjectType(internalInfo.Type),
		Dependencies: make(map[string]string),
	}

	// If a stack was detected with good confidence
	if bestConfidence > 0.5 {
		stackDef := d.definitions[bestStackID]
		externalInfo.Type = types.ProjectType(bestStackID)
		internalInfo.Metadata["confidence"] = bestConfidence
		internalInfo.Metadata["framework"] = stackDef.Name
		internalInfo.Metadata["stack_components"] = bestComponents

		// Set dependencies based on components
		for _, comp := range stackDef.Components.Frontend {
			externalInfo.Dependencies[comp] = "*"
		}
		for _, comp := range stackDef.Components.Backend {
			externalInfo.Dependencies[comp] = "*"
		}
		for _, comp := range stackDef.Components.Database {
			externalInfo.Dependencies[comp] = "*"
		}
		for _, comp := range stackDef.Components.AI {
			externalInfo.Dependencies[comp] = "*"
		}
		for _, comp := range stackDef.Components.Deployment {
			externalInfo.Dependencies[comp] = "*"
		}

		// Set language
		if contains(stackDef.Components.Frontend, "nextjs") || contains(stackDef.Components.Frontend, "react") {
			// Check for TypeScript
			if d.hasTypeScriptFiles(dir) {
				internalInfo.Language = "typescript"
				// Set LLM provider if AI components are present
				if len(stackDef.Components.AI) > 0 {
					for _, ai := range stackDef.Components.AI {
						if ai == "openai" {
							externalInfo.LLMProvider = "openai"
							break
						} else if ai == "langchain" {
							externalInfo.LLMProvider = "langchain"
							break
						} else if ai == "gemini" {
							externalInfo.LLMProvider = "google"
							break
						}
					}
				}
			} else {
				internalInfo.Language = "javascript"
			}
		} else if contains(stackDef.Components.Backend, "django") || contains(stackDef.Components.Backend, "flask") {
			internalInfo.Language = "python"
		} else if contains(stackDef.Components.Backend, "express") || contains(stackDef.Components.Backend, "node") {
			internalInfo.Language = "javascript"
		}
	}

	return externalInfo, nil
}

// evaluateStack checks if a project matches a given stack definition
func (d *StackDetector) evaluateStack(dir string, def StackDefinition) (float64, map[string]interface{}) {
	totalConfidence := 0.0
	maxConfidence := 0.0
	detectedComponents := make(map[string]interface{})

	// Check main patterns (must-haves)
	for _, pattern := range def.MainPatterns {
		maxConfidence += pattern.Confidence
		if d.matchesPattern(dir, pattern) {
			totalConfidence += pattern.Confidence
		}
	}

	// Check extra patterns (nice-to-haves)
	for _, pattern := range def.ExtraPatterns {
		maxConfidence += pattern.Confidence
		if d.matchesPattern(dir, pattern) {
			totalConfidence += pattern.Confidence
		}
	}

	// Normalize confidence score
	normalizedConfidence := 0.0
	if maxConfidence > 0 {
		normalizedConfidence = totalConfidence / maxConfidence
	}

	// Check required components
	requiredCount := 0
	for _, comp := range def.RequiredComponents {
		if d.hasComponent(dir, comp) {
			requiredCount++
			detectedComponents[comp] = true
		} else {
			detectedComponents[comp] = false
		}
	}

	// If not all required components are present, reduce confidence
	requiredRatio := 1.0
	if len(def.RequiredComponents) > 0 {
		requiredRatio = float64(requiredCount) / float64(len(def.RequiredComponents))
	}

	// Check optional components
	for _, comp := range def.OptionalComponents {
		if d.hasComponent(dir, comp) {
			detectedComponents[comp] = true
			// Bonus for optional components
			normalizedConfidence += 0.05
		} else {
			detectedComponents[comp] = false
		}
	}

	// Final confidence is based on pattern matches and required components
	finalConfidence := normalizedConfidence * requiredRatio

	// Cap at 0.95 maximum confidence
	if finalConfidence > 0.95 {
		finalConfidence = 0.95
	}

	return finalConfidence, detectedComponents
}

// matchesPattern checks if a project matches a detection pattern
func (d *StackDetector) matchesPattern(dir string, pattern DetectionPattern) bool {
	switch pattern.Type {
	case PatternDependency:
		return d.hasDependency(dir, pattern.Pattern, pattern.Path)
	case PatternFile:
		return d.hasFile(dir, pattern.Pattern)
	case PatternImport:
		return d.hasImport(dir, pattern.Pattern, pattern.Path)
	case PatternContent:
		return d.hasContent(dir, pattern.Pattern, pattern.Path)
	case PatternEnvironment:
		return d.hasEnvironmentVar(dir, pattern.Pattern)
	default:
		return false
	}
}

// hasDependency checks if a dependency exists in a package file
func (d *StackDetector) hasDependency(dir, dependency, packageFile string) bool {
	filePath := filepath.Join(dir, packageFile)
	content, exists := d.readFileCache(filePath)
	if !exists {
		return false
	}

	// Simple check for dependencies in package.json or requirements.txt
	return strings.Contains(content, fmt.Sprintf("\"%s\"", dependency)) ||
		strings.Contains(content, fmt.Sprintf("'%s'", dependency)) ||
		strings.Contains(content, fmt.Sprintf("%s==", dependency)) ||
		strings.Contains(content, fmt.Sprintf("%s>=", dependency))
}

// hasFile checks if a file matching the pattern exists
func (d *StackDetector) hasFile(dir, pattern string) bool {
	// Check if it's a direct file path or a pattern
	if !strings.Contains(pattern, "*") && !strings.Contains(pattern, "(") && !strings.Contains(pattern, "[") {
		filePath := filepath.Join(dir, pattern)
		_, exists := d.fileExistsCache(filePath)
		return exists
	}

	// It's a pattern, use filepath.Glob
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return false
	}
	return len(matches) > 0
}

// hasImport checks if a file has a specific import statement
func (d *StackDetector) hasImport(dir, importPattern, filePath string) bool {
	re, err := regexp.Compile(importPattern)
	if err != nil {
		return false
	}

	// If filePath is a pattern, check all matching files
	if strings.Contains(filePath, "*") || strings.Contains(filePath, "{") {
		matches, err := filepath.Glob(filepath.Join(dir, filePath))
		if err != nil {
			return false
		}
		for _, match := range matches {
			content, exists := d.readFileCache(match)
			if exists && re.MatchString(content) {
				return true
			}
		}
		return false
	}

	// Single file check
	fullPath := filepath.Join(dir, filePath)
	content, exists := d.readFileCache(fullPath)
	if !exists {
		return false
	}
	return re.MatchString(content)
}

// hasContent checks if any file matching the path pattern contains the content pattern
func (d *StackDetector) hasContent(dir, contentPattern, pathPattern string) bool {
	re, err := regexp.Compile(contentPattern)
	if err != nil {
		return false
	}

	// Handle glob patterns
	matches, err := filepath.Glob(filepath.Join(dir, pathPattern))
	if err != nil {
		return false
	}

	for _, match := range matches {
		content, exists := d.readFileCache(match)
		if exists && re.MatchString(content) {
			return true
		}
	}
	return false
}

// hasEnvironmentVar checks if an environment variable is defined in .env files
func (d *StackDetector) hasEnvironmentVar(dir, varName string) bool {
	envFiles := []string{".env", ".env.local", ".env.development", ".env.production"}
	for _, envFile := range envFiles {
		filePath := filepath.Join(dir, envFile)
		content, exists := d.readFileCache(filePath)
		if exists && strings.Contains(content, varName+"=") {
			return true
		}
	}
	return false
}

// hasComponent checks if a specific technology component is detected
func (d *StackDetector) hasComponent(dir, component string) bool {
	switch component {
	case "nextjs":
		return d.hasFile(dir, "next.config.js") || d.hasFile(dir, "next.config.mjs") ||
			d.hasDependency(dir, "next", "package.json")
	case "react":
		return d.hasDependency(dir, "react", "package.json") &&
			d.hasDependency(dir, "react-dom", "package.json")
	case "vue":
		return d.hasDependency(dir, "vue", "package.json") ||
			d.hasFile(dir, "vue.config.js")
	case "supabase":
		return d.hasDependency(dir, "@supabase/supabase-js", "package.json") ||
			d.hasEnvironmentVar(dir, "SUPABASE_URL")
	case "langchain":
		return d.hasDependency(dir, "langchain", "package.json") ||
			d.hasDependency(dir, "langchain", "requirements.txt")
	case "openai":
		return d.hasDependency(dir, "openai", "package.json") ||
			d.hasDependency(dir, "openai", "requirements.txt") ||
			d.hasEnvironmentVar(dir, "OPENAI_API_KEY")
	case "gemini":
		return d.hasDependency(dir, "@google/generative-ai", "package.json") ||
			d.hasDependency(dir, "google-generativeai", "requirements.txt") ||
			d.hasEnvironmentVar(dir, "GEMINI_API_KEY")
	case "postgres":
		return d.hasDependency(dir, "pg", "package.json") ||
			d.hasDependency(dir, "psycopg2", "requirements.txt") ||
			d.hasEnvironmentVar(dir, "DATABASE_URL")
	case "pgvector":
		return d.hasDependency(dir, "pgvector", "requirements.txt") ||
			d.hasContent(dir, "CREATE EXTENSION vector", "**/*.sql") ||
			d.hasContent(dir, "pgvector", "**/*.{js,ts,py}")
	case "tailwind":
		return d.hasDependency(dir, "tailwindcss", "package.json") ||
			d.hasFile(dir, "tailwind.config.js") ||
			d.hasFile(dir, "tailwind.config.ts")
	case "stripe":
		return d.hasDependency(dir, "stripe", "package.json") ||
			d.hasDependency(dir, "stripe", "requirements.txt") ||
			d.hasEnvironmentVar(dir, "STRIPE_SECRET_KEY")
	case "django":
		return d.hasFile(dir, "manage.py") &&
			d.hasContent(dir, "django", "requirements.txt")
	case "express":
		return d.hasDependency(dir, "express", "package.json")
	case "mongodb":
		return d.hasDependency(dir, "mongodb", "package.json") ||
			d.hasDependency(dir, "mongoose", "package.json") ||
			d.hasEnvironmentVar(dir, "MONGO_URI")
	default:
		return false
	}
}

// hasTypeScriptFiles checks if the project contains TypeScript files
func (d *StackDetector) hasTypeScriptFiles(dir string) bool {
	tsConfigPath := filepath.Join(dir, "tsconfig.json")
	_, tsConfigExists := d.fileExistsCache(tsConfigPath)
	if tsConfigExists {
		return true
	}

	// Look for .ts or .tsx files
	matches, err := filepath.Glob(filepath.Join(dir, "**/*.ts"))
	if err == nil && len(matches) > 0 {
		return true
	}
	matches, err = filepath.Glob(filepath.Join(dir, "**/*.tsx"))
	if err == nil && len(matches) > 0 {
		return true
	}

	return false
}

// fileExistsCache checks if a file exists with caching
func (d *StackDetector) fileExistsCache(path string) (os.FileInfo, bool) {
	if cache, ok := d.fileCache.Load(path + "_exists"); ok {
		if fileCache, ok := cache.(fileReadCache); ok {
			if fileCache.exists {
				// We don't have the FileInfo in the cache, so we need to get it again
				info, err := os.Stat(path)
				if err != nil {
					return nil, false
				}
				return info, true
			}
			return nil, false
		}
	}

	info, err := os.Stat(path)
	exists := err == nil
	d.fileCache.Store(path+"_exists", fileReadCache{exists: exists})
	return info, exists
}

// readFileCache reads a file with caching
func (d *StackDetector) readFileCache(path string) (string, bool) {
	if cache, ok := d.fileCache.Load(path); ok {
		if fileCache, ok := cache.(fileReadCache); ok {
			return fileCache.content, fileCache.exists
		}
	}

	content, err := os.ReadFile(path)
	exists := err == nil
	d.fileCache.Store(path, fileReadCache{
		exists:  exists,
		content: string(content),
	})

	return string(content), exists
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
