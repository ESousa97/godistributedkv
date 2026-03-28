package storage

import (
	"sync"
)

// Store represents an in-memory thread-safe key-value store.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
	wal  *WAL
}

// NewStore initializes and returns a new Store instance.
func NewStore(wal *WAL) *Store {
	return &Store{
		data: make(map[string]string),
		wal:  wal,
	}
}

// Recover loads the state from the WAL.
func (s *Store) Recover() error {
	if s.wal == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.wal.Recover(s.data)
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
func (s *Store) Set(key, value string) error {
	if s.wal != nil {
		if err := s.wal.Append(key, value, false); err != nil {
			return err
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return nil
}

// Delete removes a key and its associated value from the store.
func (s *Store) Delete(key string) error {
	if s.wal != nil {
		if err := s.wal.Append(key, "", true); err != nil {
			return err
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
	return nil
}
