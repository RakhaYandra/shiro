package detector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

func mkproj(t *testing.T, files map[string]int) string {
	t.Helper()
	root := t.TempDir()
	for rel, n := range files {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, make([]byte, n), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func byRule(fs []Finding, id string) *Finding {
	for i := range fs {
		if fs[i].RuleID == id {
			return &fs[i]
		}
	}
	return nil
}

func TestBuiltin_Positive(t *testing.T) {
	root := mkproj(t, map[string]int{
		"package.json":                 10,
		"node_modules/a.js":            500,
		"go.mod":                       10,
		"pkg/mod/cache/download/x.zip": 400,
		"pyproject.toml":               10,
		"pkg1/__pycache__/a.pyc":       100,
		"pkg1/__init__.py":             5,
		"pkg1/.venv/pyvenv.cfg":        10,
		"pkg1/.venv/lib/big.so":        200,
		"Cargo.toml":                   10,
		"target/debug/app":             600,
	})
	sum, err := scanner.Aggregate(root, scanner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	got := DefaultRegistry().Detect(sum)
	wantConf := map[string]string{
		"DET-003": "high", "DET-002b": "medium", "DET-004a": "high",
		"DET-004b": "high", "DET-GEN": "high",
	}
	for id, want := range wantConf {
		f := byRule(got, id)
		if f == nil {
			t.Fatalf("missing %s in %+v", id, got)
		}
		if f.Risk != "review" || f.Reason == "" {
			t.Fatalf("bad evidence %+v", f)
		}
		if f.Confidence != want {
			t.Fatalf("%s confidence = %q, want %q", id, f.Confidence, want)
		}
	}
}

func TestBuiltin_Negative(t *testing.T) {
	root := mkproj(t, map[string]int{
		"src/main.go":   50,
		"docs/read.txt": 20,
	})
	sum, err := scanner.Aggregate(root, scanner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if got := DefaultRegistry().Detect(sum); len(got) != 0 {
		t.Fatalf("false positives: %+v", got)
	}
}

func TestBuiltin_NoMarker_Low(t *testing.T) {
	root := mkproj(t, map[string]int{"lonely/target/debug/app": 100})
	sum, err := scanner.Aggregate(root, scanner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	f := byRule(DefaultRegistry().Detect(sum), "DET-GEN")
	if f == nil {
		t.Fatal("want lenient finding without marker")
	}
	if f.Confidence != "low" {
		t.Fatalf("confidence = %q, want low", f.Confidence)
	}
}
