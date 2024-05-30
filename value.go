package safe

import (
	"sync"
)

type Value[T any] struct {
	data T
	mu   sync.RWMutex
}

func (v *Value[T]) Load() (val T) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.data
}

func (v *Value[T]) Store(val T) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.data = val
}

func (v *Value[T]) Swap(new T) (old T) {
	v.mu.Lock()
	defer v.mu.Unlock()

	old, v.data = v.data, new
	return
}

func NewValue[T any]() *Value[T] {
	return &Value[T]{}
}
