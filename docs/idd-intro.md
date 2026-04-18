# IDD Link Validator Documentation

> This project follows the IDD (Intent-Driven Development) framework for documentation management.

## Project Structure

```
idd-cli/
├── cmd/validator/         # CLI entry point
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
|--------|---------|---------|
| `SPEC-` | Functionality specification | `SPEC-BE-001` |
| `CONTRACT-` | Interface/behavior contract | `CONTRACT-BE-001` |
| `TEST-` | Test case | `TEST-BE-001` |
| `DESIGN-` | Architecture design decision | `DESIGN-BE-001` |

**Regex Patterns:**
```yaml
identifier_patterns:
  spec: "SPEC-[A-Z]+-[0-9]+"
  contract: "CONTRACT-[A-Z]+-[0-9]+"
  test: "TEST-[A-Z]+-[0-9]+"
  design: "DESIGN-[A-Z]+-[0-9]+"
```

## Document Structure

Documents are organized by module in the `docs/` directory:

```
docs/
├── backend/
│   ├── spec.md       # SPEC-BE-001, SPEC-BE-002
│   ├── testing.md    # TEST-BE-001
│   ├── contract.md   # CONTRACT-BE-001
│   └── design.md     # DESIGN-BE-001
└── IDD.md            # This file
```

Each IDD document must include YAML frontmatter for tool parsing:

```yaml
---
markers:
  - id: SPEC-BE-001
    name: Feature description
  - id: CONTRACT-BE-001
    name: Contract interface description
---
```

### Identifier Index Table Format

```markdown
| ID | Title | Status | Tests |
|----|-------|--------|-------|
| [SPEC-BE-001](#spec-be-001) | Feature Name | Done | TEST-BE-001 |
```

## Code Annotation Format

```go
// @spec SPEC-BE-001: Feature description
// @contract CONTRACT-BE-001: Interface specification
// @test TEST-BE-001: Test verification
// @design DESIGN-BE-001: Architecture design
func DoSomething() {
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
    contract: "CONTRACT-[A-Z]+-[0-9]+"
    test: "TEST-[A-Z]+-[0-9]+"
    design: "DESIGN-[A-Z]+-[0-9]+"

code:
  patterns:
    - "**/*.go"
  annotations:
    - "@spec"
    - "@contract"
    - "@test"
    - "@design"

validation:
  require_bidirectional: true
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
go build -o idd-verify ./cmd/validator

# Run validation
./idd-verify run --config idd.yaml

# Using Make
make build && make run
```
