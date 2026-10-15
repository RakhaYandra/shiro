package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func fakeHome(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, d := range []string{".ssh", "Projects/proj/node_modules", "Documents"} {
		if err := os.MkdirAll(filepath.Join(home, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return home
}

func TestDeny_StrictSubtree(t *testing.T) {
	home := fakeHome(t)
	p := Default(home)
	if err := p.Check(filepath.Join(home, ".ssh", "id_ed25519")); err == nil {
		t.Fatal("want deny inside ~/.ssh")
	}
	if err := p.Check(filepath.Join(home, ".ssh")); err == nil {
		t.Fatal("want deny ~/.ssh itself")
	}
}

func TestDeny_RootItself(t *testing.T) {
	home := fakeHome(t)
	p := Default(home)
	if err := p.Check(filepath.Join(home, "Projects")); err == nil {
		t.Fatal("want deny ~/Projects itself")
	}
}

func TestAllow_ArtifactUnderProjects(t *testing.T) {
	home := fakeHome(t)
	p := Default(home)
	if err := p.Check(filepath.Join(home, "Projects", "proj", "node_modules")); err != nil {
		t.Fatalf("want allow, got %v", err)
	}
}

func TestDeny_Traversal(t *testing.T) {
	home := fakeHome(t)
	p := Default(home)
	tricky := filepath.Join(home, "Projects", "proj", "..", "..", ".ssh", "x")
	if err := p.Check(tricky); err == nil {
		t.Fatal("want deny path resolving into ~/.ssh")
	}
}

func TestDeny_OutsideRoot(t *testing.T) {
	home := fakeHome(t)
	p := Default(home)
	if err := p.CheckIn(filepath.Join(home, "Projects", "proj", "node_modules"), filepath.Join(home, "Projects")); err != nil {
		t.Fatalf("want allow inside root, got %v", err)
	}
	if err := p.CheckIn("/etc/passwd", filepath.Join(home, "Projects")); err == nil {
		t.Fatal("want deny outside scan root")
	}
}
