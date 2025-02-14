package analysis

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

// Common tree-sitter queries for Go code
const (
	goImportQuery = `
(source_file
  (import_declaration
    (import_spec_list
      (import_spec
        path: (interpreted_string_literal) @import.path))))
`

	goFunctionQuery = `
(source_file
  (function_declaration
    name: (identifier) @function.name))
`

	goAPIEndpointQuery = `
(source_file
  (call_expression
    function: (selector_expression
      operand: (identifier) @router
      field: (field_identifier) @method) @api.method
    arguments: (argument_list) @api.args))
`
)

func (p *Parser) analyzeGoFile(path string, tree *sitter.Tree, content []byte, analysis *ProjectAnalysis) error {
	// Initialize queries if needed
	if err := p.initGoQueries(); err != nil {
		return err
	}

	// Extract imports
	imports, err := p.extractGoImports(tree, content)
	if err != nil {
		return err
	}
	analysis.Imports[path] = imports

	// Extract functions
	functions, err := p.extractGoFunctions(tree, content)
	if err != nil {
		return err
	}
	analysis.Functions[path] = functions

	return nil
}

func (p *Parser) initGoQueries() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.queries[Go]; ok {
		return nil
	}

	// Initialize queries map for Go if not exists
	p.queries[Go] = make(map[string]*sitter.Query)

	// Initialize import query
	importQuery, err := sitter.NewQuery([]byte(goImportQuery), p.language[Go])
	if err != nil {
		return fmt.Errorf("failed to create Go import query: %w", err)
	}
	p.queries[Go]["import"] = importQuery

	// Initialize function query
	functionQuery, err := sitter.NewQuery([]byte(goFunctionQuery), p.language[Go])
	if err != nil {
		return fmt.Errorf("failed to create Go function query: %w", err)
	}
	p.queries[Go]["function"] = functionQuery

	return nil
}

func (p *Parser) extractGoImports(tree *sitter.Tree, content []byte) ([]string, error) {
	query := p.queries[Go]["import"]
	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, tree.RootNode())

	var imports []string
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if capture.Node != nil {
				importPath := capture.Node.Content(content)
				// Remove quotes from import path
				importPath = strings.Trim(importPath, "\"")
				imports = append(imports, importPath)
			}
		}
	}

	return imports, nil
}

func (p *Parser) extractGoFunctions(tree *sitter.Tree, content []byte) ([]FunctionInfo, error) {
	query := p.queries[Go]["function"]
	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, tree.RootNode())

	var functions []FunctionInfo
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if capture.Node != nil {
				fn := FunctionInfo{
					Name:      capture.Node.Content(content),
					StartLine: uint32(capture.Node.StartPoint().Row + 1),
					EndLine:   uint32(capture.Node.EndPoint().Row + 1),
				}
				fn.IsExported = isExported(fn.Name)
				functions = append(functions, fn)
			}
		}
	}

	return functions, nil
}

func (p *Parser) extractGoAPIEndpoints(tree *sitter.Tree, content []byte) ([]APIEndpoint, error) {
	query := p.queries[Go]["api"]
	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, tree.RootNode())

	var endpoints []APIEndpoint
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		var endpoint APIEndpoint
		for _, capture := range match.Captures {
			nodeContent := capture.Node.Content(content)
			switch {
			case strings.HasPrefix(nodeContent, "@router"):
				if isRouter(nodeContent) {
					endpoint.Handler = nodeContent
				}
			case strings.HasPrefix(nodeContent, "@method"):
				endpoint.Method = nodeContent
			case strings.HasPrefix(nodeContent, "@api.args"):
				endpoint.Path = extractPath(nodeContent)
			}
		}

		if endpoint.Method != "" && endpoint.Path != "" {
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints, nil
}

// Helper functions

func isExported(name string) bool {
	if len(name) == 0 {
		return false
	}
	// In Go, a name is exported if it begins with an uppercase letter
	return name[0] >= 'A' && name[0] <= 'Z'
}

func isRouter(name string) bool {
	routers := []string{"router", "mux", "e", "r", "app"}
	for _, r := range routers {
		if name == r {
			return true
		}
	}
	return false
}

func extractPath(args string) string {
	// Simple path extraction, in practice you'd want more robust parsing
	if len(args) == 0 {
		return ""
	}
	// Remove parentheses and get first string argument
	args = args[1 : len(args)-1]
	return args
}
