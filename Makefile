APP ?= hosts
GO_COMMON_FLAGS ?= -v -mod=vendor
OUTPUT_DIR ?= build
TEST_TIMEOUT ?= 1m

default: build

vendor:
	go mod vendor -v

build: vendor
	@mkdir $(OUTPUT_DIR) > /dev/null 2>&1 || true
	go build $(GO_COMMON_FLAGS) -o $(OUTPUT_DIR)/$(APP) cmd/$(APP)/*

install: build
	go install $(GO_COMMON_FLAGS) cmd/$(APP)/*

test: test-unit

test-unit:
	go test $(GO_COMMON_FLAGS) -timeout=$(TEST_TIMEOUT) ./...
