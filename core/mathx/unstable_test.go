package mathx

import (
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
)

func TestUnstable_AroundDuration(t *testing.T) {
	unstable := NewUnstable(0.05)
	for i := 0; i < 1000; i++ {
		val := unstable.AroundDuration(time.Second)
		assert.True(t, float64(time.Second)*0.95 <= float64(val))
		assert.True(t, float64(val) <= float64(time.Second)*1.05)
	}
}
