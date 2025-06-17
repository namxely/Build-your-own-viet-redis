# VietRedis Server Makefile
# High-Performance Redis Implementation by Namxely

# Build information
PROJECT_NAME := viet-redis-server
BINARY_NAME := vietredis
VERSION := 2.0.0-saigon
BUILD_TIME := $(shell date +"%Y-%m-%d %H:%M:%S")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_BY := Namxely-$(shell whoami)@$(shell hostname)

# Go build settings
GO_VERSION := 1.21
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
CGO_ENABLED := 0

# Directories
BUILD_DIR := build
DIST_DIR := dist
DOCS_DIR := docs
LOGS_DIR := logs
DATA_DIR := data

# Docker settings
DOCKER_REGISTRY := namxely
DOCKER_IMAGE := $(DOCKER_REGISTRY)/vietredis-server
DOCKER_TAG := $(VERSION)

# Build flags
LDFLAGS := -w -s \
	-X 'main.version=$(VERSION)' \
	-X 'main.buildTime=$(BUILD_TIME)' \
	-X 'main.gitCommit=$(GIT_COMMIT)' \
	-X 'main.buildBy=$(BUILD_BY)'

GCFLAGS := -trimpath
ASMFLAGS := -trimpath

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
RESET := \033[0m

.PHONY: all build build-prod build-all clean test test-coverage benchmark \
	run run-cluster docker-build docker-push docker-run install uninstall \
	help setup deps fmt lint vet security release dev-setup

# Default target
all: clean test build

# Help target
help:
	@echo "$(CYAN)🇻🇳 VietRedis Server Build System$(RESET)"
	@echo "$(CYAN)================================$(RESET)"
	@echo ""
	@echo "$(YELLOW)Build Commands:$(RESET)"
	@echo "  build         Build VietRedis server for current platform"
	@echo "  build-prod    Build optimized production version"
	@echo "  build-all     Build for all supported platforms"
	@echo "  clean         Clean build artifacts"
	@echo ""
	@echo "$(YELLOW)Development Commands:$(RESET)"
	@echo "  run           Run VietRedis server"
	@echo "  run-cluster   Run VietRedis in cluster mode"
	@echo "  dev-setup     Setup development environment"
	@echo "  fmt           Format Go code"
	@echo "  lint          Run linter"
	@echo "  vet           Run go vet"
	@echo ""
	@echo "$(YELLOW)Test Commands:$(RESET)"
	@echo "  test          Run all tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  benchmark     Run benchmarks"
	@echo ""
	@echo "$(YELLOW)Docker Commands:$(RESET)"
	@echo "  docker-build  Build Docker image"
	@echo "  docker-push   Push Docker image to registry"
	@echo "  docker-run    Run VietRedis in Docker container"
	@echo ""
	@echo "$(YELLOW)Installation Commands:$(RESET)"
	@echo "  install       Install VietRedis to system"
	@echo "  uninstall     Remove VietRedis from system"
	@echo ""
	@echo "$(YELLOW)Other Commands:$(RESET)"
	@echo "  setup         Setup project directories"
	@echo "  deps          Download dependencies"
	@echo "  security      Run security checks"
	@echo "  release       Create release package"
	@echo ""

# Setup project directories
setup:
	@echo "$(BLUE)Setting up project directories...$(RESET)"
	@mkdir -p $(BUILD_DIR) $(DIST_DIR) $(LOGS_DIR) $(DATA_DIR)
	@mkdir -p $(DATA_DIR)/master $(DATA_DIR)/replica
	@echo "$(GREEN)✓ Project directories created$(RESET)"

# Download dependencies
deps:
	@echo "$(BLUE)Downloading dependencies...$(RESET)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies downloaded$(RESET)"

# Format Go code
fmt:
	@echo "$(BLUE)Formatting Go code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✓ Code formatted$(RESET)"

# Run linter
lint:
	@echo "$(BLUE)Running linter...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not found, installing...$(RESET)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi
	@echo "$(GREEN)✓ Linting completed$(RESET)"

# Run go vet
vet:
	@echo "$(BLUE)Running go vet...$(RESET)"
	@go vet ./...
	@echo "$(GREEN)✓ go vet completed$(RESET)"

# Security check
security:
	@echo "$(BLUE)Running security checks...$(RESET)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "$(YELLOW)gosec not found, installing...$(RESET)"; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi
	@echo "$(GREEN)✓ Security check completed$(RESET)"

# Build for current platform
build: setup deps fmt
	@echo "$(BLUE)Building VietRedis Server v$(VERSION)...$(RESET)"
	@echo "$(PURPLE)Platform: $(GOOS)/$(GOARCH)$(RESET)"
	@echo "$(PURPLE)Go Version: $(shell go version)$(RESET)"
	@echo "$(PURPLE)Build Time: $(BUILD_TIME)$(RESET)"
	@echo "$(PURPLE)Git Commit: $(GIT_COMMIT)$(RESET)"
	@CGO_ENABLED=$(CGO_ENABLED) go build \
		-ldflags="$(LDFLAGS)" \
		-gcflags="$(GCFLAGS)" \
		-asmflags="$(ASMFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		.
	@echo "$(GREEN)✓ Build completed: $(BUILD_DIR)/$(BINARY_NAME)$(RESET)"

# Build production version with extra optimizations
build-prod: setup deps fmt vet
	@echo "$(BLUE)Building production VietRedis Server...$(RESET)"
	@CGO_ENABLED=0 go build \
		-ldflags="$(LDFLAGS)" \
		-gcflags="$(GCFLAGS)" \
		-asmflags="$(ASMFLAGS)" \
		-tags="netgo,osusergo,static_build" \
		-installsuffix netgo \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		.
	@echo "$(GREEN)✓ Production build completed$(RESET)"

# Build for all platforms
build-all: setup deps fmt
	@echo "$(BLUE)Building for all platforms...$(RESET)"
	@mkdir -p $(DIST_DIR)
	
	# Linux amd64
	@echo "$(PURPLE)Building for Linux amd64...$(RESET)"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .
	
	# Linux arm64
	@echo "$(PURPLE)Building for Linux arm64...$(RESET)"
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
		-ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 .
	
	# macOS amd64
	@echo "$(PURPLE)Building for macOS amd64...$(RESET)"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
		-ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .
	
	# macOS arm64 (Apple Silicon)
	@echo "$(PURPLE)Building for macOS arm64...$(RESET)"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build \
		-ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .
	
	# Windows amd64
	@echo "$(PURPLE)Building for Windows amd64...$(RESET)"
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
		-ldflags="$(LDFLAGS)" \
		-o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	
	@echo "$(GREEN)✓ All platform builds completed$(RESET)"
	@ls -la $(DIST_DIR)/

# Run tests
test:
	@echo "$(BLUE)Running tests...$(RESET)"
	@go test -v -race ./...
	@echo "$(GREEN)✓ Tests completed$(RESET)"

# Run tests with coverage
test-coverage:
	@echo "$(BLUE)Running tests with coverage...$(RESET)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(RESET)"

# Run benchmarks
benchmark:
	@echo "$(BLUE)Running benchmarks...$(RESET)"
	@go test -bench=. -benchmem ./...
	@echo "$(GREEN)✓ Benchmarks completed$(RESET)"

# Clean build artifacts
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(RESET)"
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@rm -f coverage.out coverage.html
	@rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	@echo "$(GREEN)✓ Clean completed$(RESET)"

# Run VietRedis server
run: build
	@echo "$(BLUE)Starting VietRedis Server...$(RESET)"
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Run VietRedis in cluster mode
run-cluster: build
	@echo "$(BLUE)Starting VietRedis Cluster...$(RESET)"
	@mkdir -p $(DATA_DIR)/node1 $(DATA_DIR)/node2 $(DATA_DIR)/node3
	@echo "Starting node 1 on port 7000..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --port 7000 --cluster-enabled yes --cluster-config-file $(DATA_DIR)/node1/nodes.conf --dir $(DATA_DIR)/node1 &
	@echo "Starting node 2 on port 7001..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --port 7001 --cluster-enabled yes --cluster-config-file $(DATA_DIR)/node2/nodes.conf --dir $(DATA_DIR)/node2 &
	@echo "Starting node 3 on port 7002..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --port 7002 --cluster-enabled yes --cluster-config-file $(DATA_DIR)/node3/nodes.conf --dir $(DATA_DIR)/node3 &
	@sleep 3
	@echo "$(GREEN)✓ Cluster nodes started$(RESET)"
	@echo "$(YELLOW)To create cluster, run: redis-cli --cluster create 127.0.0.1:7000 127.0.0.1:7001 127.0.0.1:7002$(RESET)"

# Build Docker image
docker-build:
	@echo "$(BLUE)Building Docker image...$(RESET)"
	@docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		-t $(DOCKER_IMAGE):latest \
		.
	@echo "$(GREEN)✓ Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)$(RESET)"

# Push Docker image
docker-push: docker-build
	@echo "$(BLUE)Pushing Docker image...$(RESET)"
	@docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	@docker push $(DOCKER_IMAGE):latest
	@echo "$(GREEN)✓ Docker image pushed$(RESET)"

# Run Docker container
docker-run:
	@echo "$(BLUE)Running VietRedis in Docker...$(RESET)"
	@docker run -d \
		--name vietredis \
		-p 6379:6379 \
		-v $(PWD)/$(DATA_DIR):/data \
		$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "$(GREEN)✓ VietRedis container started$(RESET)"
	@echo "$(YELLOW)Connect with: redis-cli -p 6379$(RESET)"

# Install to system
install: build-prod
	@echo "$(BLUE)Installing VietRedis to system...$(RESET)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo cp vietredis.conf /etc/vietredis.conf
	@sudo mkdir -p /var/lib/vietredis /var/log/vietredis
	@echo "$(GREEN)✓ VietRedis installed to /usr/local/bin/$(BINARY_NAME)$(RESET)"

# Uninstall from system
uninstall:
	@echo "$(BLUE)Uninstalling VietRedis from system...$(RESET)"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@sudo rm -f /etc/vietredis.conf
	@echo "$(GREEN)✓ VietRedis uninstalled$(RESET)"

# Development setup
dev-setup:
	@echo "$(BLUE)Setting up development environment...$(RESET)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/air-verse/air@latest
	@echo "$(GREEN)✓ Development tools installed$(RESET)"

# Create release package
release: build-all
	@echo "$(BLUE)Creating release package...$(RESET)"
	@mkdir -p $(DIST_DIR)/release
	@cp README.md CHANGELOG.md LICENSE $(DIST_DIR)/release/
	@cp vietredis.conf docker-compose.yml $(DIST_DIR)/release/
	@cd $(DIST_DIR) && tar -czf vietredis-$(VERSION).tar.gz release/ $(BINARY_NAME)-*
	@echo "$(GREEN)✓ Release package created: $(DIST_DIR)/vietredis-$(VERSION).tar.gz$(RESET)"

# Version info
version:
	@echo "$(CYAN)VietRedis Server Build Information$(RESET)"
	@echo "$(CYAN)================================$(RESET)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Built By: $(BUILD_BY)"
	@echo "Go Version: $(shell go version)"
	@echo "Platform: $(GOOS)/$(GOARCH)"
