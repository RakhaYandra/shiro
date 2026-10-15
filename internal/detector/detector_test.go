package detector

import (
	"testing"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

func tree() *scanner.DirNode {
	return &scanner.DirNode{
		Path:     "/x",
		Apparent: 1000,
		Children: []*scanner.DirNode{
			{Path: "/x/node_modules", Apparent: 700},
			{Path: "/x/src", Apparent: 300},
		},
	}
}

func TestRegistry_Match(t *testing.T) {
	reg := NewRegistry()
	reg.Register(Rule{
		ID:       "DET-T",
		Category: "test",
		Basename: "node_modules",
		Markers:  []string{"package.json"},
	})
	sum := scanner.Summary{Root: "/x", Tree: tree()}
	got := reg.Detect(sum)
	if len(got) != 1 {
		t.Fatalf("findings = %d, want 1", len(got))
	}
	f := got[0]
	if f.RuleID != "DET-T" || f.Path != "/x/node_modules" || f.Bytes != 700 {
		t.Fatalf("bad finding: %+v", f)
	}
	if f.Risk != "review" {
		t.Fatalf("risk = %q, want review", f.Risk)
	}
	if f.Reason == "" {
		t.Fatal("missing reason")
	}
}

func TestRegistry_NoMatch(t *testing.T) {
	reg := NewRegistry()
	reg.Register(Rule{ID: "DET-T", Category: "test", Basename: "nope"})
	sum := scanner.Summary{Root: "/x", Tree: tree()}
	if got := reg.Detect(sum); len(got) != 0 {
		t.Fatalf("findings = %d, want 0", len(got))
	}
}
