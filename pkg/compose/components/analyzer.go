// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/golang"
)

// ProjectAnalyzer uses tree-sitter to analyze project files for component detection
type ProjectAnalyzer struct {
	jsParser  *sitter.Parser
	pyParser  *sitter.Parser
	goParser  *sitter.Parser
}

// NewProjectAnalyzer creates a new project analyzer with initialized parsers
func NewProjectAnalyzer() *ProjectAnalyzer {
	return &ProjectAnalyzer{
		jsParser:  newParser(javascript.GetLanguage()),
		pyParser:  newParser(python.GetLanguage()),
		goParser:  newParser(golang.GetLanguage()),
	}
}

func newParser(lang *sitter.Language) *sitter.Parser {
	parser := sitter.NewParser()
	parser.SetLanguage(lang)
	return parser
}

// AnalyzeProject performs deep analysis of a project directory
func (a *ProjectAnalyzer) AnalyzeProject(dir string) (*ProjectInfo, error) {
	// Verify directory exists and is accessible
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory '%s' does not exist. Please verify the path and try again", dir)
	} else if err != nil {
		return nil, fmt.Errorf("cannot access directory '%s'. Please check your permissions. Error: %w", dir, err)
	}

	info := &ProjectInfo{
		Dependencies:    make(map[string][]string),
		EnvVars:        make([]EnvVar, 0),
		ExposedPorts:   make([]int, 0),
		DetectedFrameworks: make([]string, 0),
	}

	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access '%s'. Please ensure you have read permissions. Error: %w", path, err)
		}

		// Skip node_modules, vendor directories
		if f.IsDir() && (f.Name() == "node_modules" || f.Name() == "vendor") {
			return filepath.SkipDir
		}

		switch filepath.Ext(path) {
		case ".js", ".jsx", ".ts", ".tsx":
			if err := a.analyzeJavaScript(path, info); err != nil {
				return fmt.Errorf("analyzing JavaScript: %w", err)
			}
		case ".py":
			if err := a.analyzePython(path, info); err != nil {
				return fmt.Errorf("analyzing Python: %w", err)
			}
		case ".go":
			if err := a.analyzeGo(path, info); err != nil {
				return fmt.Errorf("analyzing Go: %w", err)
			}
		}

		return nil
	})

	return info, err
}

// analyzeJavaScript analyzes JavaScript/TypeScript files
func (a *ProjectAnalyzer) analyzeJavaScript(path string, info *ProjectInfo) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file '%s'. Please verify the file exists and you have read permissions. Error: %w", path, err)
	}

	tree := a.jsParser.Parse(nil, content)
	if tree == nil {
		return fmt.Errorf("failed to parse JavaScript/TypeScript file '%s'. Please ensure the file is valid and not corrupted", path)
	}
	root := tree.RootNode()

	// Find imports and dependencies
	query := `
		(import_statement) @import
		(call_expression
			function: (identifier) @require
			arguments: (arguments (string) @module)
			(#eq? @require "require"))
	`
	q, err := sitter.NewQuery([]byte(query), javascript.GetLanguage())
	if err != nil {
		return fmt.Errorf("internal error: failed to create query for JavaScript analysis. This might indicate a bug in the analyzer. Error: %w", err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, root)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			if c.Node.Type() == "string" {
				moduleName := string(content[c.Node.StartByte():c.Node.EndByte()])
				moduleName = strings.Trim(moduleName, "\"'")
				if !strings.HasPrefix(moduleName, ".") {
					info.Dependencies["npm"] = append(info.Dependencies["npm"], moduleName)
				}
			}
		}
	}

	// Find environment variables
	envQuery := `
		(member_expression
			object: (member_expression
				object: (identifier) @process
				property: (property_identifier) @env)
			property: (property_identifier) @var
			(#eq? @process "process")
			(#eq? @env "env"))
	`
	q, err = sitter.NewQuery([]byte(envQuery), javascript.GetLanguage())
	if err != nil {
		return err
	}

	qc = sitter.NewQueryCursor()
	qc.Exec(q, root)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			if c.Node.Type() == "property_identifier" && c.Node.Parent().Type() == "member_expression" {
				varName := string(content[c.Node.StartByte():c.Node.EndByte()])
				info.EnvVars = append(info.EnvVars, EnvVar{
					Key:   varName,
					Value: fmt.Sprintf("${%s}", varName),
				})
			}
		}
	}

	return nil
}

