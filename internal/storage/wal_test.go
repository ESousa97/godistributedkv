package storage

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestWAL_Lifecycle(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wal-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	walPath := filepath.Join(tmpDir, "test.wal")

	// 1. Create WAL
	wal, err := NewWAL(walPath)
	if err != nil {
		t.Fatalf("NewWAL failed: %v", err)
	}

	// 2. Append entry
	if aErr := wal.Append("key1", "val1", false); aErr != nil {
		t.Errorf("Append failed: %v", aErr)
	}

	// 3. Verify permissions (0600) - skip on Windows as it doesn't support POSIX permissions
	if runtime.GOOS != "windows" {
		info, sErr := os.Stat(walPath)
		if sErr != nil {
			t.Fatalf("Stat failed: %v", sErr)
		}
		if mode := info.Mode().Perm(); mode != 0600 {
			t.Errorf("WAL file permissions = %o, want 0600", mode)
		}
	}

	// 4. Close
	if cErr := wal.Close(); cErr != nil {
		t.Errorf("Close failed: %v", cErr)
	}

	// 5. Reopen and verify existence
	wal2, err := NewWAL(walPath)
	if err != nil {
		t.Fatalf("Reopening WAL failed: %v", err)
	}
	_ = wal2.Close()
}
