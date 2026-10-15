package output

import (
	"fmt"
	"sort"
	"strings"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

// DisplayOpts controls rendering only (never scan scope).
type DisplayOpts struct {
	Depth   int
	MinSize int64
}

// RenderText prints a human-readable summary with top children.
func RenderText(s scanner.Summary, o DisplayOpts) string {
	var b strings.Builder
	b.WriteString("Shiro · Storage Intelligence\n\n")
	fmt.Fprintf(&b, "Scanned  %s (%d files, %d dirs)\n", Human(s.TotalApparent), s.Files, s.Dirs)
	if s.AllocatedKnown {
		fmt.Fprintf(&b, "Allocated %s\n", Human(s.TotalAllocated))
	}
	if s.Partial {
		fmt.Fprintf(&b, "[PARTIAL] %d errors (report incomplete)\n", len(s.Errors))
		for _, e := range s.Errors {
			fmt.Fprintf(&b, "  ! %s: %s\n", e.Path, e.Err)
		}
	}
	if s.Tree != nil {
		b.WriteString("\nTop:\n")
		renderLevel(&b, s.Tree.Children, o, 1)
	}
	return b.String()
}

func renderLevel(b *strings.Builder, kids []*scanner.DirNode, o DisplayOpts, level int) {
	if o.Depth > 0 && level > o.Depth {
		return
	}
	sorted := append([]*scanner.DirNode(nil), kids...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Apparent > sorted[j].Apparent })
	for _, k := range sorted {
		if k.Apparent < o.MinSize {
			continue
		}
		fmt.Fprintf(b, "  %8s  %s\n", Human(k.Apparent), k.Path)
		renderLevel(b, k.Children, o, level+1)
	}
}
