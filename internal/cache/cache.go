/*
 license x
*/

package cache

import (
	"sync"
	"time"

	"github.com/nelsonstr/o801/internal/interfaces"
)

// Cache represents a KV cache with expiration time.
type Cache[T interfaces.Cacheable] struct {
	mu   sync.RWMutex
	data map[int64]cacheItem[T]
	ttl  time.Duration
}

type cacheItem[T interfaces.Cacheable] struct {
	value      T
	expiration time.Time
}

// NewCache creates a new cache with the specified ttl.
func NewCache[T interfaces.Cacheable](ttl time.Duration) *Cache[T] {
	cache := &Cache[T]{
		data: make(map[int64]cacheItem[T]),
		ttl:  ttl,
	}

	go cache.startCleanup()

	return cache
}

// Set adds a KV to the cache.
func (c *Cache[T]) Set(key int64, value *T) {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.data[key] = cacheItem[T]{
		value:      *value,
		expiration: time.Now().Add(c.ttl),
	}
}

// Get retrieves the value associated with the given key from the cache.
// returns the generic type T.
func (c *Cache[T]) Get(key int64) (T, bool) {
	c.mu.RLock()

	defer c.mu.RUnlock()

	item, exists := c.data[key]
	if !exists {
		var nilT T

		return nilT, false
	}

	return item.value, true
}

// Len returns the current number of items in the cache.
func (c *Cache[T]) Len() int {
	c.mu.RLock()

	defer c.mu.RUnlock()

	return len(c.data)
}

func (c *Cache[T]) delete(key int64) {
	delete(c.data, key)
}

func (c *Cache[T]) startCleanup() {
	ticker := time.NewTicker(c.ttl)

	defer ticker.Stop()

	for range ticker.C {
		c.cleanup()
	}
}

func (c *Cache[T]) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.data {
		if item.expiration.Before(now) {
			c.delete(key)
		}
	}
}
