package cleanup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RakhaYandra/shiro/internal/policy"
)

func mkitem(t *testing.T, rel string, n int) (string, CleanupItem) {
	t.Helper()
	root := t.TempDir()
	p := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, make([]byte, n), 0o644); err != nil {
		t.Fatal(err)
	}
	fi, _ := os.Lstat(filepath.Dir(p))
	return root, CleanupItem{
		Path: filepath.Dir(p), Action: "remove", Risk: "review",
		RuleID: "DET-003", ExpectedSize: int64(n),
		ExpectedModTime: fi.ModTime().Unix(),
	}
}

func TestRevalidate_OK(t *testing.T) {
	root, it := mkitem(t, "node_modules/a.js", 100)
	v := Revalidate(policy.Default(t.TempDir()), it, root)
	if v.Status != StatusOK {
		t.Fatalf("status = %q, want ok", v.Status)
	}
}

func TestRevalidate_StaleSize(t *testing.T) {
	root, it := mkitem(t, "node_modules/a.js", 100)
	if err := os.WriteFile(filepath.Join(root, "node_modules", "b.js"), make([]byte, 50), 0o644); err != nil {
		t.Fatal(err)
	}
	v := Revalidate(policy.Default(t.TempDir()), it, root)
	if v.Status != StatusStale {
		t.Fatalf("status = %q, want stale", v.Status)
	}
}

func TestRevalidate_StaleMissing(t *testing.T) {
	root, it := mkitem(t, "node_modules/a.js", 100)
	if err := os.RemoveAll(it.Path); err != nil {
		t.Fatal(err)
	}
	_ = root
	v := Revalidate(policy.Default(t.TempDir()), it, root)
	if v.Status != StatusStale {
		t.Fatalf("status = %q, want stale", v.Status)
	}
}

func TestRevalidate_DeniedSymlink(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	link := filepath.Join(root, "node_modules")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}
	it := CleanupItem{Path: link, Action: "remove", Risk: "review", RuleID: "DET-003"}
	v := Revalidate(policy.Default(t.TempDir()), it, root)
	if v.Status != StatusDenied {
		t.Fatalf("status = %q, want denied", v.Status)
	}
}

func TestRevalidate_DeniedProtected(t *testing.T) {
	home := t.TempDir()
	ssh := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(ssh, 0o755); err != nil {
		t.Fatal(err)
	}
	it := CleanupItem{Path: ssh, Action: "remove", Risk: "review", RuleID: "DET-X"}
	v := Revalidate(policy.Default(home), it, home)
	if v.Status != StatusDenied {
		t.Fatalf("status = %q, want denied", v.Status)
	}
}
