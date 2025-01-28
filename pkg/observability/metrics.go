package observability

import (
	"sync"
	"time"
)

// MetricType represents the type of metric being collected
type MetricType int

const (
	COUNTER MetricType = iota
	GAUGE
	HISTOGRAM
)

// Metric represents a single metric
type Metric struct {
	Name        string
	Type        MetricType
	Value       float64
	Labels      map[string]string
	LastUpdated time.Time
}

// MetricsCollector collects and manages metrics
type MetricsCollector struct {
	mu      sync.RWMutex
	metrics map[string]*Metric
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make(map[string]*Metric),
	}
}

// Counter increments a counter metric
func (c *MetricsCollector) Counter(name string, value float64, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metric, exists := c.metrics[name]
	if !exists {
		metric = &Metric{
			Name:   name,
			Type:   COUNTER,
			Labels: labels,
		}
		c.metrics[name] = metric
	}

	metric.Value += value
	metric.LastUpdated = time.Now()
}

// Gauge sets a gauge metric
func (c *MetricsCollector) Gauge(name string, value float64, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics[name] = &Metric{
		Name:        name,
		Type:        GAUGE,
		Value:       value,
		Labels:      labels,
		LastUpdated: time.Now(),
	}
}

// Histogram adds a value to a histogram metric
func (c *MetricsCollector) Histogram(name string, value float64, labels map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metric, exists := c.metrics[name]
	if !exists {
		metric = &Metric{
			Name:   name,
			Type:   HISTOGRAM,
			Labels: labels,
		}
		c.metrics[name] = metric
	}

	// For simplicity, we're just storing the latest value
	// In a real implementation, we'd store a distribution of values
	metric.Value = value
	metric.LastUpdated = time.Now()
}

// GetMetric returns a metric by name
func (c *MetricsCollector) GetMetric(name string) *Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.metrics[name]
}

// GetAllMetrics returns all metrics
func (c *MetricsCollector) GetAllMetrics() []*Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics := make([]*Metric, 0, len(c.metrics))
	for _, metric := range c.metrics {
		metrics = append(metrics, metric)
	}
	return metrics
}

// Reset resets all metrics
func (c *MetricsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = make(map[string]*Metric)
}
