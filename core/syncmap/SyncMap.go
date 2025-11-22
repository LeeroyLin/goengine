package syncmap

import (
	"sync"
	"sync/atomic"
)

type SyncMap[K interface{}, V interface{}] struct {
	m         *sync.Map
	weakCount atomic.Int32 // 不严格计数，短时间窗口计数会有错误
}

func NewSyncMap[K interface{}, V interface{}]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m:         &sync.Map{},
		weakCount: atomic.Int32{},
	}
}

func (sm *SyncMap[K, V]) Add(key K, value V) {
	_, present := sm.m.Swap(key, value)

	// 没有旧值，是新增
	if !present {
		sm.weakCount.Add(1)
	}
}

func (sm *SyncMap[K, V]) Delete(key K) {
	_, present := sm.m.LoadAndDelete(key)

	// 之前有值
	if present {
		sm.weakCount.Add(-1)
	}
}

func (sm *SyncMap[K, V]) Get(key K) (V, bool) {
	var v V

	res, ok := sm.m.Load(key)
	if ok {
		v = res.(V)
	}

	return v, ok
}

func (sm *SyncMap[K, V]) Clear() {
	sm.m.Clear()
	sm.weakCount.Store(0)
}

func (sm *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	sm.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

// Count 真实尺寸，会遍历，效率低
func (sm *SyncMap[K, V]) Count() int {
	size := 0
	sm.m.Range(func(key, value any) bool {
		size++

		return true
	})

	return size
}

// WeakCount 弱计数，短时间窗口内不太准确，大致反应数量
func (sm *SyncMap[K, V]) WeakCount() int32 {
	return sm.weakCount.Load()
}

func (sm *SyncMap[K, V]) GetAndDelete(key K) (value V, loaded bool) {
	var v V

	res, ok := sm.m.LoadAndDelete(key)
	if ok {
		v = res.(V)
	}

	return v, ok
}

func (sm *SyncMap[K, V]) GetOrAdd(key K, value V) (actual V, loaded bool) {
	var v V

	res, ok := sm.m.LoadOrStore(key, value)
	if ok {
		v = res.(V)
	}

	return v, ok
}

func (sm *SyncMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return sm.m.CompareAndDelete(key, old)
}

func (sm *SyncMap[K, V]) CompareAndSwap(key K, old V, new V) (swapped bool) {
	return sm.m.CompareAndSwap(key, old, new)
}

func (sm *SyncMap[K, V]) Swap(key K, new V) (previous V, swaped bool) {
	var v V

	res, ok := sm.m.Swap(key, new)
	if ok {
		v = res.(V)
	}

	return v, ok
}
