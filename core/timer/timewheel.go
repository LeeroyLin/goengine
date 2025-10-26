package timer

import (
	"github.com/LeeroyLin/goengine/core/elog"
	"github.com/LeeroyLin/goengine/core/syncmap"
	"sync"
)

// TimeWheel 时间轮 将Timer按时间刻度存储 时间轮可嵌套
type TimeWheel struct {
	// 时间轮名字
	name string
	// 刻度之间的间隔时间，单位毫秒
	scalesIntervalMS int64
	// 刻度数
	scales int
	// 当前时间指针指向，等待执行
	currIndex int
	// 所有timer 第一层键为时间刻度，第二层键为timer id
	timerQueue map[int]*syncmap.SyncMap[uint32, *Timer]
	// 下一层时间轮
	nextTimeWheel *TimeWheel
	// 时间叠加
	timeCntMS int64
	// 读写锁
	sync.RWMutex
}

// NewTimeWheel 创建一个时间轮
// Parameters:
//
//	name:		时间轮名字
//	scales:		时间轮刻度数
//	scalesIntervalMS:	时间轮间隔时间，单位毫秒
//
// Returns:
//
//	*TimeWheel:	创建好的时间轮
func NewTimeWheel(name string, scales int, scalesIntervalMS int64) *TimeWheel {
	wheel := &TimeWheel{
		name:             name,
		scales:           scales,
		scalesIntervalMS: scalesIntervalMS,
		timerQueue:       make(map[int]*syncmap.SyncMap[uint32, *Timer]),
	}

	for i := 0; i < scales; i++ {
		wheel.timerQueue[i] = syncmap.NewSyncMap[uint32, *Timer]()
	}

	elog.Info("[TimeWheel] init time wheel:", name, "scales:", scales, "scalesIntervalMS:", scalesIntervalMS, "is down.")

	return wheel
}

// AddTimer 添加计时器
func (tw *TimeWheel) AddTimer(timer *Timer) {
	tw.Lock()
	defer tw.Unlock()

	now := CurrMilli()

	// 获得毫秒延迟时间
	delayMS := timer.callAtMS - now

	// 如果延迟时间小于刻度间隔时间
	if delayMS <= tw.scalesIntervalMS {
		// 有下一层时间轮
		if tw.nextTimeWheel != nil {
			// 添加到下一层时间轮
			tw.nextTimeWheel.AddTimer(timer)
			return
		}

		// 直接添加到当前刻度
		m := tw.timerQueue[tw.currIndex]
		m.Add(timer.GetId(), timer)
		return
	}

	// 需要跨越几个刻度
	dn := delayMS / tw.scalesIntervalMS
	if delayMS%tw.scalesIntervalMS == 0 {
		dn--
	}
	if dn < 0 {
		dn = 0
	}
	m := tw.timerQueue[(tw.currIndex+int(dn))%tw.scales]
	m.Add(timer.GetId(), timer)
}

// RemoveTimer 移除计时器
func (tw *TimeWheel) RemoveTimer(tid uint32) {
	for _, m := range tw.timerQueue {
		if _, ok := m.Get(tid); ok {
			m.Delete(tid)
			return
		}
	}
}

func (tw *TimeWheel) ClearTimer() {
	for _, m := range tw.timerQueue {
		m.Clear()
	}
}

// AddNextTimeWheel 添加下一层时间轮
func (tw *TimeWheel) AddNextTimeWheel(next *TimeWheel) {
	tw.Lock()
	tw.Unlock()

	tw.nextTimeWheel = next
}

// DoTick 时间步进
func (tw *TimeWheel) DoTick(dtMS int64) {
	timers := tw.tickAndGetTimers(dtMS)

	// 回调
	for _, timer := range timers {
		timer.Run()
	}
}

func (tw *TimeWheel) tickAndGetTimers(dtMS int64) []*Timer {
	tw.Lock()
	defer tw.Unlock()

	tw.timeCntMS += dtMS

	var timers []*Timer

	changeToNext := tw.timeCntMS >= tw.scalesIntervalMS

	// 推进到下一个刻度
	if changeToNext {
		// 获得可回调的计时器
		timers = tw.getAvailableTimers()

		tw.timeCntMS = tw.timeCntMS - tw.scalesIntervalMS
		tw.currIndex = (tw.currIndex + 1) % tw.scales
	}

	if tw.nextTimeWheel != nil {
		// 叠加子时间轮的可回调计时器
		timers = append(timers, tw.nextTimeWheel.tickAndGetTimers(dtMS)...)
	}

	if changeToNext {
		// 尝试将计时器添加到下一层时间轮
		tw.tryPushToNext()
	}

	return timers
}

// 尝试将计时器添加到下一层时间轮
func (tw *TimeWheel) tryPushToNext() {
	if tw.nextTimeWheel == nil {
		return
	}

	m := tw.timerQueue[tw.currIndex]

	now := CurrMilli()

	m.Range(func(tid uint32, timer *Timer) bool {
		deltaMS := timer.callAtMS - now

		// 剩余时间小于刻度间隔时间
		if deltaMS <= tw.scalesIntervalMS {
			m.Delete(tid)

			// 记录到下一层时间轮
			tw.nextTimeWheel.AddTimer(timer)
		}

		return true
	})
}

// 获得可回调的计时器
func (tw *TimeWheel) getAvailableTimers() []*Timer {
	m := tw.timerQueue[tw.currIndex]

	now := CurrMilli()

	// 可以回调的计时器
	callAvailableTimers := make([]*Timer, 0)

	m.Range(func(tid uint32, timer *Timer) bool {
		deltaMS := timer.callAtMS - now

		// 时间满足
		if deltaMS <= 0 {
			m.Delete(tid)

			// 记录到可回调计时器列表
			callAvailableTimers = append(callAvailableTimers, timer)
			return true
		}

		return true
	})

	return callAvailableTimers
}
