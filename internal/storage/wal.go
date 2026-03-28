package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// LogEntry represents a single operation in the Write-Ahead Log.
// It serializes to JSON for persistence on disk.
type LogEntry struct {
	// Key is the unique identifier for the entry.
	Key string `json:"key"`
	// Value is the data associated with the key. Omitted for delete operations.
	Value string `json:"value,omitempty"`
	// IsDelete indicates if this entry represents a key deletion.
	IsDelete bool `json:"is_delete"`
}

// WAL (Write-Ahead Log) handles sequential log writing and recovery.
// It provides durability by ensuring all write operations are synced to disk
// before being applied to the in-memory [Store].
type WAL struct {
	mu   sync.Mutex
	file *os.File
}

// NewWAL opens or creates a WAL file at the specified path.
// It ensures that the parent directory exists before attempting to open the file.
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

// Append adds a new [LogEntry] to the log and flushes it to disk.
// It returns an error if the serialization or disk write fails.
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

// Recover reads the sequential log from disk and replays it into the provided
// data map. This is typically used during [Store] initialization to restore
// the previous state.
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

// Close closes the underlying WAL file.
func (w *WAL) Close() error {
	return w.file.Close()
}
