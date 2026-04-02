BINARY := pchome
MAIN := ./cmd/pchome
HOST_GOOS := $(shell go env GOOS)
BIN_EXT := $(if $(filter windows,$(HOST_GOOS)),.exe,)
GORELEASER_VERSION ?= v2.15.1
GORELEASER := ./scripts/run-goreleaser.sh

.PHONY: build install test verify release-check release-snapshot clean

build:
	mkdir -p bin
	go build -trimpath -o ./bin/$(BINARY)$(BIN_EXT) $(MAIN)

install:
	go install -trimpath $(MAIN)

test:
	go test ./...

verify: test release-check

release-check:
	GORELEASER_VERSION=$(GORELEASER_VERSION) $(GORELEASER) check

release-snapshot:
	GORELEASER_VERSION=$(GORELEASER_VERSION) $(GORELEASER) release --snapshot --skip=publish --clean

clean:
	rm -rf ./bin ./dist
