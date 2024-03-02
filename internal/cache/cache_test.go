package cache

import (
	"testing"
)

// TestCache_Set проверяет функцию Set для структуры cache.
func TestCache_Set(t *testing.T) {
	c := NewCache()

	c.Set("key", "value")

	val, ok := c.Get("key")

	if !ok || val != "value" {
		t.Errorf("Cache set value = %v, want value = value", val)
	}
}

// TestCache_Get проверяет функцию Get для структуры cache.
func TestCache_Get(t *testing.T) {
	c := NewCache()

	c.Set("key", "value")

	val, ok := c.Get("key")

	if !ok || val != "value" {
		t.Errorf("Cache get value = %v, want value = value", val)
	}
}

// TestCache_GetAll проверяет функцию GetAll для структуры cache.
func TestCache_GetAll(t *testing.T) {
	c := NewCache()

	c.Set("key", "value")

	val, ok := c.GetAll()
	if !ok {
		t.Errorf("Cache get value = %v, want value = value", val)
	}

	if len(val) != 1 {
		t.Errorf("Cache get value = %v, want value = value", val)
	}

	if val[0] != "value" {
		t.Errorf("Cache get value = %v, want value = value", val)
	}
}

// TestCache_Delete проверяет функцию Delete для структуры cache.
func TestCache_Delete(t *testing.T) {
	c := NewCache()

	c.Set("key", "value")

	c.Delete("key")

	_, ok := c.Get("key")
	if ok {
		t.Errorf("Cache value was not deleted")
	}
}
