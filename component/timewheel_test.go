package component

import (
	"testing"
	"time"
)

/*
--- FAIL: TestTimeWheel (12.02s)

	timewheel_test.go:26: start: 2024-03-11 17:37:48.7934477 +0800 CST m=+0.004245201
	timewheel_test.go:29: task1, 2024-03-11 17:37:50.8066652 +0800 CST m=+2.017462701
	timewheel_test.go:37: task2, 2024-03-11 17:37:55.80156 +0800 CST m=+7.012357501

FAIL
exit status 1
FAIL    byteurl/component       12.055s
*/
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
