# idd-cli

A CLI tool with an embedded **IDD (Intent-Driven Development)** skill. The
skill guides semantic authoring; idd-cli creates safe document skeletons,
repairs structural identity, and validates the complete documentation/code
traceability graph.

The self-describing format reduces duplicated declarations and relationships;
it does not treat fewer lines as a documentation goal. Human-readable
explanation remains part of the product contract.

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
- **Human-readable depth**: Structural deduplication preserves design rationale,
  implementation boundaries, failure behavior, and verification evidence

## Supported Languages

Source annotations are bound to real declarations with Tree-sitter for:

- Go: `.go`
- TypeScript and TSX: `.ts`, `.tsx`
- JavaScript and JSX: `.js`, `.jsx`
- C++: `.cpp`, `.cc`, `.cxx`, `.hpp`, `.hh`, `.hxx`
- Java: `.java`
- Python: `.py`

The pinned parsers use CGO. A configured extension without a pinned grammar or
a syntax error in a supported file produces a `source-parse` finding; idd-cli
does not fall back to line-oriented or regular-expression binding.

## Skill and CLI Workflow

The two parts have a strict responsibility boundary:

- The IDD skill captures intent and authors meaningful Component, Contract,
  SPEC, TEST, code, and repair changes.
- idd-cli exports that skill, initializes four-file document sets, repairs only
  safe structure, and reports graph problems at exact locations.
- `.idd.yaml` configures validation; it never owns markers or semantic records.

An approved user request is sufficient intent. Teams may keep an `intent.md`
for discussion history, but idd-cli does not require it, add an Intent field to
SPEC records, or build Intent-to-SPEC graph edges.

### Install or refresh the pair

```bash
# Install or upgrade the binary
CGO_ENABLED=1 go install github.com/jingxu9x/idd-cli/cmd/idd-cli@latest

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

# Inspect the exact generated locations that still need authored content
idd-cli docs status docs/internal/auth --format json

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

Semantic records are Markdown sections in the file that owns them. Their fixed
fields are deliberately small, but their explanatory prose should be as
detailed as the behavior requires:

The excerpt below demonstrates record identity and relationship ownership only;
it is not a complete documentation example.

```markdown
<!-- design.md -->
## Component: AuthModule

**Purpose:** Separate credential verification from transport concerns and own
the request-scoped authentication decision.

<!-- contract.md -->
## Contract: Authenticator

**Guarantees:** The boundary returns a complete identity on success and one
public invalid-credentials error for expected rejection without leaking which
credential failed.

<!-- spec.md -->
## SPEC-INTERNAL_AUTH-001: User authentication

- **Design:** `AuthModule`
- **Contract:** `Authenticator`

**Requirement:** Authenticate users with validated credentials.

**Acceptance:** Valid credentials return the expected identity; unknown users
and wrong secrets return the same public error, and no failure returns a
partial identity.

<!-- testing.md -->
## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_AUTH-001`

**Purpose:** Verify accepted and rejected credentials.

**Oracle:** The returned identity and public error category exactly match each
scenario, with no secret or partial identity in failures.

### Evidence and scenarios

Describe the positive, negative, boundary, and failure cases and what each case
proves.
```

### Semantic completeness

Single ownership answers where a fact is written; it does not reduce what must
be explained. A useful document set lets a new maintainer understand:

- why each component exists, what it owns, and why its boundaries were chosen;
- what callers can observe, including inputs, outputs, errors, invariants,
  ownership, side effects, and compatibility promises;
- complete required behavior, edge and failure cases, non-goals, rationale,
  examples, and acceptance evidence;
- how tests prove those claims, including scenarios, fixtures, isolation,
  oracles, and meaningful exclusions.

Use paragraphs, lists, diagrams, examples, and subordinate headings naturally.
Do not restore repeated backlinks or copy every source signature into the
records. Conversely, do not collapse a design or requirement to a one-sentence
summary merely because idd-cli can parse it.

idd-cli validates structure and traceability. A passing report cannot certify
that the prose is sufficient or the design is correct; the paired Skill and
human review own that judgment.

The TEST `Covers` field is the only authored SPEC/TEST relationship. Adding a
SPEC touches only `spec.md`; adding a TEST or changing coverage touches only
`testing.md`. idd-cli derives the reverse edge. Fixed filenames replace
repeated `related_files` paths.

A contract TEST also owns a `Contracts` field naming the Contract records whose
guarantees it proves. A behavior TEST must not declare `Contracts`. Components
may declare `Depends on`; local names and `<package>#<name>` cross-package
references must resolve, and the derived dependency graph must be acyclic.
Component and Contract names must not contain the reserved `#` or `,`
reference separators.

