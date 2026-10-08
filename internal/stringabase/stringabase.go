package stringabase

import "sync"

type Stringabase struct {
	database map[string]string
	mu       sync.RWMutex
}
