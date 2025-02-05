// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package sysinfo

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// SystemInfo holds information about the user's system and deployment
type SystemInfo struct {
	OS            string    `json:"os"`
	OSVersion     string    `json:"os_version"`
	Architecture  string    `json:"architecture"`
	IDE           string    `json:"ide"`
	AIModel       string    `json:"ai_model"`
	DeploymentURL string    `json:"deployment_url"`
	Timestamp     time.Time `json:"timestamp"`
}

// GetSystemInfo gathers information about the user's system
func GetSystemInfo() *SystemInfo {
	info := &SystemInfo{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		Timestamp:    time.Now(),
	}

	// Get OS version
	switch runtime.GOOS {
	case "darwin":
		if out, err := os.ReadFile("/System/Library/CoreServices/SystemVersion.plist"); err == nil {
			info.OSVersion = parseOSVersion(string(out))
		}
	case "linux":
		if out, err := os.ReadFile("/etc/os-release"); err == nil {
			info.OSVersion = parseOSVersion(string(out))
		}
	case "windows":
		info.OSVersion = os.Getenv("OS")
	}

	// Get IDE information from environment variables
	info.IDE = os.Getenv("NEXLAYER_IDE")
	info.AIModel = os.Getenv("NEXLAYER_AI_MODEL")

	return info
}

// parseOSVersion extracts the OS version from system files
func parseOSVersion(content string) string {
	// For macOS
	if strings.Contains(content, "ProductVersion") {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if strings.Contains(line, "ProductVersion") {
				parts := strings.Split(line, "<string>")
				if len(parts) > 1 {
					version := strings.Split(parts[1], "</string>")[0]
					return strings.TrimSpace(version)
				}
			}
		}
	}

	// For Linux
	if strings.Contains(content, "VERSION_ID") {
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "VERSION_ID=") {
				version := strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
				return version
			}
		}
	}

	return "unknown"
}

// FormatFeedback formats system info into a feedback message
func (si *SystemInfo) FormatFeedback(deploymentName string) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Deployment completed: %s\n", deploymentName))
	if si.DeploymentURL != "" {
		builder.WriteString(fmt.Sprintf("URL: %s\n", si.DeploymentURL))
	}
	builder.WriteString(fmt.Sprintf("Time: %s\n", si.Timestamp.Format(time.RFC3339)))
	builder.WriteString(fmt.Sprintf("IDE: %s\n", si.IDE))
	builder.WriteString(fmt.Sprintf("AI Model: %s\n", si.AIModel))
	builder.WriteString(fmt.Sprintf("System: %s %s (%s)\n", si.OS, si.OSVersion, si.Architecture))

	return builder.String()
}
