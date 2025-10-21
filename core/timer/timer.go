package timer

import (
	"sync"
	"time"
)

const (
	HOUR_NAME     = "HOUR"
	HOUR_INTERVAL = 60 * 60 * 1e3 // ms为精度
	HOUR_SCALES   = 12

	MINUTE_NAME     = "MINUTE"
	MINUTE_INTERVAL = 60 * 1e3
	MINUTE_SCALES   = 60

	SECOND_NAME     = "SECOND"
	SECOND_INTERVAL = 1e3
	SECOND_SCALES   = 60

	TIMERS_MAX_CAP = 2048 // 每个时间轮刻度挂载定时器的最大个数
)

var idGen uint32 = 0
var genMutex sync.Mutex

func getNextID() uint32 {
	genMutex.Lock()
	defer genMutex.Unlock()

	idGen++
	return idGen
}

type Timer struct {
	// 定时器唯一id
	tid uint32
	// 延迟调用的函数
	delayFunc *DelayFunc
	// 调用时间（unix 时间，单位ms）
	callAtMS int64
}

func (t *Timer) GetId() uint32 {
	return t.tid
}

// CurrMilli 获得毫秒时间戳
func CurrMilli() int64 {
	return time.Now().UnixMilli()
}

// CurrMicro 获得微秒时间戳
func CurrMicro() int64 {
	return time.Now().UnixMicro()
}

// CurrNano 获得纳秒时间戳
func CurrNano() int64 {
	return time.Now().UnixNano()
}

// NewTimerAt 在毫秒时间戳时调用回调函数
func NewTimerAt(df *DelayFunc, callAtMS int64) *Timer {
	return &Timer{
		tid:       getNextID(),
		delayFunc: df,
		callAtMS:  callAtMS,
	}
}

// NewTimerAfter 在当前时间戳后延迟毫秒时间调用回调函数
func NewTimerAfter(df *DelayFunc, afterMS int64) *Timer {
	return NewTimerAt(df, CurrMilli()+afterMS)
}

func (t *Timer) Run() {
	go func() {
		now := CurrMilli()
		if t.callAtMS > now {
			time.Sleep(time.Duration(t.callAtMS-now) * time.Millisecond)
		}

		t.delayFunc.Call()
	}()
}
