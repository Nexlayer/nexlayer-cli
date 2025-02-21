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

// CheckPathSetup ensures the nexlayer-cli binary is accessible from PATH
func CheckPathSetup() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	pathEnv := os.Getenv("PATH")
	pathDirs := strings.Split(pathEnv, string(os.PathListSeparator))

	// Handle multiple GOPATHs
	gopaths := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))
	if len(gopaths) == 0 || (len(gopaths) == 1 && gopaths[0] == "") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		gopaths = []string{filepath.Join(homeDir, "go")}
	}

	binaryInPath := false
	gopathBin := ""
	for _, gp := range gopaths {
		gopathBin = filepath.Join(gp, "bin")
		for _, dir := range pathDirs {
			if dir == gopathBin && filepath.Dir(exePath) == gopathBin {
				binaryInPath = true
				break
			}
		}
		if binaryInPath {
			break
		}
	}

	if !binaryInPath {
		shell := getShell()
		guidance := getPathGuidance(shell, gopathBin)
		if confirm("Would you like to automatically add nexlayer-cli to PATH?") {
			if err := appendToShellConfig(shell, gopathBin); err != nil {
				return fmt.Errorf("failed to update PATH: %w\nManual instructions:\n\n%s", err, guidance)
			}
			fmt.Println("PATH updated successfully. You may need to restart your terminal.")
		} else {
			return fmt.Errorf("nexlayer-cli is not in PATH. To fix this:\n\n%s", guidance)
		}
	}

	return nil
}

// getShell detects the current shell
func getShell() string {
	if runtime.GOOS == "windows" {
		if _, err := os.Stat("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"); err == nil {
			return "powershell"
		}
		if _, err := os.Stat("C:\\Program Files\\Git\\bin\\bash.exe"); err == nil {
			return "gitbash"
		}
		return "cmd"
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "/bin/sh"
	}
	return shell
}

// getPathGuidance provides shell-specific instructions to add to PATH
func getPathGuidance(shell, gopathBin string) string {
	switch {
	case strings.Contains(shell, "zsh"):
		return fmt.Sprintf(`Add the following line to your ~/.zshrc:
echo 'export PATH="$PATH:%s"' >> ~/.zshrc
source ~/.zshrc`, gopathBin)
	case strings.Contains(shell, "bash"), shell == "gitbash":
		return fmt.Sprintf(`Add the following line to your ~/.bashrc:
echo 'export PATH="$PATH:%s"' >> ~/.bashrc
source ~/.bashrc`, gopathBin)
	case strings.Contains(shell, "fish"):
		return fmt.Sprintf(`Add the following line to your ~/.config/fish/config.fish:
set -gx PATH $PATH %s`, gopathBin)
	case shell == "powershell":
		return fmt.Sprintf(`Run the following commands in PowerShell:
$env:Path += ";%s"
Set-ItemProperty -Path 'HKCU:\Environment' -Name Path -Value $env:Path`, gopathBin)
	default:
		return fmt.Sprintf(`Add %s to your PATH via System Properties.`, gopathBin)
	}
}

// appendToShellConfig appends PATH update to shell config file
func appendToShellConfig(shell, gopathBin string) error {
	var configFile string
	var command string
	switch {
	case strings.Contains(shell, "zsh"):
		configFile = filepath.Join(os.Getenv("HOME"), ".zshrc")
		command = fmt.Sprintf("\nexport PATH=\"$PATH:%s\"\n", gopathBin)
	case strings.Contains(shell, "bash"), shell == "gitbash":
		configFile = filepath.Join(os.Getenv("HOME"), ".bashrc")
		command = fmt.Sprintf("\nexport PATH=\"$PATH:%s\"\n", gopathBin)
	case strings.Contains(shell, "fish"):
		configFile = filepath.Join(os.Getenv("HOME"), ".config", "fish", "config.fish")
		command = fmt.Sprintf("\nset -gx PATH $PATH %s\n", gopathBin)
	default:
		return fmt.Errorf("automatic PATH update not supported for %s", shell)
	}

	f, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file %s: %w", configFile, err)
	}
	defer f.Close()

	_, err = f.WriteString(command)
	if err != nil {
		return fmt.Errorf("failed to write to config file %s: %w", configFile, err)
	}
	return nil
}

// confirm prompts the user for confirmation
func confirm(prompt string) bool {
	fmt.Printf("%s [y/N]: ", prompt)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

// IsFirstRun checks if this is the first run and creates config dir if needed
func IsFirstRun() bool {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Failed to get user config directory: %v\n", err)
		return true
	}
	nexlayerDir := filepath.Join(configDir, "nexlayer")
	if _, err := os.Stat(nexlayerDir); os.IsNotExist(err) {
		if err := os.MkdirAll(nexlayerDir, 0755); err != nil {
			fmt.Printf("Failed to create config directory: %v\n", err)
			return true
		}
		return true
	}
	return false
}
