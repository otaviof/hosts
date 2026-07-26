APP = hosts
OUTPUT_DIR ?= bin
VERSION ?= $(shell cat ./version)

BIN ?= $(OUTPUT_DIR)/$(APP)
CMD ?= ./cmd/$(APP)/...
PKG ?= ./pkg/$(APP)/...

GOFLAGS ?= -trimpath
LDFLAGS ?= -s -w

GOFLAGS_TEST ?= \
	-v  \
	-failfast \
	-race \
	-cover \
	-coverprofile=coverage.txt \
	-covermode=atomic

LOWER_OSTYPE ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
CPUTYPE ?= $(shell uname -m)
INSTALL_DIR ?= /usr/local/bin

ARGS ?=

.EXPORT_ALL_VARIABLES:

default: build

vendor:
	go mod vendor -v

.PHONY: $(BIN)
$(BIN):
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" $(GOFLAGS) -o $(BIN) $(CMD)

build: $(BIN)

.PHONY: run
run:
	go run $(CMD) $(ARGS)

install: build
	install -m 0755 $(BIN) $(INSTALL_DIR)/$(APP)

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR) >/dev/null

.PHONY: lint
lint:
	golangci-lint run ./...

test: test-unit

.PHONY: test-unit
test-unit:
	go test $(GOFLAGS_TEST) $(CMD) $(PKG)

snapshot:
	goreleaser --clean --snapshot --skip=publish

release:
	git tag $(VERSION)
	git push --tags origin $(VERSION)
	goreleaser --clean
