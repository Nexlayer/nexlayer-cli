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
    name: (identifier) @function.name
    parameters: (parameter_list) @function.params
    result: (result) @function.result))
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

	// Extract API endpoints
	endpoints, err := p.extractGoAPIEndpoints(tree, content)
	if err != nil {
		return err
	}
	analysis.APIEndpoints = append(analysis.APIEndpoints, endpoints...)

	return nil
}

func (p *Parser) initGoQueries() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.queries[Go]; ok {
		return nil
	}

	// Initialize all queries
	queries := map[string]string{
		"import":   goImportQuery,
		"function": goFunctionQuery,
		"api":      goAPIEndpointQuery,
	}

	for name, query := range queries {
		q, err := sitter.NewQuery([]byte(query), p.language[Go])
		if err != nil {
			return fmt.Errorf("failed to create Go %s query: %w", name, err)
		}
		p.queries[Go] = q
	}

	return nil
}

func (p *Parser) extractGoImports(tree *sitter.Tree, content []byte) ([]string, error) {
	query := p.queries[Go]
	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, tree.RootNode())

	var imports []string
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			nodeContent := capture.Node.Content(content)
			if strings.HasPrefix(nodeContent, "@import.path") {
				// Remove quotes from import path
				importPath := nodeContent[1 : len(nodeContent)-1] // Remove quotes
				imports = append(imports, importPath)
			}
		}
	}

	return imports, nil
}

func (p *Parser) extractGoFunctions(tree *sitter.Tree, content []byte) ([]FunctionInfo, error) {
	query := p.queries[Go]
	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, tree.RootNode())

	var functions []FunctionInfo
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		var fn FunctionInfo
		for _, capture := range match.Captures {
			nodeContent := capture.Node.Content(content)
			switch {
			case strings.HasPrefix(nodeContent, "@function.name"):
				fn.Name = nodeContent
			case strings.HasPrefix(nodeContent, "@function.params"):
				fn.Signature = nodeContent
			}
		}

		if fn.Name != "" {
			fn.StartLine = uint32(match.Captures[0].Node.StartPoint().Row)
			fn.EndLine = uint32(match.Captures[0].Node.EndPoint().Row)
			fn.IsExported = isExported(fn.Name)
			functions = append(functions, fn)
		}
	}

	return functions, nil
}

func (p *Parser) extractGoAPIEndpoints(tree *sitter.Tree, content []byte) ([]APIEndpoint, error) {
	query := p.queries[Go]
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
