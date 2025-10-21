package timer

import (
	"fmt"
	"testing"
	"time"
)

func func1(v ...interface{}) {
	fmt.Println(v...)
}

func TestTimer(t *testing.T) {
	fmt.Println("Timer test")

	timer := NewTimerAt(NewDelayFunc(func1, []interface{}{"1", "2", "3"}), CurrMilli()+500)
	timer.Run()

	timer2 := NewTimerAfter(NewDelayFunc(func1, []interface{}{"1", "2", "3"}), 1500)
	timer2.Run()

	time.Sleep(time.Second * 2)
}
