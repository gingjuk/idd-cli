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
│   └── walk/             # File traversal
├── examples/            # IDD identifier usage examples
├── docs/                # Project IDD documentation
├── cmd/idd-cli/skills/  # Workflow embedded into the binary
└── .github/workflows/    # CI configuration
```

## Integrated Skill and CLI Roles

The embedded IDD skill and idd-cli are one release-aligned workflow:

```text
idd-cli generate skill
        ↓
Agent loads the paired IDD workflow
        ↓
intent.md / plan.md
        ↓
docs init <package> when new
        ↓
human-readable design, contract, SPEC, and TEST records
        ↓
tests + implementation + code annotations
        ↓
docs fix for safe structure
        ↓
run . --format llm-markdown
        ↺ skill-guided semantic repairs
        ↓
repository tests + run . --format json
```

The skill owns semantic decisions. idd-cli owns these executable boundaries:

- `generate skill` exports the workflow embedded in the installed binary;
- `docs init` creates a new four-document skeleton without fake records;
- `docs fix` repairs identity and missing skeletons without rewriting prose;
- `run .` validates the complete project graph and emits exact findings.

Regenerate the installed skill after upgrading idd-cli. `.idd.yaml` is only
validator configuration and must not contain package markers or semantic
records.

Run the final validation from the project root. Supplying
`docs/<single-package>` narrows documentation collection but not source
annotation collection, so it is not a package-only validity check.

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
├── internal/auth/
│   ├── design.md     # Components and architecture decisions
│   ├── contract.md   # Contracts, interfaces, and invariants
│   ├── spec.md       # SPEC declarations and detail
│   └── testing.md    # TEST declarations, coverage, and strategy
└── IDD.md            # This file
```

Each document has an `idd` frontmatter block containing only its identity:

```yaml
# spec.md
idd:
  version: "1.0"
  package: internal/auth
  document: spec
```

The semantic declarations live in human-readable Markdown records:

```markdown
<!-- spec.md -->
## SPEC-INTERNAL_AUTH-001: User authentication

- **Design:** `AuthModule`
- **Contract:** `Authenticator`

**Requirement:** Authenticate users with validated credentials.

<!-- testing.md -->
## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_AUTH-001`

**Purpose:** Verify accepted and rejected credentials.
```

The TEST `Covers` field is the only authored SPEC/TEST relationship. idd-cli
derives the reverse SPEC-to-TEST edge. `design.md` owns `## Component:`
records and `contract.md` owns `## Contract:` records. Self-describing
documents must not repeat legacy marker, `related_files`, `Tests`, or
`Spec Coverage` metadata.

Level-two headings delimit records. The prose inside each record is ordinary
Markdown, and large registry tables are rejected as an authored declaration
format.

A SPEC represents cohesive behavior, not one function. Keep function signatures,
parameters, and returns in source code or generated API documentation.

Packages where none of the four files has an `idd` block continue to use the
legacy frontmatter format, so migration can happen one package at a time.
`docs init` refuses legacy marker metadata and central catalogs before writing;
migrate those relationships explicitly before enabling self-describing mode.

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

1. **Document Schema** — Each role owns correctly shaped Markdown records.
2. **Path Identity** — Package and role values match the Markdown path.
3. **Document Set** — All four fixed documents are self-describing.
4. **Migration Safety** — A package-local central marker catalog is rejected.
5. **Reference Integrity** — Design, contract, and coverage references resolve.
6. **Completeness** — Every SPEC has at least one derived TEST backlink.
7. **TEST Kind** — `Kind: test` uses `@test`; `Kind: contract` uses
   `@test-contract`.
8. **Doc/Code Correspondence** — Document identifiers match code annotations.
9. **Legacy Compatibility** — Legacy rules apply when no `idd` block is present.

## Configuration Example

This project configuration controls collection, validation, and reporting. It
does not own Components, Contracts, SPECs, TESTs, or coverage.

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
    spec: "@implement"
    test: "@test"
    test_contract: "@test-contract"

validation:
  require_doc_link_consistency: true
  require_spec_fields: true
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

# Export and install this binary's paired skill
./idd-cli generate skill -o idd-skill.md

# Create documents for a new package only
./idd-cli docs init internal/auth

# Normalize safe structural identity
./idd-cli docs fix docs/internal/auth

# Normalize only one independently generated document
./idd-cli docs fix docs/internal/auth/testing.md

# Let the skill repair the complete finding report
./idd-cli run . --format llm-markdown

# Final automation-facing gate
./idd-cli run . --config .idd.yaml --format json
```

An invalid graph exits non-zero after writing its report. Use the skill to edit
semantic records or source annotations; use `docs fix` only for identity and
missing-document findings.
