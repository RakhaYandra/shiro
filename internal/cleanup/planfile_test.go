package cleanup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlanFile_Roundtrip(t *testing.T) {
	p := CleanupPlan{
		ID: "deadbeef", Root: "/x", CreatedAt: 7, SchemaVersion: 1,
		Items: []CleanupItem{{
			Path: "/x/node_modules", Action: "remove", Reason: "r",
			Risk: "review", ExpectedBytes: 5, RuleID: "DET-003",
			ExpectedSize: 5, ExpectedModTime: 6,
		}},
	}
	path := filepath.Join(t.TempDir(), "plan.json")
	if err := SavePlan(path, p); err != nil {
		t.Fatal(err)
	}
	got, err := LoadPlan(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != p.ID || len(got.Items) != 1 || got.Items[0].Path != "/x/node_modules" {
		t.Fatalf("bad roundtrip: %+v", got)
	}
}

func TestPlanFile_RejectsBadSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plan.json")
	if err := os.WriteFile(path, []byte(`{"schema_version":99,"id":"x"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadPlan(path); err == nil {
		t.Fatal("want error for unknown schema")
	}
}

func TestPlanFile_RejectsCorrupt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plan.json")
	if err := os.WriteFile(path, []byte(`{not json`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadPlan(path); err == nil {
		t.Fatal("want error for corrupt file")
	}
}
