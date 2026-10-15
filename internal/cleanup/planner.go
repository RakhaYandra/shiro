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
	Path            string
	Action          string
	Reason          string
	Risk            string
	ExpectedBytes   int64
	RuleID          string
	ExpectedSize    int64
	ExpectedModTime int64
}

type Denied struct {
	Path   string
	RuleID string
	Reason string
}

type CleanupPlan struct {
	ID            string
	Root          string
	CreatedAt     int64
	SchemaVersion int
	Items         []CleanupItem
	Denied        []Denied
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
