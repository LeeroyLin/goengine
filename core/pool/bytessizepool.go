package pool

import (
	"engine/core/elog"
	"strconv"
	"strings"
)

type BytesSizePool struct {
	pools map[int]*BytesPool
}

const (
	// 最大缓存个数
	maxCacheNum int = 32
)

var BytesPools *BytesSizePool

func init() {
	BytesPools = &BytesSizePool{
		pools: make(map[int]*BytesPool),
	}
}

// Get 获得比目标尺寸相等或更大的字节数组
func (sp *BytesSizePool) Get(bytesSize int) []byte {
	// 找等于或大于该长度的池子
	key := getNearGreaterPower(bytesSize)

	p, ok := sp.pools[key]

	if ok {
		return p.Get()
	}

	return make([]byte, bytesSize)
}

// Set 归还字节数组
func (sp *BytesSizePool) Set(b []byte) {
	l := len(b)

	// 归还时，找等于或小于该长度的池子
	key := getNearLessPower(l)

	p, ok := sp.pools[key]

	if !ok {
		p = NewBytesPool(key, maxCacheNum)
		sp.pools[key] = p
	}

	p.Set(b)
}

// Sprint 获得打印文本
func (sp *BytesSizePool) Sprint() string {
	var builder strings.Builder

	builder.WriteString("[BytesSizePool]\n")

	for size, p := range sp.pools {
		builder.WriteString("    - ")
		builder.WriteString(strconv.Itoa(size))
		builder.WriteString(" : ")
		builder.WriteString(strconv.Itoa(len(p.cacheChan)))
		builder.WriteString("\n")
	}

	return builder.String()
}

// 获得等于或大于该值的 2的整数次幂 的值
func getNearGreaterPower(val int) int {
	result := 1
	for result < val {
		result <<= 1
	}
	return result
}

// 获得等于或小于该值的 2的整数次幂 的值
func getNearLessPower(val int) int {
	if val < 1 {
		elog.Error("[Pool] wrong val to getNearLessPower.", val)
		return 0
	}

	lastResult := 1
	result := 1
	for {
		if result > val {
			return lastResult
		}

		lastResult = result

		result <<= 1
	}
}
