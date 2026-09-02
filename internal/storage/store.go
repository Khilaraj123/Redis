// In-memory datastore
// Top-level store: map[string]Value + RWMutex
package storage

import "sync"

type Store struct {
	mu   sync.RWMutex
	data map[string]Value
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Value),
	}
}
