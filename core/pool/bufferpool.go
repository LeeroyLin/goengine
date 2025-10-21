package pool

import "bytes"

type BufferPool struct {
	BytesSize int
	cacheChan chan *bytes.Buffer
	bytesPool *BytesPool
}

func NewBufferPool(bytesSize, max int) *BufferPool {
	p := &BufferPool{
		BytesSize: bytesSize,
		cacheChan: make(chan *bytes.Buffer, max),
		bytesPool: NewBytesPool(bytesSize, max),
	}

	return p
}

func (p *BufferPool) Get() *bytes.Buffer {
	select {
	case b := <-p.cacheChan:
		b.Reset()
		return b
	default:
		b := bytes.NewBuffer(p.bytesPool.Get())
		return b
	}
}

func (p *BufferPool) Set(b *bytes.Buffer) {
	select {
	case p.cacheChan <- b:
	default:
	}
}

func (p *BufferPool) Dispose() {
	p.bytesPool.Dispose()
	close(p.cacheChan)
}
