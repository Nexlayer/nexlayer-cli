// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package analysis

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
)

// Language represents a supported programming language
type Language string

const (
	Go         Language = "go"
	JavaScript Language = "javascript"
	Python     Language = "python"
)

// Parser manages tree-sitter parsers for different languages
type Parser struct {
	mu       sync.RWMutex
	parsers  map[Language]*sitter.Parser
	queries  map[Language]map[string]*sitter.Query
	language map[Language]*sitter.Language
}

// NewParser creates a new tree-sitter parser manager
func NewParser() *Parser {
	p := &Parser{
		parsers:  make(map[Language]*sitter.Parser),
		queries:  make(map[Language]map[string]*sitter.Query),
		language: make(map[Language]*sitter.Language),
	}

	// Initialize supported languages
	p.language[Go] = golang.GetLanguage()
	p.language[JavaScript] = javascript.GetLanguage()
	p.language[Python] = python.GetLanguage()

	return p
}

// ProjectAnalysis contains the analysis results for a project
type ProjectAnalysis struct {
	Dependencies  map[string]string         // Direct dependencies and their versions
	Imports       map[string][]string       // Package imports by file
	Functions     map[string][]FunctionInfo // Functions by file
	Frameworks    []string                  // Detected frameworks
	DatabaseTypes []string                  // Detected database types
	APIEndpoints  []APIEndpoint             // Detected API endpoints
}

// FunctionInfo contains information about a function
type FunctionInfo struct {
	Name       string
	Signature  string
	StartLine  uint32
	EndLine    uint32
	IsExported bool
}

// APIEndpoint represents a detected API endpoint
type APIEndpoint struct {
	Path       string
	Method     string
	Handler    string
	Parameters []string
}

// AnalyzeProject performs a deep analysis of the project directory
func (p *Parser) AnalyzeProject(ctx context.Context, projectDir string) (*ProjectAnalysis, error) {
	analysis := &ProjectAnalysis{
		Dependencies:  make(map[string]string),
		Imports:       make(map[string][]string),
		Functions:     make(map[string][]FunctionInfo),
		Frameworks:    make([]string, 0),
		DatabaseTypes: make([]string, 0),
		APIEndpoints:  make([]APIEndpoint, 0),
	}

	// Walk through project files
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-source files
		if info.IsDir() || !isSupportedFile(path) {
			return nil
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Analyze file
		if err := p.analyzeFile(path, analysis); err != nil {
			return fmt.Errorf("failed to analyze %s: %w", path, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("project analysis failed: %w", err)
	}

	return analysis, nil
}

// analyzeFile analyzes a single source file
func (p *Parser) analyzeFile(path string, analysis *ProjectAnalysis) error {
	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Determine language
	lang := detectLanguage(path)
	if lang == "" {
		return nil // Unsupported language
	}

	// Get or create parser for the language
	parser, err := p.GetParser(lang)
	if err != nil {
		return err
	}

	// Parse file
	tree := parser.Parse(nil, content)
	if tree == nil {
		return fmt.Errorf("failed to parse %s", path)
	}
	defer tree.Close()

	// Call language-specific analysis
	switch lang {
	case Go:
		return p.analyzeGoFile(path, tree, content, analysis)
	case JavaScript:
		// TODO: Implement JavaScript analysis in javascript.go
		return nil
	case Python:
		// TODO: Implement Python analysis in python.go
		return nil
	}

	return nil
}

// GetParser returns a parser for the given language
func (p *Parser) GetParser(lang Language) (*sitter.Parser, error) {
	p.mu.RLock()
	parser, ok := p.parsers[lang]
	p.mu.RUnlock()

	if ok {
		return parser, nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Check again in case another goroutine created it
	if parser, ok = p.parsers[lang]; ok {
		return parser, nil
	}

	// Create new parser
	parser = sitter.NewParser()
	language, ok := p.language[lang]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	parser.SetLanguage(language)
	p.parsers[lang] = parser

	return parser, nil
}

// Helper functions

func isSupportedFile(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".go", ".js", ".jsx", ".ts", ".tsx", ".py":
		return true
	default:
		return false
	}
}

func detectLanguage(path string) Language {
	ext := filepath.Ext(path)
	switch ext {
	case ".go":
		return Go
	case ".js", ".jsx", ".ts", ".tsx":
		return JavaScript
	case ".py":
		return Python
	default:
		return ""
	}
}

// Language-specific analysis functions are implemented in their respective files:
// - go.go for Go analysis
// - javascript.go for JavaScript analysis
// - python.go for Python analysis
