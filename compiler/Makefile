build:
	@echo "Building compiler..."
	@go build -o bin/compiler src/main.go

repl: build
	@echo "Running REPL..."
	@./bin/compiler

test:
	@echo "Running tests..."
	@go test ./... -v

fmt:
	@echo "Formatting code..."
	@go fmt ./...