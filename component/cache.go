package component

import (
	"container/list"
	"sync"
)

type (
	Cache struct {
		name string
		lock sync.Mutex
		data map[string]any
		lru  lru
	}
)

func NewCache(name string, onEvict func(string)) *Cache {
	cache := &Cache{
		name: name,
		data: make(map[string]any),
		lru:  emptyLruCache,
	}

	return cache
}

func (c *Cache) Del(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.lru.del(key)
	delete(c.data, key)
}

func (c *Cache) Get(key string) (any, bool) {
	return c.get(key)
}

func (c *Cache) Set(key string, val any) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[key] = val
	c.lru.add(key)
}

func (c *Cache) get(key string) (any, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	val, ok := c.data[key]
	if ok {
		c.lru.add(key)
	}
	return val, ok
}

func (c *Cache) size() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

type (
	lru interface {
		add(key string)
		del(key string)
	}

	emptyLru struct{}

	lruCache struct {
		capacity int
		list     *list.List
		evicts   map[string]*list.Element
		onEvict  func(string) // 删除时的回调
	}
)

var _ lru = (*emptyLru)(nil)
var _ lru = (*lruCache)(nil)

var emptyLruCache = emptyLru{}

func (e emptyLru) add(string) {} // do nothing
func (e emptyLru) del(string) {} // do nothing

func newLruCache(capacity int, onEvict func(string)) lru {
	return &lruCache{
		capacity,
		list.New(),
		make(map[string]*list.Element, capacity),
		onEvict,
	}
}

func (c *lruCache) add(key string) {
	if elem, ok := c.evicts[key]; ok {
		c.list.MoveToFront(elem)
		return
	}

	c.evicts[key] = c.list.PushFront(key) // PushFront 在list前插入key并返回Element
	if c.list.Len() > c.capacity {
		c.removeOldest()
	}
}

func (c *lruCache) del(key string) {
	if elem, ok := c.evicts[key]; ok {
		c.removeElement(elem)
	}
}

// removeOldest 删除lru末项
func (c *lruCache) removeOldest() {
	elem := c.list.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}

func (c *lruCache) removeElement(elem *list.Element) {
	c.list.Remove(elem)
	key := elem.Value.(string)
	delete(c.evicts, key)
	c.onEvict(key)
}
