package cache

import (
	"crypto/md5"
	"math"
	"runtime"
)

// HashLru is LRU cache with hash sharding.
// Each shard performs its own concurrency control
type HashLru struct {
	lrus      []lru
	sliceSize int
	size      int
}

func NewHashLru(size int, onEvict func(string)) *HashLru {
	sliceSize := getNumLrus(runtime.NumCPU())
	capacity := int(math.Ceil(float64(size) / float64(sliceSize)))
	h := &HashLru{
		lrus:      make([]lru, sliceSize),
		sliceSize: sliceSize,
		size:      sliceSize * capacity,
	}
	for i := 0; i < h.sliceSize; i++ {
		h.lrus[i] = NewLruCache(capacity, onEvict)
	}
	return h
}

func (h *HashLru) Get(key string) (string, bool) {
	pos := h.getPos(key)
	return h.lrus[pos].Get(key)
}

func (h *HashLru) Add(key, value string) {
	pos := h.getPos(key)
	h.lrus[pos].Add(key, value)
}

func (h *HashLru) Del(key string) {
	pos := h.getPos(key)
	h.lrus[pos].Del(key)
}

func (h *HashLru) getPos(key string) int {
	x := int(md5.Sum([]byte(key))[0])
	return x % h.sliceSize
}

func getNumLrus(cores int) int {
	switch {
	case CPU64 <= cores:
		return CPU64
	case CPU32 <= cores:
		return CPU32
	case CPU16 <= cores:
		return CPU16
	case CPU8 <= cores:
		return CPU8
	case CPU4 <= cores:
		return CPU4
	case CPU2 <= cores:
		return CPU2
	default:
		return CPU1
	}
}

const (
	CPU1 = 1 << iota
	CPU2
	CPU4
	CPU8
	CPU16
	CPU32
	CPU64
)
