package cleanup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RakhaYandra/shiro/internal/policy"
)

func TestExecute_OnlyApprovedRemoved(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "node_modules")
	keep := filepath.Join(root, "src")
	os.MkdirAll(filepath.Join(target, "sub"), 0o755)
	os.WriteFile(filepath.Join(target, "sub", "a.js"), make([]byte, 100), 0o644)
	os.MkdirAll(keep, 0o755)
	os.WriteFile(filepath.Join(keep, "main.go"), []byte("x"), 0o644)
	fi, _ := os.Lstat(target)
	sum := int64(100)
	plan := CleanupPlan{
		ID: "p1", Root: root, CreatedAt: 1, SchemaVersion: 1,
		Items: []CleanupItem{{
			Path: target, Action: "remove", Reason: "t", Risk: "review",
			ExpectedBytes: sum, RuleID: "DET-003",
			ExpectedSize: sum, ExpectedModTime: fi.ModTime().Unix(),
		}},
	}
	rep := Execute(policy.Default(t.TempDir()), plan)
	if len(rep.Results) != 1 {
		t.Fatalf("results = %d", len(rep.Results))
	}
	r := rep.Results[0]
	if r.Status != "success" {
		t.Fatalf("status = %q (%s)", r.Status, r.Reason)
	}
	if _, err := os.Lstat(target); !os.IsNotExist(err) {
		t.Fatal("target not removed")
	}
	if _, err := os.Lstat(filepath.Join(keep, "main.go")); err != nil {
		t.Fatal("neighbor removed!")
	}
}

func TestExecute_StaleSkipped(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "node_modules")
	os.MkdirAll(target, 0o755)
	os.WriteFile(filepath.Join(target, "a.js"), make([]byte, 10), 0o644)
	plan := CleanupPlan{
		ID: "p2", Root: root, CreatedAt: 1, SchemaVersion: 1,
		Items: []CleanupItem{{
			Path: target, Action: "remove", Reason: "t", Risk: "review",
			ExpectedBytes: 10, RuleID: "DET-003",
			ExpectedSize: 9999, ExpectedModTime: 1,
		}},
	}
	rep := Execute(policy.Default(t.TempDir()), plan)
	if rep.Results[0].Status != "skipped" {
		t.Fatalf("status = %q, want skipped", rep.Results[0].Status)
	}
	if _, err := os.Lstat(target); err != nil {
		t.Fatal("stale target must not be removed")
	}
}

func TestExecute_ProtectedDenied(t *testing.T) {
	home := t.TempDir()
	ssh := filepath.Join(home, ".ssh")
	os.MkdirAll(ssh, 0o755)
	plan := CleanupPlan{
		ID: "p3", Root: home, CreatedAt: 1, SchemaVersion: 1,
		Items: []CleanupItem{{
			Path: ssh, Action: "remove", Reason: "t", Risk: "review",
			RuleID: "DET-X", ExpectedSize: -1,
		}},
	}
	rep := Execute(policy.Default(home), plan)
	if rep.Results[0].Status != "denied" {
		t.Fatalf("status = %q, want denied", rep.Results[0].Status)
	}
	if _, err := os.Lstat(ssh); err != nil {
		t.Fatal("protected path removed!")
	}
}
