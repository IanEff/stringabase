package stringabase

import "sync"

// Stringabase is just the in-memory key-value store and its attendant accessors.   Hot papa!
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
