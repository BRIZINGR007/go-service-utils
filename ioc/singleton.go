package ioc

import "sync"

// Singleton is a generic holder for any type that should only be initialized once.
type Singleton[T any] struct {
	instance *T
	once     sync.Once
}

// Get initializes the instance only once and returns it.
func (s *Singleton[T]) Get(initializer func() *T) *T {
	s.once.Do(func() {
		s.instance = initializer()
	})
	return s.instance
}
