package snapshot

import "sync"

type Snapshot struct {
	start           bool     // 标记是否开始快照
	readyChan       chan int // 准备完毕队列
	snapshotHandler func()   // 快照函数
	sync.Mutex
}

func (ss *Snapshot) SetSnapshotHandler(snapshotHandler func()) {
	ss.snapshotHandler = snapshotHandler
}

func (ss *Snapshot) MarkStart(needReadyNum int) {
	ss.Lock()
	defer ss.Unlock()

	if ss.start {
		return
	}

	ss.start = true

	for i := 0; i < needReadyNum; i++ {
		ss.readyChan <- i
	}

	// 用于判断全部准备完毕
	ss.readyChan <- -1
}

func (ss *Snapshot) ReadyOne() {
	select {
	case v := <-ss.readyChan:
		// 全部准备完毕
		if v == -1 {
			if ss.snapshotHandler != nil {
				ss.snapshotHandler()
			}
		}
	default:
		return
	}
}
