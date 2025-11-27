package pool

import (
	"errors"
	"fmt"
	"math"
)

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

// 添加一些uint值
func (b *PFMBuffer) putSameUintVal(l int, littleEndian bool, handler func(m int) byte) error {
	err := b.MakeSureCap(l)
	if err != nil {
		return err
	}

	for i := 0; i < l; i++ {
		m := i * 8
		if !littleEndian {
			m = (l - i - 1) * 8
		}

		err = b.WriteByte(handler(m))
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *PFMBuffer) PutUint16(v uint16, littleEndian bool) error {
	return b.putSameUintVal(2, littleEndian, func(m int) byte {
		return byte(v >> m)
	})
}

func (b *PFMBuffer) PutUint32(v uint32, littleEndian bool) error {
	return b.putSameUintVal(4, littleEndian, func(m int) byte {
		return byte(v >> m)
	})
}

func (b *PFMBuffer) PutUint64(v uint64, littleEndian bool) error {
	return b.putSameUintVal(4, littleEndian, func(m int) byte {
		return byte(v >> m)
	})
}

func (b *PFMBuffer) WriteBasicVal(data any, littleEndian bool) error {
	switch v := data.(type) {
	case *bool:
		if *v {
			return b.WriteByte(1)
		} else {
			return b.WriteByte(0)
		}
	case bool:
		if v {
			return b.WriteByte(1)
		} else {
			return b.WriteByte(0)
		}
	case []bool:
		err := b.MakeSureCap(len(v))
		if err != nil {
			return err
		}
		for _, x := range v {
			if x {
				err = b.WriteByte(1)
				if err != nil {
					return err
				}
			} else {
				err = b.WriteByte(0)
				if err != nil {
					return err
				}
			}
		}
	case *int8:
		return b.WriteByte(byte(*v))
	case int8:
		return b.WriteByte(byte(v))
	case []int8:
		err := b.MakeSureCap(len(v))
		if err != nil {
			return err
		}
		for _, x := range v {
			err = b.WriteByte(byte(x))
			if err != nil {
				return err
			}
		}
	case *uint8:
		return b.WriteByte(*v)
	case uint8:
		return b.WriteByte(v)
	case []uint8:
		_, err := b.Write(v)
		return err
	case *int16:
		return b.PutUint16(uint16(*v), littleEndian)
	case int16:
		return b.PutUint16(uint16(v), littleEndian)
	case []int16:
		for _, x := range v {
			err := b.PutUint16(uint16(x), littleEndian)
			if err != nil {
				return err
			}
		}
	case *uint16:
		return b.PutUint16(*v, littleEndian)
	case uint16:
		return b.PutUint16(v, littleEndian)
	case []uint16:
		for _, x := range v {
			err := b.PutUint16(x, littleEndian)
			if err != nil {
				return err
			}
		}
	case *int32:
		return b.PutUint32(uint32(*v), littleEndian)
	case int32:
		return b.PutUint32(uint32(v), littleEndian)
	case []int32:
		for _, x := range v {
			err := b.PutUint32(uint32(x), littleEndian)
			if err != nil {
				return err
			}
		}
	case *uint32:
		return b.PutUint32(*v, littleEndian)
	case uint32:
		return b.PutUint32(v, littleEndian)
	case []uint32:
		for _, x := range v {
			err := b.PutUint32(x, littleEndian)
			if err != nil {
				return err
			}
		}
	case *int64:
		return b.PutUint64(uint64(*v), littleEndian)
	case int64:
		return b.PutUint64(uint64(v), littleEndian)
	case []int64:
		for _, x := range v {
			err := b.PutUint64(uint64(x), littleEndian)
			if err != nil {
				return err
			}
		}
	case *uint64:
		return b.PutUint64(*v, littleEndian)
	case uint64:
		return b.PutUint64(v, littleEndian)
	case []uint64:
		for _, x := range v {
			err := b.PutUint64(x, littleEndian)
			if err != nil {
				return err
			}
		}
	case *float32:
		return b.PutUint32(math.Float32bits(*v), littleEndian)
	case float32:
		return b.PutUint32(math.Float32bits(v), littleEndian)
	case []float32:
		for _, x := range v {
			err := b.PutUint32(math.Float32bits(x), littleEndian)
			if err != nil {
				return err
			}
		}
	case *float64:
		return b.PutUint64(math.Float64bits(*v), littleEndian)
	case float64:
		return b.PutUint64(math.Float64bits(v), littleEndian)
	case []float64:
		for _, x := range v {
			err := b.PutUint64(math.Float64bits(x), littleEndian)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New(fmt.Sprintf("pfm buffer not match val type: %v", v))
	}

	return nil
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

// MakeSureCap 确保容量足够，不够提前开辟
func (b *PFMBuffer) MakeSureCap(addLen int) error {
	delta := addLen + b.Len() - b.Cap()

	if delta > 0 {
		err := b.grow(delta)
		if err != nil {
			return err
		}
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
