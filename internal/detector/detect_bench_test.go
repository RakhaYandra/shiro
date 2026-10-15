package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

func BenchmarkDetect_1k(b *testing.B) {
	root := b.TempDir()
	for d := 0; d < 100; d++ {
		dir := filepath.Join(root, fmt.Sprintf("proj%03d", d), "node_modules")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			b.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "a.js"), make([]byte, 64), 0o644); err != nil {
			b.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, fmt.Sprintf("proj%03d", d), "package.json"), []byte("{}"), 0o644); err != nil {
			b.Fatal(err)
		}
	}
	sum, err := scanner.Aggregate(root, scanner.Options{})
	if err != nil {
		b.Fatal(err)
	}
	reg := DefaultRegistry()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = reg.Detect(sum)
	}
}
