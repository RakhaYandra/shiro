package output

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/RakhaYandra/shiro/internal/analysis"
	"github.com/RakhaYandra/shiro/internal/scanner"
)

func explainSummary() scanner.Summary {
	return scanner.Summary{
		Root:          "/x",
		TotalApparent: 1000,
		Files:         4,
		Dirs:          2,
		Tree: &scanner.DirNode{
			Path:     "/x",
			Apparent: 1000,
			Children: []*scanner.DirNode{{Path: "/x/sub", Apparent: 600}},
		},
		TopFiles:  []scanner.FileRef{{Path: "/x/sub/big.mp4", ApparentBytes: 500}},
		ExtTotals: map[string]int64{".mp4": 500, ".go": 300},
		ExtCounts: map[string]int64{".mp4": 1, ".go": 3},
	}
}

func TestRenderExplain_HasReasons(t *testing.T) {
	out := RenderExplainText(explainSummary(), "dev")
	for _, want := range []string{"media", "CAT-EXT", "big.mp4", "/x/sub"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
}

func TestRenderExplainJSON_Schema(t *testing.T) {
	raw := RenderExplainJSON(explainSummary(), analysis.Reclaimable(explainSummary()), "dev")
	var v map[string]any
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatal(err)
	}
	for _, k := range []string{"schema_version", "root", "categories", "top_dirs", "top_files", "reclaimable"} {
		if _, ok := v[k]; !ok {
			t.Fatalf("missing key %q", k)
		}
	}
	cats := v["categories"].([]any)
	c0 := cats[0].(map[string]any)
	if _, ok := c0["rule_id"]; !ok {
		t.Fatal("category missing rule_id")
	}
}
