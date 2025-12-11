package utils

import (
	"fmt"
	"sync"
	"time"
)

type ShortSnowflake struct {
	mu        sync.Mutex
	timestamp int64 // 上次生成ID的时间戳（毫秒）
	sequence  int64
}

func NewShortSnowflake() *ShortSnowflake {
	return &ShortSnowflake{}
}

func (s *ShortSnowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	// 处理时钟回拨：等待时间追平
	if now < s.timestamp {
		// 计算需要等待的毫秒数
		waitMs := s.timestamp - now
		fmt.Printf("检测到时钟回拨，等待%d毫秒...\n", waitMs)
		time.Sleep(time.Duration(waitMs) * time.Millisecond)
		// 等待后重新获取当前时间
		now = time.Now().UnixMilli()
	}

	// 处理序列号（同一毫秒内递增，超过则等待下一毫秒）
	if now == s.timestamp {
		s.sequence++
		if s.sequence > 2047 { // 11位序列号最大为2047
			time.Sleep(1 * time.Millisecond)
			now = time.Now().UnixMilli()
			s.sequence = 0
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now

	// 组合ID：时间戳(42位) + 序列号(11位)
	return (now << 11) | s.sequence
}
