/*
 license x
*/

package cache

import (
	userModel "github.com/nelsonstr/o801/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCache_Set_Get_Expired(t *testing.T) {
	ttl := 100 * time.Millisecond
	cache := NewCache[userModel.UserView](ttl)

	key := int64(1)
	value := &userModel.UserView{Name: "nelson"}

	// Test Set and Read methods.
	cache.Set(key, value)
	cachedValue, exists := cache.Get(key)

	assert.True(t, exists)
	assert.Equal(t, value.Name, cachedValue.Name)

	// Sleep for a duration longer than the TTL to ensure item expires.
	time.Sleep(2 * ttl)

	// Test item expiration.
	cachedValue, exists = cache.Get(key)
	assert.False(t, exists)

	assert.Equal(t, 0, cache.Len())
}

func TestCache_Delete(t *testing.T) {
	ttl := 100 * time.Millisecond
	cache := NewCache[userModel.UserView](ttl)

	key := int64(1)
	value := &userModel.UserView{Name: "nelson"}

	// Test Set and Read methods.
	cache.Set(key, value)
	cachedValue, exists := cache.Get(key)
	assert.True(t, exists)
	assert.Equal(t, value.Name, cachedValue.Name)

	// Test Delete method.
	cache.delete(key)
	cachedValue, exists = cache.Get(key)
	assert.False(t, exists)

	assert.Equal(t, 0, cache.Len())
}
