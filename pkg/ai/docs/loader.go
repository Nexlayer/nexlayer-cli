package docs

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Content holds documentation content with metadata
type Content struct {
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	Type     string            `json:"type"`
	Tags     []string          `json:"tags"`
	Metadata map[string]string `json:"metadata"`
}

// Store manages documentation content
type Store struct {
	content map[string]*Content
}

// NewStore creates a new documentation store
func NewStore() *Store {
	return &Store{
		content: make(map[string]*Content),
	}
}

// LoadContent loads documentation content from embedded files
func (s *Store) LoadContent(fsys fs.FS) error {
	// Walk through all files in the filesystem
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Process only markdown and JSON files
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".md" && ext != ".json" {
			return nil
		}

		// Read file content
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %v", path, err)
		}

		var content *Content

		switch ext {
		case ".json":
			// Parse JSON content
			if err := json.Unmarshal(data, &content); err != nil {
				return fmt.Errorf("failed to parse JSON file %s: %v", path, err)
			}
		case ".md":
			// Create content from markdown
			content = &Content{
				Title:   strings.TrimSuffix(filepath.Base(path), ext),
				Content: string(data),
				Type:    "markdown",
				Tags:    extractTags(string(data)),
			}
		}

		// Store content
		s.content[path] = content
		return nil
	})
}

// LoadPluginDocs loads documentation from plugin directories
func (s *Store) LoadPluginDocs(pluginPath string) error {
	// Walk through all plugin directories
	entries, err := os.ReadDir(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Look for README.md in each plugin directory
		readmePath := filepath.Join(pluginPath, entry.Name(), "README.md")
		data, err := os.ReadFile(readmePath)
		if err != nil {
			// Skip if README doesn't exist
			continue
		}

		// Create content from markdown
		content := &Content{
			Title:   entry.Name(),
			Content: string(data),
			Type:    "plugin",
			Tags:    extractTags(string(data)),
		}

		// Store content
		s.content[fmt.Sprintf("plugins/%s/README.md", entry.Name())] = content
	}

	return nil
}

// GetContent retrieves content by path
func (s *Store) GetContent(path string) (*Content, bool) {
	content, ok := s.content[path]
	return content, ok
}

// Search searches for content matching the query
func (s *Store) Search(query string) []*Content {
	var results []*Content
	query = strings.ToLower(query)

	for _, content := range s.content {
		// Search in title
		if strings.Contains(strings.ToLower(content.Title), query) {
			results = append(results, content)
			continue
		}

		// Search in content
		if strings.Contains(strings.ToLower(content.Content), query) {
			results = append(results, content)
			continue
		}

		// Search in tags
		for _, tag := range content.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, content)
				break
			}
		}
	}

	return results
}

// GetByType returns all content of a specific type
func (s *Store) GetByType(docType string) []*Content {
	var results []*Content
	for _, content := range s.content {
		if content.Type == docType {
			results = append(results, content)
		}
	}
	return results
}

// GetByTag returns all content with a specific tag
func (s *Store) GetByTag(tag string) []*Content {
	var results []*Content
	tag = strings.ToLower(tag)
	for _, content := range s.content {
		for _, t := range content.Tags {
			if strings.ToLower(t) == tag {
				results = append(results, content)
				break
			}
		}
	}
	return results
}

// extractTags extracts tags from markdown content
func extractTags(content string) []string {
	var tags []string
	lines := strings.Split(content, "\n")
	
	// Look for tags section
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Tags:") || strings.HasPrefix(line, "tags:") {
			// Extract tags from the line
			tagPart := strings.TrimPrefix(strings.TrimPrefix(line, "Tags:"), "tags:")
			tagPart = strings.TrimSpace(tagPart)
			
			// Split tags by comma
			for _, tag := range strings.Split(tagPart, ",") {
				tag = strings.TrimSpace(tag)
				if tag != "" {
					tags = append(tags, tag)
				}
			}
			break
		}
	}
	
	return tags
}
