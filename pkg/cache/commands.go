package cache

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// AddCacheCommands adds cache-related commands to the root command
func AddCacheCommands(rootCmd *cobra.Command, manager *Manager) {
	cacheCmd := &cobra.Command{
		Use:   "cache",
		Short: "Manage the Nexlayer CLI cache",
	}

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all cached data",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := manager.ClearCache(); err != nil {
				return fmt.Errorf("failed to clear cache: %w", err)
			}
			fmt.Println("Cache cleared successfully")
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all cached items",
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := manager.ListCache()
			if err != nil {
				return fmt.Errorf("failed to list cache: %w", err)
			}

			if len(items) == 0 {
				fmt.Println("Cache is empty")
				return nil
			}

			fmt.Printf("Found %d cached items:\n\n", len(items))
			for _, item := range items {
				fmt.Printf("Key: %s\n", item.Key)
				fmt.Printf("  Size: %d bytes\n\n", item.Size)
			}

			return nil
		},
	}

	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize pending offline operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := manager.OfflineConfig.SyncPendingOperations(&http.Client{
				Timeout: 30 * time.Second,
			}); err != nil {
				return fmt.Errorf("sync failed: %w", err)
			}
			fmt.Println("Successfully synchronized pending operations")
			return nil
		},
	}

	pendingCmd := &cobra.Command{
		Use:   "pending",
		Short: "List pending offline operations",
		RunE: func(cmd *cobra.Command, args []string) error {
			ops := manager.OfflineConfig.GetPendingOperations()
			if len(ops) == 0 {
				fmt.Println("No pending operations")
				return nil
			}

			fmt.Printf("Found %d pending operations:\n\n", len(ops))
			for i, op := range ops {
				fmt.Printf("%d. %s %s\n", i+1, op.Method, op.Endpoint)
				fmt.Printf("   Created: %s\n", op.CreatedAt.Format(time.RFC3339))
				if len(op.Payload) > 0 {
					fmt.Printf("   Has payload: yes\n")
				}
				fmt.Println()
			}
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show offline mode status",
		RunE: func(cmd *cobra.Command, args []string) error {
			isOnline := manager.OfflineConfig.IsOnline()
			mode := manager.OfflineConfig.Mode

			fmt.Printf("Network Status: %s\n", map[bool]string{true: "Online", false: "Offline"}[isOnline])
			fmt.Printf("Operation Mode: %s\n", map[OfflineMode]string{
				OnlineMode:        "Online",
				AutoOfflineMode:   "Auto-Offline",
				StrictOfflineMode: "Strict-Offline",
			}[mode])

			pendingOps := len(manager.OfflineConfig.GetPendingOperations())
			fmt.Printf("Pending Operations: %d\n", pendingOps)

			return nil
		},
	}

	cacheCmd.AddCommand(clearCmd)
	cacheCmd.AddCommand(listCmd)
	cacheCmd.AddCommand(syncCmd)
	cacheCmd.AddCommand(pendingCmd)
	cacheCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(cacheCmd)
}
