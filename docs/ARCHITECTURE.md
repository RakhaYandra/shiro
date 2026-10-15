# Shiro Architecture (M1)

```
Cobra CLI (internal/cli)
  -> scanner (internal/scanner): Walk + Aggregate, stdlib only
  -> output (internal/output): text + JSON from same Summary
  -> store (internal/store): interface; sqlite impl (modernc.org/sqlite)
```

## Boundaries

- `scanner` has no CLI, SQLite, or output deps.
- `output` is pure rendering; display flags never change scan scope.
- `store` persists metadata only; scans work with `NoopStore`.
- Detector / analysis / cleanup / TUI do not exist yet (post-M1).

## Safety

- Read-only by default: `Lstat`, no follow, no delete/move/chmod.
- Canonical root (`Abs` + `EvalSymlinks`); exclusions are boundary-checked
  prefixes, not substring matches.
- Per-path errors collected (`Partial=true`, exit 3), never fatal.
