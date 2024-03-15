package cache

import (
	"container/list"
	"sync"
)

type (
	lru interface {
		Add(key string)
		Del(key string)
	}

	emptyLru struct{}

	lruCache struct {
		capacity int
		list     *list.List
		items    map[string]*list.Element
		lock     sync.RWMutex
		onEvict  func(string) // Callback when deleting
	}
)

var _ lru = (*emptyLru)(nil)
var _ lru = (*lruCache)(nil)

var emptyLruCache = emptyLru{}

func (e emptyLru) Get(string) (any, bool) { // do nothing
	return nil, false
}
func (e emptyLru) Add(string) {} // do nothing
func (e emptyLru) Del(string) {} // do nothing

func newLruCache(capacity int, onEvict func(string)) lru {
	return &lruCache{
		capacity,
		list.New(),
		make(map[string]*list.Element, capacity),
		sync.RWMutex{},
		onEvict,
	}
}

func (lru *lruCache) Get(key string) (any, bool) {
	lru.lock.RLock()
	defer lru.lock.RUnlock()
	if elem, ok := lru.getItem(key); ok {
		lru.list.MoveToFront(elem)
		return elem.Value, true
	}
	return nil, false
}

func (lru *lruCache) Add(key string) {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	if elem, ok := lru.items[key]; ok {
		lru.list.MoveToFront(elem)
		return
	}

	lru.items[key] = lru.list.PushFront(key)
	if lru.list.Len() > lru.capacity {
		lru.removeOldest()
	}
}

func (lru *lruCache) Del(key string) {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	if elem, ok := lru.items[key]; ok {
		lru.removeElement(elem)
	}
}

func (lru *lruCache) getItem(key string) (*list.Element, bool) {
	// already locked
	if elem, ok := lru.items[key]; ok {
		return elem, true
	}
	return nil, false
}

// removeOldest delete the last item of lru
func (lru *lruCache) removeOldest() {
	elem := lru.list.Back()
	if elem != nil {
		lru.removeElement(elem)
	}
}

func (lru *lruCache) removeElement(elem *list.Element) {
	lru.list.Remove(elem)
	key := elem.Value.(string)
	delete(lru.items, key)
	lru.onEvict(key)
}
