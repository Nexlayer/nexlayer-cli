package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	"github.com/sahilm/fuzzy"
)

// DocSearch provides fuzzy search over documentation files
type DocSearch struct {
	docs     []string
	contents map[string]string
}

// NewDocSearch creates a new DocSearch instance
func NewDocSearch(docsPath, templatesPath string) (*DocSearch, error) {
	ds := &DocSearch{
		contents: make(map[string]string),
	}

	// Load documentation files
	if docsPath != "" {
		if err := ds.loadFiles(docsPath); err != nil {
			return nil, fmt.Errorf("failed to load docs: %w", err)
		}
	}

	// Load template files
	if templatesPath != "" {
		if err := ds.loadFiles(templatesPath); err != nil {
			return nil, fmt.Errorf("failed to load templates: %w", err)
		}
	}

	return ds, nil
}

// loadFiles recursively loads files from the given path
func (ds *DocSearch) loadFiles(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Only load text files
		switch filepath.Ext(path) {
		case ".md", ".txt", ".yaml", ".yml":
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", path, err)
			}

			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", path, err)
			}

			ds.docs = append(ds.docs, relPath)
			ds.contents[relPath] = string(content)
		}

		return nil
	})
}

// Search performs a fuzzy search over documentation files
func (ds *DocSearch) Search(query string) []string {
	// Perform fuzzy search over doc paths
	matches := fuzzy.Find(query, ds.docs)

	// Sort matches by score
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	// Extract matched paths
	results := make([]string, len(matches))
	for i, match := range matches {
		results[i] = match.Str
	}

	return results
}

// GetContent returns the content of a documentation file
func (ds *DocSearch) GetContent(path string) string {
	return ds.contents[path]
}
