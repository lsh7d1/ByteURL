package component

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestSingleFlightCallDo 测试Do的正常工作
func TestSingleFlightCallDo(t *testing.T) {
	g := NewSingleFlight()
	res, err := g.Do("key", func() (any, error) {
		return "value", nil
	})
	assert.Equal(t, "value", res)
	assert.Nil(t, err)
}

// TestSingleFlightCallDoWithError 测试Do能返回Error
func TestSingleFlightCallDoWithError(t *testing.T) {
	g := NewSingleFlight()
	res, err := g.Do("key", func() (any, error) {
		return "value", assert.AnError
	})
	assert.Equal(t, "value", res)
	assert.Equal(t, assert.AnError, err)
}

// TestSingleFlightCallDoRestrict 测试Do的并发抑制
func TestSingleFlightCallDoRestrict(t *testing.T) {
	g := NewSingleFlight()
	ch := make(chan string)
	var times int32 = 0
	fn := func() (any, error) {
		atomic.AddInt32(&times, 1)
		return <-ch, nil
	}

	const n int = 100
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			res, err := g.Do("key", fn)
			assert.Equal(t, "value", res)
			assert.Nil(t, err)
		}()
		wg.Done()
	}

	time.Sleep(time.Millisecond * 100)
	ch <- "value"
	wg.Wait()
	assert.Equal(t, int32(1), times)
}

// TestSingleFlightCallDoDiffRestrict 测试差异化的并发抑制
func TestSingleFlightCallDoDiffRestrict(t *testing.T) {
	g := NewSingleFlight()
	ch := make(chan struct{})
	var times int32 = 0
	testCases := []string{"key1", "key2", "key3", "key4", "key1", "key2", "key3", "key4", "key2", "key4", "key2", "key1", "key2"}

	fn := func() (any, error) {
		atomic.AddInt32(&times, 1)
		time.Sleep(time.Microsecond * 100)
		return nil, nil
	}

	var wg sync.WaitGroup
	for _, key := range testCases {
		wg.Add(1)
		go func(key string) {
			<-ch // 阻塞所有请求
			_, err := g.Do(key, fn)
			assert.Nil(t, err)
			wg.Done()
		}(key)
	}

	time.Sleep(time.Millisecond * 100) // 阻塞100ms，统一启动
	close(ch)
	wg.Wait()

	assert.Equal(t, int32(4), times)
}

// TestSingleFlightCallDoExRestrict 测试DoEx的并发抑制
func TestSingleFlightCallDoExRestrict(t *testing.T) {
	g := NewSingleFlight()
	ch := make(chan string)

	fn := func() (any, error) {
		return <-ch, nil
	}

	const n int = 100
	var wg sync.WaitGroup
	var gotFresh int32 = 0
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			res, fresh, err := g.DoEx("key", fn)
			if fresh {
				atomic.AddInt32(&gotFresh, 1)
			}
			assert.Nil(t, err)
			assert.Equal(t, "value", res)
			wg.Done()
		}()
	}

	time.Sleep(time.Millisecond * 100)
	ch <- "value"
	wg.Wait()
	assert.Equal(t, int32(1), gotFresh)
}
