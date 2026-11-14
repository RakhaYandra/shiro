# Shiro — Storage Intelligence

Know your storage. Clean with confidence.

Read-only by default. No silent deletion. Every recommendation carries a
reason, a size, and a risk.

## Requirements

- Linux
- Go 1.27+ (to build)
- Docker optional — only for `shiro docker` (inventory via the Engine API, GET-only)

## Installation

```sh
git clone https://github.com/RakhaYandra/shiro && cd shiro
make install        # builds bin/shiro and installs to ~/.local/bin/shiro (override with PREFIX)
```

Release binaries (linux amd64/arm64 + `SHA256SUMS`) can be built with `make release`.

## Quickstart

```sh
shiro scan ~/Projects --depth 3
shiro explain ~/Projects
shiro clean ~/Projects --dry-run
shiro clean ~/Projects --dry-run --json > plan.json
# review plan.json, then:
shiro clean --plan plan.json --execute
```

## Commands

| Command | Effect |
|---|---|
| `shiro scan PATH` | Read-only scan, text summary. Never modifies the target. |
| `shiro scan PATH --json` | Stable JSON (`schema_version: 1`). |
| `shiro explain [PATH]` | Categories, top dirs/files, reclaimable estimate. Read-only. |
| `shiro explain [PATH] --json` | Machine-readable explanation (`schema_version: 1`). |
| `shiro clean PATH --dry-run` | Print-only plan. `--dry-run` is required. |
| `shiro clean PATH --dry-run --json` | Machine-readable plan (`schema_version: 1`). |
| `shiro clean --plan plan.json --execute` | Revalidate the plan, prompt `yes`, then remove approved targets. |
| `shiro clean --plan plan.json --execute --yes` | Non-interactive (policy still enforced). |
| `shiro duplicates PATH` | Report identical content (sha256, two-stage). Read-only. |
| `shiro duplicates PATH --min-size 10MiB --json` | Tune hashing scope; machine output. |
| `shiro history [PATH]` | List recorded scans. |
| `shiro diff [PATH]` | Compare the two newest scans for PATH. |
| `shiro diff --from ID --to ID` | Compare explicit scans (must share root). |
| `shiro docker [--json]` | Docker inventory via the Engine API (GET-only). |
| `shiro tui [PATH]` | Interactive tree + findings + dry-run preview. Never executes cleanup. |
| `shiro --version` / `--help` | Always available. |

Useful `scan` flags: `--depth N` (display only, default 3), `--exclude PATH`
(repeatable), `--one-file-system`, `--min-size SIZE` (e.g. `100MiB`, `2GiB`).

## How cleaning works

1. `--dry-run` produces a reviewable plan — it never deletes.
2. The plan records each approved target with a reason, a size, and a risk.
3. `--execute` revalidates the plan against the filesystem before deleting,
   and asks for the literal word `yes` (`--yes` never skips policy checks).

The TUI is read-only by construction: a lint rule fails the build if cleanup
execution is reachable from the TUI.

## Design principles

- Read-only by default; deletion only through a reviewed plan.
- Every output format can say the same thing: text for humans, JSON
  (`schema_version: 1`) for automation.
- No telemetry, no network calls except optional Docker Engine inventory
  (GET-only, against the local socket).

## Docs

`docs/CLI.md` commands, flags, exit codes · `docs/SAFETY.md` threat model +
policy · `docs/THREAT_REVIEW.md` v0.2.0 review · `docs/DATA_MODEL.md` size
semantics · `docs/DETECTORS.md` rules + limits · `docs/BENCHMARKS.md` baseline.

## License

No `LICENSE` file yet — all rights reserved for now; adding one is an open
decision (an explicit LICENSE choice is pending).
