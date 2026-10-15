package output

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/RakhaYandra/shiro/internal/analysis"
)

func dupReport() analysis.DupReport {
	return analysis.DupReport{
		Root: "/x",
		Groups: []analysis.DupGroup{{
			Hash: "abc", Size: 2048,
			Files: []string{"/x/a.bin", "/x/b.bin"},
		}},
		Reclaimable: 2048,
		FilesHashed: 2,
	}
}

func TestRenderDupText(t *testing.T) {
	out := RenderDupText(dupReport())
	for _, want := range []string{"/x/a.bin", "/x/b.bin", "2.0 KiB", "abc"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q:\n%s", want, out)
		}
	}
}

func TestRenderDupJSON(t *testing.T) {
	var v map[string]any
	if err := json.Unmarshal(RenderDupJSON(dupReport()), &v); err != nil {
		t.Fatal(err)
	}
	if v["schema_version"] != float64(1) {
		t.Fatalf("bad schema: %v", v)
	}
	if _, ok := v["groups"]; !ok {
		t.Fatal("missing groups")
	}
}
