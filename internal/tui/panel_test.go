package tui

import (
	"strings"
	"testing"

	"github.com/RakhaYandra/shiro/internal/detector"
)

func TestNodeDetail(t *testing.T) {
	m := NewModel(testSummary(), []detector.Finding{{
		RuleID: "DET-003", Path: "/x/b", Bytes: 700, Category: "node",
		Risk: "review", Confidence: "high", Reason: "marker present",
	}})
	out := m.nodeDetail("/x/b")
	for _, want := range []string{"700", "DET-003", "high", "review"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
	out = m.nodeDetail("/x/a")
	if strings.Contains(out, "DET-003") {
		t.Fatalf("finding leaked to unrelated node:\n%s", out)
	}
}
