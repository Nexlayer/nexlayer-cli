// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package starters

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
)

// Template represents a project starter template
type Template struct {
	Name  string            `json:"name"`
	Desc  string            `json:"description"`
	Files map[string]string `json:"files"`
	Stack []string          `json:"stack"`
}

// TemplateItem implements list.Item for bubbletea
type TemplateItem struct {
	Name string
	Desc string
}

func (t TemplateItem) Title() string       { return t.Name }
func (t TemplateItem) Description() string { return t.Desc }
func (t TemplateItem) FilterValue() string { return t.Name }

var defaultTemplates = []Template{
	{
		Name:  "Full Stack App",
		Desc:  "React + FastAPI + PostgreSQL",
		Stack: []string{"react", "fastapi", "postgres"},
		Files: map[string]string{
			"frontend/package.json": `{
				"name": "{{.Name}}-frontend",
				"version": "0.1.0",
				"private": true,
				"dependencies": {
					"react": "^18.2.0",
					"react-dom": "^18.2.0",
					"react-scripts": "5.0.1"
				}
			}`,
			"backend/requirements.txt": `fastapi==0.104.1
uvicorn==0.24.0
sqlalchemy==2.0.23
psycopg2-binary==2.9.9`,
		},
	},
	{
		Name:  "API Service",
		Desc:  "FastAPI + PostgreSQL",
		Stack: []string{"fastapi", "postgres"},
		Files: map[string]string{
			"requirements.txt": `fastapi==0.104.1
uvicorn==0.24.0
sqlalchemy==2.0.23
psycopg2-binary==2.9.9`,
		},
	},
	{
		Name:  "Static Website",
		Desc:  "React Single Page App",
		Stack: []string{"react"},
		Files: map[string]string{
			"package.json": `{
				"name": "{{.Name}}",
				"version": "0.1.0",
				"private": true,
				"dependencies": {
					"react": "^18.2.0",
					"react-dom": "^18.2.0",
					"react-scripts": "5.0.1"
				}
			}`,
		},
	},
}

// CreateProject creates a new project from a template
func CreateProject(projectName, templateName string) error {
	var selectedTemplate *Template
	for _, t := range defaultTemplates {
		if t.Name == templateName {
			selectedTemplate = &t
			break
		}
	}
	if selectedTemplate == nil {
		return fmt.Errorf("template %q not found", templateName)
	}

	// Create project directory
	if err := os.MkdirAll(projectName, 0o755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	// Create files from template
	for path, content := range selectedTemplate.Files {
		fullPath := filepath.Join(projectName, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", filepath.Dir(fullPath), err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", fullPath, err)
		}
	}

	return nil
}

// GetTemplateItems returns a list of template items for bubbletea
func GetTemplateItems() []list.Item {
	var items []list.Item
	for _, t := range defaultTemplates {
		items = append(items, TemplateItem{
			Name: t.Name,
			Desc: t.Desc,
		})
	}
	return items
}
