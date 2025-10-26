package syncmap

import "sync"

type SyncMap[K interface{}, V interface{}] struct {
	m *sync.Map
}

func NewSyncMap[K interface{}, V interface{}]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m: &sync.Map{},
	}
}

func (sm *SyncMap[K, V]) Add(key K, value V) {
	sm.m.Store(key, value)
}

func (sm *SyncMap[K, V]) Delete(key K) {
	sm.m.Delete(key)
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
}

func (sm *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	sm.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}

func (sm *SyncMap[K, V]) Size() int {
	size := 0
	sm.m.Range(func(key, value any) bool {
		size++

		return true
	})

	return size
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
