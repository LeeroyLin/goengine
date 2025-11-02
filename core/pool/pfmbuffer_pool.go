package pool

type PFMBufferPool struct {
	cap   int
	cache chan *PFMBuffer
}

func NewPFMBufferPool(cap, max int) *PFMBufferPool {
	return &PFMBufferPool{
		cap:   cap,
		cache: make(chan *PFMBuffer, max),
	}
}

func (p *PFMBufferPool) Get() *PFMBuffer {
	select {
	case b := <-p.cache:
		b.Reset()
		return b
	default:
		return NewPFMBuffer(p.cap)
	}
}

func (p *PFMBufferPool) Set(b *PFMBuffer) {
	select {
	case p.cache <- b:
		return
	default:
		return
	}
}
