package cache

import (
	"fmt"
	"reflect"
	"testing"
)

// type exHashLru struct {
// 	HashLru
// 	data map[string]string
// }

// func (h *exHashLru) Add(key string, value string) {
// 	h.data[key] = value
// }

func TestLruCacheAdd(t *testing.T) {
	h := NewHashLru(100, func(s string) {
		fmt.Println("onEvited")
	})

	h.Add("114514", "1919810")
	x, _ := h.Get("114514")
	fmt.Println(reflect.TypeOf(x))
	fmt.Printf("x: %#v\n", x)
}

func TestLruCacheDel(t *testing.T) {
	h := NewHashLru(100, func(s string) {
		fmt.Println(s)
	})
	h.Add("key1", "value1")
	h.Add("key2", "value22")

	h.Del("key3")
	h.Del("key2")
}
