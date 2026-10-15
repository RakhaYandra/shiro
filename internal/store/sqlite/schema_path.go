package sqlite

import (
	"os"
	"path/filepath"
)

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
