package component

import (
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheGet(t *testing.T) {
	cache := NewCache("myCache", func(s string) {})
	cache.Set("key", "value")
	value, ok := cache.Get("key")
	assert.Equal(t, "value", value)
	assert.Equal(t, true, ok)
}

func BenchmarkCache(b *testing.B) {
	cache := NewCache("myCache", func(s string) {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		num := rand.Int()
		cache.Set(strconv.Itoa(num), num)
		cache.Get(strconv.Itoa(rand.Int()))
	}
}

func TestCacheTake(t *testing.T) {
	cache := NewCache("myCache", func(s string) {})

	const n int = 100
	var wg sync.WaitGroup
	var count int32 = 0
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			val, err := cache.Take("key", func() (any, error) {
				atomic.AddInt32(&count, 1)
				time.Sleep(time.Millisecond * 100)
				return "value", nil
			})
			assert.Equal(t, "value", val)
			assert.Nil(t, err)
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(1), count)
}
