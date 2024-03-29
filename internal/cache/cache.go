// Package cache provides functionality for caching data.
package cache

import (
	"context"
	"fmt"
	"sync"

	"L0/internal/repository"
)

// Cache represents a structure for caching data.
type cache struct {
	data  map[string]interface{}
	mutex sync.RWMutex
}

// NewCache creates a new instance of Cache.
func NewCache() *cache {
	return &cache{
		data: make(map[string]interface{}),
	}
}

// Set sets a value in the cache for the specified key.
func (c *cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

// Get returns the value from the cache for the specified key.
func (c *cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

// GetAll returns all values in the cache.
func (c *cache) GetAll() ([]interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	values := make([]interface{}, 0, len(c.data))
	for _, value := range c.data {
		values = append(values, value)
	}
	return values, len(values) > 0
}

// Delete deletes a value from the cache for the specified key.
func (c *cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}

// Load loads data into the cache from the repository.
func (c *cache) Load(ctx context.Context, orderRepository repository.OrderRepository) error {
	orders, err := orderRepository.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("can't get orders from database: %w", err)
	}

	for _, order := range orders {
		c.Set(order.OrderUID, order)
	}

	return nil
}
