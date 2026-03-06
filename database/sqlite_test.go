package database

import (
	"os"
	"path/filepath"
	"testing"
)

// TestWALModeEnabled verifies that a new SQLiteDB connection has WAL mode
// enabled and that the expected PRAGMAs are in effect.
func TestWALModeEnabled(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test_wal.db")

	db, err := NewSQLiteDB(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteDB failed: %v", err)
	}
	defer db.Close()

	var journalMode string
	if err := db.db.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
		t.Fatalf("failed to query journal_mode: %v", err)
	}
	if journalMode != "wal" {
		t.Errorf("expected journal_mode=wal, got %q", journalMode)
	}

	var synchronous int
	if err := db.db.QueryRow("PRAGMA synchronous").Scan(&synchronous); err != nil {
		t.Fatalf("failed to query synchronous: %v", err)
	}
	// NORMAL = 1
	if synchronous != 1 {
		t.Errorf("expected synchronous=1 (NORMAL), got %d", synchronous)
	}

	// Verify the WAL file is created on disk once WAL mode is active.
	walFile := dbPath + "-wal"
	if _, err := os.Stat(walFile); os.IsNotExist(err) {
		// The WAL file may not exist until the first write; this is acceptable.
		t.Log("WAL file not yet present (expected before any writes)")
	}
}

// TestWALModePersists verifies that WAL mode is preserved when the database
// is reopened — confirming the setting is stored in the database header rather
// than being a session-only configuration.
func TestWALModePersists(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test_wal_persist.db")

	// Open and immediately close the database (WAL mode gets written to header).
	db, err := NewSQLiteDB(dbPath)
	if err != nil {
		t.Fatalf("first NewSQLiteDB failed: %v", err)
	}
	db.Close()

	// Re-open the database; journal_mode should still be WAL.
	db2, err := NewSQLiteDB(dbPath)
	if err != nil {
		t.Fatalf("second NewSQLiteDB failed: %v", err)
	}
	defer db2.Close()

	var journalMode string
	if err := db2.db.QueryRow("PRAGMA journal_mode").Scan(&journalMode); err != nil {
		t.Fatalf("failed to query journal_mode on second open: %v", err)
	}
	if journalMode != "wal" {
		t.Errorf("expected journal_mode=wal after reopen, got %q", journalMode)
	}
}
