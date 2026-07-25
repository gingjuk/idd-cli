# idd-cli

A CLI tool with an embedded **IDD (Intent-Driven Development)** skill. The
skill guides semantic authoring; idd-cli creates safe document skeletons,
repairs structural identity, and validates the complete documentation/code
traceability graph.

## Purpose

idd-cli validates that a project adheres to the IDD documentation system:

- **Consistent identifiers**: All IDD markers use path-based module naming (e.g., `SPEC-INTERNAL_AUTH-001` for `internal/auth/`)
- **Derived traceability**: `testing.md` records coverage once and idd-cli derives reverse links
- **Package documentation**: Every Go package has proper doc comments with Spec/Contract paths
- **Code-doc traceability**: Code annotations (`@implement`, `@test`, `@test-contract`) link to formal docs
- **Test coverage**: Every SPEC is covered through a TEST record's `Covers`
  field; contract tests use `@test-contract`
- **Version-aligned workflow**: The binary exports the IDD skill that describes
  the document model it validates

## Supported Languages

**Code**: Go (extensible in future)

## Skill and CLI Workflow

The two parts have a strict responsibility boundary:

- The IDD skill captures intent and authors meaningful Component, Contract,
  SPEC, TEST, code, and repair changes.
- idd-cli exports that skill, initializes four-file document sets, repairs only
  safe structure, and reports graph problems at exact locations.
- `.idd.yaml` configures validation; it never owns markers or semantic records.

### Install or refresh the pair

```bash
# Install or upgrade the binary
go install github.com/jingxu9x/idd-cli/cmd/idd-cli@latest

# Export the exact skill embedded in this binary
idd-cli generate skill -o idd-skill.md
```

Install `idd-skill.md` as `SKILL.md` through your agent's skill mechanism.
Regenerate it after every idd-cli upgrade so the agent and validator use the
same document contract.

### Develop with the integrated loop

Run commands from the project root:

```bash
# New package only; existing self-describing packages skip this command
idd-cli docs init internal/auth

# The IDD skill now authors design, contract, SPEC, TEST, tests, and code.

# Normalize safe identity after independent document generation
idd-cli docs fix docs/internal/auth
idd-cli docs fix docs/internal/auth/testing.md

# Agent-oriented repair report
idd-cli run . --format llm-markdown

# After the skill repairs each finding and repository tests pass
idd-cli run . --format json
```

`docs fix` never writes semantic content. Identity or missing-file findings can
use `docs fix`; schema, reference, Markdown, migration, TEST-kind, and
annotation findings must be repaired by editing their canonical Markdown
record or source declaration.

Always use `idd-cli run .` for the validity gate. A documentation subdirectory
argument narrows documentation collection, but source annotations still come
from the current project working tree.

### Example Prompts

With the IDD skill installed, use these prompts with your AI agent:

**Implement a feature:**

```bash
Use the intent-driven-development skill to implement user authentication with
email/password. Use idd-cli for document bootstrap, structural repair, and the
final project-level validation loop.
```

**Check compliance before committing:**

```bash
Use the IDD skill and run `idd-cli run . --format llm-markdown`. Repair every
finding at its canonical owner, rerun repository tests, then finish with the
JSON project-level gate.
```

The agent uses the skill to decide what the documents mean and uses idd-cli to
prove that their structure and traceability are complete.

## Configuration

See `examples/idd-config-example.yaml` for full configuration options.

## Identifier Format

Format: `TYPE-MODULE-NUMBER`

- `TYPE`: `SPEC` or `TEST` for self-describing records. `CONTRACT` and
  `DESIGN` identifiers are accepted only by the legacy metadata model.
- `MODULE`: Path-based abbreviation following directory structure
- `ID`: Sequential number within the module

Examples:

| Identifier | Meaning |
|------------|---------|
| `SPEC-INTERNAL_AUTH-001` | Spec 001 for `internal/auth/` module |
| `TEST-INTERNAL_ENGINE-001` | Test 001 for `internal/engine/` module |
| `SPEC-PKG_WALK-001` | Spec 001 for `pkg/walk/` module |

