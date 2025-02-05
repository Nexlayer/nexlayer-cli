package di

import (
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/config"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/observability"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

// Container handles dependency injection for the application.
type Container struct {
	mu        sync.RWMutex
	config    *config.Config
	apiClient api.APIClient
	uiManager ui.Manager
	logger    *observability.Logger
}

// NewContainer creates a new dependency injection container.
func NewContainer() *Container {
	return &Container{}
}

// GetAPIClient returns the API client instance.
func (c *Container) GetAPIClient() api.APIClient {
	// First, acquire a read lock.
	c.mu.RLock()
	if c.apiClient != nil {
		defer c.mu.RUnlock()
		return c.apiClient
	}
	c.mu.RUnlock()

	// Acquire a write lock and double-check.
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.apiClient == nil {
		c.apiClient = api.NewClient(c.GetConfig().GetAPIEndpoint(""))
	}
	return c.apiClient
}

// GetConfig returns the configuration instance.
func (c *Container) GetConfig() *config.Config {
	c.mu.RLock()
	if c.config != nil {
		defer c.mu.RUnlock()
		return c.config
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.config == nil {
		c.config = config.GetConfig()
	}
	return c.config
}

// GetUIManager returns the UI manager instance.
func (c *Container) GetUIManager() ui.Manager {
	c.mu.RLock()
	if c.uiManager != nil {
		defer c.mu.RUnlock()
		return c.uiManager
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.uiManager == nil {
		c.uiManager = ui.NewManager()
	}
	return c.uiManager
}

// GetLogger returns the logger instance.
func (c *Container) GetLogger() *observability.Logger {
	c.mu.RLock()
	if c.logger != nil {
		defer c.mu.RUnlock()
		return c.logger
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.logger == nil {
		c.logger = observability.NewLogger(observability.INFO)
	}
	return c.logger
}


