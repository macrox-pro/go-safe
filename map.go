package safe

import "sync"

type MapIterFn[K comparable, V any] func(key K, value V) bool

type Map[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.data)
}

func (m *Map[K, V]) Get(key K) (v V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.data != nil {
		v, ok = m.data[key]
	}

	return
}

func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data == nil {
		m.data = make(map[K]V)
	}

	m.data[key] = value
}

func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data != nil {
		delete(m.data, key)
	}
}

func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data != nil {
		clear(m.data)
	}
}

func (m *Map[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.data == nil {
		return nil
	}

	keys := make([]K, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}

	return keys
}

func (m *Map[K, V]) Values() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.data == nil {
		return nil
	}

	values := make([]V, 0, len(m.data))
	for _, value := range m.data {
		values = append(values, value)
	}

	return values
}

func (m *Map[K, V]) Scan(iter MapIterFn[K, V]) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.data == nil {
		return
	}

	for key, value := range m.data {
		if !iter(key, value) {
			return
		}
	}
}

func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{}
}
