package cache

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
=== RUN   TestCacheGet
--- PASS: TestCacheGet (0.00s)
=== RUN   TestCacheDel
Delete ~~!!!!
--- PASS: TestCacheDel (0.00s)
=== RUN   TestCacheTake
Delete ~~!!!!
--- PASS: TestCacheTake (1.01s)
=== RUN   TestCacheTakeWithError
Delete ~~!!!!
--- PASS: TestCacheTakeWithError (1.01s)
=== RUN   TestCacheWithLimit
--- PASS: TestCacheWithLimit (0.00s)
PASS
ok      command-line-arguments  2.135s
*/

func TestCacheGet(t *testing.T) {
	cache := NewCache("myCache", time.Second)
	cache.Set("key", "value")
	cache.SetWithExpire("key2", "value2", time.Second*3)

	value, ok := cache.Get("key")
	assert.Equal(t, "value", value)
	assert.True(t, ok)

	value, ok = cache.Get("key2")
	assert.Equal(t, "value2", value)
	assert.True(t, ok)
}

func TestCacheDel(t *testing.T) {
	cache := NewCache("myCache", time.Second*2)

	cache.Set("key", "value")
	cache.Set("key2", "value2")
	cache.Del("key")

	_, ok := cache.Get("key")
	assert.False(t, ok)

	value2, ok := cache.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "value2", value2)
}

func TestCacheTake(t *testing.T) {
	const n int = 1000
	cache := NewCache("myCache", time.Minute, WithAroundCapLimit(n))
	var wg sync.WaitGroup
	var count int32 = 0
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			val, err := cache.Take("key1", func() (string, error) {
				atomic.AddInt32(&count, 1)
				time.Sleep(time.Second)
				return "value1", nil
			})
			assert.Equal(t, "value1", val)
			assert.Nil(t, err)
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(1), count)
}

func TestCacheTakeWithError(t *testing.T) {
	const n int = 100
	cache := NewCache("myCache", time.Minute, WithAroundCapLimit(n))
	var count int32
	var wg sync.WaitGroup
	errSame := errors.New("SameError")
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			val, err := cache.Take("key", func() (string, error) {
				atomic.AddInt32(&count, 1)
				time.Sleep(time.Second)
				return "", errSame
			})
			assert.Nil(t, val)
			assert.Equal(t, errSame, err)
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, int32(1), count)
}

func TestCacheWithLimit(t *testing.T) {
	cache := NewCache("myCache", time.Second*3, WithCapLimit(3))
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")

	// Set key-value pairs that exceed capacity limits
	_, ok := cache.Get("key1")
	assert.False(t, ok)
	value2, ok := cache.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "value2", value2)
	value3, ok := cache.Get("key3")
	assert.True(t, ok)
	assert.Equal(t, "value3", value3)
	value4, ok := cache.Get("key4")
	assert.True(t, ok)
	assert.Equal(t, "value4", value4)

	// Refresh the position of key2 in lru
	// and add key5, key6
	_, _ = cache.Get("key2")
	cache.Set("key5", "value5")
	cache.Set("key6", "value6")

	value2, ok = cache.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "value2", value2)
	value5, ok := cache.Get("key5")
	assert.True(t, ok)
	assert.Equal(t, "value5", value5)
	value6, ok := cache.Get("key6")
	assert.True(t, ok)
	assert.Equal(t, "value6", value6)
}
