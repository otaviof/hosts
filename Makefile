APP ?= hosts
GO_COMMON_FLAGS ?= -v -mod=vendor
OUTPUT_DIR ?= build
TEST_TIMEOUT ?= 1m
RUN_ARGS ?=

default: build

vendor:
	go mod vendor -v

$(OUTPUT_DIR)/$(APP):
	@mkdir $(OUTPUT_DIR) > /dev/null 2>&1 || true
	go build $(GO_COMMON_FLAGS) -o $(OUTPUT_DIR)/$(APP) cmd/$(APP)/*

build: vendor $(OUTPUT_DIR)/$(APP)

clean:
	rm -rf "$(OUTPUT_DIR)" || true

test: test-unit

test-unit:
	go test $(GO_COMMON_FLAGS) -timeout=$(TEST_TIMEOUT) ./...

install: build
	cp -v $(OUTPUT_DIR)/$(APP) $(GOPATH)/bin/

run:
	go run $(GO_COMMON_FLAGS) cmd/$(APP)/* $(RUN_ARGS)