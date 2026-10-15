package store

import "github.com/RakhaYandra/shiro/internal/scanner"

type Store interface {
	SaveScan(s scanner.Summary, scannerVersion string) (id string, err error)
	Close() error
}
