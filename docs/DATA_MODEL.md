# Shiro Data Model (M1)

## Size semantics

- **Apparent size**: logical bytes (`FileInfo.Size()`). Always present.
- **Allocated size**: `Stat_t.Blocks * 512` on Linux. Valid only when
  `AllocatedKnown=true`; otherwise report 0 and treat as unknown.
- **Reclaimable estimate**: not computed in M1 (no cleanup planner yet).

## Rules

1. Parent size = sum of children; totals never double-count parent+child.
2. Symlinks are entries (`kind=symlink`), never followed, size 0.
3. Hardlinks (same dev+ino) count once.
4. Sparse files: apparent large, allocated small — show both.
5. Scan results are snapshots; paths may change after the scan.

## Entry

`{path, kind: file|directory|symlink|other, apparent_bytes,
allocated_bytes, allocated_known, modtime_unix, device_id, inode}`

## SQLite (optional)

Tables `scans(id, root_path, started_at, finished_at, status, scanner_version)`
and `entries(scan_id, path, kind, apparent_bytes, allocated_bytes, modified_at)`.
FK enforced. No file contents or secrets stored. `--no-store` disables.

History stores every directory node (full tree), not just top level —
roughly one row per directory. Individual files are never stored. Diff
compares per-path bytes between two scans of the same root; the total is
the shallowest common path's delta (avoids double-counting nested dirs).
No auto-prune: history grows one scan row-set per `scan`; prune deferred.
