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
		kids := append([]*scanner.DirNode(nil), s.Tree.Children...)
		sort.Slice(kids, func(i, j int) bool { return kids[i].Apparent > kids[j].Apparent })
		for _, k := range kids {
			if k.Apparent < o.MinSize {
				continue
			}
			fmt.Fprintf(&b, "  %8s  %s\n", Human(k.Apparent), k.Path)
		}
	}
	return b.String()
}
