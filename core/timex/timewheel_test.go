package timex

import (
	"testing"
	"time"
)

func TestTimeWheel(t *testing.T) {
	tw, err := NewTimeWheel(8, time.Millisecond*500) // 4 seconds a circle
	if err != nil {
		panic(err)
	}
	defer tw.Stop()

	t.Errorf("start: %v", time.Now())

	_ = tw.AddTask("task1", time.Second*2, func() {
		t.Errorf("task1, %v", time.Now())
	})

	_ = tw.AddTask("task2", time.Second*11, func() {
		t.Errorf("task2, %v", time.Now())
	})

	_ = tw.AddTask("task2", time.Second*7, func() {
		t.Errorf("task2, %v", time.Now())
	})

	<-time.After(time.Second * 12)
}
