package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func mkTree(b *testing.B, dirs, filesPer int) string {
	b.Helper()
	root := b.TempDir()
	for d := 0; d < dirs; d++ {
		dir := filepath.Join(root, fmt.Sprintf("d%04d", d))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			b.Fatal(err)
		}
		for f := 0; f < filesPer; f++ {
			p := filepath.Join(dir, fmt.Sprintf("f%04d.bin", f))
			if err := os.WriteFile(p, make([]byte, 128), 0o644); err != nil {
				b.Fatal(err)
			}
		}
	}
	return root
}

func BenchmarkAggregate_1k(b *testing.B) {
	root := mkTree(b, 10, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Aggregate(root, Options{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAggregate_10k(b *testing.B) {
	root := mkTree(b, 100, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Aggregate(root, Options{}); err != nil {
			b.Fatal(err)
		}
	}
}
