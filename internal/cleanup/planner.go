package cleanup

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"

	"github.com/RakhaYandra/shiro/internal/detector"
	"github.com/RakhaYandra/shiro/internal/policy"
	"github.com/RakhaYandra/shiro/internal/scanner"
)

type CleanupItem struct {
	Path            string `json:"path"`
	Action          string `json:"action"`
	Reason          string `json:"reason"`
	Risk            string `json:"risk"`
	ExpectedBytes   int64  `json:"expected_bytes"`
	RuleID          string `json:"rule_id"`
	ExpectedSize    int64  `json:"expected_size"`
	ExpectedModTime int64  `json:"expected_mod_time"`
}

type Denied struct {
	Path   string `json:"path"`
	RuleID string `json:"rule_id"`
	Reason string `json:"reason"`
}

type CleanupPlan struct {
	ID            string        `json:"id"`
	Root          string        `json:"root"`
	CreatedAt     int64         `json:"created_at"`
	SchemaVersion int           `json:"schema_version"`
	DryRun        bool          `json:"dry_run,omitempty"`
	Items         []CleanupItem `json:"items"`
	Denied        []Denied      `json:"denied"`
}

type Filter struct {
	Category string
	MinSize  int64
}

func Plan(pol policy.Policy, findings []detector.Finding, root string, f Filter) (CleanupPlan, error) {
	plan := CleanupPlan{
		ID:            newID(),
		Root:          root,
		CreatedAt:     time.Now().Unix(),
		SchemaVersion: 1,
		DryRun:        true,
	}
	for _, fd := range findings {
		if f.Category != "" && fd.Category != f.Category {
			continue
		}
		if fd.Bytes < f.MinSize {
			continue
		}
		dec := policy.Decide(pol, fd, root)
		if !dec.Allowed {
			plan.Denied = append(plan.Denied, Denied{
				Path: fd.Path, RuleID: fd.RuleID, Reason: dec.Reason,
			})
			continue
		}
		size, mtime := snapshot(fd.Path)
		plan.Items = append(plan.Items, CleanupItem{
			Path:            fd.Path,
			Action:          "remove",
			Reason:          fd.Reason + " | policy: " + dec.Reason,
			Risk:            fd.Risk,
			ExpectedBytes:   fd.Bytes,
			RuleID:          fd.RuleID,
			ExpectedSize:    size,
			ExpectedModTime: mtime,
		})
	}
	return plan, nil
}

func snapshot(path string) (int64, int64) {
	fi, err := os.Lstat(path)
	if err != nil {
		return -1, -1
	}
	if !fi.IsDir() {
		return fi.Size(), fi.ModTime().Unix()
	}
	sum, err := scanner.Aggregate(path, scanner.Options{})
	if err != nil {
		return -1, fi.ModTime().Unix()
	}
	return sum.TotalApparent, fi.ModTime().Unix()
}

func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
