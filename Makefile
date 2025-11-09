.PHONY: build build-linux build-macos build-windows run clean test

# Variables
BINARY_NAME=server
MAIN_PATH=main.go
BIN_DIR=./bin

# Default build
build:
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)"

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)-linux"

# Build for Linux (static binary with musl)
build-linux-static:
	@echo "Building static Linux binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	CC=musl-gcc \
	go build -ldflags="-linkmode external -extldflags '-static'" \
	-o $(BIN_DIR)/$(BINARY_NAME)-linux-static $(MAIN_PATH)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)-linux-static"

# Build for macOS
build-macos:
	@echo "Building for macOS..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-macos $(MAIN_PATH)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)-macos"

# Build for macOS ARM64 (M1/M2)
build-macos-arm:
	@echo "Building for macOS ARM64..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-macos-arm $(MAIN_PATH)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)-macos-arm"

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-windows.exe $(MAIN_PATH)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)-windows.exe"

# Build all platforms
build-all: build-linux build-macos build-macos-arm build-windows
	@echo "All builds complete!"

# Run the server
run:
	@echo "Running server..."
	go run $(MAIN_PATH)

ensure-compile-daemon:
	@which go > /dev/null || (echo "Error: Go is not installed or not in PATH" && exit 1)
	@which CompileDaemon > /dev/null || (echo "Installing CompileDaemon..." && go install github.com/githubnemo/CompileDaemon@latest)

# Run with live reload (requires air: go install github.com/githubnemo/CompileDaemon@latest)
dev:
	@echo "Starting development server with live reload..."
	CompileDaemon -build="go build -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PATH)" -command="$(BIN_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Create necessary directories
setup:
	@echo "Setting up project structure..."
	mkdir -p $(BIN_DIR)
	@echo "Setup complete!"

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t smart-spore-hub:latest .

# Run in Docker
docker-run:
	@echo "Running in Docker..."
	docker run -p 8080:8080 --env-file .env smart-spore-hub:latest

# Help
help:
	@echo "Available commands:"
	@echo "  make build              - Build for current platform"
	@echo "  make build-linux        - Build for Linux"
	@echo "  make build-macos        - Build for macOS (Intel)"
	@echo "  make build-macos-arm    - Build for macOS (M1/M2)"
	@echo "  make build-windows      - Build for Windows"
	@echo "  make build-all          - Build for all platforms"
	@echo "  make run                - Run the server"
	@echo "  make dev                - Run with live reload"
	@echo "  make clean              - Remove build artifacts"
	@echo "  make test               - Run tests"
	@echo "  make test-coverage      - Run tests with coverage"
	@echo "  make deps               - Install dependencies"
	@echo "  make fmt                - Format code"
	@echo "  make lint               - Lint code"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-run         - Run in Docker"
	@echo "  make help               - Show this help message"
