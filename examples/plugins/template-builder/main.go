// Formatted with gofmt -s
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/detector"
	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/generator"
	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/types"
	"gopkg.in/yaml.v3"
)

// PluginMetadata for --describe
type PluginMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: nexlayer template:generate [--dry-run]")
		os.Exit(1)
	}

	if os.Args[1] == "--describe" {
		describePlugin()
		return
	}

	// Get project name from current directory
	projectName := getProjectName()

	// Detect project stack
	stack := detector.DetectStack()

	// Generate template
	template := generator.GenerateTemplate(projectName, stack)

	// Refine template if AI is available
	if err := refineTemplate(template, stack); err != nil {
		fmt.Printf("Warning: Could not refine template: %v\n", err)
	}

	// Convert to YAML
	yamlData, err := yaml.Marshal(template)
	if err != nil {
		fmt.Printf("Error marshaling template: %v\n", err)
		os.Exit(1)
	}

	// Check if dry run
	if len(os.Args) > 2 && os.Args[2] == "--dry-run" {
		fmt.Println(string(yamlData))
		return
	}

	// Write template to file
	filename := fmt.Sprintf("%s-nexlayer-template.yaml", projectName)
	if err := writeTemplateFile(filename, yamlData); err != nil {
		fmt.Printf("Error writing template: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[SUCCESS] Template generated successfully: %s\n", filename)
}

func describePlugin() {
	metadata := PluginMetadata{
		Name:        "template-builder",
		Version:     "1.0.0",
		Description: "Generates Nexlayer deployment templates based on project stack",
		Usage:       "nexlayer template:generate [--dry-run]",
	}

	json.NewEncoder(os.Stdout).Encode(metadata)
}

func getProjectName() string {
	dir, err := os.Getwd()
	if err != nil {
		return "template-builder"
	}
	return filepath.Base(dir)
}

func writeTemplateFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

func refineTemplate(template *types.NexlayerTemplate, stack *types.ProjectStack) error {
	refiner := NewAIRefiner()
	if refiner == nil {
		return fmt.Errorf("no AI refiner available")
	}

	return refiner.RefineTemplate(template, *stack)
}
