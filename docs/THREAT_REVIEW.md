# Shiro Threat Review — v0.2.0 (2026-10-16)

Reviewer: implementer self-review + prior fresh-reviewer pass (M1) + fuzz.
Method: code grep for destructive calls, test evidence per threat.

## Results

| # | Threat | Verdict | Evidence |
|---|---|---|---|
| 1 | Silent deletion during scan/explain/dry-run | PASS | Only `RemoveAll` in codebase is `cleanup/executor.go` (post-validation). Scan/plan paths have zero mutating calls. Integration tests assert metadata unchanged. |
| 2 | Path traversal (`a/../b`) into protected/out-of-root | PASS | Canonical (`Abs`+`EvalSymlinks`) + boundary compare. `TestDeny_Traversal` green. Fuzz `FuzzContainment` 3.3M execs, no failure. |
| 3 | Symlink escape (link → `~/.ssh`, link as target) | PASS | Walk never follows; validator denies symlink targets; `TestIntegration_SymlinkEscapeDenied` green. |
| 4 | Stale plan (changed after preview) | PASS | Exact size+mtime revalidation; E2E stale case skipped with exit 4, target kept. |
| 5 | Protected bypass via `--yes`/flags | PASS | `--yes` skips prompt only; denials enforced in `Execute` regardless. `TestExecute_ProtectedDenied` green. |
| 6 | Malformed plan file → wrong deletion | PASS | Strict JSON (`DisallowUnknownFields`, schema==1, id+root required). `FuzzDecodePlan` 1M execs, no failure. |
| 7 | Size-filter overflow (`--min-size 999...X`) | FIXED | Fuzz found float→int64 wrap to negative (would widen scope, not delete). Fixed with NaN/Inf/overflow rejection; 4.4M execs clean after. |
| 8 | Secrets in SQLite inventory | PASS | Only path/kind/bytes stored; contents never read (lstat-only scanner). |
| 9 | Shell injection via paths | PASS | No `exec`, no shell string building anywhere in non-test code. |
| 10 | FK bypass → orphan inventory rows | FIXED (M1) | `SetMaxOpenConns(1)`; multi-conn test green. Integrity-only impact. |

## Known residual (accepted)

- `deviceOf` uses following `syscall.Stat`: a symlink under
  `--one-file-system` is dropped entirely instead of counted. Read-only
  effect (under-report), no safety impact. Deferred.
- Mtime granularity: filesystems with 1s+ granularity could theoretically
  miss sub-second swaps; size check covers content changes. Accepted.
- No quarantine/undo: deletions are final; mitigated by preview + prompt +
  audit. Documented in SAFETY.md.
