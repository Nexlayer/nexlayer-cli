package di

import (
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/config"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/observability"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

// Container handles dependency injection for the application
type Container struct {
	mu sync.RWMutex

	config    *config.Config
	apiClient api.APIClient
	uiManager ui.Manager
	logger    *observability.Logger
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return &Container{}
}

// GetAPIClient returns the API client instance
func (c *Container) GetAPIClient() api.APIClient {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.apiClient == nil {
		c.mu.RUnlock()
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.apiClient == nil {
			c.apiClient = api.NewClient(c.GetConfig().GetAPIEndpoint(""))
		}
	}

	return c.apiClient
}

// GetConfig returns the configuration instance
func (c *Container) GetConfig() *config.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.config == nil {
		c.mu.RUnlock()
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.config == nil {
			c.config = config.GetConfig()
		}
	}

	return c.config
}

// GetUIManager returns the UI manager instance
func (c *Container) GetUIManager() ui.Manager {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.uiManager == nil {
		c.mu.RUnlock()
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.uiManager == nil {
			c.uiManager = ui.NewManager()
		}
	}

	return c.uiManager
}

// GetLogger returns the logger instance
func (c *Container) GetLogger() *observability.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.logger == nil {
		c.mu.RUnlock()
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.logger == nil {
			c.logger = observability.NewLogger(observability.INFO)
		}
	}

	return c.logger
}

// GetMetricsCollector returns the metrics collector instance
func (c *Container) GetMetricsCollector() *observability.MetricsCollector {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.metricsCollector == nil {
		c.mu.RUnlock()
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.metricsCollector == nil {
			c.metricsCollector = observability.NewMetricsCollector()
		}
	}

	return c.metricsCollector
}
