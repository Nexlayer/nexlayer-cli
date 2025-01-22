package cache

import (
	"container/list"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slices"
)

const (
	defaultTTL = 5 * time.Minute
	cacheDir   = ".nexlayer/cache"
)

// CachedResponse represents a cached HTTP response
type CachedResponse struct {
	StatusCode int                 `json:"status_code"`
	Body       []byte              `json:"body"`
	Headers    map[string][]string `json:"headers"`
	CachedAt   time.Time           `json:"cached_at"`
	ExpiresAt  time.Time           `json:"expires_at"`
}

// OfflineMode represents the offline operation mode
type OfflineMode int

const (
	// OnlineMode normal operation with network access
	OnlineMode OfflineMode = iota
	// AutoOfflineMode automatically switches to offline mode when network is unavailable
	AutoOfflineMode
	// StrictOfflineMode forces offline-only operation
	StrictOfflineMode
)

// Operation represents a pending operation to be synced
type Operation struct {
	Method    string          `json:"method"`
	Endpoint  string          `json:"endpoint"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

// OfflineConfig stores offline operation settings
type OfflineConfig struct {
	Mode              OfflineMode   `json:"mode"`
	LastOnlineCheck   time.Time     `json:"last_online_check"`
	NetworkTimeout    time.Duration `json:"network_timeout"`
	AutoSyncInterval  time.Duration `json:"auto_sync_interval"`
	PendingOperations []Operation   `json:"pending_operations"`
	mu                sync.RWMutex
}

// CacheConfig holds configuration for the cache manager
type CacheConfig struct {
	MaxSize    int64         // Maximum size in bytes
	MaxEntries int           // Maximum number of entries
	TTL        time.Duration // Time-to-live for cache entries
}

// Manager handles caching operations
type Manager struct {
	baseDir       string
	config        CacheConfig
	mu            sync.RWMutex
	entries       map[string]*cacheEntry
	lru           *list.List
	currentSize   int64
	OfflineConfig *OfflineConfig
	baseURL       string
}

type cacheEntry struct {
	key        string
	size       int64
	lastUsed   time.Time
	lruElement *list.Element
	response   *CachedResponse
}

// NewManager creates a new cache manager with configuration
func NewManager(config CacheConfig) (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, cacheDir)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	offlineConfig := &OfflineConfig{
		Mode:             OnlineMode,
		NetworkTimeout:   5 * time.Second,
		AutoSyncInterval: 30 * time.Minute,
	}

	if err := offlineConfig.loadPendingOperations(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load pending operations: %v\n", err)
	}

	return &Manager{
		baseDir:       cacheDir,
		config:        config,
		entries:       make(map[string]*cacheEntry),
		lru:           list.New(),
		OfflineConfig: offlineConfig,
		baseURL:       "https://service.api.nexlayer.ai",
	}, nil
}

// IsOnline checks if network connectivity is available
func (c *OfflineConfig) IsOnline() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we've checked recently (within last minute)
	if time.Since(c.LastOnlineCheck) < time.Minute {
		return c.Mode == OnlineMode
	}

	// Try to connect to the API endpoint
	client := &http.Client{Timeout: c.NetworkTimeout}
	resp, err := client.Get("https://service.api.nexlayer.ai/health")
	c.LastOnlineCheck = time.Now()

	if err != nil || resp.StatusCode != http.StatusOK {
		if c.Mode == OnlineMode {
			c.Mode = AutoOfflineMode
		}
		return false
	}

	if c.Mode == AutoOfflineMode {
		c.Mode = OnlineMode
	}
	return true
}

// QueueOperation adds an operation to be executed when online
func (c *OfflineConfig) QueueOperation(method, endpoint string, payload []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	op := Operation{
		Method:    method,
		Endpoint:  endpoint,
		CreatedAt: time.Now(),
	}
	if payload != nil {
		op.Payload = json.RawMessage(payload)
	}
	c.PendingOperations = append(c.PendingOperations, op)
	c.savePendingOperations()
}

// GetPendingOperations returns all queued operations
func (c *OfflineConfig) GetPendingOperations() []Operation {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]Operation{}, c.PendingOperations...)
}

// ClearPendingOperations removes all queued operations
func (c *OfflineConfig) ClearPendingOperations() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.PendingOperations = nil
	c.savePendingOperations()
}

// savePendingOperations persists pending operations to disk
func (c *OfflineConfig) savePendingOperations() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(homeDir, ".nexlayer", "pending_operations.json")
	data, err := json.Marshal(c.PendingOperations)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// loadPendingOperations loads pending operations from disk
func (c *OfflineConfig) loadPendingOperations() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(homeDir, ".nexlayer", "pending_operations.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &c.PendingOperations)
}

// SyncPendingOperations attempts to execute queued operations
func (c *OfflineConfig) SyncPendingOperations(client *http.Client) error {
	if !c.IsOnline() {
		return fmt.Errorf("cannot sync while offline")
	}

	operations := c.GetPendingOperations()
	var failed []Operation

	for _, op := range operations {
		req, err := http.NewRequest(op.Method, "https://service.api.nexlayer.ai"+op.Endpoint, nil)
		if err != nil {
			failed = append(failed, op)
			continue
		}

		if len(op.Payload) > 0 {
			req.Body = http.NoBody
		}

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode >= 400 {
			failed = append(failed, op)
			continue
		}
	}

	c.mu.Lock()
	c.PendingOperations = failed
	c.mu.Unlock()

	return c.savePendingOperations()
}

// Set stores a response in cache with size management
func (m *Manager) Set(key string, statusCode int, body []byte, headers map[string][]string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	size := int64(len(body))

	// Check if we need to evict entries
	for m.currentSize+size > m.config.MaxSize || m.lru.Len() >= m.config.MaxEntries {
		if !m.evictOldest() {
			break
		}
	}

	entry := &cacheEntry{
		key:      key,
		size:     size,
		lastUsed: time.Now(),
		response: &CachedResponse{
			StatusCode: statusCode,
			Body:       body,
			Headers:    headers,
			CachedAt:   time.Now(),
		},
	}

	entry.lruElement = m.lru.PushFront(key)
	m.entries[key] = entry
	m.currentSize += size

	return m.persistEntry(key, entry)
}

// Get retrieves a cached response with LRU update
func (m *Manager) Get(key string) (*CachedResponse, bool, error) {
	m.mu.RLock()
	entry, exists := m.entries[key]
	if !exists {
		m.mu.RUnlock()
		return nil, false, nil
	}

	if time.Since(entry.lastUsed) > m.config.TTL {
		m.mu.RUnlock()
		m.Remove(key)
		return nil, false, nil
	}

	resp := entry.response
	m.mu.RUnlock()

	// Update LRU in a separate write lock
	m.mu.Lock()
	entry.lastUsed = time.Now()
	m.lru.MoveToFront(entry.lruElement)
	m.mu.Unlock()

	return resp, true, nil
}

// evictOldest removes the least recently used entry
func (m *Manager) evictOldest() bool {
	element := m.lru.Back()
	if element == nil {
		return false
	}

	key := element.Value.(string)
	entry := m.entries[key]

	m.lru.Remove(element)
	delete(m.entries, key)
	m.currentSize -= entry.size

	// Clean up the persisted cache file
	os.Remove(m.getCachePath(key))
	return true
}

// BatchGet retrieves multiple cached responses efficiently
func (m *Manager) BatchGet(keys []string) map[string]*CachedResponse {
	m.mu.RLock()
	results := make(map[string]*CachedResponse, len(keys))
	var expiredKeys []string

	for _, key := range keys {
		if entry, exists := m.entries[key]; exists {
			if time.Since(entry.lastUsed) <= m.config.TTL {
				results[key] = entry.response
			} else {
				expiredKeys = append(expiredKeys, key)
			}
		}
	}
	m.mu.RUnlock()

	// Clean up expired entries
	if len(expiredKeys) > 0 {
		m.mu.Lock()
		for _, key := range expiredKeys {
			m.Remove(key)
		}
		m.mu.Unlock()
	}

	return results
}

// Remove removes a cache entry
func (m *Manager) Remove(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entry, exists := m.entries[key]; exists {
		m.lru.Remove(entry.lruElement)
		delete(m.entries, key)
		m.currentSize -= entry.size
		return os.Remove(m.getCachePath(key))
	}
	return nil
}

// persistEntry persists a cache entry to disk
func (m *Manager) persistEntry(key string, entry *cacheEntry) error {
	data, err := json.Marshal(entry.response)
	if err != nil {
		return fmt.Errorf("failed to marshal cache response: %w", err)
	}

	path := m.getCachePath(key)
	return os.WriteFile(path, data, 0644)
}

// getCachePath returns the full path for a cache key
func (m *Manager) getCachePath(key string) string {
	return filepath.Join(m.baseDir, fmt.Sprintf("%x.json", key))
}

// ClearCache clears all cached items
func (m *Manager) ClearCache() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	dir, err := os.ReadDir(m.baseDir)
	if err != nil {
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	for _, entry := range dir {
		if err := os.Remove(filepath.Join(m.baseDir, entry.Name())); err != nil {
			return fmt.Errorf("failed to remove cache file %s: %w", entry.Name(), err)
		}
	}

	m.entries = make(map[string]*cacheEntry)
	m.lru = list.New()
	m.currentSize = 0

	return nil
}

// CacheInfo represents information about a cached item
type CacheInfo struct {
	Key  string `json:"key"`
	Size int    `json:"size"`
}

// ListCache returns information about all cached items
func (m *Manager) ListCache() ([]CacheInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var items []CacheInfo
	for _, entry := range m.entries {
		items = append(items, CacheInfo{
			Key:  entry.key,
			Size: int(entry.size),
		})
	}

	slices.SortFunc(items, func(a, b CacheInfo) int {
		return strings.Compare(a.Key, b.Key)
	})

	return items, nil
}

// GenerateCacheKey creates a unique key for caching
func (m *Manager) GenerateCacheKey(method, endpoint string) string {
	return fmt.Sprintf("%s:%s", method, endpoint)
}
