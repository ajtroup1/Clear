# Variables
BIN := bin/clear

# Targets
.PHONY: all build run test clean fresh

build:
	@go build -o ${BIN} ./src/cmd/main.go

run: build
	@./${BIN}

test:
	@go test -v ./...
	
fmt:
	@go fmt ./...

clean:
	@go clean
	@rm -f ${BIN}

# A combined target to clean, build, and run the application
fresh: clean build run
