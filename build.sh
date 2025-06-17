#!/bin/bash

# Build script for My Redis Server
# Author: [Your Name]
# Date: June 2025

set -e

PROJECT_NAME="my-redis-server"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date +"%Y-%m-%d %H:%M:%S")
GO_VERSION=$(go version | awk '{print $3}')

echo "==================================="
echo "Building $PROJECT_NAME"
echo "Version: $VERSION"
echo "Build Time: $BUILD_TIME"
echo "Go Version: $GO_VERSION"
echo "==================================="

# Clean previous builds
echo "Cleaning previous builds..."
rm -f $PROJECT_NAME $PROJECT_NAME.exe $PROJECT_NAME-* 

# Build for current platform
echo "Building for current platform..."
go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o $PROJECT_NAME

# Build for multiple platforms
echo "Building for multiple platforms..."

# Linux amd64
echo "Building for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o $PROJECT_NAME-linux-amd64

# Linux arm64  
echo "Building for Linux arm64..."
GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o $PROJECT_NAME-linux-arm64

# macOS amd64
echo "Building for macOS amd64..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o $PROJECT_NAME-darwin-amd64

# macOS arm64 (Apple Silicon)
echo "Building for macOS arm64..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o $PROJECT_NAME-darwin-arm64

# Windows amd64
echo "Building for Windows amd64..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o $PROJECT_NAME-windows-amd64.exe

echo "==================================="
echo "Build completed successfully!"
echo ""
echo "Binaries created:"
ls -la $PROJECT_NAME*
echo ""
echo "To run locally: ./$PROJECT_NAME"
echo "==================================="
