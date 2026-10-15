package cleanup

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RakhaYandra/shiro/internal/detector"
	"github.com/RakhaYandra/shiro/internal/policy"
)

func mkfiles(t *testing.T, files map[string]int) string {
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

func TestPlan_ItemsCarryLiveStat(t *testing.T) {
	root := mkfiles(t, map[string]int{"node_modules/a.js": 500})
	pol := policy.Policy{}
	fs := []detector.Finding{{
		RuleID: "DET-003", Path: filepath.Join(root, "node_modules"),
		Bytes: 500, Category: "node", Risk: "review",
		Confidence: "high", Reason: "test",
	}}
	plan, err := Plan(policy.Default(t.TempDir()), fs, root, Filter{})
	_ = pol
	if err != nil {
		t.Fatal(err)
	}
	if plan.SchemaVersion != 1 || plan.ID == "" {
		t.Fatalf("bad plan meta: %+v", plan)
	}
	if len(plan.Items) != 1 {
		t.Fatalf("items = %d, want 1", len(plan.Items))
	}
	it := plan.Items[0]
	if it.ExpectedSize != 500 || it.RuleID != "DET-003" || it.Risk != "review" {
		t.Fatalf("bad item: %+v", it)
	}
	if it.ExpectedModTime <= 0 {
		t.Fatal("missing mtime snapshot")
	}
}

func TestPlan_DeniedExcluded(t *testing.T) {
	home := t.TempDir()
	ssh := filepath.Join(home, ".ssh")
	if err := os.MkdirAll(ssh, 0o755); err != nil {
		t.Fatal(err)
	}
	fs := []detector.Finding{{
		RuleID: "DET-X", Path: ssh, Bytes: 10,
		Category: "x", Risk: "review", Confidence: "low", Reason: "test",
	}}
	plan, err := Plan(policy.Default(home), fs, home, Filter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Items) != 0 || len(plan.Denied) != 1 {
		t.Fatalf("want 0 items 1 denied, got %+v", plan)
	}
}

func TestPlan_Filters(t *testing.T) {
	root := mkfiles(t, map[string]int{
		"node_modules/a.js": 500,
		"target/app":        100,
	})
	pol := policy.Default(t.TempDir())
	fs := []detector.Finding{
		{RuleID: "DET-003", Path: filepath.Join(root, "node_modules"), Bytes: 500, Category: "node", Risk: "review", Confidence: "high", Reason: "t"},
		{RuleID: "DET-GEN", Path: filepath.Join(root, "target"), Bytes: 100, Category: "build", Risk: "review", Confidence: "low", Reason: "t"},
	}
	plan, err := Plan(pol, fs, root, Filter{Category: "node"})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Items) != 1 || plan.Items[0].RuleID != "DET-003" {
		t.Fatalf("category filter failed: %+v", plan.Items)
	}
	plan, err = Plan(pol, fs, root, Filter{MinSize: 200})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Items) != 1 || plan.Items[0].RuleID != "DET-003" {
		t.Fatalf("min-size filter failed: %+v", plan.Items)
	}
}

func TestPlan_UniqueIDs(t *testing.T) {
	pol := policy.Default(t.TempDir())
	a, _ := Plan(pol, nil, "", Filter{})
	b, _ := Plan(pol, nil, "", Filter{})
	if a.ID == "" || a.ID == b.ID {
		t.Fatal("IDs must be unique non-empty")
	}
}
