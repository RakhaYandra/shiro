package output

import (
	"encoding/json"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

// RenderJSON emits stable machine-readable inventory (schema_version 1).
func RenderJSON(s scanner.Summary, scannerVersion string) []byte {
	type child struct {
		Path     string `json:"path"`
		Apparent int64  `json:"apparent_bytes"`
	}
	type doc struct {
		SchemaVersion      int                 `json:"schema_version"`
		ScannerVersion     string              `json:"scanner_version"`
		Root               string              `json:"root"`
		TotalApparentBytes int64               `json:"total_apparent_bytes"`
		TotalAllocBytes    int64               `json:"total_allocated_bytes"`
		AllocatedKnown     bool                `json:"allocated_known"`
		Files              int64               `json:"files"`
		Dirs               int64               `json:"dirs"`
		Symlinks           int64               `json:"symlinks"`
		Partial            bool                `json:"partial"`
		Errors             []scanner.ScanError `json:"errors"`
		TopChildren        []child             `json:"top_children"`
	}
	d := doc{
		SchemaVersion:      1,
		ScannerVersion:     scannerVersion,
		Root:               s.Root,
		TotalApparentBytes: s.TotalApparent,
		TotalAllocBytes:    s.TotalAllocated,
		AllocatedKnown:     s.AllocatedKnown,
		Files:              s.Files,
		Dirs:               s.Dirs,
		Symlinks:           s.Symlinks,
		Partial:            s.Partial,
		Errors:             s.Errors,
	}
	if d.Errors == nil {
		d.Errors = []scanner.ScanError{}
	}
	if s.Tree != nil {
		for _, k := range s.Tree.Children {
			d.TopChildren = append(d.TopChildren, child{k.Path, k.Apparent})
		}
	}
	raw, _ := json.MarshalIndent(d, "", "  ")
	return append(raw, '\n')
}
