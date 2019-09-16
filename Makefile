APP = hosts
BUILD_DIR ?= build

.PHONY: default bootstrap build clean test

default: build

bootstrap:
	go mod vendor

build: clean
	GO111MODULE=on go build -mod=vendor -v -o $(BUILD_DIR)/$(APP) cmd/$(APP)/*

clean:
	rm -rf $(BUILD_DIR) > /dev/null

clean-vendor:
	rm -rf ./vendor > /dev/null

test:
	go test -cover -v pkg/$(APP)/*

install:
	GO111MODULE=on go install ./cmd/$(APP)
