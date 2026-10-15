package analysis

import (
	"sort"
	"strings"
)

type Change struct {
	Path  string `json:"path"`
	Old   int64  `json:"old_bytes"`
	New   int64  `json:"new_bytes"`
	Delta int64  `json:"delta"`
}

type HistoryDiff struct {
	Added      []Change
	Removed    []Change
	Grown      []Change
	Shrunk     []Change
	TotalDelta int64
}

func Diff(old, newm map[string]int64) HistoryDiff {
	var d HistoryDiff
	for path, nb := range newm {
		ob, ok := old[path]
		if !ok {
			d.Added = append(d.Added, Change{path, 0, nb, nb})
			continue
		}
		if nb > ob {
			d.Grown = append(d.Grown, Change{path, ob, nb, nb - ob})
		} else if nb < ob {
			d.Shrunk = append(d.Shrunk, Change{path, ob, nb, nb - ob})
		}
	}
	for path, ob := range old {
		if _, ok := newm[path]; !ok {
			d.Removed = append(d.Removed, Change{path, ob, 0, -ob})
		}
	}
	byDelta := func(cs []Change) {
		sort.Slice(cs, func(i, j int) bool {
			if cs[i].Delta != cs[j].Delta {
				return cs[i].Delta > cs[j].Delta
			}
			return cs[i].Path < cs[j].Path
		})
	}
	byDelta(d.Added)
	byDelta(d.Grown)
	sort.Slice(d.Shrunk, func(i, j int) bool {
		if d.Shrunk[i].Delta != d.Shrunk[j].Delta {
			return d.Shrunk[i].Delta < d.Shrunk[j].Delta
		}
		return d.Shrunk[i].Path < d.Shrunk[j].Path
	})
	sort.Slice(d.Removed, func(i, j int) bool { return d.Removed[i].Path < d.Removed[j].Path })
	d.TotalDelta = rootDelta(old, newm)
	return d
}

func rootDelta(old, newm map[string]int64) int64 {
	best := ""
	for path := range old {
		if _, ok := newm[path]; !ok {
			continue
		}
		if best == "" || depth(path) < depth(best) {
			best = path
		}
	}
	if best == "" {
		return 0
	}
	return newm[best] - old[best]
}

func depth(path string) int {
	if path == "" {
		return 0
	}
	return strings.Count(path, "/")
}
