// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detection

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
)

// AIAssistantConfig defines the configuration for detecting an AI assistant
type AIAssistantConfig struct {
	Name            string
	ConfigPaths     []string
	EnvVars         []string
	DefaultModel    string
	DetectionMethod func() (string, string)
}

var (
	// assistantConfigs maps IDE names to their detection configurations
	assistantConfigs = map[string]AIAssistantConfig{
		"Cursor": {
			Name: "Cursor",
			ConfigPaths: []string{
				"Library/Application Support/Cursor/User/settings.json",
				".config/Cursor/User/settings.json",
			},
			EnvVars:      []string{"CURSOR_TRACE_ID", "CURSOR_LLM_MODEL"},
			DefaultModel: "claude-3-sonnet",
		},
		"VSCode": {
			Name:        "VSCode",
			ConfigPaths: []string{".vscode/extensions"},
			EnvVars:     []string{"VSCODE_GIT_IPC_HANDLE", "VSCODE_PID"},
		},
		"Windsurf": {
			Name:         "Windsurf",
			ConfigPaths:  []string{".windsurf/config.json"},
			EnvVars:      []string{"WINDSURF", "WINDSURF_LLM"},
			DefaultModel: "gpt-4",
		},
		"Zed": {
			Name:         "Zed",
			ConfigPaths:  []string{".zed/settings.json"},
			EnvVars:      []string{"ZED_ROOT", "ZED_LLM"},
			DefaultModel: "gpt-4",
		},
	}

	// Cache for detection results
	detectionCache struct {
		assistant string
		model     string
		timestamp time.Time
		sync.RWMutex
	}
	cacheDuration = 5 * time.Minute
)

// LLMDetector detects AI-powered IDEs and LLM-based coding assistants
type LLMDetector struct{}

func (d *LLMDetector) Priority() int { return 250 }

func (d *LLMDetector) Detect(dir string) (*types.ProjectInfo, error) {
	assistant, model := detectAIAssistantAndModel()
	if assistant == "Unknown" {
		return nil, nil
	}
	return &types.ProjectInfo{
		Type:        types.TypeUnknown,
		LLMProvider: assistant,
		LLMModel:    model,
	}, nil
}

// detectAIAssistantAndModel detects the AI assistant and its model with caching
func detectAIAssistantAndModel() (string, string) {
	detectionCache.RLock()
	if time.Since(detectionCache.timestamp) < cacheDuration {
		assistant, model := detectionCache.assistant, detectionCache.model
		detectionCache.RUnlock()
		return assistant, model
	}
	detectionCache.RUnlock()

	assistant, model := performDetection()

	detectionCache.Lock()
	detectionCache.assistant = assistant
	detectionCache.model = model
	detectionCache.timestamp = time.Now()
	detectionCache.Unlock()

	return assistant, model
}

// performDetection performs the actual detection of AI assistant and model
func performDetection() (string, string) {
	// First check running processes
	if processes := detectRunningAIProcesses(); len(processes) > 0 {
		for _, process := range processes {
			if assistant, model := matchProcessToAssistant(process); assistant != "Unknown" {
				return assistant, model
			}
		}
	}

	// Check environment variables
	for name, config := range assistantConfigs {
		for _, envVar := range config.EnvVars {
			if os.Getenv(envVar) != "" {
				return name, getAssistantModel(name)
			}
		}
	}

	// Check configuration files
	if assistant, model := detectFromConfigFiles(); assistant != "Unknown" {
		return assistant, model
	}

	return "Unknown", "Unknown"
}

// detectRunningAIProcesses detects running AI assistant processes
func detectRunningAIProcesses() []string {
	var processes []string

	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("ps", "-ax", "-o", "comm=")
		output, err := cmd.Output()
		if err == nil {
			for _, line := range strings.Split(string(output), "\n") {
				if strings.Contains(line, "Cursor.app") ||
					strings.Contains(line, "Code.app") ||
					strings.Contains(line, "Zed") {
					processes = append(processes, line)
				}
			}
		}
	case "linux":
		cmd := exec.Command("ps", "-e", "-o", "comm=")
		output, err := cmd.Output()
		if err == nil {
			for _, line := range strings.Split(string(output), "\n") {
				if strings.Contains(line, "cursor") ||
					strings.Contains(line, "code") ||
					strings.Contains(line, "zed") {
					processes = append(processes, line)
				}
			}
		}
	case "windows":
		cmd := exec.Command("tasklist", "/FO", "CSV")
		output, err := cmd.Output()
		if err == nil {
			for _, line := range strings.Split(string(output), "\n") {
				if strings.Contains(strings.ToLower(line), "cursor.exe") ||
					strings.Contains(strings.ToLower(line), "code.exe") ||
					strings.Contains(strings.ToLower(line), "zed.exe") {
					processes = append(processes, line)
				}
			}
		}
	}

	return processes
}

// matchProcessToAssistant matches a process name to an AI assistant
func matchProcessToAssistant(process string) (string, string) {
	process = strings.ToLower(process)
	if strings.Contains(process, "cursor") {
		return "Cursor", getAssistantModel("Cursor")
	}
	if strings.Contains(process, "code") && !strings.Contains(process, "cursor") {
		return "VSCode", getAssistantModel("VSCode")
	}
	if strings.Contains(process, "zed") {
		return "Zed", getAssistantModel("Zed")
	}
	return "Unknown", "Unknown"
}

// detectFromConfigFiles checks for AI assistant configuration files
func detectFromConfigFiles() (string, string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "Unknown", "Unknown"
	}

	for name, config := range assistantConfigs {
		for _, configPath := range config.ConfigPaths {
			path := filepath.Join(home, configPath)
			if _, err := os.Stat(path); err == nil {
				return name, getAssistantModel(name)
			}
		}
	}

	return "Unknown", "Unknown"
}

// getAssistantModel gets the model for a specific AI assistant
func getAssistantModel(assistant string) string {
	config, ok := assistantConfigs[assistant]
	if !ok {
		return "Unknown"
	}

	// Check environment variables first
	for _, envVar := range config.EnvVars {
		if strings.Contains(envVar, "MODEL") || strings.Contains(envVar, "LLM") {
			if model := os.Getenv(envVar); model != "" {
				return model
			}
		}
	}

	// Check configuration files
	if model := getModelFromConfig(assistant); model != "" {
		return model
	}

	// Return default model if available
	if config.DefaultModel != "" {
		return config.DefaultModel
	}

	return "Unknown"
}

// getModelFromConfig reads the model from assistant's configuration file
func getModelFromConfig(assistant string) string {
	config, ok := assistantConfigs[assistant]
	if !ok {
		return ""
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	for _, configPath := range config.ConfigPaths {
		path := filepath.Join(home, configPath)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var settings map[string]interface{}
		if err := json.Unmarshal(data, &settings); err != nil {
			continue
		}

		// Check common model setting keys
		for _, key := range []string{
			"llmModel",
			"ai.model",
			"model",
			assistant + ".model",
			"ai.llm.model",
		} {
			if model, ok := settings[key].(string); ok {
				return model
			}
		}
	}

	return ""
}
