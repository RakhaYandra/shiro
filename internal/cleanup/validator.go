package cleanup

import (
	"fmt"
	"os"

	"github.com/RakhaYandra/shiro/internal/policy"
	"github.com/RakhaYandra/shiro/internal/scanner"
)

type ItemStatus string

const (
	StatusOK     ItemStatus = "ok"
	StatusStale  ItemStatus = "stale"
	StatusDenied ItemStatus = "denied"
)

type Verdict struct {
	Status ItemStatus
	Reason string
}

func Revalidate(pol policy.Policy, it CleanupItem, root string) Verdict {
	if err := pol.CheckIn(it.Path, root); err != nil {
		return Verdict{StatusDenied, err.Error()}
	}
	fi, err := os.Lstat(it.Path)
	if err != nil {
		return Verdict{StatusStale, fmt.Sprintf("target gone or unreadable: %v", err)}
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return Verdict{StatusDenied, "target is a symlink; refusing"}
	}
	var size int64
	if fi.IsDir() {
		sum, err := scanner.Aggregate(it.Path, scanner.Options{})
		if err != nil {
			return Verdict{StatusStale, fmt.Sprintf("cannot re-measure: %v", err)}
		}
		size = sum.TotalApparent
	} else {
		size = fi.Size()
	}
	if it.ExpectedSize >= 0 && size != it.ExpectedSize {
		return Verdict{StatusStale, fmt.Sprintf("size changed %d -> %d", it.ExpectedSize, size)}
	}
	if it.ExpectedModTime > 0 && fi.ModTime().Unix() != it.ExpectedModTime {
		return Verdict{StatusStale, "modification time changed since plan"}
	}
	return Verdict{StatusOK, "unchanged since plan"}
}
