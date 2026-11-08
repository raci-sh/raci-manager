BINARY := raci-manager
PKG := github.com/raci-sh/raci-manager
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X $(PKG)/src/internal/version.Version=$(VERSION) \
	-X $(PKG)/src/internal/version.Commit=$(COMMIT) \
	-X $(PKG)/src/internal/version.Date=$(DATE)

.PHONY: build test cover lint fmt clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./src/cmd/raci-manager

test:
	go test -race ./...

cover:
	scripts/coverage-gate.sh

lint:
	golangci-lint run

fmt:
	golangci-lint fmt

clean:
	rm -f $(BINARY) coverage.out
	rm -rf dist/
