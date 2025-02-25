// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package config provides a centralized configuration management system for the Nexlayer CLI.
// It supports loading configuration from multiple sources (files, environment variables, flags)
// and provides a unified interface for accessing configuration values.
package config

import (
	"time"
)

// Provider defines the interface for accessing configuration values
type Provider interface {
	// Get retrieves a configuration value by key
	Get(key string) interface{}

	// GetString retrieves a string configuration value
	GetString(key string) string

	// GetInt retrieves an integer configuration value
	GetInt(key string) int

	// GetBool retrieves a boolean configuration value
	GetBool(key string) bool

	// GetFloat64 retrieves a float64 configuration value
	GetFloat64(key string) float64

	// GetTime retrieves a time.Time configuration value
	GetTime(key string) time.Time

	// GetDuration retrieves a time.Duration configuration value
	GetDuration(key string) time.Duration

	// GetStringSlice retrieves a string slice configuration value
	GetStringSlice(key string) []string

	// GetStringMap retrieves a map of string configuration values
	GetStringMap(key string) map[string]interface{}

	// GetStringMapString retrieves a map of string configuration values
	GetStringMapString(key string) map[string]string

	// IsSet checks if a configuration value is set
	IsSet(key string) bool

	// SetDefault sets a default value for a configuration key
	SetDefault(key string, value interface{})

	// Set sets a configuration value
	Set(key string, value interface{})

	// AllSettings returns all settings as a map
	AllSettings() map[string]interface{}

	// ConfigFileUsed returns the config file used for loading configuration
	ConfigFileUsed() string
}

// ProviderKey is the context key for storing the configuration provider
type ProviderKey struct{}

// DefaultProviderKey is the default context key for the configuration provider
var DefaultProviderKey = ProviderKey{}
