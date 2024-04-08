package timex

import (
	"container/list"
	"errors"
	"log"
	"sync"
	"time"
)

var (
	ErrArgs   = errors.New("[timewheel]: invalid args")
	ErrClosed = errors.New("[timewheel]: has been closed")
)

type (
	taskEntry struct {
		task   func()
		key    string
		pos    int
		circle int
	}

	TimeWheel struct {
		sync.Once
		interval   time.Duration
		ticker     *time.Ticker
		cursorPos  int
		numSlots   int
		slots      []*list.List
		stopChan   chan struct{}
		addChan    chan *taskEntry
		removeChan chan string
		// The map here is thread-safe, ensuring the thread-safety
		// of the time wheel during concurrent reading and writing.
		// TODO: Leave it as interface, you can use non-thread-safe map, and the caller ensures thread safety.
		keyToEntry sync.Map // key --> val:(taskEntry->list.Element->any)
	}
)

func NewTimeWheel(numSlots int, interval time.Duration) (*TimeWheel, error) {
	tw := &TimeWheel{
		interval:   interval,
		ticker:     time.NewTicker(interval),
		numSlots:   numSlots,
		slots:      make([]*list.List, numSlots),
		stopChan:   make(chan struct{}),
		addChan:    make(chan *taskEntry),
		removeChan: make(chan string),
	}

	for i := 0; i < numSlots; i++ {
		tw.slots[i] = list.New()
	}
	go tw.run()
	return tw, nil
}

// func newTimeWheelWithTicker(numSlots int, interval time.Duration, ticker *time.Ticker) (*TimeWheel, error) {
// 	panic("unimplemented")
// }

func (tw *TimeWheel) Stop() {
	tw.Do(func() {
		tw.ticker.Stop()
		close(tw.stopChan)
	})
}

func (tw *TimeWheel) AddTask(key string, delay time.Duration, task func()) error {
	if key == "" || delay <= 0 || task == nil {
		return ErrArgs
	}
	pos, circle := tw.getPositionAndCircle(delay)
	select {
	case tw.addChan <- &taskEntry{
		task:   task,
		key:    key,
		pos:    pos,
		circle: circle,
	}:
		return nil
	case <-tw.stopChan:
		return ErrClosed
	}
}

func (tw *TimeWheel) RemoveTask(key string) error {
	if key == "" {
		return ErrArgs
	}

	select {
	case tw.removeChan <- key:
		return nil
	case <-tw.stopChan:
		return ErrClosed
	}
}

func (tw *TimeWheel) onTick() {
	tw.cursorPos = (tw.cursorPos + 1) % tw.numSlots
	slot := tw.slots[tw.cursorPos]
	tw.scanAndRunTask(slot)
}

func (tw *TimeWheel) run() {
	for {
		select {
		case <-tw.stopChan:
			return
		case <-tw.ticker.C:
			tw.onTick()
		case taskEntry := <-tw.addChan:
			tw.addTask(taskEntry)
		case key := <-tw.removeChan:
			tw.removeTask(key)
		}
	}
}

func (tw *TimeWheel) scanAndRunTask(slot *list.List) {
	for e := slot.Front(); e != nil; {
		taskEntry := e.Value.(*taskEntry)
		if taskEntry.circle > 0 {
			taskEntry.circle--
			e = e.Next()
			continue
		}

		// execute real task function
		// TODO: use coroutine pool
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("Panic in task execution, key: %v, err: %v\n", taskEntry.key, err)
				}
			}()
			taskEntry.task()
		}()

		// after execution, delete from TimeWheel
		next := e.Next()
		slot.Remove(e)                      // delete from slot
		tw.keyToEntry.Delete(taskEntry.key) // delete from Mapping Table
		e = next
	}
}

func (tw *TimeWheel) getPositionAndCircle(delay time.Duration) (pos, circle int) {
	step := int(delay / tw.interval)
	pos = (tw.cursorPos + step) % tw.numSlots
	circle = (step - 1) / tw.numSlots
	return
}

func (tw *TimeWheel) addTask(taskEntry *taskEntry) {
	slot := tw.slots[taskEntry.pos]
	if _, ok := tw.keyToEntry.Load(taskEntry.key); ok { // Duplicate task, delete old one
		tw.removeTask(taskEntry.key)
	}
	eTask := slot.PushBack(taskEntry)
	tw.keyToEntry.Store(taskEntry.key, eTask)
}

func (tw *TimeWheel) removeTask(key string) {
	val, ok := tw.keyToEntry.Load(key)
	if !ok {
		return
	}
	elementTask := val.(*list.Element)
	tw.keyToEntry.Delete(key) // delete from Mapping Table

	taskEntry := elementTask.Value.(*taskEntry)
	tw.slots[taskEntry.pos].Remove(elementTask) // delete from TimeWheel slots
}
