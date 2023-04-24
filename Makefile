APP ?= hosts
OUTPUT_DIR ?= bin
BIN ?= $(OUTPUT_DIR)/$(APP)

CMD ?= ./cmd/...

GOFLAGS ?= -v -mod=vendor
GOFLAGS_TEST ?= -v -race -cover -timeout=1m

ARGS ?=

default: $(BIN)

.PHONY: $(BIN)
$(BIN):
	go build -o $(BIN) $(CMD)

clean:
	rm -rf "$(OUTPUT_DIR)" || true

test: test-unit

test-unit:
	go test $(GOFLAGS_TEST) ./...

.PHONY: install
install:
	go install $(CMD)

run:
	go run $(CMD) $(ARGS)
