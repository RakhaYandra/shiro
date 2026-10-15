package analysis

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkFindDuplicates(b *testing.B) {
	root := b.TempDir()
	buf := make([]byte, 2<<20)
	if _, err := rand.Read(buf); err != nil {
		b.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		p := filepath.Join(root, string(rune('a'+i))+".bin")
		if err := os.WriteFile(p, buf, 0o644); err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := FindDuplicates(root, DupOptions{MinSize: 1 << 20}); err != nil {
			b.Fatal(err)
		}
	}
}
