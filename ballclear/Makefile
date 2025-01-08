build:
	@echo "Building the Clear executable..."
	@go build -o bin/clear ./src/main.go

repl: build
	@echo "Starting REPL..."
	@./bin/clear

run: build
	@echo "Running script..."
	@./bin/clear $(ARGS)

test:
	@echo "Running tests..."
	@go test ./... -v

fmt:
	@echo "Formatting..."
	@go fmt ./...
