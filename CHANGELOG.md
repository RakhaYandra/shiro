# Changelog

## v0.3.0 (2026-10-16)

Post-MVP intelligence batch.

- `duplicates`: two-stage sha256 report (size → prefix → full), hardlink
  excluded, changed-during-hash unverified, report-only, default min 1MiB.
- `history` / `diff`: full-tree directory inventory in SQLite, added /
  removed / grown / shrunk with growth ranking, explicit `--from/--to`.
- `docker`: read-only Engine API inventory (GET-only client, no new deps);
  in-use volumes flagged high; never prunes.
- `tui`: interactive tree + findings + dry-run preview; no execute path
  (enforced by lint guard); refuses non-terminal.
- No auto-prune for history; duplicate cleanup stays manual.

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
