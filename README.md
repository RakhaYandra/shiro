# Shiro — Storage Intelligence

Know your storage. Clean with confidence.

M1 scope: `shiro scan PATH` read-only + text/JSON inventory.

```sh
go build -o bin/shiro ./cmd/shiro
./bin/shiro scan ~/Projects --depth 3
./bin/shiro scan ~ --json
```
