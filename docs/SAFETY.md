# Shiro Safety (M4)

## Threat model

Cleanup deletes user data. Threats: false-positive detectors, path
traversal / symlink escape, stale plans (file changed after preview),
over-broad protected matching (prefix bugs), silent over-scan.

## Two-tier protection

- **Strict** (`~/.ssh`, `~/.gnupg`, `~/.password-store`): the path itself
  and everything beneath it is denied.
- **Roots** (`~/Documents`, `~/Projects`, `~/.config`, `~/.local/share`,
  `~/.local/state`): only the path itself is denied; detector findings
  beneath them stay eligible (the main use case).

Both tiers compare canonical paths (`Abs` + `EvalSymlinks`) with
separator-boundary prefixes — never substring or naive prefix match.

## Containment

Every candidate must sit inside the scan root (canonical). Traversal
(`a/../b`) and symlinks resolving outside the root or into protected
subtrees are denied with reasons.

## Two-phase rule

Plan (print-only, with live size+mtime snapshot per item) → save as file →
review → `clean --plan FILE --execute`. `--yes` never bypasses denials or
stale skips. `clean` without `--dry-run` (and `--plan` without `--execute`)
is refused. Non-terminal stdin without `--yes` is refused.

## Executor guarantees (M5)

- Revalidates every item: policy recheck, no-follow Lstat, exact size+mtime
  match. Changed, missing, symlinked, protected, or out-of-root targets are
  skipped/denied, never removed.
- Removal via `os.RemoveAll` directly (no shell). Post-removal Lstat must
  show the target gone, else the item is reported `unknown`/`failed`.
- Audit per item: success / skipped / denied / failed with reasons and
  verified reclaimed bytes. Partial failure continues and reports exit 5.

## Proven by tests

- `tests/integration/clean_test.go`: dry-run changes zero files/metadata;
  symlink-into-`~/.ssh` denied.
- `internal/policy`: strict subtree, root-itself, traversal, out-of-root.
- No `Remove/Chmod/Chown/Rename/exec` exists in the scan/plan path.
- Docker client exposes GET only (`internal/docker` has no POST/DELETE
  path); volumes in use are `high` risk and the report never prunes.
