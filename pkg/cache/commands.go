package cache

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/config"
)

// CacheCmd represents the cache command
var CacheCmd = &cobra.Command{
	Use:   "cache [command]",
	Short: "Manage cache",
	Long: `Manage the Nexlayer CLI cache.
Example: nexlayer cache clear`,
	Args: cobra.MinimumNArgs(1),
	RunE: runCache,
}

func runCache(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}
	command := args[0]

	switch command {
	case "clear":
		return clearCache()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func clearCache() error {
	cache := Get()
	cache.Clear()

	cacheDir := filepath.Join(config.GetConfigDir(), "cache")
	if err := os.RemoveAll(cacheDir); err != nil {
		return fmt.Errorf("failed to remove cache directory: %w", err)
	}

	fmt.Println("Cache cleared successfully")
	return nil
}
