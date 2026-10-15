# Changelog

## v0.2.0 (2026-10-16)

Full workflow: scan → explain → dry-run → execute.

- `scan`: read-only inventory, text + stable JSON, exit 0/2/3, SQLite store.
- `explain`: categories with rule evidence, top dirs/files, sparse estimate.
- detectors: node, go, python, generic build outputs (all review-risk).
- `clean --dry-run`: print-only plans with policy denials.
- `clean --plan FILE --execute`: revalidation, prompt/`--yes`, per-item audit.
- policy: two-tier protected paths, canonical containment, no shell anywhere.
- hardening: benchmarks, fuzz (fixed size-overflow bug), threat review.

## v0.1.0-scan

M1 scanner + SQLite inventory + text/JSON output.
