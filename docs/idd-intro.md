# IDD CLI Documentation

> This project follows the IDD (Intent-Driven Development) framework for documentation management.

## Project Structure

```text
idd-cli/
├── cmd/idd-cli/         # CLI entry point
├── internal/
│   ├── collector/         # Document/code identifier collection
│   ├── config/           # Configuration loading
│   ├── engine/          # Validation engine
│   ├── graph/           # Linkage graph
│   ├── model/           # Data models
│   ├── reporter/        # Report generation
│   └── auth/             # Example code
├── pkg/
│   ├── pattern/          # IDD identifier regex patterns
│   ├── walk/            # File traversal
│   └── annotation/       # Code annotation parsing
├── examples/            # IDD identifier usage examples
├── docs/                # Project IDD documentation
├── skills/              # IDD skill definitions
└── .github/workflows/    # CI configuration
```

## IDD Identifier Format

**Format:** `<TYPE>-<MODULE>-<NUMBER>`

| Prefix | Meaning | Example |
| ------ | ------- | ------- |
| `SPEC-` | Functionality specification | SPEC-XX-001 |
| `TEST-` | Test case | TEST-XX-001 |

Note: CONTRACT and DESIGN are no longer primary identifiers. contract.md describes contracts that implement SPECs. design.md describes architecture decisions.

**Regex Patterns:**

```yaml
identifier_patterns:
  spec: "SPEC-[A-Z]+-[0-9]+"
  test: "TEST-[A-Z]+-[0-9]+"
  test_contract: "TEST-[A-Z]+-[0-9]+"
```

## Ignoring Annotations

In some cases, you may want idd-cli to ignore certain annotations. For example, test data strings containing IDD identifiers should not be treated as real annotations.

### Nolint Directives

```go
// idd:ignore                    // Ignore this line only

// idd:ignore-start              // Start ignoring (alternative syntax)
// idd:ignore-end                // Stop ignoring (alternative syntax)
```

### Examples

**Single line ignore:**

```go
code := `// idd:ignore
// @implement SPEC-XX-001          // This annotation will NOT be collected
`
func TestRealAnnotation() {
    // @test TEST-XX-001              // This annotation WILL be collected
}
```

**Scope ignore:**

```go
// idd:ignore start
code := `
// @implement SPEC-XX-001      // Ignored
// @test TEST-XX-001         // Ignored
`
// idd:ignore end

func TestCode() {
    // @test TEST-XX-001              // Collected normally
}
```

### When to Use

- **Test data**: Annotations inside backtick strings in test files should be wrapped with `// idd:ignore` or `// idd:ignore start/end`
- **Example code**: Code examples in comments that show annotations but aren't meant to be collected
- **Temporary annotations**: Annotations you're not ready to link yet

## Document Structure

Documents are organized by module in the `docs/` directory:

```text
docs/
├── backend/
│   ├── spec.md       # SPEC-XX-001, SPEC-XX-002
│   ├── testing.md    # TEST-XX-001
│   ├── contract.md   # Contract implementations (no CONTRACT identifiers)
│   └── design.md     # Architecture decisions (no DESIGN identifiers)
└── IDD.md            # This file
```

Each IDD document must include YAML frontmatter for tool parsing:

```yaml
---
markers:
  - id: SPEC-XX-001
    name: Feature description

related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---
```

### Identifier Index Table Format

```markdown
| ID | Title | Status | Tests |
|----|-------|--------|-------|
| [SPEC-XX-001](#spec-xx-001) | Feature Name | Done | TEST-XX-001 |
```

## Code Annotation Format

```go
// @implement SPEC-XX-001: Feature description
func DoSomething() {
    // implementation
}

// @test TEST-XX-001: Test verification
func TestDoSomething() {
    // implementation
}

// @test-contract TEST-XX-002: Contract test for Login interface
func TestDoSomethingContract() {
    // implementation
}
```

## Validation Rules

1. **Completeness** — Every SPEC must have at least one TEST link
2. **Bidirectional** — If SPEC→TEST exists, TEST→SPEC must also exist
3. **Orphan Detection** — No identifiers with zero connections
4. **Consistency** — Cross-reference chain consistency

## Configuration Example

```yaml
version: "1.0"

docs:
  patterns:
    - "docs/**/*.md"
  identifier_patterns:
    spec: "SPEC-[A-Z]+-[0-9]+"
    test: "TEST-[A-Z]+-[0-9]+"
    test_contract: "TEST-[A-Z]+-[0-9]+"

code:
  patterns:
    - "**/*.go"
  annotations:
    - "@implement"
    - "@test"
    - "@test-contract"

validation:
  require_doc_link_consistency: true
  allow_orphans: false
  require_spec_test_coverage: true

output:
  file: "idd-report.json"
  include_graph: true
  verbose: false
```

## Usage

```bash
# Build
go build -o idd-cli ./cmd/idd-cli

# Run validation
./idd-cli run --config .idd.yaml

# Using Make
make build && make run
```
