package storage

import (
	"testing"
)

func TestStore_SetAndGet(t *testing.T) {
	s := NewStore(nil)
	key := "test-key"
	value := "test-value"

	s.Set(key, value)

	got, found := s.Get(key)
	if !found {
		t.Errorf("Get(%q) found = false, want true", key)
	}
	if got != value {
		t.Errorf("Get(%q) = %q, want %q", key, got, value)
	}
}

func TestStore_GetNotFound(t *testing.T) {
	s := NewStore(nil)
	key := "non-existent"

	_, found := s.Get(key)
	if found {
		t.Errorf("Get(%q) found = true, want false", key)
	}
}

func TestStore_Delete(t *testing.T) {
	s := NewStore(nil)
	key := "to-delete"
	value := "some-value"

	s.Set(key, value)
	s.Delete(key)

	_, found := s.Get(key)
	if found {
		t.Errorf("Get(%q) found = true after Delete, want false", key)
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	s := NewStore(nil)
	const iterations = 1000

	// Concurrent writes
	done := make(chan bool)
	go func() {
		for i := 0; i < iterations; i++ {
			s.Set("key", "value")
		}
		done <- true
	}()

	// Concurrent reads
	go func() {
		for i := 0; i < iterations; i++ {
			s.Get("key")
		}
		done <- true
	}()

	<-done
	<-done
}
