# Shiro CLI (M1)

## Commands

| Command | Effect |
|---|---|
| `shiro scan PATH` | Read-only scan, text summary. Never modifies target. |
| `shiro scan PATH --json` | Stable JSON (`schema_version: 1`). |
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
