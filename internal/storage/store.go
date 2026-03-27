package storage

import (
	"sync"
)

// Store represents an in-memory thread-safe key-value store.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewStore initializes and returns a new Store instance.
func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

// Get retrieves a value from the store by its key.
// Returns the value and a boolean indicating if the key was found.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	return val, ok
}

// Set stores a value associated with a key in the store.
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

// Delete removes a key and its associated value from the store.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
