package utils

import (
	"fmt"
	"sync"
	"time"
)

type Snowflake struct {
	mu        sync.Mutex
	timestamp int64 // 上次生成ID的时间戳（毫秒）
	machineID int64
	sequence  int64
}

func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > 1023 {
		return nil, fmt.Errorf("机器ID必须在0-1023之间")
	}
	return &Snowflake{machineID: machineID}, nil
}

func (s *Snowflake) Generate() int64 {
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
		if s.sequence > 4095 { // 12位序列号最大为4095
			time.Sleep(1 * time.Millisecond)
			now = time.Now().UnixMilli()
			s.sequence = 0
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now

	// 组合ID：时间戳(41位) + 机器ID(10位) + 序列号(12位)
	return (now << 22) | (s.machineID << 12) | s.sequence
}
