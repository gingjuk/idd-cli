# idd-cli

A CLI tool that validates **IDD (Intent-Driven Development)** documentation compliance. Ensures project documentation follows IDD conventions with consistent identifier formats, doc-link-consistency, and proper code-doc traceability—accelerating human comprehension of code through structured documentation.

## Purpose

idd-cli validates that a project adheres to the IDD documentation system:

- **Consistent identifiers**: All IDD markers use path-based module naming (e.g., `SPEC-INT_AUTH-001` for `internal/auth/`)
- **Bidirectional links**: Specs, contracts, tests, and code all reference each other
- **Package documentation**: Every Go package has proper doc comments with Spec/Contract paths
- **Code-doc traceability**: Code annotations (`@implement`, `@test`, `@test-contract`) link to formal docs
- **Test coverage**: Every spec has corresponding tests via `**Tests:**`/`**Spec Coverage:**` linkage; every contract has validation via `@test-contract` annotations

## Supported Languages

**Code**: Go (extensible in future)

## Quick Start

```bash
# Install (once)
go install github.com/jingxu9x/idd-cli/cmd/idd-cli@latest

# Run validation
idd-cli run .

# With verbose output
idd-cli run . -v
```

### Install IDD Skill

```bash
idd-cli generate skill --output idd-skill.md
cp idd-skill.md ~/.claude/skills/
```

### Example Prompts

With the IDD skill installed, use these prompts with your AI agent:

**Implement a feature:**

```bash
Implement user authentication with email/password.
```

**Check compliance before committing:**

```bash
Run idd-cli to validate and fix any IDD documentation errors.
```

The agent will create SPEC, CONTRACT, TEST, DESIGN docs with appropriate IDD identifiers, add code annotations, and ensure doc-link-consistency.

## Configuration

See `examples/idd-config-example.yaml` for full configuration options.

## Identifier Format

Format: `TYPE-MODULE_ID`

- `TYPE`: SPEC (specification), TEST (test case), CONTRACT, DESIGN
- `MODULE`: Path-based abbreviation following directory structure
- `ID`: Sequential number within the module

Examples:

| Identifier | Meaning |
|------------|---------|
| `SPEC-INT_AUTH-001` | Spec 001 for `internal/auth/` module |
| `TEST-INT_ENG-001` | Test 001 for `internal/engine/` module |
| `SPEC-PKG_WALK-001` | Spec 001 for `pkg/walk/` module |

### Module Abbreviation Rules

Directories are abbreviated to ensure readable identifiers:

| Directory | Abbreviation |
|----------|--------------|
| `internal/<mod>` | `INT_<MOD>` (e.g., `internal/auth` → `INT_AUTH`) |
| `pkg/<mod>` | `PKG_<MOD>` (e.g., `pkg/walk` → `PKG_WALK`) |
| `cmd/<mod>` | `CMD_<MOD>` (e.g., `cmd/idd-cli` → `CMD_IDD`) |

Directories ≤4 characters are kept uppercase (e.g., `auth` → `AUTH`).

## Documentation Structure

Each module should have documentation under `docs/<path>/`:

```text
docs/
├── internal/
│   ├── auth/           # INT_AUTH
│   │   ├── spec.md    # SPEC-INT_AUTH-001
│   │   ├── contract.md
│   │   ├── testing.md
│   │   └── design.md
│   ├── engine/         # INT_ENG
│   └── model/          # INT_MOD
└── pkg/
    └── walk/          # PKG_WALK
```

## Code Annotations

```go
// Package auth provides authentication utilities.

// Spec: docs/internal/auth/spec.md
// Contract: docs/internal/auth/contract.md
package auth

// @implement SPEC-INT_AUTH-001
func Login(email, password string) error { ... }

// @test TEST-INT_AUTH-001
func TestLogin(t *testing.T) { ... }

// @test-contract TEST-INT_AUTH-001
func TestLoginContract(t *testing.T) { ... }
```

The `@test-contract` annotation links contract test functions to their corresponding test identifiers, not to SPEC identifiers.

## Validation Rules

idd-cli enforces these rules:

| Rule | Description |
|------|-------------|
| `package-doc-comment` | Every package has doc comment with Spec/Contract paths |
| `doc-link-consistency` | All forward links have corresponding backlinks |
| `orphan-detection` | No undefined or unreferenced identifiers |
| `doc-code-correspondence` | Doc references match actual code |
| `test-annotation` | Test functions have proper `@test` annotations |
| `contract-test-coverage` | Specs have corresponding contract tests |
| `design-sections` | Docs have required ## Design sections |
| `frontmatter-mismatch` | YAML frontmatter markers match content |
| `module-prefix-mismatch` | Identifiers use correct path-based module prefix |

## Project Structure

```text
idd-cli/
├── cmd/idd-cli/           # CLI entry point
├── internal/
│   ├── collector/         # Doc/code identifier collection
│   ├── engine/            # Validation engine
│   ├── graph/             # Linkage graph
│   └── model/             # Data models
├── docs/
│   ├── internal/          # Internal packages (INT_*)
│   └── pkg/               # External packages (PKG_*)
└── skills/                # IDD skill definitions
```
