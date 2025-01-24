package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	instance *Cache
	once     sync.Once
)

// Cache represents a cache for storing deployment information
type Cache struct {
	cache *cache.Cache
}

// Get returns the singleton cache instance
func Get() *Cache {
	once.Do(func() {
		instance = &Cache{
			cache: cache.New(5*time.Minute, 10*time.Minute),
		}
	})
	return instance
}

// Set stores a value in the cache
func (c *Cache) Set(key string, value interface{}) {
	c.cache.Set(key, value, cache.DefaultExpiration)
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.cache.Delete(key)
}

// Clear removes all values from the cache
func (c *Cache) Clear() {
	c.cache.Flush()
}

// SaveToFile saves the cache to a file
func (c *Cache) SaveToFile(filename string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	items := c.cache.Items()
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	for key, item := range items {
		if err := enc.Encode(map[string]interface{}{
			"key":   key,
			"value": item.Object,
		}); err != nil {
			return fmt.Errorf("failed to encode cache item: %w", err)
		}
	}

	return nil
}

// LoadFromFile loads the cache from a file
func (c *Cache) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	for {
		var item map[string]interface{}
		if err := dec.Decode(&item); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to decode cache item: %w", err)
		}
		c.cache.Set(item["key"].(string), item["value"], cache.DefaultExpiration)
	}

	return nil
}
