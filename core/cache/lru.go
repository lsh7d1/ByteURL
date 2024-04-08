package cache

import (
	"container/list"
	"sync"
)

type (
	lru interface {
		Get(key string) (string, bool)
		Add(key, value string)
		Del(key string)
	}

	emptyLru struct{}

	entry struct {
		key, value string
	}

	LruCache struct {
		capacity int
		list     *list.List
		items    map[string]*list.Element
		lock     sync.RWMutex
		onEvict  func(string) // Callback when deleting
	}
)

var _ lru = (*emptyLru)(nil)
var _ lru = (*LruCache)(nil)

var emptyLruCache = emptyLru{}

func (e emptyLru) Get(string) (string, bool) { // do nothing
	return emptyString, false
}
func (e emptyLru) Add(string, string) {} // do nothing
func (e emptyLru) Del(string)         {} // do nothing

func NewLruCache(capacity int, onEvict func(string)) lru {
	return &LruCache{
		capacity,
		list.New(),
		make(map[string]*list.Element, capacity),
		sync.RWMutex{},
		onEvict,
	}
}

func (lru *LruCache) Get(key string) (string, bool) {
	lru.lock.RLock()
	defer lru.lock.RUnlock()
	if elem, ok := lru.getItem(key); ok {
		lru.list.MoveToFront(elem)
		return elem.Value.(entry).value, true
	}
	return emptyString, false
}

func (lru *LruCache) Add(key, value string) {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	if elem, ok := lru.items[key]; ok {
		lru.list.MoveToFront(elem)
		return
	}

	lru.items[key] = lru.list.PushFront(entry{key, value})
	if lru.list.Len() > lru.capacity {
		lru.removeOldest()
	}
}

func (lru *LruCache) Del(key string) {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	if elem, ok := lru.items[key]; ok {
		lru.removeElement(elem)
	}
}

func (lru *LruCache) getItem(key string) (*list.Element, bool) {
	// already locked
	if elem, ok := lru.items[key]; ok {
		return elem, true
	}
	return nil, false
}

// removeOldest delete the last item of lru
func (lru *LruCache) removeOldest() {
	elem := lru.list.Back()
	if elem != nil {
		lru.removeElement(elem)
	}
}

// removeElement really delete a [k, v] object, (list.element form)
func (lru *LruCache) removeElement(elem *list.Element) {
	lru.list.Remove(elem)
	key := elem.Value.(entry).key
	delete(lru.items, key)
	if lru.onEvict != nil {
		lru.onEvict(key)
	}
}

const (
	emptyString string = ""
)
