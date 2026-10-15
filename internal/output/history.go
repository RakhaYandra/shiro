package output

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/RakhaYandra/shiro/internal/analysis"
	"github.com/RakhaYandra/shiro/internal/store/sqlite"
)

func RenderHistoryText(root string, scans []sqlite.ScanMeta) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Shiro · History for %s\n\n", root)
	if len(scans) == 0 {
		b.WriteString("No scans recorded. Run `shiro scan` without --no-store.\n")
		return b.String()
	}
	for _, s := range scans {
		short := s.ID
		if len(short) > 12 {
			short = short[:12]
		}
		fmt.Fprintf(&b, "  %s  %8s  %-8s  %s\n",
			short, Human(s.TotalBytes), s.Status,
			time.Unix(s.StartedAt, 0).Format("2006-01-02 15:04"))
	}
	return b.String()
}

func RenderHistoryJSON(root string, scans []sqlite.ScanMeta) []byte {
	if scans == nil {
		scans = []sqlite.ScanMeta{}
	}
	raw, _ := json.MarshalIndent(map[string]any{
		"schema_version": 1,
		"root":           root,
		"scans":          scans,
	}, "", "  ")
	return append(raw, '\n')
}

func RenderDiffText(from, to string, d analysis.HistoryDiff) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Shiro · Diff %s → %s  (total %+d)\n", shortID(from), shortID(to), d.TotalDelta)
	printChanges(&b, "Grown", d.Grown, true)
	printChanges(&b, "Added", d.Added, true)
	printChanges(&b, "Shrunk", d.Shrunk, false)
	if len(d.Removed) > 0 {
		b.WriteString("\nRemoved\n")
		for _, c := range d.Removed {
			fmt.Fprintf(&b, "  %8s  %s\n", Human(c.Old), c.Path)
		}
	}
	if len(d.Added)+len(d.Removed)+len(d.Grown)+len(d.Shrunk) == 0 {
		b.WriteString("\nNo changes.\n")
	}
	return b.String()
}

func printChanges(b *strings.Builder, title string, cs []analysis.Change, _ bool) {
	if len(cs) == 0 {
		return
	}
	fmt.Fprintf(b, "\n%s\n", title)
	for _, c := range cs {
		fmt.Fprintf(b, "  %+d  %s\n", c.Delta, c.Path)
	}
}

func shortID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

func RenderDiffJSON(from, to string, d analysis.HistoryDiff) []byte {
	raw, _ := json.MarshalIndent(map[string]any{
		"schema_version": 1,
		"from":           from,
		"to":             to,
		"total_delta":    d.TotalDelta,
		"added":          d.Added,
		"removed":        d.Removed,
		"grown":          d.Grown,
		"shrunk":         d.Shrunk,
	}, "", "  ")
	return append(raw, '\n')
}
