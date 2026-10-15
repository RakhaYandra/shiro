package store

import "github.com/RakhaYandra/shiro/internal/scanner"

// Store persists scan inventory. No file contents, ever.
type Store interface {
	SaveScan(s scanner.Summary, scannerVersion string) (id string, err error)
	Close() error
}
