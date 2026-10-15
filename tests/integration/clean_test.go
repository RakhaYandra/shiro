package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RakhaYandra/shiro/internal/cleanup"
	"github.com/RakhaYandra/shiro/internal/detector"
	"github.com/RakhaYandra/shiro/internal/policy"
	"github.com/RakhaYandra/shiro/internal/scanner"
)

func TestIntegration_DryRunMutatesNothing(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "package.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	nm := filepath.Join(root, "node_modules")
	if err := os.MkdirAll(nm, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nm, "a.js"), make([]byte, 500), 0o640); err != nil {
		t.Fatal(err)
	}
	before := snapshot(t, root)

	sum, err := scanner.Aggregate(root, scanner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := cleanup.Plan(policy.Default(t.TempDir()),
		detector.DefaultRegistry().Detect(sum), sum.Root, cleanup.Filter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(plan.Items))
	}
	after := snapshot(t, root)
	if len(before) != len(after) {
		t.Fatal("dry-run changed file count")
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

func TestIntegration_SymlinkEscapeDenied(t *testing.T) {
	home := t.TempDir()
	secret := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(secret, 0o755); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(home, "work")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "node_modules")
	if err := os.Symlink(secret, link); err != nil {
		t.Fatal(err)
	}
	pol := policy.Default(home)
	if err := pol.CheckIn(link, root); err == nil {
		t.Fatal("symlink pointing into ~/.ssh must be denied")
	}
}
