# Shiro CLI (M1)

## Commands

| Command | Effect |
|---|---|
| `shiro scan PATH` | Read-only scan, text summary. Never modifies target. |
| `shiro scan PATH --json` | Stable JSON (`schema_version: 1`). |
| `shiro explain [PATH]` | Categories + top dirs/files + reclaimable estimate. Read-only, no store. |
| `shiro explain [PATH] --json` | Machine-readable explanation (`schema_version: 1`). |
| `shiro clean PATH --dry-run` | Print-only cleanup plan. Never deletes. `--dry-run` required. |
| `shiro clean PATH --dry-run --json` | Machine-readable plan (`schema_version: 1`). |
| `shiro clean PATH --dry-run --json > plan.json` | Save plan file for review. |
| `shiro clean --plan plan.json --execute` | Revalidate + prompt `yes` + remove approved targets. |
| `shiro clean --plan plan.json --execute --yes` | Non-interactive (policy still enforced). |
| `shiro duplicates PATH` | Report identical content (sha256, two-stage). Read-only, never deletes. |
| `shiro duplicates PATH --min-size 10MiB --json` | Tune hashing scope; machine output. |
| `shiro --version` / `--help` | Always available. |

## Flags (`scan`)

- `--depth N` — display depth only (default 3). Scan scope is always full.
- `--exclude PATH` — repeatable path exclusion (prefix-boundary match).
- `--one-file-system` — do not cross filesystem boundaries.
- `--min-size SIZE` — display filter, e.g. `100MiB`, `2GiB`, `512KiB`.
- `--json` — machine-readable output.
- `--no-store` — skip SQLite inventory.
- `--db PATH` — SQLite path override (default `~/.local/share/shiro/shiro.db`).

## Exit codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | General error |
| 2 | Invalid CLI usage |
| 3 | Partial scan (errors recorded in output, report incomplete) |
| 4 | Cleanup denied by policy/validation (or nothing executable) |
| 5 | Cleanup partially failed |
