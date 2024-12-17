GO=go
SRC_DIR=src
OUTPUT_DIR=bin
GO_FILES=$(wildcard $(SRC_DIR)/*.go)
OUTPUT=$(OUTPUT_DIR)/compiler

all: build

build: $(GO_FILES)
	$(GO) build -o $(OUTPUT) $(GO_FILES)

run: build
	./$(OUTPUT)

clean:
	del $(OUTPUT)

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run

test:
	$(GO) test ./...

install:
	$(GO) mod tidy

.PHONY: all build run clean fmt lint test install
