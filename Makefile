APP ?= hosts
GO_COMMON_FLAGS ?= -v -mod=vendor
GO_TEST_FLAGS ?= -cover
OUTPUT_DIR ?= _output
TEST_TIMEOUT ?= 1m
RUN_ARGS ?=

default: build

vendor:
	go mod vendor -v

.PHONY: $(OUTPUT_DIR)/$(APP)
$(OUTPUT_DIR)/$(APP):
	@mkdir $(OUTPUT_DIR) > /dev/null 2>&1 || true
	go build $(GO_COMMON_FLAGS) -o $(OUTPUT_DIR)/$(APP) cmd/$(APP)/*

build: vendor $(OUTPUT_DIR)/$(APP)

clean:
	rm -rf "$(OUTPUT_DIR)" || true

test: test-unit

test-unit:
	go test $(GO_COMMON_FLAGS) $(GO_TEST_FLAGS) -timeout=$(TEST_TIMEOUT) ./...

.PHONY: install
install: build
	install -m +x $(OUTPUT_DIR)/$(APP) $(GOPATH)/bin/$(APP)

run:
	go run $(GO_COMMON_FLAGS) cmd/$(APP)/* $(RUN_ARGS)