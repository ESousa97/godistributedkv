// Package storage provides an in-memory thread-safe key-value store with
// persistence support via a Write-Ahead Log (WAL).
//
// The primary entry point is [Store], which handles all data operations.
package storage

import (
	"sync"
)

// Store represents an in-memory thread-safe key-value store.
// It supports persistence by delegating write operations to a [WAL] instance.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
	wal  *WAL
}

// NewStore initializes and returns a new [Store] instance.
// It requires a [WAL] to handle data persistence and recovery.
func NewStore(wal *WAL) *Store {
	return &Store{
		data: make(map[string]string),
		wal:  wal,
	}
}

// Recover loads the state from the associated [WAL].
// If no WAL is configured, it returns nil immediately.
func (s *Store) Recover() error {
	if s.wal == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.wal.Recover(s.data)
}

// Get retrieves a value from the store by its key.
// It returns the value and a boolean indicating if the key was found.
// This operation is thread-safe and utilizes a read lock.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	return val, ok
}

// Set stores a value associated with a key in the store.
// If a [WAL] is configured, it synchronously persists the operation before
// updating the in-memory state.
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
// If a [WAL] is configured, it synchronously persists the deletion before
// updating the in-memory state.
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
