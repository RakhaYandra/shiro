package analysis

import "github.com/RakhaYandra/shiro/internal/scanner"

type Estimate struct {
	Bytes  int64
	Known  bool
	Reason string
}

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
