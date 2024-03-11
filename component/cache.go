package component

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

const (
	defaultSlots = 300 // Default TimeWheel slot number
)

type (
	CacheOption func(*Cache)

	Cache struct {
		name        string
		lock        sync.Mutex
		data        map[string]any
		lru         lru
		barrier     SingleFlight
		expire      time.Duration
		timingWheel *TimeWheel
	}
)

func NewCache(name string, expire time.Duration, opts ...CacheOption) *Cache {
	cache := &Cache{
		name:    name,
		data:    make(map[string]any),
		lru:     emptyLruCache,
		barrier: NewSingleFlight(),
		expire:  expire,
	}

	for _, opt := range opts {
		opt(cache)
	}

	timeWheel, err := NewTimeWheel(defaultSlots, time.Second)
	if err != nil {
		panic(fmt.Errorf("[cache]: NewTimeWheel(defaultSlots, time.Second) failed... err: %#v", err))
	}

	cache.timingWheel = timeWheel

	return cache
}

func (c *Cache) Del(key string) {
	c.lock.Lock()
	c.lru.del(key)
	delete(c.data, key)
	c.lock.Unlock()
	c.timingWheel.RemoveTask(key)
}

func (c *Cache) Get(key string) (any, bool) {
	return c.get(key)
}

func (c *Cache) Set(key string, value any) {
	c.SetWithExpire(key, value, c.expire)
}

func (c *Cache) SetWithExpire(key string, value any, expire time.Duration) {
	c.lock.Lock()
	c.data[key] = value
	c.lru.add(key)
	c.lock.Unlock()

	c.timingWheel.AddTask(key, expire, func() {
		c.Del(key)
	})
}

// Take Try to get the value of the specified key in the Cache
// If there is a direct return, otherwise call fetch
// to obtain it and return after storing it in the Cache.
func (c *Cache) Take(key string, fetch func() (any, error)) (any, error) {
	if value, ok := c.get(key); ok {
		return value, nil
	}

	value, err := c.barrier.Do(key, func() (any, error) {
		v, e := fetch()
		if e != nil {
			return nil, e
		}
		c.Set(key, v)
		return v, nil
	})

	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *Cache) get(key string) (any, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.data[key]
	if ok {
		c.lru.add(key)
	}
	return value, ok
}

func (c *Cache) onEvict(key string) {
	// already locked
	delete(c.data, key)
	c.timingWheel.RemoveTask(key)
}

func (c *Cache) size() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

func WithCapLimit(limit int) CacheOption {
	return func(c *Cache) {
		if limit <= 0 {
			return
		}
		c.lru = newLruCache(limit, c.onEvict)
	}
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
		onEvict  func(string) // Callback when deleting
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

func (lru *lruCache) add(key string) {
	if elem, ok := lru.evicts[key]; ok {
		lru.list.MoveToFront(elem)
		return
	}

	lru.evicts[key] = lru.list.PushFront(key)
	if lru.list.Len() > lru.capacity {
		lru.removeOldest()
	}
}

func (lru *lruCache) del(key string) {
	if elem, ok := lru.evicts[key]; ok {
		lru.removeElement(elem)
	}
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
	delete(lru.evicts, key)
	lru.onEvict(key)
}
