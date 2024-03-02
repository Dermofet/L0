package cache

import (
	"L0/internal/repository"
	"context"
)

//go:generate mockgen -source=./interfaces.go -destination=cache_mock.go -package=cache

// Cache represents a structure for caching data.
type Cache interface {
	// Set sets a value in the cache for the specified key.
	// It takes a key and a value as input parameters.
	// Returns nothing.
	Set(key string, value interface{})

	// Get returns the value from the cache for the specified key.
	// It takes a key as input parameter.
	// Returns the value and a boolean indicating whether the key was found in the cache.
	Get(key string) (interface{}, bool)

	// GetAll returns all values in the cache.
	// Returns a slice of values.
	// Returns nothing.
	GetAll() ([]interface{}, bool)

	// Delete deletes a value from the cache for the specified key.
	// It takes a key as input parameter.
	// Returns nothing.
	Delete(key string)

	// Load loads data into the cache from the repository.
	// It takes a context and an OrderRepository as input parameters.
	// Returns an error if the operation fails.
	Load(ctx context.Context, orderRepository repository.OrderRepository) error
}
