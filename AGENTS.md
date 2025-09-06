# AGENTS.md

## Project Overview

This project is a Go-based **TUI CLI tool** (`carca`) for creating and managing **Carcassonne tournaments** on [BoardGameArena](https://boardgamearena.com).

- Tournaments are **1v1 matches**, **best-of-3**, **Swiss system**.
- Fixtures are stored in **CSV files under `./data/`** (one file per division).
- The CLI uses the **Charm ecosystem** ([Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), etc.) for a pleasant, interactive user experience.
- The tool automates tournament creation only when explicitly requested.

---

## Agents and Responsibilities

### 1. **Data Agent**

**Role:** Handle reading, parsing, and validating fixtures.

- Input: CSV files in `./data/`.
- Output: Go structs (`Division`, `Match`, `Player`).
- Use stdlib `encoding/csv` for parsing.
- Validation: unique players, correct best-of-3 format, no empty fields.

---

### 2. **Tournament Logic Agent**

**Role:** Model rules and progression.

- Represent divisions, rounds, pairings, results.
- Handle best-of-3 resolution.
- Support Swiss tournament rules (pairing logic, score tracking).
- Provide APIs for current standings and fixtures.

---

### 3. **BGA Integration Agent**

**Role:** Handle communication with BoardGameArena.

- Use stdlib `net/http` for API requests.
- Manage authentication.
- Create tournaments and matches on demand.
- Update tournament state from BGA responses.
- Ensure integration is mocked for testing.

---

### 4. **CLI/TUI Agent**

**Role:** Present interactive CLI menus using Charm libraries.

- Entry command: `carca`.
- Main menu options:
  1. **Create Tournament** (trigger BGA integration).
  2. **View Fixture** (display parsed CSV fixtures).
  3. **View Positions** (show Swiss standings).
  4. **Exit**.
- Use **Bubble Tea** for stateful views and navigation.
- Use **Lip Gloss** for styling.

---

### 5. **Testing & QA Agent**

**Role:** Maintain code quality.

- Unit tests for CSV parsing, Swiss logic, CLI interactions.
- Mock BGA API for deterministic tests.
- Use `testing` + `testify` for assertions.
- Run `go test ./...` in CI pipeline.
- Ensure style and linting via `golangci-lint`.

---

## Development Practices

- **Language:** Go 1.22+
- **CLI Framework:** Charm Bubble Tea + Lip Gloss.
- **CSV Parsing:** stdlib `encoding/csv`.
- **HTTP:** stdlib `net/http`.
- **Versioning:** Git with [Conventional Commits](https://www.conventionalcommits.org/).
- **Branching:** `main` stable, feature branches for dev.
- **Code Style:** idiomatic Go (`gofmt`, `goimports`).
- **Tests:** standard library testing.
- **CI/CD:** GitHub Actions for linting, tests, release builds.

---

## Naming Conventions

- Packages: short, lowercase (`tournament`, `fixtures`, `bga`, `cli`).
- Types: PascalCase (`Match`, `Player`, `Division`).
- Functions: PascalCase for exported, camelCase for internal.
- Files: snake_case, descriptive (`fixtures_parser.go`, `bga_client.go`).
- CLI commands: kebab-case (`carca create`).

---

## Suggested Folder Structure

carca/
├── cmd/
│ └── carca/ # main.go (entrypoint for CLI/TUI)
├── internal/
│ ├── cli/ # Bubble Tea models, views, menus
│ │ ├── menu.go
│ │ ├── create.go
│ │ ├── fixture.go
│ │ └── positions.go
│ ├── fixtures/ # CSV parsing and validation
│ │ └── parser.go
│ ├── tournament/ # Core tournament logic (Swiss, best-of-3)
│ │ └── tournament.go
│ ├── bga/ # BGA integration (HTTP client)
│ │ ├── client.go
│ │ └── mock_client.go
│ └── utils/ # Helpers (logging, formatting, etc.)
├── data/ # Example CSV fixture files
│ ├── division_a.csv
│ └── division_b.csv
├── tests/ # Integration & end-to-end tests
├── go.mod
├── go.sum
└── AGENTS.md

---

## Example CLI Flow

```bash
$ carca
Welcome to Carcassonne Tournament Manager
Please select an option:
1. Create Tournament
2. View Fixture
3. View Positions
4. Exit
```
