package utils

import "sync"

type SafeMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: make(map[K]V),
	}
}

func NewSafeMapWithValues[K comparable, V any](values map[K]V) *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: values,
	}
}

func (s *SafeMap[K, V]) Set(k K, v V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[k] = v
}

func (s *SafeMap[K, V]) Get(k K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[k]
	return val, ok
}

func (s *SafeMap[K, V]) Delete(k K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, k)
}

func (s *SafeMap[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *SafeMap[K, V]) ForEach(f func(K, V)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, val := range s.data {
		f(key, val)
	}
}

func (s *SafeMap[K, V]) All() map[K]V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	results := make(map[K]V, 0)
	for k, val := range s.data {
		results[k] = val
	}

	return results
}
