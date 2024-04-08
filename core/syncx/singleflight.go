package syncx

import "sync"

type (
	// SingleFlight allows concurrent calls to the same key to share the call results
	SingleFlight interface {
		// Do execute the given function fn and returns
		// the result according to the specified key.
		// If there are multiple requests using the same key,
		// only one request will execute fn. Other requests
		// wait and share the results.
		Do(key string, fn func() (any, error)) (any, error)
		// DoEx is similar to the Do method, but it returns
		// an additional bool value indicating whether the
		// result is fresh (whether it is the result obtained
		// by executing fn recently.)
		DoEx(key string, fn func() (any, error)) (any, bool, error)
	}

	// call blocks a group of calls for the same key
	call struct {
		wg  sync.WaitGroup
		res any
		err error
	}

	// flightGroup is the implementation class of SingleFlight
	// and saves ongoing calls
	flightGroup struct {
		lock  sync.Mutex
		calls map[string]*call
	}
)

func NewSingleFlight() SingleFlight {
	return &flightGroup{
		calls: make(map[string]*call),
	}
}

func (f *flightGroup) Do(key string, fn func() (any, error)) (any, error) {
	c, created := f.createCall(key)
	if created {
		return c.res, c.err
	}

	f.doCall(c, key, fn)
	return c.res, c.err
}

func (f *flightGroup) DoEx(key string, fn func() (any, error)) (any, bool, error) {
	c, created := f.createCall(key)
	if created {
		return c.res, false, c.err
	}

	f.doCall(c, key, fn)
	return c.res, true, c.err
}

func (f *flightGroup) createCall(key string) (c *call, created bool) {
	f.lock.Lock()
	if c, ok := f.calls[key]; ok {
		f.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	c = new(call)
	c.wg.Add(1)
	f.calls[key] = c
	f.lock.Unlock()
	return c, false
}

func (f *flightGroup) doCall(c *call, key string, fn func() (any, error)) {
	defer func() {
		f.lock.Lock()
		delete(f.calls, key)
		f.lock.Unlock()
		c.wg.Done()
	}()
	c.res, c.err = fn()
}
