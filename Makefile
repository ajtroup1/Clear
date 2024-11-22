# Variables
BIN := bin/clear
MAIN := ./src/cmd/main.go

# Targets
.PHONY: all build run test clean fresh

build:
	@go build -o ${BIN} ${MAIN}

repl: build
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
