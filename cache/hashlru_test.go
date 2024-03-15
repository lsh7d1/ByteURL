package cache

import "testing"

type exHashLru struct {
	HashLru
	data map[string]string
}

func (h *exHashLru) Add(key string, value string) {
	h.data[key] = value
}

func TestLruCacheAdd(t *testing.T) {

}
