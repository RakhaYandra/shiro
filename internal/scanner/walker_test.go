package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path string, n int) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, make([]byte, n), 0o644); err != nil {
		t.Fatal(err)
	}
}

func collect(t *testing.T, root string, opts Options) ([]Entry, []ScanError) {
	t.Helper()
	entries, errs := Walk(root, opts)
	var out []Entry
	var serrs []ScanError
	errDone := make(chan struct{})
	go func() {
		for e := range errs {
			serrs = append(serrs, e)
		}
		close(errDone)
	}()
	for e := range entries {
		out = append(out, e)
	}
	<-errDone
	return out, serrs
}

func TestWalk_BasicTree(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.txt"), 10)
	writeFile(t, filepath.Join(root, "sub", "b.txt"), 20)
	writeFile(t, filepath.Join(root, "sub", "deep", "c.txt"), 30)

	entries, serrs := collect(t, root, Options{})
	if len(serrs) != 0 {
		t.Fatalf("unexpected errors: %v", serrs)
	}
	var files int
	var bytes int64
	for _, e := range entries {
		if e.Kind == KindFile {
			files++
			bytes += e.ApparentBytes
		}
	}
	if files != 3 {
		t.Fatalf("want 3 files, got %d", files)
	}
	if bytes != 60 {
		t.Fatalf("want 60 bytes, got %d", bytes)
	}
}

func TestWalk_NoFollowSymlink(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	writeFile(t, filepath.Join(outside, "secret.txt"), 100)
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "real.txt"), 5)

	entries, _ := collect(t, root, Options{})
	for _, e := range entries {
		if e.Kind == KindFile && e.Path == filepath.Join(outside, "secret.txt") {
			t.Fatal("walked outside root via symlink")
		}
	}
	foundLink := false
	for _, e := range entries {
		if e.Path == filepath.Join(root, "link") {
			foundLink = true
			if e.Kind != KindSymlink {
				t.Fatalf("link kind = %q, want symlink", e.Kind)
			}
		}
	}
	if !foundLink {
		t.Fatal("symlink entry missing")
	}
}

func TestWalk_BrokenSymlink(t *testing.T) {
	root := t.TempDir()
	if err := os.Symlink(filepath.Join(root, "nope"), filepath.Join(root, "broken")); err != nil {
		t.Fatal(err)
	}
	entries, serrs := collect(t, root, Options{})
	if len(serrs) != 0 {
		t.Fatalf("broken symlink must not be fatal: %v", serrs)
	}
	var links int
	for _, e := range entries {
		if e.Kind == KindSymlink {
			links++
		}
	}
	if links != 1 {
		t.Fatalf("want 1 symlink entry, got %+v", entries)
	}
}

func TestWalk_Exclude(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "keep.txt"), 5)
	writeFile(t, filepath.Join(root, "skip", "x.txt"), 50)
	writeFile(t, filepath.Join(root, "skip2", "y.txt"), 7)

	entries, _ := collect(t, root, Options{Excludes: []string{filepath.Join(root, "skip")}})
	var bytes int64
	for _, e := range entries {
		if e.Kind == KindFile {
			bytes += e.ApparentBytes
		}
	}
	if bytes != 12 {
		t.Fatalf("want 12 bytes after exclude, got %d", bytes)
	}
}

func TestWalk_PermissionDenied_Partial(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "ok.txt"), 5)
	denied := filepath.Join(root, "denied")
	if err := os.MkdirAll(filepath.Join(denied, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(denied, "sub", "x.txt"), 10)
	if err := os.Chmod(denied, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(denied, 0o755)

	_, serrs := collect(t, root, Options{})
	if os.Geteuid() == 0 {
		t.Skip("running as root, permission errors do not occur")
	}
	if len(serrs) == 0 {
		t.Fatal("want partial errors recorded, got none")
	}
}
