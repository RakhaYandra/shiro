# Shiro Detectors (M3)

Detectors answer "what is this?", never "may this be deleted?".
Every finding is `risk=review` and carries `rule_id` + `reason`.

## Rules

| Rule | Target | Markers (confidence high) | No marker |
|---|---|---|---|
| DET-003 | `node_modules` | package.json / lockfiles | low |
| DET-002a | `go-build` | — (medium, distinctive name) | medium |
| DET-002b | `mod` under `pkg/` | — (medium) | medium |
| DET-004a | `__pycache__` | `__init__.py` / setup.py / pyproject.toml | low |
| DET-004b/c | `.venv` / `venv` | `pyvenv.cfg` (inside or parent) | low |
| DET-GEN | `target` / `build` / `dist` / `.next` | Cargo.toml / gradle / pom / package.json | low |

Markers are checked inside the candidate and in its parent (1 stat each,
read-only). Missing marker never suppresses a finding (lenient by design);
it lowers confidence instead.

## Limitations

- Basename matching: a personal directory literally named `build` with no
  marker still reports (low confidence). Read the reason, not just the row.
- No `low` risk exists yet: nothing is proven safe. `low` confidence ≠ low risk.
- Reclaimable per-finding is not estimated (planner phase).
- No Docker, no orphans, no duplicates yet.
