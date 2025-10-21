package timer

import (
	"fmt"
	"testing"
	"time"
)

func secFunc(v ...interface{}) {
	fmt.Println(v...)
}

func minFunc(v ...interface{}) {
	fmt.Println(v...)
}

func TestTimeWheel(t *testing.T) {
	twMinute := NewTimeWheel("Minute", 60, 1000)
	twSecond := NewTimeWheel("Second", 5, 200)

	twMinute.AddNextTimeWheel(twSecond)

	dfSec := NewDelayFunc(secFunc, []interface{}{"SecFunc"})
	twSecond.AddTimer(NewTimerAfter(dfSec, 800))
	twSecond.AddTimer(NewTimerAfter(dfSec, 1300))

	dfMin := NewDelayFunc(minFunc, []interface{}{"MinFunc"})

	twMinute.AddTimer(NewTimerAfter(dfMin, 10000))
	twMinute.AddTimer(NewTimerAfter(dfMin, 62000))

	go func() {
		for {
			lastMS := CurrMilli()

			time.Sleep(time.Millisecond * 100)

			dt := CurrMilli() - lastMS

			twMinute.DoTick(dt)
		}
	}()

	select {
	case <-time.After(time.Hour * 1):
	}
}
