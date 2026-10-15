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

Plan (print-only, with live size+mtime snapshot per item) → review →
execution (M5, not implemented). `--yes` will never bypass denials.
`clean` without `--dry-run` is refused.

## Proven by tests

- `tests/integration/clean_test.go`: dry-run changes zero files/metadata;
  symlink-into-`~/.ssh` denied.
- `internal/policy`: strict subtree, root-itself, traversal, out-of-root.
- No `Remove/Chmod/Chown/Rename/exec` exists in the scan/plan path.
