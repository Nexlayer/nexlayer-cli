// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// CheckPathSetup verifies if the nexlayer-cli binary is properly accessible from PATH
func CheckPathSetup() error {
	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Get PATH directories
	pathEnv := os.Getenv("PATH")
	pathDirs := strings.Split(pathEnv, string(os.PathListSeparator))

	// Get GOPATH/bin
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		gopath = filepath.Join(homeDir, "go")
	}
	gopathBin := filepath.Join(gopath, "bin")

	// Check if binary is in PATH
	binaryInPath := false
	for _, dir := range pathDirs {
		if dir == filepath.Dir(exePath) {
			binaryInPath = true
			break
		}
	}

	if !binaryInPath {
		// Provide guidance based on shell and OS
		shell := os.Getenv("SHELL")
		if shell == "" {
			if runtime.GOOS == "windows" {
				shell = "cmd.exe"
			} else {
				shell = "/bin/bash"
			}
		}

		var guidance string
		switch {
		case strings.Contains(shell, "zsh"):
			guidance = fmt.Sprintf(`Add the following line to your ~/.zshrc:
echo 'export PATH="$PATH:%s"' >> ~/.zshrc
source ~/.zshrc`, gopathBin)
		case strings.Contains(shell, "bash"):
			guidance = fmt.Sprintf(`Add the following line to your ~/.bashrc:
echo 'export PATH="$PATH:%s"' >> ~/.bashrc
source ~/.bashrc`, gopathBin)
		case strings.Contains(shell, "fish"):
			guidance = fmt.Sprintf(`Add the following line to your ~/.config/fish/config.fish:
set -gx PATH $PATH %s`, gopathBin)
		case runtime.GOOS == "windows":
			guidance = fmt.Sprintf(`Add %s to your PATH by:
1. Press Win + X and select "System"
2. Click "Advanced system settings"
3. Click "Environment Variables"
4. Under "User variables", select "Path" and click "Edit"
5. Click "New" and add: %s`, gopathBin, gopathBin)
		}

		return fmt.Errorf("nexlayer-cli is not in PATH. To fix this:\n\n%s", guidance)
	}

	return nil
}

// IsFirstRun checks if this is the first time running nexlayer-cli
func IsFirstRun() bool {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return true
	}
	
	nexlayerDir := filepath.Join(configDir, "nexlayer")
	_, err = os.Stat(nexlayerDir)
	return os.IsNotExist(err)
}
