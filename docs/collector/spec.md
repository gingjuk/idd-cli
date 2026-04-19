---
markers:
  - id: SPEC-COLLECTOR-001
    name: Document Collector
  - id: SPEC-COLLECTOR-002
    name: Code Collector
---

# Specification (collector)

## SPEC-COLLECTOR-001: Document Collector

**Status:** Done

**Requirement:**

The Document Collector (`DocCollector`) must collect IDD identifiers from markdown documentation files, parsing frontmatter markers and extracting identifier references from content.

**Implementation:** `internal/collector/doc_collector.go`

**Key Types:**

- `DocCollector` — Main collector struct with config reference
- `Frontmatter` — YAML frontmatter with marker definitions
- `Marker` — Individual marker with ID, name, and describe fields

**Key Functionality:**

- Recursively walk directories to find `.md` files
- Parse YAML frontmatter to extract defined markers
- Extract IDD references from markdown content using regex
- Validate frontmatter markers match actual content headings
- Validate module prefix matches directory structure
- Report errors for malformed markers (not wrapped in backticks)

**Acceptance Criteria:**

- [x] Collects identifiers from markdown files in directory trees
- [x] Parses frontmatter markers correctly
- [x] Detects defined markers vs referenced markers
- [x] Extracts title from heading containing identifier
- [x] Validates module prefix consistency
- [x] Ignores paths configured in `ignore_paths`
- [x] Reports errors for bare markers (not wrapped in backticks)

**Tests:** `TEST-COLLECTOR-001`

**Related:** `CONTRACT-BE-001`

---

## SPEC-COLLECTOR-002: Code Collector

**Status:** Done

**Requirement:**

The Code Collector (`CodeCollector`) must collect IDD annotations from source code files (Go, TypeScript, JavaScript), extracting `@spec`, `@contract`, `@test`, and `@design` annotations.

**Implementation:** `internal/collector/code_collector.go`

**Key Types:**

- `CodeCollector` — Main collector struct with config reference

**Key Functionality:**

- Walk directory trees to find source files (`.go`, `.ts`, `.tsx`, `.js`)
- Extract annotations using configured patterns
- Extract function context (function name, preceding comments)
- Build identifiers from annotations with code location
- Set origin to `model.OriginCode` for code-based identifiers

**Acceptance Criteria:**

- [x] Supports Go, TypeScript, and JavaScript files
- [x] Extracts `@spec`, `@contract`, `@test`, `@design` annotations
- [x] Captures function name and preceding comments as context
- [x] Reports file path and line number for each annotation
- [x] Ignores paths configured in `ignore_paths`
- [x] Handles multiple annotations on same line

**Tests:** `TEST-COLLECTOR-002`

**Related:** `CONTRACT-BE-001`