### Module Naming Rules

The module is derived without a lookup table: uppercase every package-path
component, replace hyphens with underscores, and join components with
underscores.

| Package path | Module |
| --- | --- |
| `internal/auth` | `INTERNAL_AUTH` |
| `pkg/walk` | `PKG_WALK` |
| `cmd/idd-cli` | `CMD_IDD_CLI` |

## Documentation Structure

Each module should have documentation under `docs/<path>/`:

```text
docs/
├── internal/
│   ├── auth/           # INTERNAL_AUTH
│   │   ├── design.md   # components + architecture
│   │   ├── contract.md # contracts + observable boundaries
│   │   ├── spec.md     # SPEC records + specification detail
│   │   └── testing.md  # TEST records, coverage, and strategy
│   ├── engine/         # INTERNAL_ENGINE
│   └── model/          # INTERNAL_MODEL
└── pkg/
    └── walk/          # PKG_WALK
```

Every file is self-describing. YAML frontmatter contains only identity:

```yaml
idd:
  version: "1.0"
  package: internal/auth
  document: spec # design, contract, spec, or testing
```

Semantic records are short Markdown sections in the file that owns them:

```markdown
<!-- design.md -->
## Component: AuthModule

The module separates credential verification from transport concerns.

<!-- contract.md -->
## Contract: Authenticator

The boundary returns one public invalid-credentials error.

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

The TEST `Covers` field is the only authored SPEC/TEST relationship. Adding a
SPEC touches only `spec.md`; adding a TEST or changing coverage touches only
`testing.md`. idd-cli derives the reverse edge. Fixed filenames replace
repeated `related_files` paths.

Large Markdown tables are not a declaration format. Separate read-only reports
may summarize the records, while level-two record sections remain the
canonical, validated source.

If none of the four files has an `idd` block, idd-cli continues to use the
legacy frontmatter format. `docs init` refuses legacy marker metadata and old
central catalogs before writing, so migration remains explicit.

See [`examples/compact-module-docs/`](examples/compact-module-docs/) for a
complete four-file example.

## Code Annotations

```go
// Package auth provides authentication utilities.

// Spec: docs/internal/auth/spec.md
// Contract: docs/internal/auth/contract.md
package auth

// @implement SPEC-INTERNAL_AUTH-001
func Login(email, password string) error { ... }

// @test TEST-INTERNAL_AUTH-001
func TestLogin(t *testing.T) { ... }

// @test-contract TEST-INTERNAL_AUTH-002
func TestLoginContract(t *testing.T) { ... }
```

The `@test-contract` annotation links contract test functions to their corresponding test identifiers, not to SPEC identifiers.

## Validation Rules

idd-cli enforces these rules:

| Rule | Description |
|------|-------------|
| `idd-document-schema` | Each file contains valid role-owned Markdown records |
| `idd-document-identity` | Package and role match `docs/<package>/<role>.md` |
| `idd-document-set` | All four self-describing documents exist |
| `idd-document-reference` | Design, contract, and coverage references resolve |
| `idd-document-markdown` | Records are well-formed and do not duplicate legacy metadata |
| `idd-document-migration` | Central catalogs and semantic YAML fields require migration |
| `idd-document-test-kind` | TEST kind agrees with `@test` or `@test-contract` |
| `package-doc-comment` | Every package has doc comment with Spec/Contract paths |
| `spec-required-fields` | Legacy SPEC sections contain their required fields |
| `doc-link-consistency` | Self-described or legacy coverage relationships are consistent |
| `orphan-detection` | No undefined or unreferenced identifiers |
| `doc-code-correspondence` | Doc references match actual code |
| `test-annotation` | Test functions use the correct test annotation |
| `contract-test-coverage` | Legacy CONTRACT identifiers have contract-test coverage |
| `design-sections` | Docs have required ## Design sections |
| `frontmatter-mismatch` | Legacy YAML frontmatter markers match content |

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
│   ├── internal/          # Internal packages (INTERNAL_*)
│   └── pkg/               # External packages (PKG_*)
└── cmd/idd-cli/skills/    # Workflow embedded in the binary
```
