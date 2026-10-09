package stringabase

import "sync"

// TODO: Ousterhout: This is where a thorough and complete decription of
// your "class" belongs.
// Stringbase is a blah blah blah blah
type Stringabase struct {
	database map[string]string
	mu       *sync.RWMutex
}

func NewStringabase() *Stringabase {
	return &Stringabase{
		database: make(map[string]string),
		mu:       &sync.RWMutex{},
	}
}

func (s Stringabase) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.database[key]
	return value, ok
}

func (s *Stringabase) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.database[key] = value // clobber existing
}
