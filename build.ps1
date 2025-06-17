# VietRedis Server Build Script for Windows
# 🇻🇳 High-Performance Redis Implementation Made in Vietnam
# Version: 2.0.0 Saigon Edition

$PROJECT_NAME = "viet-redis-server"
$BINARY_NAME = "vietredis"
$VERSION = "2.0.0-saigon"
$BUILD_TIME = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
$GIT_COMMIT = try { (git rev-parse --short HEAD 2>$null) } catch { "unknown" }
$BUILD_BY = "Namxely-$env:USERNAME@$env:COMPUTERNAME"
$GO_VERSION = (go version).Split()[2]

Write-Host "🇻🇳 ===============================================" -ForegroundColor Red
Write-Host "   VietRedis Server Build System" -ForegroundColor Yellow
Write-Host "   High-Performance Redis Made in Vietnam" -ForegroundColor Yellow
Write-Host "===============================================" -ForegroundColor Red
Write-Host ""
Write-Host "📦 Project: $PROJECT_NAME" -ForegroundColor Cyan
Write-Host "🚀 Version: $VERSION" -ForegroundColor Yellow
Write-Host "⏰ Build Time: $BUILD_TIME" -ForegroundColor Yellow
Write-Host "🔧 Go Version: $GO_VERSION" -ForegroundColor Yellow
Write-Host "📝 Git Commit: $GIT_COMMIT" -ForegroundColor Yellow
Write-Host "👤 Built By: $BUILD_BY" -ForegroundColor Yellow
Write-Host ""

# Clean previous builds
Write-Host "🧹 Cleaning previous builds..." -ForegroundColor Cyan
Remove-Item -Path "$BINARY_NAME*", "build", "dist" -Recurse -Force -ErrorAction SilentlyContinue

# Create directories
Write-Host "📁 Creating build directories..." -ForegroundColor Cyan
New-Item -ItemType Directory -Path "build", "dist" -Force | Out-Null

# Build flags
$LDFLAGS = "-w -s -X 'main.version=$VERSION' -X 'main.buildTime=$BUILD_TIME' -X 'main.gitCommit=$GIT_COMMIT' -X 'main.buildBy=$BUILD_BY'"

# Build for current platform
Write-Host "🔨 Building for current platform..." -ForegroundColor Green
$env:CGO_ENABLED = "0"
go build -ldflags "$LDFLAGS" -gcflags "-trimpath" -asmflags "-trimpath" -o "build\$BINARY_NAME.exe" .

# Build for multiple platforms
Write-Host "🌍 Building for multiple platforms..." -ForegroundColor Green

# Linux amd64
Write-Host "Building for Linux amd64..." -ForegroundColor Blue
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o "$PROJECT_NAME-linux-amd64"

# Linux arm64
Write-Host "Building for Linux arm64..." -ForegroundColor Blue
$env:GOOS = "linux"
$env:GOARCH = "arm64"
go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o "$PROJECT_NAME-linux-arm64"

# macOS amd64
Write-Host "Building for macOS amd64..." -ForegroundColor Blue
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o "$PROJECT_NAME-darwin-amd64"

# macOS arm64 (Apple Silicon)
Write-Host "Building for macOS arm64..." -ForegroundColor Blue
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o "$PROJECT_NAME-darwin-arm64"

# Windows amd64
Write-Host "Building for Windows amd64..." -ForegroundColor Blue
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags "-X main.version=$VERSION -X 'main.buildTime=$BUILD_TIME'" -o "$PROJECT_NAME-windows-amd64.exe"

# Reset environment variables
Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue

Write-Host "===================================" -ForegroundColor Green
Write-Host "Build completed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Binaries created:" -ForegroundColor Yellow
Get-ChildItem -Name "$PROJECT_NAME*"
Write-Host ""
Write-Host "To run locally: .\$PROJECT_NAME.exe" -ForegroundColor Cyan
Write-Host "===================================" -ForegroundColor Green
