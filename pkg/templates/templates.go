// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package templates

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
)

// Template represents a project starter template
type Template struct {
	Name        string            `json:"name"`
	Desc        string            `json:"description"`
	Files       map[string]string `json:"files"`
	Stack       []string          `json:"stack"`
}

// TemplateItem implements list.Item for bubbletea
type TemplateItem struct {
	Name string
	Desc string
}

func (i TemplateItem) Title() string       { return i.Name }
func (i TemplateItem) Description() string { return i.Desc }
func (i TemplateItem) FilterValue() string { return i.Name }

var defaultTemplates = []Template{
	{
		Name: "Full Stack App",
		Desc: "React + FastAPI + PostgreSQL",
		Stack: []string{"react", "fastapi", "postgres"},
		Files: map[string]string{
			"frontend/package.json": `{
  "name": "frontend",
  "version": "0.1.0",
  "private": true,
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-scripts": "5.0.1"
  }
}`,
			"backend/requirements.txt": `fastapi>=0.104.1
uvicorn>=0.24.0
sqlalchemy>=2.0.23
psycopg2-binary>=2.9.9
python-dotenv>=1.0.0`,
			"backend/main.py": `from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def root():
    return {"message": "Hello World"}`,
		},
	},
	{
		Name: "Backend Only",
		Desc: "FastAPI + PostgreSQL",
		Stack: []string{"fastapi", "postgres"},
		Files: map[string]string{
			"requirements.txt": `fastapi>=0.104.1
uvicorn>=0.24.0
sqlalchemy>=2.0.23
psycopg2-binary>=2.9.9
python-dotenv>=1.0.0`,
			"main.py": `from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def root():
    return {"message": "Hello World"}`,
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
		return fmt.Errorf("template %s not found", templateName)
	}

	// Create project directory
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create files from template
	for path, content := range selectedTemplate.Files {
		fullPath := filepath.Join(projectName, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(fullPath), err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
	}

	return nil
}

// GetTemplateItems returns a list of template items for bubbletea
func GetTemplateItems() []list.Item {
	items := make([]list.Item, len(defaultTemplates))
	for i, t := range defaultTemplates {
		items[i] = TemplateItem{
			Name: t.Name,
			Desc: t.Desc,
		}
	}
	return items
}
