# An unofficial Go library for interacting with the Trading212 API.
#
# Copyright (c) 2025 Finbarrs Oketunji
# Written by Finbarrs Oketunji <f@finbarrs.eu>
#
# This file is part of trading212.
#
# trading212 is an open-source software: you are free to redistribute
# and/or modify it under the terms of the MIT License.
#
# trading212 is made available with the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# MIT License for more details.
#
# You should have received a copy of the MIT License
# along with trading212. If not, see <https://opensource.org/licenses/MIT>.

.PHONY: run test build clean complexity lint fmt vet

# Variables
BINARY_NAME=trading212-demo
BUILD_DIR=./build
DEMO_FILE=demo/main.go
TRADING_FILE=demo/nvidia.go
MULTISTOCK_FILE=demo/multistock.go

# Default target
all: fmt vet test build

# Run the demo application
run:
	go run $(DEMO_FILE)

# Run the trading application
trade:
	go run $(TRADING_FILE)

# Run the multi-stock trading application
multistock:
	go run $(MULTISTOCK_FILE)

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Build the application
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(DEMO_FILE)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Check cyclomatic complexity
complexity:
	gocyclo -avg .
	gocyclo -over 10 .

# Lint the code
lint:
	golangci-lint run

# Format the code
fmt:
	go fmt ./...

# Vet the code
vet:
	go vet ./...

# Install development tools
dev-tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all quality checks
quality: fmt vet lint complexity test

# Watch and run tests on file changes (requires entr)
watch-test:
	find . -name "*.go" | entr -c go test ./...

# Show current project structure
show-structure:
	@echo "Current project files:"
	@ls -la *.go

# Help
help:
	@echo "Available targets:"
	@echo "  run           - Run the demo application"
	@echo "  trade         - Run the trading example application"
	@echo "  multistock    - Run the multi-stock trading application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  complexity    - Check cyclomatic complexity"
	@echo "  lint          - Lint the code"
	@echo "  fmt           - Format the code"
	@echo "  vet           - Vet the code"
	@echo "  dev-tools     - Install development tools"
	@echo "  quality       - Run all quality checks"
	@echo "  watch-test    - Watch files and run tests on changes"
	@echo "  show-structure - Show current project files"
	@echo "  help          - Show this help message"