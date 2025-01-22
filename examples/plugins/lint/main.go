package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	fix bool
)

func init() {
	rootCmd.Flags().BoolVarP(&fix, "fix", "f", false, "Automatically fix issues")
}

var rootCmd = &cobra.Command{
	Use:   "lint [dir]",
	Short: "Lint Go code",
	Long: `A linting plugin for Nexlayer CLI.
Example: nexlayer lint ./...`,
	Args: cobra.ExactArgs(1),
	RunE: runLint,
}

func runLint(cmd *cobra.Command, args []string) error {
	dir := args[0]

	// Get all Go files in directory
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	// Lint each file
	for _, file := range files {
		fmt.Printf("Linting %s...\n", file)
		if fix {
			fmt.Printf("Fixing %s...\n", file)
		}
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
