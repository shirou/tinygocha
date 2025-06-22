# Makefile for ゴチャキャラバトル

# Default target platform
GOOS ?= windows
GOARCH ?= amd64

# Binary name
BINARY_NAME = tinygocha
ifeq ($(GOOS),windows)
	BINARY_NAME := $(BINARY_NAME).exe
endif

# Build directory
BUILD_DIR = build

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application (for development)
.PHONY: run
run:
	go run .

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/tinygocha-windows-amd64.exe .
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/tinygocha-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/tinygocha-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/tinygocha-darwin-arm64 .
	@echo "Multi-platform build complete"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
	@echo "Dependencies installed"

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...
	@echo "Tests complete"

# Development build (with race detection)
.PHONY: dev
dev:
	@echo "Building development version..."
	go build -race -o $(BUILD_DIR)/$(BINARY_NAME)-dev .
	@echo "Development build complete"

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build for default platform ($(GOOS)/$(GOARCH))"
	@echo "  run        - Run the application for development"
	@echo "  clean      - Clean build artifacts"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  deps       - Install dependencies"
	@echo "  fmt        - Format code"
	@echo "  test       - Run tests"
	@echo "  dev        - Build development version with race detection"
	@echo "  help       - Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  GOOS       - Target OS (default: windows)"
	@echo "  GOARCH     - Target architecture (default: amd64)"
