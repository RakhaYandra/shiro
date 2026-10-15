package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAggregate_NoDoubleCount(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.txt"), 10)
	writeFile(t, filepath.Join(root, "sub", "b.txt"), 20)

	sum, err := Aggregate(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if sum.TotalApparent != 30 {
		t.Fatalf("total = %d, want 30", sum.TotalApparent)
	}
	if sum.Files != 2 || sum.Dirs < 2 {
		t.Fatalf("files=%d dirs=%d, want 2 files", sum.Files, sum.Dirs)
	}
	if sum.Partial {
		t.Fatal("should not be partial")
	}
}

func TestAggregate_Hardlink(t *testing.T) {
	root := t.TempDir()
	a := filepath.Join(root, "a.bin")
	writeFile(t, a, 100)
	if err := os.Link(a, filepath.Join(root, "b.bin")); err != nil {
		t.Skip("hardlinks unsupported:", err)
	}
	sum, err := Aggregate(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if sum.TotalApparent != 100 {
		t.Fatalf("hardlink double-counted: total = %d, want 100", sum.TotalApparent)
	}
}

func TestAggregate_EmptyDir(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}
	sum, err := Aggregate(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if sum.TotalApparent != 0 {
		t.Fatalf("total = %d, want 0", sum.TotalApparent)
	}
	if sum.Tree == nil || sum.Tree.Path == "" {
		t.Fatal("tree root missing")
	}
}

func TestAggregate_SparseFlag(t *testing.T) {
	root := t.TempDir()
	p := filepath.Join(root, "sparse.bin")
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(10 << 20); err != nil {
		t.Fatal(err)
	}
	f.Close()
	sum, err := Aggregate(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if sum.TotalApparent != 10<<20 {
		t.Fatalf("apparent = %d, want %d", sum.TotalApparent, 10<<20)
	}
	// AllocatedKnown reflects platform data; just require consistency.
	if !sum.AllocatedKnown && sum.TotalAllocated != 0 {
		t.Fatal("unknown allocated must report 0")
	}
}

func TestAggregate_TopFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "small.txt"), 10)
	writeFile(t, filepath.Join(root, "big.txt"), 300)
	writeFile(t, filepath.Join(root, "mid.txt"), 100)
	sum, err := Aggregate(root, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(sum.TopFiles) != 3 {
		t.Fatalf("top files = %d, want 3", len(sum.TopFiles))
	}
	if filepath.Base(sum.TopFiles[0].Path) != "big.txt" || sum.TopFiles[0].ApparentBytes != 300 {
		t.Fatalf("top[0] = %+v, want big.txt", sum.TopFiles[0])
	}
}
