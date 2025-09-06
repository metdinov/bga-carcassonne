# Carcassonne Tournament Manager - Just Commands
# Fast and Physical development workflow

# Default recipe to display available commands
default:
    @just --list

# Build the main application
build:
    @echo "🏗️  Building carca CLI..."
    go build -buildvcs=false -o carca ./cmd/carca
    @echo "✅ Build complete: ./carca"

# Build with optimizations for release
build-release:
    @echo "🚀 Building release version..."
    go build -buildvcs=false -ldflags="-s -w" -o carca ./cmd/carca
    @echo "✅ Release build complete: ./carca"

# Run all tests
test:
    @echo "🧪 Running all tests..."
    go test ./...
    @echo "✅ All tests passed"

# Run tests with verbose output
test-verbose:
    @echo "🧪 Running tests with verbose output..."
    go test -v ./...

# Run tests with coverage
test-coverage:
    @echo "🧪 Running tests with coverage..."
    go test -cover ./...
    @echo "✅ Coverage analysis complete"

# Run tests for specific package
test-package package:
    @echo "🧪 Testing package: {{package}}"
    go test -v ./{{package}}/

# Run the linter (golangci-lint)
lint:
    @echo "🔍 Running golangci-lint..."
    golangci-lint run
    @echo "✅ Linting complete"

# Fix linting issues automatically where possible
lint-fix:
    @echo "🔧 Auto-fixing lint issues..."
    golangci-lint run --fix
    @echo "✅ Auto-fix complete"

# Format Go code
fmt:
    @echo "📝 Formatting Go code..."
    go fmt ./...
    @echo "✅ Code formatted"

# Tidy up go.mod and go.sum
tidy:
    @echo "🧹 Tidying go modules..."
    go mod tidy
    @echo "✅ Modules tidied"

# Run the application with sample credentials
run:
    @echo "🎮 Running carca with sample credentials..."
    BGA_USER=sample BGA_PASS=sample ./carca

# Clean build artifacts
clean:
    @echo "🧹 Cleaning build artifacts..."
    rm -f carca
    @echo "✅ Clean complete"

# Full development cycle: clean, tidy, fmt, lint, test, build
all: clean tidy fmt lint test build
    @echo "🏆 Full development cycle complete!"

# Full development cycle without linting (for working code)
all-no-lint: clean tidy fmt test build
    @echo "🏆 Development cycle complete (no linting)!"

# Check if golangci-lint is installed
check-tools:
    @echo "🔧 Checking development tools..."
    @command -v golangci-lint >/dev/null 2>&1 || { echo "❌ golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
    @echo "✅ All tools available"

# Install development dependencies
install-tools:
    @echo "🔧 Installing development tools..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    @echo "✅ Development tools installed"

# Comprehensive development environment setup
setup:
    @echo "🚀 Setting up development environment..."
    @if command -v asdf >/dev/null 2>&1; then \
        echo "📦 asdf found, setting up managed environment..."; \
        just setup-asdf; \
    else \
        echo "📦 Using system tools..."; \
        just check-tools || just install-tools; \
    fi
    @echo "🧪 Running initial tests..."
    just test
    @echo "🏗️  Building application..."
    just build
    @echo "🎉 Setup complete! Run 'just check-versions' to verify."

# Setup development environment with asdf (if available)
setup-asdf:
    @echo "🔧 Setting up asdf environment..."
    @command -v asdf >/dev/null 2>&1 || { echo "❌ asdf not found. Install from https://asdf-vm.com/"; exit 1; }
    asdf plugin add golang || echo "golang plugin already installed"
    asdf plugin add just || echo "just plugin already installed"
    asdf install
    @echo "✅ asdf environment setup complete"

# Check asdf versions match .tool-versions
check-versions:
    @echo "🔍 Checking tool versions..."
    @echo "Expected versions from .tool-versions:"
    @cat .tool-versions
    @echo "\nCurrent versions:"
    @go version 2>/dev/null || echo "❌ Go not found"
    @just --version 2>/dev/null || echo "❌ just not found"

# Watch for changes and run tests (requires entr)
watch-test:
    @echo "👀 Watching for changes and running tests..."
    find . -name "*.go" | entr -c just test

# Generate test coverage report
coverage-html:
    @echo "📊 Generating HTML coverage report..."
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "✅ Coverage report generated: coverage.html"

# Benchmark tests
bench:
    @echo "⚡ Running benchmarks..."
    go test -bench=. ./...
    @echo "✅ Benchmarks complete"

# Verify build works on multiple platforms
build-all:
    @echo "🌐 Building for multiple platforms..."
    GOOS=linux GOARCH=amd64 go build -buildvcs=false -o carca-linux-amd64 ./cmd/carca
    GOOS=darwin GOARCH=amd64 go build -buildvcs=false -o carca-darwin-amd64 ./cmd/carca
    GOOS=darwin GOARCH=arm64 go build -buildvcs=false -o carca-darwin-arm64 ./cmd/carca
    GOOS=windows GOARCH=amd64 go build -buildvcs=false -o carca-windows-amd64.exe ./cmd/carca
    @echo "✅ Multi-platform builds complete"

# Clean multi-platform build artifacts
clean-all:
    @echo "🧹 Cleaning all build artifacts..."
    rm -f carca carca-* coverage.out coverage.html
    @echo "✅ All artifacts cleaned"
