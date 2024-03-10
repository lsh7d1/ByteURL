package component

import (
	"testing"
	"time"
)

/*
--- FAIL: TestTimeWheel (5.02s)

	timewheel_test.go:15: start: 2024-03-10 15:29:28.8803944 +0800 CST m=+0.007561001
	timewheel_test.go:18: task1, 2024-03-10 15:29:30.8872625 +0800 CST m=+2.014429101
	timewheel_test.go:26: task2, 2024-03-10 15:29:31.8918554 +0800 CST m=+3.019022001

FAIL
exit status 1
FAIL    byteurl/component       5.050s
*/
func TestTimeWheel(t *testing.T) {
	tw, err := NewTimeWheel(8, time.Millisecond*500)
	if err != nil {
		panic(err)
	}
	defer tw.Stop()

	t.Errorf("start: %v", time.Now())

	tw.AddTask("task1", time.Second*2, func() {
		t.Errorf("task1, %v", time.Now())
	})

	tw.AddTask("task2", time.Second*10, func() {
		t.Errorf("task2, %v", time.Now())
	})

	tw.AddTask("task2", time.Second*3, func() {
		t.Errorf("task2, %v", time.Now())
	})

	<-time.After(time.Second * 5)
}
