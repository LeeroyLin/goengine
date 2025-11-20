package pool

import "errors"

const (
	maxInt = int(^uint(0) >> 1)
)

var errTooLarge = errors.New("buf length too large")

type PFMBuffer struct {
	buf    []byte
	offset int
}

func NewPFMBuffer(cap int) *PFMBuffer {
	return &PFMBuffer{
		buf:    make([]byte, 0, cap),
		offset: 0,
	}
}

func (b *PFMBuffer) Len() int {
	return b.offset
}

func (b *PFMBuffer) Cap() int {
	return cap(b.buf)
}

func (b *PFMBuffer) Write(p []byte) (int, error) {
	l := len(p)

	delta := l + b.Len() - b.Cap()

	if delta > 0 {
		err := b.grow(delta)
		if err != nil {
			return 0, err
		}
	}

	for i := 0; i < l; i++ {
		if b.offset < len(b.buf) {
			b.buf[b.offset] = p[i]
		} else {
			b.buf = append(b.buf, p[i])
		}
		b.offset++
	}

	return l, nil
}

func (b *PFMBuffer) WriteByte(bt byte) error {
	delta := 1 + b.Len() - b.Cap()

	if delta > 0 {
		err := b.grow(delta)
		if err != nil {
			return err
		}
	}

	if b.offset < len(b.buf) {
		b.buf[b.offset] = bt
	} else {
		b.buf = append(b.buf, bt)
	}
	b.offset++

	return nil
}

func (b *PFMBuffer) WriteUtil(p []byte, l int) error {
	delta := l + b.Len() - b.Cap()

	if delta > 0 {
		err := b.grow(delta)
		if err != nil {
			return err
		}
	}

	for i := 0; i < l; i++ {
		if b.offset < len(b.buf) {
			b.buf[b.offset] = p[i]
		} else {
			b.buf = append(b.buf, p[i])
		}
		b.offset++
	}

	return nil
}

func (b *PFMBuffer) Bytes() ([]byte, int) {
	return b.buf, b.Len()
}

func (b *PFMBuffer) BytesClamp(l int) []byte {
	return b.buf[:l]
}

func (b *PFMBuffer) AvailableBytes() []byte {
	return b.BytesClamp(b.Len())
}

func (b *PFMBuffer) Reset() {
	b.offset = 0
}

func (b *PFMBuffer) grow(delta int) error {
	if delta+b.Cap() > maxInt {
		return errTooLarge
	}

	t := make([]byte, delta)
	b.buf = append(b.buf, t...)

	return nil
}
