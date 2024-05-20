package sync

import "sync"

type SafeMap[K comparable, V any] struct {
	data  map[K]V
	mutex sync.RWMutex
}

func (s *SafeMap[K, V]) Put(key K, val V) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data[key] = val
}
func (s *SafeMap[K, V]) Get(key K) (val V) {
	s.mutex.RLock()
	defer s.mutex.Unlock()
	return s.data[key]
}
func (s SafeMap[K, V]) LoadOrStore(key K, newVal V) (val V, loaded bool) {
	s.mutex.RLock()
	res, ok := s.data[key]
	if ok {
		return res, true
	}
	return res, false

}
