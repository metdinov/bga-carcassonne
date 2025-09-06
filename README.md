# Carcassonne Tournament Manager

Tournament management for [BoardGameArena](https://boardgamearena.com) Carcassonne tournaments.

## Overview

- **Tournament Fixtures** - Parse CSV files with match schedules
- **Interactive Navigation** - Browse tournaments by division and round
- **Match Selection** - Copy tournament links or create new tournaments
- **Swiss System Support** - View tournament standings and progression

## Quick Start

### Prerequisites

- Go 1.25.1+ (or use asdf for version management)
- [just](https://github.com/casey/just) command runner
- [golangci-lint](https://golangci-lint.run/) (optional, for linting)
- [asdf](https://asdf-vm.com/) (optional, for version management)

#### Optional: Version Management with asdf

For consistent development environments, use [asdf](https://asdf-vm.com/) to manage Go and just versions:

```bash
# Install asdf (see https://asdf-vm.com/guide/getting-started.html)
# Add required plugins
asdf plugin add golang
asdf plugin add just

# Install versions specified in .tool-versions file
asdf install
```

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd bga-carcassonne

# Comprehensive setup (recommended)
just setup                  # Auto-detects asdf, installs tools, tests, and builds
```

### Configuration

Set your BoardGameArena credentials:

```bash
# Environment variables (recommended)
export BGA_USER="your-username"
export BGA_PASS="your-password"

# Or create a .env file
echo "BGA_USER=your-username" > .env
echo "BGA_PASS=your-password" >> .env
```

### Usage

```bash
# Run the application
./carca

# Or with explicit credentials
BGA_USER=username BGA_PASS=password ./carca
```

## Development Workflow

This project uses [just](https://github.com/casey/just) for development automation.

### Available Commands

```bash
just                    # Show all available commands
just build              # Build the main application
just test               # Run all tests
just lint               # Run golangci-lint
just fmt                # Format Go code
just all                # Full development cycle
```

### Common Tasks

```bash
# Initial setup
just setup                  # Complete environment setup (auto-detects asdf)
just check-versions         # Verify tool versions match .tool-versions

# Development cycle
just fmt tidy lint test build

# Testing
just test                    # All tests
just test-verbose            # Verbose output
just test-coverage           # With coverage
just test-package cli        # Specific package

# Code quality
just lint                    # Run linter
just lint-fix               # Auto-fix issues
just fmt                    # Format code

# Building
just build                  # Development build
just build-release          # Optimized release
just build-all              # Multi-platform builds

# Utilities
just clean                  # Clean artifacts
just run                    # Run with sample creds
just coverage-html          # Generate HTML coverage
```

### Full Development Cycle

```bash
# One command for the complete workflow
just all
```

This runs: `clean` → `tidy` → `fmt` → `lint` → `test` → `build`

## Project Structure

```
carca/
├── cmd/carca/           # Main application entry point
├── internal/
│   ├── cli/            # TUI interface (Bubble Tea models)
│   ├── fixtures/       # CSV parsing and tournament data
│   └── utils/          # Shared utilities
├── data/               # Example tournament CSV files
├── justfile            # Development automation
└── .golangci.yml       # Linter configuration
```

## Features

### 🎮 Interactive TUI

- **Main Menu** - Navigate between options with arrow keys or vim keys (j/k)
- **Division Selection** - Choose from Elite, Platinum A/B, Oro A/B/C/D
- **Fixture Display** - Professional table format with match details
- **Match Selection** - Interactive selection with ↑/↓ or j/k keys

### ⌨️ Navigation

**Menu Navigation:**

- `↑/↓` or `j/k` - Navigate menu items
- `Enter` - Select option
- `q/Ctrl+C` - Quit

**Fixture Navigation:**

- `←/→`, `h/l`, or `PgUp/PgDown` - Navigate rounds
- `↑/↓` or `j/k` - Select matches
- `Enter` - Copy tournament link (played) or show create prompt (unplayed)
- `c` - Create tournament for unplayed match
- `Esc/q` - Go back

### 📊 Tournament Data

- **CSV Parsing** - Read tournament fixtures from CSV files
- **Match Status** - Visual indicators for played (✓) vs unplayed (○) matches
- **Consistent Layout** - Professional table formatting across all rounds
- **Tournament Links** - Extract and copy BGA tournament URLs
- **Player Information** - Handle variable-length player names with consistent alignment

### 🔗 Clipboard Integration

- **Link Copying** - Instantly copy tournament URLs to clipboard
- **Status Messages** - User feedback for all actions
- **Auto-clear** - Messages disappear after 3 seconds

## Testing

The project maintains comprehensive test coverage:

```bash
# Run all tests
just test

# Specific test suites
just test-package fixtures    # CSV parsing tests
just test-package cli         # TUI interaction tests

# Coverage analysis
just test-coverage
just coverage-html           # Generate HTML report
```

## Code Quality

### Linting

```bash
just lint                    # Check code quality
just lint-fix               # Auto-fix where possible
```

The project uses `golangci-lint` with TUI-appropriate settings that account for the complexity of interactive terminal applications.

### Formatting

```bash
just fmt                    # Format all Go files
```

Uses standard `go fmt` with consistent import organization.
