package output

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/RakhaYandra/shiro/internal/cleanup"
)

func dryPlan() cleanup.CleanupPlan {
	return cleanup.CleanupPlan{
		ID: "abc", Root: "/x", CreatedAt: 1, SchemaVersion: 1,
		Items: []cleanup.CleanupItem{{
			Path: "/x/node_modules", Action: "remove", Reason: "r",
			Risk: "review", ExpectedBytes: 500, RuleID: "DET-003",
			ExpectedSize: 500, ExpectedModTime: 2,
		}},
		Denied: []cleanup.Denied{{Path: "/x/.ssh", RuleID: "DET-X", Reason: "protected"}},
	}
}

func TestRenderPlanText(t *testing.T) {
	out := RenderPlanText(dryPlan())
	for _, want := range []string{"DRY-RUN", "DET-003", "/x/node_modules", "DENIED", "/x/.ssh"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
}

func TestRenderPlanJSON(t *testing.T) {
	var v map[string]any
	if err := json.Unmarshal(RenderPlanJSON(dryPlan()), &v); err != nil {
		t.Fatal(err)
	}
	if v["schema_version"] != float64(1) || v["id"] != "abc" {
		t.Fatalf("bad meta: %v", v)
	}
	if _, ok := v["items"]; !ok {
		t.Fatal("missing items")
	}
	if _, ok := v["denied"]; !ok {
		t.Fatal("missing denied")
	}
}
