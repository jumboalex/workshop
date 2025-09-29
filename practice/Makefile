.PHONY: test build run clean coverage fmt vet lint help

# Default target
all: test build

# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./algorithms/...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./algorithms/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the runner executable
build:
	@echo "Building runner..."
	@go build -o bin/runner ./cmd/runner

# Run the runner
run: build
	@echo "Running application..."
	@./bin/runner

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Run golangci-lint (if installed)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./algorithms/...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Check for updates
check:
	@echo "Checking code..."
	@$(MAKE) fmt
	@$(MAKE) vet
	@$(MAKE) test

# Help command
help:
	@echo "Available targets:"
	@echo "  make test      - Run all tests"
	@echo "  make coverage  - Generate test coverage report"
	@echo "  make build     - Build the runner executable"
	@echo "  make run       - Build and run the application"
	@echo "  make fmt       - Format code"
	@echo "  make vet       - Run go vet"
	@echo "  make lint      - Run golangci-lint"
	@echo "  make clean     - Remove build artifacts"
	@echo "  make bench     - Run benchmarks"
	@echo "  make tidy      - Tidy dependencies"
	@echo "  make check     - Format, vet, and test code"
	@echo "  make help      - Show this help message"