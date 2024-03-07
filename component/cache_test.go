package component

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestCache(t *testing.T) {
	cache := NewCache("myCache", func(s string) {})
	for i := 0; i < 1000; i++ {
		num := rand.Int()
		go cache.Set(strconv.Itoa(num), num)
		go cache.Get(strconv.Itoa(rand.Int()))
	}
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
