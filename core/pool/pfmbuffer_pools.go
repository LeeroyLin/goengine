package pool

import "github.com/LeeroyLin/goengine/core/syncmap"

var PFMBufferCtl = NewPFMBufferPools(4096)

type PFMBufferPools struct {
	cacheMax int
	pools    *syncmap.SyncMap[int, *PFMBufferPool]
}

func NewPFMBufferPools(cacheMax int) *PFMBufferPools {
	return &PFMBufferPools{
		cacheMax: cacheMax,
		pools:    syncmap.NewSyncMap[int, *PFMBufferPool](),
	}
}

// Get 获得比目标尺寸相等或更大的字节数组
func (p *PFMBufferPools) Get(cap int) *PFMBuffer {
	// 找等于或大于该长度的池子
	key := getNearGreaterPower(cap)

	pool, ok := p.pools.Get(key)

	if ok {
		return pool.Get()
	}

	return NewPFMBuffer(cap)
}

// Set 归还字节数组
func (p *PFMBufferPools) Set(buffer *PFMBuffer) {
	l := buffer.Cap()

	// 归还时，找等于或小于该长度的池子
	key := getNearLessPower(l)

	pool, ok := p.pools.Get(key)

	if !ok {
		pool = NewPFMBufferPool(key, p.cacheMax)
		p.pools.Add(key, pool)
	}

	pool.Set(buffer)
}
