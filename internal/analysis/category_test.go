package analysis

import (
	"testing"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

func TestCategorize_MixedTree(t *testing.T) {
	sum := scanner.Summary{
		Root:          "/x",
		TotalApparent: 1000,
		ExtTotals: map[string]int64{
			".mp4": 500,
			".go":  300,
			".zip": 150,
			"":     50,
		},
	}
	cats := Categorize(sum)
	if len(cats) == 0 {
		t.Fatal("no categories")
	}
	if cats[0].Category != "media" || cats[0].Bytes != 500 {
		t.Fatalf("top = %+v, want media/500", cats[0])
	}
	for _, c := range cats {
		if c.RuleID == "" || c.Reason == "" {
			t.Fatalf("missing evidence: %+v", c)
		}
	}
}

func TestCategorize_UnknownExt(t *testing.T) {
	sum := scanner.Summary{Root: "/x", ExtTotals: map[string]int64{".xyz123": 10}}
	cats := Categorize(sum)
	found := false
	for _, c := range cats {
		if c.Category == "other" && c.Bytes == 10 {
			found = true
		}
	}
	if !found {
		t.Fatalf("unknown ext must fall to other: %+v", cats)
	}
}
