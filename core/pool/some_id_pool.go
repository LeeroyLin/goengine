package pool

func NewInt8IdPool(max int8) *IdPool[int8] {
	return NewIdPool[int8](int64(max), func(idx int64) int8 {
		return int8(idx + 1)
	})
}

func NewInt16IdPool(max int16) *IdPool[int16] {
	return NewIdPool[int16](int64(max), func(idx int64) int16 {
		return int16(idx + 1)
	})
}

func NewInt32IdPool(max int) *IdPool[int32] {
	return NewIdPool[int32](int64(max), func(idx int64) int32 {
		return int32(idx + 1)
	})
}

func NewUint8IdPool(max uint8) *IdPool[uint8] {
	return NewIdPool[uint8](int64(max), func(idx int64) uint8 {
		return uint8(idx + 1)
	})
}

func NewUint16IdPool(max uint16) *IdPool[uint16] {
	return NewIdPool[uint16](int64(max), func(idx int64) uint16 {
		return uint16(idx + 1)
	})
}

func NewUint32IdPool(max uint32) *IdPool[uint32] {
	return NewIdPool[uint32](int64(max), func(idx int64) uint32 {
		return uint32(idx + 1)
	})
}
