.PHONY: build run vet test race lint clean install ci release bench

PREFIX ?= $(HOME)/.local
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X github.com/RakhaYandra/shiro/internal/cli.version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/shiro ./cmd/shiro

release:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/shiro-linux-amd64 ./cmd/shiro
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/shiro-linux-arm64 ./cmd/shiro
	sha256sum dist/shiro-linux-* > dist/SHA256SUMS

lint:
	test -z "$$(gofmt -l cmd internal tests)" || (gofmt -l cmd internal tests && exit 1)
	go vet ./...
	test -z "$$(rg -l "cleanup\.Execute" internal/tui || true)" || (echo "tui must not execute cleanup" && exit 1)

install: build
	install -Dm755 bin/shiro $(PREFIX)/bin/shiro

run: build
	./bin/shiro --help

vet:
	go vet ./...

test:
	go test ./...

race:
	go test -race ./...

bench:
	go test -bench=. -benchmem ./internal/scanner/

ci: lint test
	go test -race ./internal/... ./tests/...

clean:
	rm -rf bin/ dist/
