.PHONY: test build clean lint fmt vet coverage benchmark

# Default target
all: fmt vet test

# Run tests
test:
	go test -v ./...

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
benchmark:
	go test -bench=. -benchmem ./...

# Build the package
build:
	go build ./...

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Clean build artifacts
clean:
	go clean ./...
	rm -f coverage.out coverage.html

# Run all checks
check: fmt vet lint test

# Install dependencies
deps:
	go mod tidy
	go mod download

# Generate documentation
docs:
	godoc -http=:6060

# Run example
example:
	go run cmd/test/main.go
