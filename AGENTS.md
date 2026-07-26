# IDD CLI — Agent Guide

## Build & Run Commands

```bash
# Build (Tree-sitter grammars require CGO; outputs to bin/idd-cli)
CGO_ENABLED=1 go build -o bin/idd-cli ./cmd/idd-cli

# Or use go install
CGO_ENABLED=1 go install github.com/jingxu9x/idd-cli/cmd/idd-cli@latest

# Or use Make (binary name in Makefile is "idd-cli")
make build

# Run validation
./bin/idd-cli run .
./bin/idd-cli run . -v          # verbose
./bin/idd-cli --config .idd.yaml run .

# Run tests (CI requires CGO_ENABLED=1)
CGO_ENABLED=1 go test -v -race ./...

# Lint
CGO_ENABLED=1 golangci-lint run ./...

# Format
gofmt -w .
goimports -w .
```

## Binary Name

The binary is built as `idd-cli` (see Makefile's `BINARY_NAME`).

## Project Structure

```text
cmd/idd-cli/main.go        # CLI entry, uses cobra
internal/
  engine/engine.go           # Core validation orchestration
  engine/source_validation.go # AST-backed source rules
  engine/filesystem.go       # Configured source/doc traversal
  collector/                 # Doc and code collectors
  collector/source_ast.go    # Seven-language Tree-sitter binding
  collector/document_schema.go # Shared scaffold/completion schema
  collector/idd_document_validation.go # Self-describing document rules
  graph/graph.go             # LinkageGraph structure
  model/identifier.go        # Identifier, IdentifierSet types
  config/config.go           # .idd.yaml loading
  reporter/reporter.go       # JSON/markdown output
  auth/auth.go               # Stub types (not real auth implementation)
pkg/pattern/idd.go           # IDD regex patterns (SPEC-*, TEST-*, etc.)
```

## Identifier Format

Format: `TYPE-MODULE-NUMBER` (e.g., `SPEC-BE-001`, `TEST-BE-001`)

Pattern in `.idd.yaml`: `SPEC-[A-Z]+-[0-9]+`

Code annotations: `@implement`, `@test`, `@test-contract`

Code-to-document association is identifier-derived. Do not add repeated
`Spec:`, `Contract:`, or `Test:` document paths to source-file headers;
`idd-cli run` joins declaration annotations to the owning self-describing
document records.

## Config Loading

Config search order:

1. `--config` flag value
2. `./.idd.yaml`
3. `./config/.idd.yaml`

## Code Architecture Note

The **README.md architecture diagrams show interfaces that don't exist in code**. The actual implementation:

- No `validator.ValidationRule` interface — rules are explicit `Engine`
  methods split across `engine.go`, `source_validation.go`,
  `annotation_validation.go`, and `filesystem.go`
- No `collector.Collector` interface — just `DocCollector` and `CodeCollector` structs
- No `linker.Linker` interface — graph building is in `engine.buildGraph()`

Trust code over README for implementation details.

## Testing

- Tests exist in `internal/graph/graph_test.go` and `internal/model/*_test.go`
- Most packages lack test coverage
- Run with `-race` flag (CI requirement)
- **Use table-driven tests** to organize test cases grouped by scenario

## IDD Documentation System

This tool validates the IDD (Intent-Driven Development) documentation system
defined in `cmd/idd-cli/skills/SKILL.md`. Key conventions:

- Every scanned source package and nested sub-package maps to
  `docs/<package>/` and owns `design.md`, `contract.md`, `spec.md`, and
  `testing.md`
- The exact basename is the only document-role authority. IDD frontmatter
  contains `version` and `package`; `idd.document` is invalid
- IDD role files have no line-count limit and must not be split into
  `design-*`, `contract-*`, `spec-*`, or `testing-*` files
- `examples/self-describing-module-docs/` is the one normative four-file
  document example
- `examples/idd-config-example.yaml` mirrors `config.Default()`
- The tool itself uses IDD annotations in `internal/auth/auth.go` as an example

## Pre-Commit Rules

Before committing any changes, the following checks must pass:

### 1. Config Defaults Sync Check

If `internal/config/config.go` `Default()` function is modified, `examples/idd-config-example.yaml` must be updated to match the new defaults.

**Check command**: Diff `Default()` values in `config.go` with `examples/idd-config-example.yaml`

### 2. Example Document Check

If `examples/self-describing-module-docs/*.md` files are modified, validate
both completion and the collector-owned schema/reference graph.

**Check commands**:

```bash
./bin/idd-cli docs status examples/self-describing-module-docs --format json
CGO_ENABLED=1 go test ./internal/collector -run TestReferenceExampleDocuments
```

### 3. IDD Validation Check

Before any commit, run idd-cli to ensure all docs are compliant:

```bash
./bin/idd-cli run .
```

If validation fails, fix errors before committing. Common issues:

- Generated scaffold markers or required fields remain incomplete
- Document version/package identity does not match its intended `docs/`
  location
- A source package or sub-package lacks one of its four canonical files, or a
  role was placed in a split/non-canonical filename
- SPEC, TEST, Component, or Contract references do not resolve
- Source annotations have no matching document record
