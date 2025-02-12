// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PackageInfo struct {
	Name           string            `json:"name"`
	Path           string            `json:"path"`
	Files          []string         `json:"files"`
	Functions      []FunctionInfo    `json:"functions"`
	Interfaces     []InterfaceInfo  `json:"interfaces"`
	Structs        []StructInfo     `json:"structs"`
	Doc           string           `json:"doc"`
	Responsibility string           `json:"responsibility"`
	Imports        []string         `json:"imports"`
	DependedOnBy   []string         `json:"dependedOnBy"`
}

type FunctionInfo struct {
	Name       string   `json:"name"`
	Doc        string   `json:"doc"`
	Parameters []string `json:"parameters"`
	Returns    []string `json:"returns"`
}

type InterfaceInfo struct {
	Name    string   `json:"name"`
	Doc     string   `json:"doc"`
	Methods []string `json:"methods"`
}

type StructInfo struct {
	Name   string     `json:"name"`
	Doc    string     `json:"doc"`
	Fields []string   `json:"fields"`
}

type DependencyInfo struct {
	Module    string   `json:"module"`
	Version   string   `json:"version"`
	Requires  []string `json:"requires"`
}

type Metadata struct {
	ProjectName     string                 `json:"projectName"`
	Version         string                 `json:"version"`
	Packages        []PackageInfo          `json:"packages"`
	Dependencies    []DependencyInfo       `json:"dependencies"`
	CallGraphs      map[string][]string    `json:"callGraphs"`
	BuildInfo       map[string]interface{} `json:"buildInfo"`
}

func main() {
	if err := generateMetadata(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating metadata: %v\n", err)
		os.Exit(1)
	}
}

func generateMetadata() error {
	// Create AI training analysis directory if it doesn't exist
	aiDir := filepath.Join("ai_training", "analysis")
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		return fmt.Errorf("failed to create AI training directory: %w", err)
	}

	metadata := &Metadata{
		ProjectName: "nexlayer-cli",
		Packages:    make([]PackageInfo, 0),
		CallGraphs:  make(map[string][]string),
		BuildInfo:   make(map[string]interface{}),
	}

	// Parse Go modules info
	if err := parseModuleInfo(metadata); err != nil {
		return fmt.Errorf("failed to parse module info: %w", err)
	}

	// Parse package structure
	if err := parsePackages(metadata); err != nil {
		return fmt.Errorf("failed to parse packages: %w", err)
	}

	// Generate call graphs for key packages
	if err := generateCallGraphs(metadata); err != nil {
		return fmt.Errorf("failed to generate call graphs: %w", err)
	}

	// Write metadata to file
	file, err := os.Create(filepath.Join(aiDir, "metadata.json"))
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metadata); err != nil {
		return fmt.Errorf("failed to encode metadata: %w", err)
	}

	return nil
}

func parseModuleInfo(metadata *Metadata) error {
	cmd := exec.Command("go", "mod", "graph")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run go mod graph: %w", err)
	}

	// Parse module dependencies
	lines := strings.Split(string(output), "\n")
	deps := make(map[string]DependencyInfo)
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		moduleParts := strings.Split(parts[0], "@")
		module := moduleParts[0]
		version := ""
		if len(moduleParts) > 1 {
			version = moduleParts[1]
		}
		
		dep, exists := deps[module]
		if !exists {
			dep = DependencyInfo{
				Module:   module,
				Version:  version,
				Requires: make([]string, 0),
			}
		}
		dep.Requires = append(dep.Requires, parts[1])
		deps[module] = dep
	}

	for _, dep := range deps {
		metadata.Dependencies = append(metadata.Dependencies, dep)
	}

	return nil
}

