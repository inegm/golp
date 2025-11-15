# Makefile for golp - Launchpad Mini Go Library

# Build directory
BUILD_DIR := build

# Example programs
EXAMPLES := basic buttons rainbow animation gameoflife

# Go build flags
GO_BUILD := go build
GO_FLAGS := -v

# Default target
.PHONY: all
all: build

# Build all examples
.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	@echo "Building all examples..."
	@$(GO_BUILD) $(GO_FLAGS) -o $(BUILD_DIR)/basic ./examples/basic/main.go
	@$(GO_BUILD) $(GO_FLAGS) -o $(BUILD_DIR)/buttons ./examples/buttons/main.go
	@$(GO_BUILD) $(GO_FLAGS) -o $(BUILD_DIR)/rainbow ./examples/rainbow/main.go
	@$(GO_BUILD) $(GO_FLAGS) -o $(BUILD_DIR)/animation ./examples/animation/main.go
	@$(GO_BUILD) $(GO_FLAGS) -o $(BUILD_DIR)/gameoflife ./examples/gameoflife/main.go
	@echo "Build complete. Binaries in $(BUILD_DIR)/"

# Build individual examples
.PHONY: basic
basic:
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) $(GO_FLAGS) -o $(BUILD_DIR)/basic ./examples/basic/main.go

.PHONY: buttons
buttons:
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) $(GO_FLAGS) -o $(BUILD_DIR)/buttons ./examples/buttons/main.go

.PHONY: rainbow
rainbow:
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) $(GO_FLAGS) -o $(BUILD_DIR)/rainbow ./examples/rainbow/main.go

.PHONY: animation
animation:
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) $(GO_FLAGS) -o $(BUILD_DIR)/animation ./examples/animation/main.go

.PHONY: gameoflife
gameoflife:
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) $(GO_FLAGS) -o $(BUILD_DIR)/gameoflife ./examples/gameoflife/main.go

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Run tests
.PHONY: test
test:
	go test -v ./pkg/...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Run go vet
.PHONY: vet
vet:
	go vet ./...

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Generate documentation
.PHONY: docs
docs:
	@echo "Generating API documentation..."
	@./generate_docs.sh > docs/API.md
	@echo "Documentation generated in docs/API.md"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build       - Build all example programs"
	@echo "  make basic       - Build basic example"
	@echo "  make buttons     - Build buttons example"
	@echo "  make rainbow     - Build rainbow example"
	@echo "  make animation   - Build animation example"
	@echo "  make gameoflife  - Build game of life example"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format code"
	@echo "  make vet         - Run go vet"
	@echo "  make deps        - Install/update dependencies"
	@echo "  make docs        - Generate API documentation"
	@echo "  make help        - Show this help message"
