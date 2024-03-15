package cache

import (
	"fmt"
	"sync"
	"time"

	"byteurl/syncx"
	"byteurl/timex"
)

const (
	defaultSlots = 300 // Number of default time wheel slots
)

type (
	CacheOption func(*Cache)

	Cache struct {
		name        string
		lock        sync.Mutex
		data        map[string]any
		lru         lru
		barrier     syncx.SingleFlight
		expire      time.Duration
		timingWheel *timex.TimeWheel
	}
)

func NewCache(name string, expire time.Duration, opts ...CacheOption) *Cache {
	cache := &Cache{
		name:    name,
		data:    make(map[string]any),
		lru:     emptyLruCache,
		barrier: syncx.NewSingleFlight(),
		expire:  expire,
	}

	for _, opt := range opts {
		opt(cache)
	}

	timeWheel, err := timex.NewTimeWheel(defaultSlots, time.Second)
	if err != nil {
		panic(fmt.Errorf("[cache]: NewTimeWheel(defaultSlots, time.Second) failed... err: %#v", err))
	}

	cache.timingWheel = timeWheel

	return cache
}

func (c *Cache) Del(key string) {
	c.lock.Lock()
	c.lru.Del(key)
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
	c.lru.Add(key)
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
		c.lru.Add(key)
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
