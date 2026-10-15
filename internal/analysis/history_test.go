package analysis

import "testing"

func TestDiff_GrownShrunkAddedRemoved(t *testing.T) {
	old := map[string]int64{"/x": 1000, "/x/grow": 100, "/x/shrink": 400, "/x/gone": 50}
	newm := map[string]int64{"/x": 1200, "/x/grow": 500, "/x/shrink": 100, "/x/fresh": 30}
	d := Diff(old, newm)
	if len(d.Added) != 1 || d.Added[0].Path != "/x/fresh" {
		t.Fatalf("added = %+v", d.Added)
	}
	if len(d.Removed) != 1 || d.Removed[0].Path != "/x/gone" {
		t.Fatalf("removed = %+v", d.Removed)
	}
	if len(d.Grown) != 1 || d.Grown[0].Path != "/x/grow" || d.Grown[0].Delta != 400 {
		t.Fatalf("grown = %+v", d.Grown)
	}
	if len(d.Shrunk) != 1 || d.Shrunk[0].Path != "/x/shrink" {
		t.Fatalf("shrunk = %+v", d.Shrunk)
	}
	if d.TotalDelta != 200 {
		t.Fatalf("total delta = %d, want 200", d.TotalDelta)
	}
}

func TestDiff_TopGrownFirst(t *testing.T) {
	old := map[string]int64{"/x/a": 10, "/x/b": 10}
	newm := map[string]int64{"/x/a": 20, "/x/b": 1010}
	d := Diff(old, newm)
	if len(d.Grown) != 2 || d.Grown[0].Path != "/x/b" {
		t.Fatalf("not sorted desc: %+v", d.Grown)
	}
}

func TestDiff_Empty(t *testing.T) {
	d := Diff(map[string]int64{}, map[string]int64{})
	if len(d.Added) != 0 || len(d.Removed) != 0 || d.TotalDelta != 0 {
		t.Fatalf("want empty diff: %+v", d)
	}
}
