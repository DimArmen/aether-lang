.PHONY: build install test clean fmt vet run help

# Build variables
BINARY_NAME=aether
VERSION=0.1.0-alpha
BUILD_DIR=bin
INSTALL_DIR=/usr/local/bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=gofmt
GOVET=$(GOCMD) vet

all: build

## build: Build the Aether CLI binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v ./cmd/aether
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## install: Install Aether CLI to system
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installation complete!"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@echo "Format complete!"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Vet complete!"

## run: Build and run Aether CLI
run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

## run-example: Run with simple-web-server example
run-example: build
	@$(BUILD_DIR)/$(BINARY_NAME) validate examples/simple-web-server/main.ae

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOGET) -v ./...
	@echo "Dependencies downloaded!"

## help: Show this help message
help:
	@echo "Aether - Infrastructure as Code with AI"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# Default target
.DEFAULT_GOAL := help
