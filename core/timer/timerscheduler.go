package timer

import (
	"context"
	"sync"
	"time"
)

type TimerScheduler struct {
	twHour    *TimeWheel
	twMinute  *TimeWheel
	twSecond  *TimeWheel
	isRunning bool

	sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTimerScheduler() *TimerScheduler {
	ts := &TimerScheduler{
		twHour:   NewTimeWheel("hour", 60, 1000000),
		twMinute: NewTimeWheel("minute", 60, 1000),
		twSecond: NewTimeWheel("second", 10, 100),
	}

	ts.twHour.AddNextTimeWheel(ts.twMinute)
	ts.twMinute.AddNextTimeWheel(ts.twSecond)

	return ts
}

func (ts *TimerScheduler) IsRunning() bool {
	return ts.isRunning
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
}

func (ts *TimerScheduler) Run() {
	ts.Lock()
	defer ts.Unlock()

	if ts.isRunning {
		return
	}

	ts.ctx, ts.cancel = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ts.ctx.Done():
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
	ts.Lock()
	defer ts.Unlock()

	if !ts.isRunning {
		return
	}

	ts.isRunning = false
	ts.cancel()
}
