package output

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

func fixture() scanner.Summary {
	return scanner.Summary{
		Root:          "/tmp/x",
		TotalApparent: 3 << 30,
		Files:         10,
		Dirs:          3,
		Tree: &scanner.DirNode{
			Path:     "/tmp/x",
			Apparent: 3 << 30,
			Children: []*scanner.DirNode{
				{Path: "/tmp/x/b", Apparent: 2 << 30},
				{Path: "/tmp/x/a", Apparent: 1 << 30},
			},
		},
	}
}

func TestHuman(t *testing.T) {
	if got := Human(1536); got != "1.5 KiB" {
		t.Fatalf("got %q", got)
	}
	if got := Human(3 << 30); got != "3.0 GiB" {
		t.Fatalf("got %q", got)
	}
}

func TestRenderText_Sorted(t *testing.T) {
	out := RenderText(fixture(), DisplayOpts{Depth: 2})
	ib := strings.Index(out, "/tmp/x/b")
	ia := strings.Index(out, "/tmp/x/a")
	if ib < 0 || ia < 0 || ib > ia {
		t.Fatalf("not sorted desc:\n%s", out)
	}
}

func TestRenderText_Partial(t *testing.T) {
	s := fixture()
	s.Partial = true
	s.Errors = []scanner.ScanError{{Path: "/tmp/x/nope", Err: "permission denied"}}
	out := RenderText(s, DisplayOpts{Depth: 1})
	if !strings.Contains(out, "PARTIAL") {
		t.Fatalf("missing PARTIAL marker:\n%s", out)
	}
}

func TestRenderJSON_Schema(t *testing.T) {
	raw := RenderJSON(fixture(), "test-v")
	var v map[string]any
	if err := json.Unmarshal(raw, &v); err != nil {
		t.Fatal(err)
	}
	if v["schema_version"] != float64(1) {
		t.Fatalf("schema_version = %v", v["schema_version"])
	}
	for _, k := range []string{"root", "total_apparent_bytes", "files", "dirs", "partial", "scanner_version"} {
		if _, ok := v[k]; !ok {
			t.Fatalf("missing key %q", k)
		}
	}
}
