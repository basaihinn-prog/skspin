package lazy

import (
	"sync"
)

type IMap[K comparable, V any] interface {
	Load(k K) (V, bool)
	Store(k K, v V)
	Delete(k K)
	Clear()
	Len() int
	Each(fn func(k K, v V))
}

type SyncMap[K comparable, V any] struct {
	m    map[K]V
	lock sync.Mutex
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{m: map[K]V{}}
}

func (m *SyncMap[K, V]) Load(k K) (V, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	v, ok := m.m[k]
	return v, ok
}

func (m *SyncMap[K, V]) Store(k K, v V) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.m[k] = v
}

func (m *SyncMap[K, V]) Delete(k K) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.m, k)
}

func (m *SyncMap[K, V]) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.m = map[K]V{}
}

func (m *SyncMap[K, V]) Len() int {
	m.lock.Lock()
	defer m.lock.Unlock()

	return len(m.m)
}

func (m *SyncMap[K, V]) Each(fn func(k K, v V)) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for k, v := range m.m {
		fn(k, v)
	}
}
