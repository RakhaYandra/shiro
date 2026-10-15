package analysis

import (
	"testing"

	"github.com/RakhaYandra/shiro/internal/scanner"
)

func TestReclaimable_Sparse(t *testing.T) {
	sum := scanner.Summary{
		Root:           "/x",
		AllocatedKnown: true,
		SparseSavings:  900,
		SparseFiles:    2,
	}
	r := Reclaimable(sum)
	if !r.Known || r.Bytes != 900 {
		t.Fatalf("got %+v, want known/900", r)
	}
	if r.Reason == "" {
		t.Fatal("missing reason")
	}
}

func TestReclaimable_Unknown(t *testing.T) {
	sum := scanner.Summary{Root: "/x"}
	r := Reclaimable(sum)
	if r.Known {
		t.Fatalf("must be unknown without allocated data: %+v", r)
	}
}
