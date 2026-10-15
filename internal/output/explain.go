package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/RakhaYandra/shiro/internal/analysis"
	"github.com/RakhaYandra/shiro/internal/detector"
	"github.com/RakhaYandra/shiro/internal/scanner"
)

func RenderExplainText(s scanner.Summary, _ string) string {
	var b strings.Builder
	b.WriteString("Shiro · Explain\n\n")
	fmt.Fprintf(&b, "Root     %s\n", s.Root)
	fmt.Fprintf(&b, "Scanned  %s (%d files, %d dirs)\n", Human(s.TotalApparent), s.Files, s.Dirs)
	if s.Partial {
		fmt.Fprintf(&b, "[PARTIAL] %d errors (report incomplete)\n", len(s.Errors))
	}
	b.WriteString("\nCategories\n")
	for _, c := range analysis.Categorize(s) {
		fmt.Fprintf(&b, "  %8s  %-10s  [%s] %s\n", Human(c.Bytes), c.Category, c.RuleID, c.Reason)
	}
	b.WriteString("\nTop directories\n")
	if s.Tree != nil {
		kids := append([]*scanner.DirNode(nil), s.Tree.Children...)
		sort.Slice(kids, func(i, j int) bool { return kids[i].Apparent > kids[j].Apparent })
		for _, k := range kids {
			fmt.Fprintf(&b, "  %8s  %s\n", Human(k.Apparent), k.Path)
		}
	}
	b.WriteString("\nTop files\n")
	for _, f := range s.TopFiles {
		fmt.Fprintf(&b, "  %8s  %s\n", Human(f.ApparentBytes), f.Path)
	}
	b.WriteString("\nDeveloper artifacts\n")
	findings := detector.DefaultRegistry().Detect(s)
	if len(findings) == 0 {
		b.WriteString("  (none detected)\n")
	}
	for _, f := range findings {
		fmt.Fprintf(&b, "  %8s  %-12s  [%s/%s] %s\n  %s\n",
			Human(f.Bytes), f.Category, f.RuleID, f.Confidence, f.Path, "  "+f.Reason)
	}
	r := analysis.Reclaimable(s)
	b.WriteString("\nReclaimable estimate\n")
	if r.Known {
		fmt.Fprintf(&b, "  %s\n  %s\n", Human(r.Bytes), r.Reason)
	} else {
		fmt.Fprintf(&b, "  unknown\n  %s\n", r.Reason)
	}
	b.WriteString("\nNote: cleanup candidates arrive with detectors (next phase).\n")
	return b.String()
}

func RenderExplainJSON(s scanner.Summary, r analysis.Estimate, scannerVersion string) []byte {
	type cat struct {
		Category string `json:"category"`
		Bytes    int64  `json:"bytes"`
		Files    int64  `json:"files"`
		RuleID   string `json:"rule_id"`
		Reason   string `json:"reason"`
	}
	type item struct {
		Path  string `json:"path"`
		Bytes int64  `json:"bytes"`
	}
	type finding struct {
		RuleID     string `json:"rule_id"`
		Path       string `json:"path"`
		Bytes      int64  `json:"bytes"`
		Category   string `json:"category"`
		Risk       string `json:"risk"`
		Confidence string `json:"confidence"`
		Reason     string `json:"reason"`
	}
	doc := map[string]any{
		"schema_version":       1,
		"scanner_version":      scannerVersion,
		"root":                 s.Root,
		"total_apparent_bytes": s.TotalApparent,
		"files":                s.Files,
		"dirs":                 s.Dirs,
		"partial":              s.Partial,
		"errors":               s.Errors,
		"top_dirs":             []item{},
		"top_files":            []item{},
		"categories":           []cat{},
		"findings":             []finding{},
		"reclaimable": map[string]any{
			"bytes":  r.Bytes,
			"known":  r.Known,
			"reason": r.Reason,
		},
	}
	if s.Errors == nil {
		doc["errors"] = []scanner.ScanError{}
	}
	if s.Tree != nil {
		for _, k := range s.Tree.Children {
			doc["top_dirs"] = append(doc["top_dirs"].([]item), item{k.Path, k.Apparent})
		}
	}
	for _, f := range s.TopFiles {
		doc["top_files"] = append(doc["top_files"].([]item), item{f.Path, f.ApparentBytes})
	}
	for _, c := range analysis.Categorize(s) {
		doc["categories"] = append(doc["categories"].([]cat),
			cat{c.Category, c.Bytes, c.Files, c.RuleID, c.Reason})
	}
	for _, f := range detector.DefaultRegistry().Detect(s) {
		doc["findings"] = append(doc["findings"].([]finding),
			finding{f.RuleID, f.Path, f.Bytes, f.Category, f.Risk, f.Confidence, f.Reason})
	}
	raw, _ := json.MarshalIndent(doc, "", "  ")
	return append(raw, '\n')
}
