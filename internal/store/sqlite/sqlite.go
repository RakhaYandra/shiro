package sqlite

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"os"
	"time"

	_ "modernc.org/sqlite"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

type Store struct{ db *sql.DB }

type NoopStore struct{}

func (NoopStore) SaveScan(scanner.Summary, string) (string, error) { return "", nil }
func (NoopStore) Close() error                                     { return nil }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		return nil, err
	}
	schema, err := os.ReadFile(schemaPath())
	if err != nil {
		schema = []byte(fallbackSchema)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) SaveScan(sum scanner.Summary, scannerVersion string) (string, error) {
	id := newID()
	now := time.Now().Unix()
	status := "complete"
	if sum.Partial {
		status = "partial"
	}
	tx, err := s.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(
		`INSERT INTO scans(id, root_path, started_at, finished_at, status, scanner_version)
		 VALUES(?, ?, ?, ?, ?, ?)`,
		id, sum.Root, now, now, status, scannerVersion,
	); err != nil {
		return "", err
	}
	put := func(path, kind string, apparent int64, alloc *int64) error {
		_, err := tx.Exec(
			`INSERT INTO entries(scan_id, path, kind, apparent_bytes, allocated_bytes)
			 VALUES(?, ?, ?, ?, ?)`,
			id, path, kind, apparent, alloc,
		)
		return err
	}
	if sum.Tree != nil {
		var alloc *int64
		if sum.AllocatedKnown {
			alloc = &sum.TotalAllocated
		}
		if err := put(sum.Tree.Path, "directory", sum.TotalApparent, alloc); err != nil {
			return "", err
		}
		for _, k := range sum.Tree.Children {
			if err := put(k.Path, "directory", k.Apparent, nil); err != nil {
				return "", err
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return id, nil
}

func (s *Store) EntryCount(scanID string) (int, error) {
	var n int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM entries WHERE scan_id = ?`, scanID).Scan(&n)
	return n, err
}

func (s *Store) InsertEntryOrphan() error {
	_, err := s.db.Exec(
		`INSERT INTO entries(scan_id, path, kind, apparent_bytes) VALUES(?, ?, ?, ?)`,
		"no-such-scan", "/x", "file", 1,
	)
	return err
}

func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

const fallbackSchema = `
CREATE TABLE IF NOT EXISTS scans (
    id TEXT PRIMARY KEY,
    root_path TEXT NOT NULL,
    started_at INTEGER NOT NULL,
    finished_at INTEGER,
    status TEXT NOT NULL,
    scanner_version TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS entries (
    scan_id TEXT NOT NULL,
    path TEXT NOT NULL,
    kind TEXT NOT NULL,
    apparent_bytes INTEGER NOT NULL,
    allocated_bytes INTEGER,
    modified_at INTEGER,
    PRIMARY KEY (scan_id, path),
    FOREIGN KEY (scan_id) REFERENCES scans(id)
);
`
