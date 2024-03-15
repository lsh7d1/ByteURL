package cache

import (
	"crypto/md5"
	"math"
	"runtime"
)

type HashLru struct {
	lrus      []lru
	sliceSize int
	size      int
}

func NewHashLru(size int, onEvict func(string)) *HashLru {
	sliceSize := getNumLrus(runtime.NumCPU())
	capacity := int(math.Ceil(float64(size) / float64(sliceSize)))
	h := &HashLru{
		lrus: make([]lru, 0),
	}
	for i := 0; i < h.sliceSize; i++ {
		h.lrus[i] = newLruCache(capacity, onEvict)
	}
	return h
}

func (h *HashLru) Add(key string) {
	pos := h.getPos(key)
	one := h.lrus[pos]
	one.Add(key)
}

func (h *HashLru) Del(key string) {
	pos := h.getPos(key)
	one := h.lrus[pos]
	one.Del(key)
}

func (h *HashLru) getPos(key string) int {
	x := int(md5.Sum([]byte(key))[0])
	return x % h.sliceSize
}

func getNumLrus(size int) int {
	switch {
	case CPU64 <= size:
		return CPU64
	case CPU32 <= size:
		return CPU32
	case CPU16 <= size:
		return CPU16
	case CPU8 <= size:
		return CPU8
	case CPU4 <= size:
		return CPU4
	case CPU2 <= size:
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
