package output

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/RakhaYandra/shiro/internal/analysis"
	"github.com/RakhaYandra/shiro/internal/store/sqlite"
)

func TestRenderHistory(t *testing.T) {
	out := RenderHistoryText("/x", []sqlite.ScanMeta{
		{ID: "abc123", StartedAt: 7, Status: "complete", TotalBytes: 1000},
	})
	for _, want := range []string{"/x", "abc123", "1000 B", "complete"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
	var v map[string]any
	if err := json.Unmarshal(RenderHistoryJSON("/x", nil), &v); err != nil {
		t.Fatal(err)
	}
	if v["schema_version"] != float64(1) {
		t.Fatalf("bad schema: %v", v)
	}
}

func TestRenderDiff(t *testing.T) {
	d := analysis.HistoryDiff{
		Added:      []analysis.Change{{Path: "/x/fresh", New: 30, Delta: 30}},
		Grown:      []analysis.Change{{Path: "/x/grow", Old: 100, New: 500, Delta: 400}},
		TotalDelta: 430,
	}
	out := RenderDiffText("a", "b", d)
	for _, want := range []string{"+400", "/x/grow", "/x/fresh", "+430"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
	var v map[string]any
	if err := json.Unmarshal(RenderDiffJSON("a", "b", d), &v); err != nil {
		t.Fatal(err)
	}
	if v["schema_version"] != float64(1) || v["total_delta"] != float64(430) {
		t.Fatalf("bad doc: %v", v)
	}
}
