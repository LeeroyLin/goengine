package pool

import (
	"github.com/LeeroyLin/goengine/core/syncmap"
	"strconv"
	"strings"
)

type BytesPools struct {
	cacheMax int
	pools    *syncmap.SyncMap[int, *BytesPool]
}

func NewBytesPools(max int) *BytesPools {
	return &BytesPools{
		cacheMax: max,
		pools:    syncmap.NewSyncMap[int, *BytesPool](),
	}
}

// Get 获得比目标尺寸相等或更大的字节数组
func (sp *BytesPools) Get(bytesSize int) []byte {
	// 找等于或大于该长度的池子
	key := getNearGreaterPower(bytesSize)

	p, ok := sp.pools.Get(key)

	if ok {
		return p.Get()
	}

	return make([]byte, bytesSize)
}

// Set 归还字节数组
func (sp *BytesPools) Set(b []byte) {
	l := len(b)

	// 归还时，找等于或小于该长度的池子
	key := getNearLessPower(l)

	p, ok := sp.pools.Get(key)

	if !ok {
		p = NewBytesPool(key, sp.cacheMax)
		sp.pools.Add(key, p)
	}

	p.Set(b)
}

// Sprint 获得打印文本
func (sp *BytesPools) Sprint() string {
	var builder strings.Builder

	builder.WriteString("[BytesPools]\n")

	sp.pools.Range(func(key int, p *BytesPool) bool {
		builder.WriteString("    - ")
		builder.WriteString(strconv.Itoa(key))
		builder.WriteString(" : ")
		builder.WriteString(strconv.Itoa(len(p.cacheChan)))
		builder.WriteString("\n")

		return true
	})

	return builder.String()
}
