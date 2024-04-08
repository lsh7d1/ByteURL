package cache

import (
	"fmt"
	"sync"
	"time"

	"byteurl/core/mathx"
	"byteurl/core/syncx"
	"byteurl/core/timex"
)

const (
	// Default number of TimeWheel slots
	defaultSlots = 300
	// Default TimeWheel check interval
	defaultInterval = time.Second
	// Default expiration time perturbation to prevent
	// cache from expiring at the same time
	// Example: this will make the expiration time distributed in [0.9, 1.1] seconds
	defaultExpireDeviation = 0.1
)

type (
	DataSource interface {
		Get(key string) (string, bool)
		Add(key, value string)
		Del(key string)
	}

	mappingTable struct {
		items map[string]string
		lock  sync.Mutex
	}

	CacheOption func(*Cache)

	Cache struct {
		name           string
		lock           sync.Mutex
		data           DataSource
		barrier        syncx.SingleFlight
		unstableExpire mathx.Unstable
		expire         time.Duration
		timingWheel    *timex.TimeWheel
	}
)

func NewCache(name string, expire time.Duration, opts ...CacheOption) *Cache {
	cache := &Cache{
		name:           name,
		lock:           sync.Mutex{},
		data:           NewMappingTable(),
		barrier:        syncx.NewSingleFlight(),
		unstableExpire: mathx.NewUnstable(defaultExpireDeviation),
		expire:         expire,
	}

	timeWheel, err := timex.NewTimeWheel(defaultSlots, defaultInterval)
	if err != nil {
		panic(fmt.Errorf("[cache]: NewTimeWheel failed... err: %#v", err))
	}

	cache.timingWheel = timeWheel

	for _, opt := range opts {
		opt(cache)
	}

	return cache
}

func (c *Cache) Get(key string) (string, bool) {
	return c.get(key)
}

func (c *Cache) Del(key string) {
	c.data.Del(key)
}

func (c *Cache) Set(key string, value string) {
	//fmt.Printf("Set: [%s,%s]\n", key, value)
	c.SetWithExpire(key, value, c.expire)
}

func (c *Cache) SetWithExpire(key string, value string, expire time.Duration) {
	c.data.Add(key, value)
}

// Take Try to get the value of the specified key in the Cache
// If there is a direct return, otherwise call fetch
// to obtain it and return after storing it in the Cache.
// TODO: I can't seem to avoid repeated writes other than using an atomic lock
func (c *Cache) Take(key string, fetch func() (string, error)) (any, error) {
	if value, ok := c.get(key); ok {
		return value, nil
	}

	value, err := c.barrier.Do(key, func() (any, error) {
		// because O(1) on map search in memory, and fetch is an IO query,
		// so we do double-check, cache might be taken by another call
		if val, ok := c.get(key); ok {
			return val, nil
		}

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

func (c *Cache) get(key string) (string, bool) {
	value, ok := c.data.Get(key)
	if ok {
		return value, true
	}
	return value, false
}

func (c *Cache) onEvict(key string) {
	// already locked by lru
	// panic("2~!!!")
	// c.data.Del(key)
	// c.timingWheel.RemoveTask(key)
}

// WithAroundCapLimit Because we use HashLru, we cannot
// guarantee the exact quantity limit. We can only guarantee
// that the capacity is not less than limit.
func WithAroundCapLimit(limit int) CacheOption {
	return func(c *Cache) {
		if limit <= 0 {
			return
		}
		c.data = NewHashLru(limit, c.onEvict)
	}
}

// WithCapLimit Accurate capacity limits can be guaranteed
func WithCapLimit(limit int) CacheOption {
	return func(c *Cache) {
		if limit <= 0 {
			return
		}
		c.data = NewLruCache(limit, c.onEvict)
	}
}

func NewMappingTable() DataSource {
	return &mappingTable{
		items: make(map[string]string),
		lock:  sync.Mutex{},
	}
}

func (mp *mappingTable) Get(key string) (string, bool) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	val, ok := mp.items[key]
	return val, ok
}

func (mp *mappingTable) Add(key, value string) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	mp.items[key] = value
}

func (mp *mappingTable) Del(key string) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	delete(mp.items, key)
}
