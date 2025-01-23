package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// DocsContent stores documentation content by category
type DocsContent struct {
	Deployment map[string]string
	Domain     map[string]string
	Status     map[string]string
	Template   map[string]string
}

// LoadDocumentation loads documentation from the local Nexlayer docs and templates repositories
func LoadDocumentation(docsPath, templatesPath string) (*DocsContent, error) {
	docs := &DocsContent{
		Deployment: make(map[string]string),
		Domain:     make(map[string]string),
		Status:     make(map[string]string),
		Template:   make(map[string]string),
	}

	// Load docs
	if err := filepath.Walk(docsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-markdown files and mintlify configuration
		if !strings.HasSuffix(path, ".md") || strings.Contains(path, "mintlify") {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Categorize content based on path and content
		contentStr := string(content)
		switch {
		case strings.Contains(path, "deployment") || strings.Contains(contentStr, "deployment"):
			docs.Deployment[filepath.Base(path)] = contentStr
		case strings.Contains(path, "domain") || strings.Contains(contentStr, "domain"):
			docs.Domain[filepath.Base(path)] = contentStr
		case strings.Contains(path, "status") || strings.Contains(path, "monitoring"):
			docs.Status[filepath.Base(path)] = contentStr
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error loading docs: %w", err)
	}

	// Load templates
	if err := filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only process template configuration files
		if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".md") {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		docs.Template[filepath.Base(path)] = string(content)
		return nil
	}); err != nil {
		return nil, fmt.Errorf("error loading templates: %w", err)
	}

	return docs, nil
}

// GetRelevantDocs returns relevant documentation content for a specific category
func (d *DocsContent) GetRelevantDocs(category string) []string {
	var content []string
	
	switch category {
	case "deployment":
		for _, doc := range d.Deployment {
			content = append(content, doc)
		}
	case "domain":
		for _, doc := range d.Domain {
			content = append(content, doc)
		}
	case "status":
		for _, doc := range d.Status {
			content = append(content, doc)
		}
	case "template":
		for _, doc := range d.Template {
			content = append(content, doc)
		}
	}

	return content
}

// GetTemplateExamples returns example configurations from templates
func (d *DocsContent) GetTemplateExamples() []string {
	var examples []string
	for _, content := range d.Template {
		examples = append(examples, content)
	}
	return examples
}

// ExtractCommands extracts command examples from documentation
func (d *DocsContent) ExtractCommands(category string) []string {
	var commands []string
	docs := d.GetRelevantDocs(category)

	for _, doc := range docs {
		// Look for command blocks in markdown
		lines := strings.Split(doc, "\n")
		inCodeBlock := false
		var currentBlock strings.Builder

		for _, line := range lines {
			if strings.HasPrefix(line, "```") {
				if inCodeBlock {
					// End of code block
					block := currentBlock.String()
					if strings.Contains(block, "nexlayer") {
						commands = append(commands, block)
					}
					currentBlock.Reset()
				}
				inCodeBlock = !inCodeBlock
				continue
			}

			if inCodeBlock {
				currentBlock.WriteString(line + "\n")
			}
		}
	}

	return commands
}
