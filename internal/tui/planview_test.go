package tui

import (
	"strings"
	"testing"

	"github.com/RakhaYandra/shiro/internal/cleanup"
)

func TestPlanTab_Preview(t *testing.T) {
	m := NewModel(testSummary(), nil)
	m.SetPlan(cleanup.CleanupPlan{
		ID: "p1", Root: "/x", CreatedAt: 1, SchemaVersion: 1,
		Items: []cleanup.CleanupItem{{
			Path: "/x/b", Action: "remove", Reason: "r", Risk: "review",
			ExpectedSize: 700, RuleID: "DET-003",
		}},
	})
	out := m.plan.View()
	for _, want := range []string{"DRY-RUN", "/x/b", "DET-003"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "execute") || strings.Contains(out, "Execute") {
		t.Fatalf("plan tab must not mention execution:\n%s", out)
	}
}
