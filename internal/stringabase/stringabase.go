// Package stringabase implements the basic in-memory key-value store.
package stringabase

import (
	"sync"
)

// Stringabase is a concurrency-safe in-memory string key-value store. It is accessed through Get(), and
// Set(). Set() will clobber existing stored values and allows for empty string values. The zero value is not
// usable: build with NewStringabase().
type Stringabase struct {
	database map[string]string
	mu       *sync.RWMutex
}

// NewStringabase returns the pointer to a freshly initialized Stringabase.
func NewStringabase() *Stringabase {
	return &Stringabase{
		database: make(map[string]string),
		mu:       &sync.RWMutex{},
	}
}

// Get returns a key's value from the Stringabase. If the key isn't present, it returns an empty string and
// false.
func (s *Stringabase) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.database[key]
	return value, ok
}

// Set sets the key value pair in the database. If the key is present, its existing value is clobbered.
func (s *Stringabase) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.database[key] = value // clobber existing
}
