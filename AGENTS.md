# IDD Link Validator — Agent Guide

## Build & Run Commands

```bash
# Build (outputs to bin/idd-verify, NOT idd-validator)
go build -o bin/idd-verify ./cmd/validator

# Or use Make (binary name in Makefile is "idd-verify")
make build

# Run validation
./bin/idd-verify run .
./bin/idd-verify run . -v          # verbose
./bin/idd-verify --config idd.yaml run .

# Run tests (CI requires CGO_ENABLED=1)
CGO_ENABLED=1 go test -v -race ./...

# Lint
golangci-lint run ./...

# Format
gofmt -w .
goimports -w .
```

## Binary Name Discrepancy

- **Makefile** builds to `bin/idd-verify`
- **README.md** incorrectly says `idd-validator`
- Always use `idd-verify` or check Makefile's `BINARY_NAME`

## Project Structure

```
cmd/validator/main.go        # CLI entry, uses cobra
internal/
  engine/engine.go           # Core validation logic (NOT in separate validator/ package)
  collector/                 # Doc and code collectors
  graph/graph.go             # LinkageGraph structure
  model/identifier.go        # Identifier, IdentifierSet types
  config/config.go           # idd.yaml loading
  reporter/reporter.go       # JSON/markdown output
  auth/auth.go               # Stub types (not real auth implementation)
pkg/pattern/idd.go           # IDD regex patterns (SPEC-*, TEST-*, etc.)
```

## Identifier Format

Format: `TYPE-MODULE-NUMBER` (e.g., `SPEC-BE-001`, `TEST-BE-001`)

Pattern in `idd.yaml`: `SPEC-[A-Z]+-[0-9]+`

Code annotations: `@spec`, `@contract`, `@test`, `@design`

## Config Loading

Config search order:
1. `--config` flag value
2. `./idd.yaml`
3. `./config/idd.yaml`

## Code Architecture Note

The **README.md architecture diagrams show interfaces that don't exist in code**. The actual implementation:
- No `validator.ValidationRule` interface — rules are hardcoded in `engine.go`
- No `collector.Collector` interface — just `DocCollector` and `CodeCollector` structs
- No `linker.Linker` interface — graph building is in `engine.buildGraph()`

Trust code over README for implementation details.

## Testing

- Tests exist in `internal/graph/graph_test.go` and `internal/model/*_test.go`
- Most packages lack test coverage
- Run with `-race` flag (CI requirement)

## IDD Documentation System

This tool validates the IDD (Intent-Driven Development) documentation system defined in `skills/SKILL.md`. Key conventions:
- `docs/` contains SPEC, CONTRACT, TEST, DESIGN markdown files
- `examples/` contains reference doc examples
- The tool itself uses IDD annotations in `internal/auth/auth.go` as an example
