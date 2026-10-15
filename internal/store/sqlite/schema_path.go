package sqlite

import (
	"os"
	"path/filepath"
)

// schemaPath locates migrations/001_init.sql relative to CWD or binary.
// Returns a non-existent path when not found so Open falls back to embedded schema.
func schemaPath() string {
	candidates := []string{
		filepath.Join("migrations", "001_init.sql"),
		filepath.Join("..", "..", "..", "migrations", "001_init.sql"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return filepath.Join("migrations", "001_init.sql")
}
