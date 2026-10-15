package analysis

import (
	"os"
	"path/filepath"
	"testing"
)

func mkdup(t *testing.T, files map[string][]byte) string {
	t.Helper()
	root := t.TempDir()
	for rel, data := range files {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func pad(b byte, n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = b
	}
	return out
}

func TestDuplicates_PairGrouped(t *testing.T) {
	data := pad('a', 2<<20)
	root := mkdup(t, map[string][]byte{"a.bin": data, "sub/b.bin": data})
	rep, err := FindDuplicates(root, DupOptions{MinSize: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Groups) != 1 || len(rep.Groups[0].Files) != 2 {
		t.Fatalf("groups = %+v, want 1 group of 2", rep.Groups)
	}
	if rep.Reclaimable != int64(len(data)) {
		t.Fatalf("reclaimable = %d, want %d", rep.Reclaimable, len(data))
	}
}

func TestDuplicates_SameSizeDifferentContent(t *testing.T) {
	root := mkdup(t, map[string][]byte{"a.bin": pad('a', 2<<20), "b.bin": pad('b', 2<<20)})
	rep, err := FindDuplicates(root, DupOptions{MinSize: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Groups) != 0 {
		t.Fatalf("groups = %+v, want none", rep.Groups)
	}
}

func TestDuplicates_HardlinkExcluded(t *testing.T) {
	root := mkdup(t, map[string][]byte{"a.bin": pad('a', 2<<20)})
	if err := os.Link(filepath.Join(root, "a.bin"), filepath.Join(root, "b.bin")); err != nil {
		t.Skip("hardlinks unsupported:", err)
	}
	rep, err := FindDuplicates(root, DupOptions{MinSize: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Groups) != 0 {
		t.Fatalf("hardlink must not count: %+v", rep.Groups)
	}
}

func TestDuplicates_MinSizeFilters(t *testing.T) {
	data := pad('a', 100)
	root := mkdup(t, map[string][]byte{"a.bin": data, "b.bin": data})
	rep, err := FindDuplicates(root, DupOptions{MinSize: 1 << 20})
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Groups) != 0 {
		t.Fatalf("small files must be skipped: %+v", rep.Groups)
	}
}
