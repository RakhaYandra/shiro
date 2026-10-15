package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

func snapshot(t *testing.T, root string) map[string]os.FileInfo {
	t.Helper()
	out := map[string]os.FileInfo{}
	_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err == nil {
			out[p] = info
		}
		return nil
	})
	return out
}

func TestIntegration_ScanReadOnly(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sub", "b.txt"), []byte("world"), 0o640); err != nil {
		t.Fatal(err)
	}
	before := snapshot(t, root)

	sum, err := scanner.Aggregate(root, scanner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if sum.Files != 2 {
		t.Fatalf("files = %d, want 2", sum.Files)
	}
	after := snapshot(t, root)
	if len(before) != len(after) {
		t.Fatalf("file count changed: %d -> %d", len(before), len(after))
	}
	for p, bi := range before {
		ai, ok := after[p]
		if !ok {
			t.Fatalf("path vanished: %s", p)
		}
		if !ai.ModTime().Equal(bi.ModTime()) || ai.Mode() != bi.Mode() || ai.Size() != bi.Size() {
			t.Fatalf("metadata changed: %s", p)
		}
	}
}

func TestIntegration_PartialReported(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root, permission errors do not occur")
	}
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "ok.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	denied := filepath.Join(root, "denied")
	if err := os.MkdirAll(filepath.Join(denied, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(denied, "sub", "x.txt"), []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(denied, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(denied, 0o755)

	sum, err := scanner.Aggregate(root, scanner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !sum.Partial || len(sum.Errors) == 0 {
		t.Fatal("want partial=true with recorded errors")
	}
}
