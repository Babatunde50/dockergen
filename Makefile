# Build settings
BINARY_NAME=dockergen
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
GO_FILES=$(shell find . -type f -name "*.go" -not -path "./vendor/*")

.PHONY: all build clean test lint vet fmt run help docker

# Default target
all: test build

# Build the binary
build:
	@echo "Building ${BINARY_NAME}..."
	@go build ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/cli

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f ${BINARY_NAME}
	@go clean

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Lint the code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Run go vet
vet:
	@echo "Vetting code..."
	@go vet ./...

# Format code
fmt:
	@echo "Formatting code..."
	@gofmt -s -w ${GO_FILES}

# Run the binary
run:
	@echo "Running binary ${BINARY_NAME}..."
	@go run ./cmd/cli

# Install the binary
install: build
	@echo "Installing ${BINARY_NAME}..."
	@mv ${BINARY_NAME} ${GOPATH}/bin/



