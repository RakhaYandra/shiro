package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RakhaYandra/shiro/internal/analysis"
)

func RenderDupText(rep analysis.DupReport) string {
	var b strings.Builder
	b.WriteString("Shiro · Duplicates (report only, nothing removed)\n\n")
	fmt.Fprintf(&b, "Root  %s\nHashed %d files\n", rep.Root, rep.FilesHashed)
	if rep.Partial {
		fmt.Fprintf(&b, "[PARTIAL] %d errors\n", len(rep.Errors))
	}
	if len(rep.Groups) == 0 {
		b.WriteString("\nNo duplicate content found.\n")
		return b.String()
	}
	for _, g := range rep.Groups {
		short := g.Hash
		if len(short) > 12 {
			short = short[:12]
		}
		fmt.Fprintf(&b, "\n%s  %s each (%d copies)\n", Human(g.Size), short, len(g.Files))
		for _, f := range g.Files {
			fmt.Fprintf(&b, "  %s\n", f)
		}
	}
	fmt.Fprintf(&b, "\nReclaimable: %s (keeping one copy per group; not guaranteed)\n",
		Human(rep.Reclaimable))
	return b.String()
}

func RenderDupJSON(rep analysis.DupReport) []byte {
	doc := map[string]any{
		"schema_version": 1,
		"root":           rep.Root,
		"groups":         rep.Groups,
		"reclaimable":    rep.Reclaimable,
		"files_hashed":   rep.FilesHashed,
		"partial":        rep.Partial,
		"errors":         rep.Errors,
	}
	if doc["groups"] == nil {
		doc["groups"] = []analysis.DupGroup{}
	}
	if doc["errors"] == nil {
		doc["errors"] = []struct{}{}
	}
	raw, _ := json.MarshalIndent(doc, "", "  ")
	return append(raw, '\n')
}
