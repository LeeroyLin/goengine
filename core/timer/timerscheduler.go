package timer

import (
	"sync"
	"time"
)

type TimerScheduler struct {
	twHour   *TimeWheel
	twMinute *TimeWheel
	twSecond *TimeWheel

	closeChan chan interface{}

	sync.Mutex
}

func NewTimerScheduler() *TimerScheduler {
	ts := &TimerScheduler{
		twHour:    NewTimeWheel("hour", 60, 1000000),
		twMinute:  NewTimeWheel("minute", 60, 1000),
		twSecond:  NewTimeWheel("second", 10, 100),
		closeChan: make(chan interface{}),
	}

	ts.twHour.AddNextTimeWheel(ts.twMinute)
	ts.twMinute.AddNextTimeWheel(ts.twSecond)

	return ts
}

func (ts *TimerScheduler) AddTimer(t *Timer) {
	ts.twHour.AddTimer(t)
}

func (ts *TimerScheduler) RemoveTimer(tid uint32) {
	ts.twHour.RemoveTimer(tid)
	ts.twMinute.RemoveTimer(tid)
	ts.twSecond.RemoveTimer(tid)
}

func (ts *TimerScheduler) ClearTimer() {
	ts.twHour.ClearTimer()
	ts.twMinute.ClearTimer()
	ts.twSecond.ClearTimer()
}

func (ts *TimerScheduler) Run() {
	ts.Lock()
	defer ts.Unlock()

	go func() {
		for {
			select {
			case <-ts.closeChan:
				return
			default:
				lastMS := CurrMilli()

				time.Sleep(time.Millisecond * 100)

				dt := CurrMilli() - lastMS

				ts.twHour.DoTick(dt)
			}
		}
	}()
}

func (ts *TimerScheduler) Stop() {
	select {
	case <-ts.closeChan:
		return
	default:
		close(ts.closeChan)
	}
}
