// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ViperProvider is a Provider implementation that uses Viper
type ViperProvider struct {
	viper *viper.Viper
	mu    sync.RWMutex
}

// NewViperProvider creates a new ViperProvider
func NewViperProvider() *ViperProvider {
	return &ViperProvider{
		viper: viper.New(),
	}
}

// Get retrieves a configuration value by key
func (p *ViperProvider) Get(key string) interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.Get(key)
}

// GetString retrieves a string configuration value
func (p *ViperProvider) GetString(key string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetString(key)
}

// GetInt retrieves an integer configuration value
func (p *ViperProvider) GetInt(key string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetInt(key)
}

// GetBool retrieves a boolean configuration value
func (p *ViperProvider) GetBool(key string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetBool(key)
}

// GetFloat64 retrieves a float64 configuration value
func (p *ViperProvider) GetFloat64(key string) float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetFloat64(key)
}

// GetTime retrieves a time.Time configuration value
func (p *ViperProvider) GetTime(key string) time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetTime(key)
}

// GetDuration retrieves a time.Duration configuration value
func (p *ViperProvider) GetDuration(key string) time.Duration {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetDuration(key)
}

// GetStringSlice retrieves a string slice configuration value
func (p *ViperProvider) GetStringSlice(key string) []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetStringSlice(key)
}

// GetStringMap retrieves a map of string configuration values
func (p *ViperProvider) GetStringMap(key string) map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetStringMap(key)
}

// GetStringMapString retrieves a map of string configuration values
func (p *ViperProvider) GetStringMapString(key string) map[string]string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.GetStringMapString(key)
}

// IsSet checks if a configuration value is set
func (p *ViperProvider) IsSet(key string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.IsSet(key)
}

// SetDefault sets a default value for a configuration key
func (p *ViperProvider) SetDefault(key string, value interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.viper.SetDefault(key, value)
}

// Set sets a configuration value
func (p *ViperProvider) Set(key string, value interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.viper.Set(key, value)
}

// AllSettings returns all settings as a map
func (p *ViperProvider) AllSettings() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.AllSettings()
}

// ConfigFileUsed returns the config file used for loading configuration
func (p *ViperProvider) ConfigFileUsed() string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.viper.ConfigFileUsed()
}

// AddConfigPath adds a path for Viper to search for the config file in
func (p *ViperProvider) AddConfigPath(path string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.viper.AddConfigPath(path)
}

// SetConfigName sets the name of the config file
func (p *ViperProvider) SetConfigName(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.viper.SetConfigName(name)
}

// SetConfigType sets the type of the config file
func (p *ViperProvider) SetConfigType(configType string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.viper.SetConfigType(configType)
}

// SetConfigFile sets the config file explicitly
func (p *ViperProvider) SetConfigFile(file string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.viper.SetConfigFile(file)
}

// ReadInConfig reads in the config file
func (p *ViperProvider) ReadInConfig() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.viper.ReadInConfig()
}

// WriteConfig writes the current configuration to the config file
func (p *ViperProvider) WriteConfig() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.viper.WriteConfig()
}

// AutomaticEnv tells Viper to check for environment variables
func (p *ViperProvider) AutomaticEnv() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.viper.AutomaticEnv()
}

// SetEnvPrefix sets the prefix for environment variables
func (p *ViperProvider) SetEnvPrefix(prefix string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.viper.SetEnvPrefix(prefix)
}

// BindPFlag binds a flag to a configuration key
func (p *ViperProvider) BindPFlag(key string, flag *pflag.Flag) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.viper.BindPFlag(key, flag)
}
