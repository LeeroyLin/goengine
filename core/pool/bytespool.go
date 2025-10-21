package pool

type BytesPool struct {
	BytesSize int
	cacheChan chan []byte
}

func NewBytesPool(bytesSize, max int) *BytesPool {
	p := &BytesPool{
		BytesSize: bytesSize,
		cacheChan: make(chan []byte, max),
	}

	return p
}

func (p *BytesPool) Get() []byte {
	select {
	case b := <-p.cacheChan:
		return b
	default:
		b := make([]byte, p.BytesSize)
		return b
	}
}

func (p *BytesPool) Set(d []byte) {
	select {
	case p.cacheChan <- d:
	default:
	}
}

func (p *BytesPool) Dispose() {
	close(p.cacheChan)
}
