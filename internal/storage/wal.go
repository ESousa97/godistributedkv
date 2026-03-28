package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// LogEntry represents a single operation in the WAL.
type LogEntry struct {
	Key      string `json:"key"`
	Value    string `json:"value,omitempty"`
	IsDelete bool   `json:"is_delete"`
}

// WAL handles sequential log writing and recovery.
type WAL struct {
	mu   sync.Mutex
	file *os.File
}

// NewWAL opens or creates a WAL file.
func NewWAL(path string) (*WAL, error) {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create WAL directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL file: %w", err)
	}

	return &WAL{
		file: f,
	}, nil
}

// Append adds a new entry to the log and syncs to disk.
func (w *WAL) Append(key, value string, isDelete bool) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	entry := LogEntry{
		Key:      key,
		Value:    value,
		IsDelete: isDelete,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	if _, err := w.file.Write(append(data, '\n')); err != nil {
		return err
	}

	return w.file.Sync()
}

// Recover reads the log and replays it into the provided store.
func (w *WAL) Recover(data map[string]string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Move to start of file
	if _, err := w.file.Seek(0, 0); err != nil {
		return err
	}

	scanner := bufio.NewScanner(w.file)
	for scanner.Scan() {
		var entry LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue // Skip corrupted lines
		}

		if entry.IsDelete {
			delete(data, entry.Key)
		} else {
			data[entry.Key] = entry.Value
		}
	}

	return scanner.Err()
}

// Close closes the underlying file.
func (w *WAL) Close() error {
	return w.file.Close()
}