// analyzePython analyzes Python files
func (a *ProjectAnalyzer) analyzePython(path string, info *ProjectInfo) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file '%s'. Please verify the file exists and you have read permissions. Error: %w", path, err)
	}

	tree := a.pyParser.Parse(nil, content)
	root := tree.RootNode()

	// Find imports
	query := `
		(import_statement) @import
		(import_from_statement) @import_from
	`
	q, err := sitter.NewQuery([]byte(query), python.GetLanguage())
	if err != nil {
		return err
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, root)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			moduleName := string(content[c.Node.StartByte():c.Node.EndByte()])
			if strings.Contains(moduleName, "import") {
				parts := strings.Fields(moduleName)
				for _, part := range parts {
					if part != "import" && part != "from" {
						info.Dependencies["pip"] = append(info.Dependencies["pip"], part)
					}
				}
			}
		}
	}

	// Find environment variables
	envQuery := `
		(call 
			function: (attribute 
				object: (identifier) @os
				attribute: (identifier) @getenv)
			arguments: (argument_list (string) @var)
			(#eq? @os "os")
			(#eq? @getenv "getenv"))
	`
	q, err = sitter.NewQuery([]byte(envQuery), python.GetLanguage())
	if err != nil {
		return err
	}

	qc = sitter.NewQueryCursor()
	qc.Exec(q, root)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			if c.Node.Type() == "string" {
				varName := string(content[c.Node.StartByte():c.Node.EndByte()])
				varName = strings.Trim(varName, "\"'")
				info.EnvVars = append(info.EnvVars, EnvVar{
					Key:   varName,
					Value: fmt.Sprintf("${%s}", varName),
				})
			}
		}
	}

	return nil
}

// analyzeGo analyzes Go files
func (a *ProjectAnalyzer) analyzeGo(path string, info *ProjectInfo) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file '%s'. Please verify the file exists and you have read permissions. Error: %w", path, err)
	}

	tree := a.goParser.Parse(nil, content)
	root := tree.RootNode()

	// Find imports
	query := `(import_declaration) @import`
	q, err := sitter.NewQuery([]byte(query), golang.GetLanguage())
	if err != nil {
		return err
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, root)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			importSpec := string(content[c.Node.StartByte():c.Node.EndByte()])
			if strings.Contains(importSpec, "\"") {
				parts := strings.Split(importSpec, "\"")
				for _, part := range parts {
					if part != "" && !strings.Contains(part, "import") {
						info.Dependencies["go"] = append(info.Dependencies["go"], strings.TrimSpace(part))
					}
				}
			}
		}
	}

	// Find environment variables
	envQuery := `
		(call_expression
			function: (selector_expression
				operand: (identifier) @os
				field: (field_identifier) @getenv)
			arguments: (argument_list (interpreted_string_literal) @var)
			(#eq? @os "os")
			(#eq? @getenv "Getenv"))
	`
	q, err = sitter.NewQuery([]byte(envQuery), golang.GetLanguage())
	if err != nil {
		return err
	}

	qc = sitter.NewQueryCursor()
	qc.Exec(q, root)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			if c.Node.Type() == "interpreted_string_literal" {
				varName := string(content[c.Node.StartByte():c.Node.EndByte()])
				varName = strings.Trim(varName, "\"")
				info.EnvVars = append(info.EnvVars, EnvVar{
					Key:   varName,
					Value: fmt.Sprintf("${%s}", varName),
				})
			}
		}
	}

	return nil
}

// ProjectInfo holds analyzed information about a project
type ProjectInfo struct {
	Dependencies      map[string][]string // Map of package manager to dependencies
	EnvVars          []EnvVar            // Environment variables used in the code
	ExposedPorts     []int               // Ports that the application listens on
	DetectedFrameworks []string           // Detected frameworks (e.g., "react", "express", "fastapi")
}
