package analysis

import "github.com/RakhaYandra/shiro/internal/scanner"

// Estimate is a conservative reclaimable-space hint.
type Estimate struct {
	Bytes  int64
	Known  bool
	Reason string
}

// Reclaimable reports sparse-file savings only. Cleanup candidates arrive
// with the detector + planner phases; until then the rest is unknown.
func Reclaimable(sum scanner.Summary) Estimate {
	if !sum.AllocatedKnown {
		return Estimate{
			Reason: "allocated size unknown on this filesystem; no estimate",
		}
	}
	return Estimate{
		Bytes:  sum.SparseSavings,
		Known:  true,
		Reason: "sum of (apparent - allocated) over sparse files; excludes cleanup candidates (later phase)",
	}
}
