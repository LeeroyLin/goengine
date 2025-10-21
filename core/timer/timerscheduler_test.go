package timer

import (
	"fmt"
	"testing"
	"time"
)

func callFunc(v ...interface{}) {
	fmt.Println(v...)
}

func TestTimerScheduler(t *testing.T) {
	ts := NewTimerScheduler()
	ts.Run()

	ts.AddTimer(NewTimerAfter(NewDelayFunc(callFunc, []interface{}{"call1"}), 800))
	ts.AddTimer(NewTimerAfter(NewDelayFunc(callFunc, []interface{}{"call2"}), 1300))
	ts.AddTimer(NewTimerAfter(NewDelayFunc(callFunc, []interface{}{"call3"}), 10000))
	ts.AddTimer(NewTimerAfter(NewDelayFunc(callFunc, []interface{}{"call4"}), 62000))

	select {
	case <-time.After(time.Hour * 1):
	}
}