func getPackageDoc(path string) string {
	cmd := exec.Command("go", "doc", path)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func parsePackages(metadata *Metadata) error {
	fset := token.NewFileSet()
	
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || strings.Contains(path, "vendor") || strings.Contains(path, "ai_training") {
			return nil
		}

		pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil // Skip directories with parse errors
		}

		for name, pkg := range pkgs {
			pkgInfo := PackageInfo{
				Name:       name,
				Path:       path,
				Files:      make([]string, 0),
				Functions:  make([]FunctionInfo, 0),
				Interfaces: make([]InterfaceInfo, 0),
				Structs:    make([]StructInfo, 0),
				Doc:        getPackageDoc(path),
			}

			ast.Inspect(pkg, func(n ast.Node) bool {
				switch x := n.(type) {
				case *ast.File:
					pkgInfo.Files = append(pkgInfo.Files, x.Name.Name)
				case *ast.FuncDecl:
					fn := FunctionInfo{
						Name:       x.Name.Name,
						Doc:        docString(x.Doc),
						Parameters: make([]string, 0),
						Returns:    make([]string, 0),
					}
					pkgInfo.Functions = append(pkgInfo.Functions, fn)
				case *ast.InterfaceType:
					if parent, ok := n.(*ast.TypeSpec); ok {
						iface := InterfaceInfo{
							Name:    parent.Name.Name,
							Methods: make([]string, 0),
						}
						pkgInfo.Interfaces = append(pkgInfo.Interfaces, iface)
					}
				case *ast.StructType:
					if parent, ok := n.(*ast.TypeSpec); ok {
						st := StructInfo{
							Name:   parent.Name.Name,
							Fields: make([]string, 0),
						}
						pkgInfo.Structs = append(pkgInfo.Structs, st)
					}
				}
				return true
			})

			// Infer package responsibility
			pkgInfo.Responsibility = inferPackageResponsibility(&pkgInfo)

			// Add imports
			for name, file := range pkg.Files {
				for _, imp := range file.Imports {
					if imp.Path != nil {
						pkgInfo.Imports = append(pkgInfo.Imports, imp.Path.Value)
					}
				}
				pkgInfo.Files = append(pkgInfo.Files, name)
			}

			metadata.Packages = append(metadata.Packages, pkgInfo)
		}

		// Build dependency graph
		for i := range metadata.Packages {
			for _, imp := range metadata.Packages[i].Imports {
				for j := range metadata.Packages {
					if strings.Contains(imp, metadata.Packages[j].Name) {
						metadata.Packages[j].DependedOnBy = append(
							metadata.Packages[j].DependedOnBy,
							metadata.Packages[i].Path,
						)
					}
				}
			}
		}

		return nil
	})

	return err
}

func inferPackageResponsibility(pkg *PackageInfo) string {
	// Infer package responsibility based on name, contents and documentation
	responsibility := "Unknown"

	// Check package name patterns
	switch {
	case strings.Contains(pkg.Path, "/cmd"):
		responsibility = "Command-line interface and entry points"
	case strings.Contains(pkg.Path, "/pkg/commands"):
		responsibility = "CLI command implementations"
	case strings.Contains(pkg.Path, "/pkg/core"):
		responsibility = "Core functionality and business logic"
	case strings.Contains(pkg.Path, "/pkg/plugins"):
		responsibility = "Plugin system and extensions"
	case strings.Contains(pkg.Path, "/pkg/api"):
		responsibility = "API client and interfaces"
	}

	// Look for key interfaces/structs that indicate responsibility
	for _, iface := range pkg.Interfaces {
		if strings.Contains(strings.ToLower(iface.Name), "service") {
			responsibility = "Service layer implementation"
		} else if strings.Contains(strings.ToLower(iface.Name), "repository") {
			responsibility = "Data access and storage"
		}
	}

	return responsibility
}

func generateCallGraphs(metadata *Metadata) error {
	// Verify go-callvis is installed
	if _, err := exec.LookPath("go-callvis"); err != nil {
		return fmt.Errorf("go-callvis not found in PATH: %w", err)
	}

	// Key packages to generate call graphs for
	keyPackages := []string{
		"./cmd",
		"./pkg/commands",
		"./pkg/core/api",
		"./pkg/plugins",
	}

	// Get project root directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Create AI training analysis directory if it doesn't exist
	aiDir := filepath.Join(wd, "ai_training", "analysis")
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		return fmt.Errorf("failed to create AI training directory: %w", err)
	}

	for _, pkg := range keyPackages {
		fmt.Printf("Generating call graph for package: %s\n", pkg)

		// Clean package path for filename
		cleanPkg := strings.Replace(pkg, "/", "_", -1)
		outputFile := filepath.Join(aiDir, fmt.Sprintf("callgraph_%s.svg", cleanPkg))

		// Get module name from go.mod
		modCmd := exec.Command("go", "list", "-m")
		modCmd.Dir = wd
		moduleBytes, err := modCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get module name: %w", err)
		}
		moduleName := strings.TrimSpace(string(moduleBytes))

		// Convert relative package path to absolute import path
		absPkg := strings.TrimPrefix(filepath.Join(moduleName, pkg), "./")

		// Build go-callvis command with focus on the package
		cmd := exec.Command("go-callvis",
			"-format=svg",
			"-focus="+absPkg,
			"-group=pkg",
			"-limit="+absPkg,
			"-nostd",
			"-file="+outputFile,
			absPkg)

		// Set working directory to project root
		cmd.Dir = wd

		// Capture both stdout and stderr
		output, err := cmd.CombinedOutput()
		if err != nil {
			// Store error output in metadata for debugging
			metadata.CallGraphs[pkg] = []string{
				"", // No valid file path
				fmt.Sprintf("Error generating call graph: %v\n%s", err, output),
			}
			fmt.Printf("Warning: failed to generate call graph for %s:\n%s\n", pkg, output)
			continue
		}

		// Verify the output file exists
		if _, err := os.Stat(outputFile); err != nil {
			metadata.CallGraphs[pkg] = []string{
				"",
				fmt.Sprintf("Error: call graph file not generated: %v", err),
			}
			continue
		}

		// Store the call graph file path and any output in metadata
		metadata.CallGraphs[pkg] = []string{
			outputFile,
			string(output),
		}
	}

	return nil
}

func docString(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	return doc.Text()
}
