# Specification: SPEC-BE-001

## IDD Link Validator Overview

**Version:** 1.0.0
**Last Updated:** 2026-04-18

## Requirement

IDD Link Validator is a CLI tool that validates bidirectional linkage consistency between IDD (Intent-Driven Development) identifiers across documentation and source code.

## Functionality

### Core Features

1. **Document Scanning** — Scans `.md` files for IDD identifiers using configurable regex patterns
2. **Code Annotation Detection** — Detects `@spec`, `@contract`, `@test`, `@design` annotations in source code
3. **Linkage Graph Construction** — Builds a directed graph of identifier references
4. **Bidirectional Validation** — Ensures all forward links have corresponding backlinks
5. **Completeness Check** — Ensures every SPEC has at least one TEST link and vice versa
6. **Orphan Detection** — Identifies identifiers with zero connections
7. **JSON Report Generation** — Outputs validation results in structured JSON format

### IDD Identifier Format

```
<TYPE>-<MODULE>-<NUMBER>
```

Examples:
- `SPEC-BE-001` — Specification for backend module 001
- `TEST-BE-001` — Test case for backend module 001
- `CONTRACT-BE-001` — Contract for backend module 001
- `DESIGN-BE-001` — Design decision for backend module 001

### User Interface

```bash
# Build
go build -o idd-verify ./cmd/validator

# Run validation
idd-verify run --config idd.yaml

# Run with verbose output
idd-verify run --config idd.yaml -v
```

### Configuration

See `idd.yaml.example` for configuration options:

- `docs.patterns` — Glob patterns for documentation files
- `docs.identifier_patterns` — Regex patterns for IDD identifiers
- `code.patterns` — Glob patterns for source files
- `code.annotations` — Annotation markers to scan
- `validation.*` — Validation rule toggles
- `output.*` — Output configuration

## Implementation

**Implementation Location:** `cmd/validator/main.go`, `internal/engine/engine.go`

**Key Modules:**
- `internal/collector/` — Document and code collection
- `internal/graph/` — Linkage graph structure
- `internal/validator/` — Validation rules
- `internal/reporter/` — Report generation

## Acceptance Criteria

- [x] CLI tool accepts `--config` flag for configuration
- [x] Scans documentation files matching configured patterns
- [x] Extracts IDD identifiers using configurable regex
- [x] Builds linkage graph from collected identifiers
- [x] Validates bidirectional links exist
- [x] Detects orphan identifiers (unless `allow_orphans: true`)
- [x] Generates JSON report with validation results
- [x] Exits with non-zero code when validation fails

## Related Documents

- **SPEC-BE-002** — Graph Linkage Structure
- **CONTRACT-BE-001** — Collector Interface Contracts
- **TEST-BE-001** — Core Validation Tests
- **DESIGN-BE-001** — Graph-First Architecture Design

## Tests

- **TEST-BE-001** — Core Validation Tests
