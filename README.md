# Shiro — Storage Intelligence

Know your storage. Clean with confidence.

Read-only by default. No silent deletion. Every recommendation carries a
reason, a size, and a risk.

## Quickstart

```sh
go build -o bin/shiro ./cmd/shiro

./bin/shiro scan ~/Projects --depth 3
./bin/shiro explain ~/Projects
./bin/shiro clean ~/Projects --dry-run
./bin/shiro clean ~/Projects --dry-run --json > plan.json
# review plan.json, then:
./bin/shiro clean --plan plan.json --execute
```

## Docs

`docs/CLI.md` commands, flags, exit codes · `docs/SAFETY.md` threat model +
policy · `docs/THREAT_REVIEW.md` v0.2.0 review · `docs/DATA_MODEL.md` size
semantics · `docs/DETECTORS.md` rules + limits · `docs/BENCHMARKS.md` baseline.
