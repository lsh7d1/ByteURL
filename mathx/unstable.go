package mathx

import (
	"math/rand"
	"sync"
	"time"
)

type Unstable struct {
	deviation float64
	r         *rand.Rand
	lock      sync.Mutex
}

func NewUnstable(deviation float64) Unstable {
	if deviation < 0 {
		deviation = 0
	} else if deviation > 1 {
		deviation = 1
	}

	return Unstable{
		deviation: deviation,
		r:         rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:      sync.Mutex{},
	}
}

func (u Unstable) AroundDuration(base time.Duration) time.Duration {
	u.lock.Lock()
	val := time.Duration((1 - u.deviation + 2*u.deviation*rand.Float64()) * float64(base))
	u.lock.Unlock()
	return val
}

func (u Unstable) AroundInt(base int64) int64 {
	u.lock.Lock()
	val := int64((1 - u.deviation + 2*u.deviation*rand.Float64()) * float64(base))
	u.lock.Unlock()
	return val
}
