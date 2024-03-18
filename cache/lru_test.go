package cache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLruCacheGet(t *testing.T) {
	lru := NewLruCache(114514, nil)

	lru.Add("key1", "value1")
	lru.Add("key2", "value2")
	lru.Add("key3", "value3")

	lru.Del("key2")

	res, ok := lru.Get("key1")
	assert.Equal(t, "value1", res)
	assert.True(t, ok)

	res, ok = lru.Get("key3")
	assert.Equal(t, "value3", res)
	assert.True(t, ok)

	res, ok = lru.Get("key2")
	assert.Equal(t, "", res)
	assert.False(t, ok)
}

func TestLruCacheonEvited(t *testing.T) {
	lru := NewLruCache(114514, func(s string) {
		fmt.Printf("key %s is deleted...\n", s)
	})

	lru.Add("key1", "value1")
	lru.Del("key1")
}
