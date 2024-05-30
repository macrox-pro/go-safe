package safe

import (
	"slices"
	"sync"
)

type SliceIterFn[T any] func(index int, elem T) bool

type Slice[T any] struct {
	data []T
	mu   sync.RWMutex
}

func (s *Slice[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

func (s *Slice[T]) Get(index int) (elem T, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if index >= 0 && index < len(s.data) {
		elem, ok = s.data[index], true
	}

	return
}

func (s *Slice[T]) Set(index int, elem T) (ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index >= 0 && index < len(s.data) {
		s.data[index], ok = elem, true
	}

	return
}

func (s *Slice[T]) Append(elems ...T) *Slice[T] {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make([]T, 0)
	}

	s.data = append(s.data, elems...)

	return s
}

func (s *Slice[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data != nil {
		clear(s.data)
	}
}

func (s *Slice[T]) Scan(iter SliceIterFn[T]) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.data == nil {
		return
	}

	for index, elem := range s.data {
		if !iter(index, elem) {
			return
		}
	}
}

func (s *Slice[T]) Sort(cmp func(a, b T) int) *Slice[T] {
	s.mu.Lock()
	defer s.mu.Unlock()

	slices.SortFunc(s.data, cmp)

	return s
}

func NewSlice[T any]() *Slice[T] {
	return &Slice[T]{}
}
