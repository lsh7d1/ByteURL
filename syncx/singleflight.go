package syncx

import "sync"

type (
	// SingleFlight 让对同一个key的并发调用共享调用结果
	SingleFlight interface {
		// Do 执行给定的函数fn，并按key返回结果
		// 如果有多个请求使用同样的key，只有一个请求会执行fn
		// 其他请求等待并共享结果
		Do(key string, fn func() (any, error)) (any, error)
		// DoEx 与Do方法类似，但它多返回一个bool值
		// 表示返回的结果是否新鲜（是否是最近执行fn得到的结果
		DoEx(key string, fn func() (any, error)) (any, bool, error)
	}

	// call 阻塞对于同一个 key 的一组调用
	call struct {
		wg  sync.WaitGroup
		res any
		err error
	}

	// flightGroup 是SingleFlight的实现类
	// 保存正在进行的调用
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