Records may declare `Status`, `Supersedes`, `Deprecated by`, and `Concerns`.
Lifecycle replacements must be same-kind, single-valued after normalization,
and acyclic. A selected concern requires a same-named, non-empty level-three
section for `security`, `concurrency`, `persistence`, `performance`, or
`compatibility`.

Large Markdown tables are not a declaration format. Separate read-only reports
may summarize the records, while level-two record sections remain the
canonical, validated source.

If none of the four files has an `idd` block, idd-cli continues to use the
legacy frontmatter format. `docs init` refuses legacy marker metadata and old
central catalogs before writing, so migration remains explicit.

Migration changes ownership, not documentation depth. Preserve useful
requirements, design reasoning, boundary explanations, failure cases,
examples, and test strategy while removing only duplicated fields, reverse
links, stale API inventories, and central catalogs.

`docs init` deliberately writes visible
`<!-- idd:scaffold slot="..." -->` markers. Generated files remain invalid
until every marker is removed and its bounded section contains real,
non-placeholder prose. `docs status` returns a deterministic work list with
file, line, role, slot, and reason. Once a record exists, the same status
command also lists any missing or placeholder required record field, such as
`Acceptance` or `Oracle`; adding only a heading therefore cannot produce a
false `complete` result. `docs fix` preserves scaffold state and never invents
requirements to make findings disappear.

See
[`examples/self-describing-module-docs/`](examples/self-describing-module-docs/)
for a complete four-file example with detailed human-readable narratives.

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

Equivalent declaration-attached comments are recognized in all supported
languages (`//` or `/* */`, and `#` in Python). The parser ignores annotation
text in strings and does not bind comments inside function bodies or detached
comments. Public declarations require `@implement`; recognized test
declarations require `@test` or `@test-contract`. The annotation and declaration
must be adjacent apart from comments and limited whitespace.

## Validation Rules

idd-cli enforces these rules:

| Rule | Description |
|------|-------------|
| `idd-document-schema` | Each file contains valid role-owned records, required fields, lifecycle state, and concern sections |
| `idd-document-identity` | Package and role match `docs/<package>/<role>.md` |
| `idd-document-set` | All four self-describing documents exist |
| `idd-document-reference` | Design, contract, and coverage references resolve |
| `idd-document-incomplete` | Generated scaffold slots still need authored content |
| `idd-document-markdown` | Records are well-formed and do not duplicate legacy metadata |
| `idd-document-migration` | Central catalogs and semantic YAML fields require migration |
| `idd-document-test-kind` | TEST kind agrees with `@test` or `@test-contract` |
| `component-dependency-cycle` | The derived Component dependency graph remains acyclic |
| `source-parse` | A configured source lacks a pinned grammar or has an invalid syntax tree |
| `package-doc-comment` | Every package has doc comment with Spec/Contract paths |
| `spec-required-fields` | Legacy SPEC sections contain their required fields |
| `doc-link-consistency` | Self-described or legacy coverage relationships are consistent |
| `orphan-detection` | No undefined or unreferenced identifiers |
| `doc-code-correspondence` | Doc references match actual code |
| `public-func-annotation` | Public source declarations use `@implement` |
| `test-annotation` | Recognized test declarations use the correct test annotation |
| `contract-test-coverage` | Legacy markers and self-describing named Contracts have contract-test evidence |
| `design-sections` | Docs have required ## Design sections |
| `frontmatter-mismatch` | Legacy YAML frontmatter markers match content |

Narrative validation rejects structural emptiness and obvious placeholders; it
does not use length thresholds or claim to verify semantic completeness.

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
