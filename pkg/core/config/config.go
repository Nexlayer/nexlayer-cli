// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package config provides a centralized configuration management system for the Nexlayer CLI.
// It supports loading configuration from multiple sources (files, environment variables, flags)
// and provides a unified interface for accessing configuration values.
package config

import (
	"context"
	"sync"
)

var (
	// defaultManager is the default configuration manager
	defaultManager *Manager
	// managerMu protects defaultManager
	managerMu sync.RWMutex
)

// init initializes the default configuration manager
func init() {
	defaultManager = DefaultManager()
}

// GetDefaultManager returns the default configuration manager
func GetDefaultManager() *Manager {
	managerMu.RLock()
	defer managerMu.RUnlock()
	return defaultManager
}

// SetDefaultManager sets the default configuration manager
func SetDefaultManager(manager *Manager) {
	managerMu.Lock()
	defer managerMu.Unlock()
	defaultManager = manager
}

// InitConfig initializes the configuration with default values and paths
func InitConfig(configPath string) error {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	return manager.InitConfig(configPath)
}

// WithContext returns a new context with the configuration provider
func WithContext(ctx context.Context) context.Context {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	return manager.WithContext(ctx)
}

// GetAPIURL returns the API URL from the configuration
func GetAPIURL() string {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	return manager.GetAPIURL()
}

// SetAPIURL sets the API URL in the configuration
func SetAPIURL(url string) {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	manager.SetAPIURL(url)
}

// GetToken returns the authentication token from the configuration
func GetToken() string {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	return manager.GetToken()
}

// SetToken sets the authentication token in the configuration
func SetToken(token string) {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	manager.SetToken(token)
}

// GetDefaultNamespace returns the default namespace from the configuration
func GetDefaultNamespace() string {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	return manager.GetDefaultNamespace()
}

// SetDefaultNamespace sets the default namespace in the configuration
func SetDefaultNamespace(namespace string) {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	manager.SetDefaultNamespace(namespace)
}

// SaveConfig saves the current configuration to disk
func SaveConfig() error {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	return manager.SaveConfig()
}

// GetConfigDir returns the directory where the configuration file is located
func GetConfigDir() (string, error) {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	return manager.GetConfigDir()
}

// GetConfigProvider returns the configuration provider
func GetConfigProvider() Provider {
	managerMu.RLock()
	manager := defaultManager
	managerMu.RUnlock()
	return manager.Provider()
}
