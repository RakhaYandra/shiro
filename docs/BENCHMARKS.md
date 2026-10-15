# Shiro Benchmarks

Machine: AMD Ryzen 7 7840HS, linux/amd64, Go 1.27.2. Date: 2026-10-16.
Reproduce: `go test -bench=. -benchmem -run=NONE ./internal/scanner/ ./internal/detector/`

## Baseline (v0.2.0)

| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| Aggregate 1k files (10 dirs) | ~1.98ms | ~686KB | ~6.5k |
| Aggregate 10k files (100 dirs) | ~22.4ms | ~6.3MB | ~62k |
| Detect 100 node_modules | ~0.44ms | ~153KB | ~1.6k |

Notes:

- Cost scales linearly with file count; TopFiles heap (cap 100) and ExtTotals
  map keep memory bounded per entry, not per tree.
- Content hashing: N/A (no hashing phase yet; M7 duplicates will add a
  separate benchmark — metadata scan vs hashing must stay split).
- Store writes are out of the hot path (one transaction per scan).

## Duplicates (v0.3.0-dev)

| Benchmark | ns/op | B/op | allocs/op |
|---|---|---|---|
| FindDuplicates 10×2MiB identical | ~13.4ms | ~704KB | ~303 |

Hashing dominates: 20MiB read+sha256 in ~13ms (~1.5GiB/s). Metadata scan
vs hashing stay split — hashing only touches `--min-size` candidates
(default 1MiB); unchanged files are never re-read within a run.
