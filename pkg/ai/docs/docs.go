package docs

import (
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed ai.md
var DocFiles embed.FS

// GetDocStore returns a new documentation store with loaded content
func GetDocStore() (*Store, error) {
	store := NewStore()

	// Load embedded documentation
	if err := store.LoadContent(DocFiles); err != nil {
		return nil, fmt.Errorf("failed to load embedded docs: %v", err)
	}

	// Load plugin documentation from examples directory
	pluginPath := filepath.Join("..", "..", "..", "examples", "plugins")
	if err := store.LoadPluginDocs(pluginPath); err != nil {
		// Just log error but continue - plugin docs are optional
		fmt.Printf("Warning: failed to load plugin docs: %v\n", err)
	}

	return store, nil
}
