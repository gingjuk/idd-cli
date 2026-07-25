---
name: intent-driven-development
description: Intent-Driven Development workflow paired with idd-cli for authoring, validation, and repair
license: MIT
compatibility: opencode
metadata:
  audience: agents
  workflow: development
  protected: true  # Modification requires explicit approval
---

# Intent-Driven Development (IDD)

> ⚠️ **Protected Skill** — Changes to this skill require explicit user approval. Do not modify without consulting the user first.

IDD is a framework-agnostic development workflow that transforms user intent into implemented, tested, and committed code through defined phases.

## Core Philosophy

**Four Pillars of IDD:**

1. **Intent is Truth** — Every deliverable originates from an approved `intent.md`. Nothing exists unless it is defined in intent and approved. Implementation without intent is guesswork.

2. **Specs are Code** — Specifications must be formalized, precise, and unambiguous. Like production code, specs undergo review, version control, and validation. Vague specs produce vague implementations.

3. **Tests Prove Intent** — Tests verify behavior against specs, not implementation details. A test proving "the function returns X for input Y" is valid; a test proving "the function calls dependency Z" is not.

4. **Nothing Guessed** — When anything is unclear, CLARIFY. Never guess, assume, or fill gaps with assumptions. Ambiguity resolved by the implementor (not the author) introduces interpretation risk. Ask, don't guess.

## Skill and CLI Contract

Use this skill and idd-cli as one workflow with separate responsibilities:

- The skill turns approved intent into meaningful design, contract, SPEC, TEST,
  code, and repair decisions.
- `idd-cli generate skill` exports the exact workflow embedded in the current
  binary. Regenerate the installed skill after upgrading idd-cli.
- `idd-cli docs init` creates a new package's four-document skeleton.
- `idd-cli docs fix` repairs identity and missing skeletons only. It never
  authors requirements, contracts, designs, purposes, or coverage.
- `idd-cli run .` checks the complete documentation/code graph and reports the
  exact record or source location that the skill must repair.
- `.idd.yaml` configures collection, validation, and output. It is never a
  marker or semantic-record source.

The skill is the semantic author; idd-cli is the executable guardrail. Neither
replaces the other.

## Workflow

```text
Install or upgrade idd-cli
    ↓
idd-cli generate skill -o <agent-skill-path>/SKILL.md
    ↓
Intent Capture → intent.md → plan.md
    ↓
idd-cli docs init <pkg_path>                 (new package only)
    ↓
Skill authors design.md + contract.md → spec.md → testing.md
    ↓
Tests and code annotations → RED → GREEN → IMPROVE
    ↓
idd-cli docs fix docs/<pkg_path>             (safe structure only)
    ↓
idd-cli run . --format llm-markdown
    ↺ Skill repairs the owning Markdown record or source annotation
    ↓
Repository tests and review
    ↓
idd-cli run . --format json                  (final project gate)
    ↓
Commit
```

**Phase Details:**

| Phase | Skill responsibility | idd-cli checkpoint |
| --- | --- | --- |
| 0. Tool alignment | Load the workflow exported by the installed binary | `idd-cli generate skill -o <path>` |
| 1. Intent and plan | Capture scope, success criteria, constraints, and work order | No document mutation |
| 2. Package bootstrap | Decide whether this is new, self-describing, or legacy documentation | `idd-cli docs init <pkg_path>` for a new package only |
| 3. Document authoring | Write role-owned, human-readable Markdown records | Optional `docs fix <file>` for identity only |
| 4. TDD implementation | Write TEST records/tests first, then implementation and annotations | No partial-project validity claim |
| 5. Repair loop | Interpret findings and edit the single owning record or source | `idd-cli run . --format llm-markdown` |
| 6. Final review | Run repository gates and review the complete change | `idd-cli run . --format json` |
| 7. Knowledge capture | Preserve approved intent, decisions, and durable documentation | Re-run the final gate after changes |

Always run validation from the project root. A narrower documentation target
does not narrow source annotation collection, so
`idd-cli run docs/<single-package>` is not a package-only validity check.

**TDD Workflow:**

1. Write test first (RED) — test should FAIL
2. Write minimal implementation (GREEN) — test should PASS
3. Refactor (IMPROVE) — verify coverage 80%+

**Troubleshoot:** check test isolation → verify mocks → fix implementation (not tests, unless tests are wrong).

## Document Dependency Chain

Documents are created in dependency order. Every package document is
self-describing and owns only the declarations for its role.

```text
User Intent
     ↓
intent.md / plan.md
     ↓
design.md + contract.md       ← architecture and observable boundaries
     ↓
spec.md                       ← SPEC declarations, rationale, and examples
     ↓
testing.md                    ← TEST declarations, coverage, and strategy
     ↓
code annotations             ← @implement / @test / @test-contract
```

For a new package, run `idd-cli docs init <pkg_path>` before authoring these
records. For an existing self-describing package, edit the owning file directly
and do not rerun initialization.

**Dependency Rules:**

| File | Owns | Must not duplicate |
| --- | --- | --- |
| `design.md` | `## Component:` records plus architecture decisions | SPEC or TEST records |
| `contract.md` | `## Contract:` records plus interfaces, errors, and invariants | SPEC backlinks |
| `spec.md` | `## SPEC-...:` records plus examples, edge cases, and rationale | TEST coverage |
| `testing.md` | `## TEST-...:` records, coverage, strategy, and fixtures | SPEC backlinks |
| Code | Executable behavior and IDD annotations | Copied specification prose |

### Specification granularity

A SPEC is a cohesive, externally meaningful behavior or invariant. It is not a
required one-to-one mirror of a function, method, type, or test case.

- One SPEC may annotate multiple declarations.
- A public declaration may reference the same behavioral SPEC as related
  declarations.
- Function signatures, parameter lists, return lists, and comments stay in
  source code or generated API documentation.
- Add a separate SPEC only when the requirement, contract, design ownership, or
  acceptance behavior is independently meaningful.

### Canonical relationship ownership

- `testing.md` TEST records' `Covers` fields own TEST → SPEC coverage.
- idd-cli derives the reverse SPEC → TEST edge; authors do not repeat a
  `Tests` field.
- `spec.md` SPEC records' `Contract` fields own SPEC → contract selection;
  `contract.md` does not repeat `Implements` backlinks.
- `spec.md` SPEC records' `Design` fields own SPEC → design component
  selection.
- `@implement` references a SPEC.
- `@test` and `@test-contract` reference a TEST.

## Document Hierarchy

All IDD documents are stored under the package path:

```text
docs/
└── <pkg_path>/           # mirrors internal/auth, pkg/pattern, cmd/idd-cli, ...
    ├── design.md         # components and architecture
    ├── contract.md       # contracts and observable boundaries
    ├── spec.md           # SPEC records and specification detail
    └── testing.md        # TEST records, coverage, and strategy
```

**Rules:**

- `<pkg_path>` mirrors the package directory relative to project root (e.g., `internal/auth` → `docs/internal/auth/`)
- `<MODULE>` is derived deterministically from `<pkg_path>`.
- The four Markdown filenames are fixed, so `related_files` paths are not
  repeated.
- Every file has an `idd` block containing only `version`, `package`, and
  `document`.
- Self-describing Markdown must not also contain legacy `markers` or
  `related_files` frontmatter.
- Level-two role headings are the authoritative declarations. Required fields
  stay compact; the rest of each section is ordinary Markdown.
- Large Markdown tables are not declaration sources. A separate generated
  report may summarize records, but the four owned documents keep level-two
  record sections as their only canonical source.
- If none of the four files has an `idd` block, the directory uses the legacy
  template. If one has it, all four must. Do not mix metadata models.

## Document Templates

### Planning Templates

**`intent.md`** (in `<module>/.planning/`):

```markdown
# Intent: <feature name>

## User Request
[What the user asked for — verbatim or paraphrased]

## Success Criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Constraints
- [ ] Constraint 1
- [ ] Constraint 2

## Out of Scope
- Item 1
- Item 2
```

**`plan.md`** (in `<module>/.planning/`):

```markdown
# Plan: <feature name>

## Status
- [ ] Not started

## Phase 1: [Name]
| Task | Status | Notes |
|------|--------|-------|
| Task description | pending | |

## Risks
| Risk | Impact | Mitigation |
|------|--------|------------|
| Risk description | High/Med/Low | Mitigation approach |

## Next Action
[Immediate next step to take]
```

### IDD Document Templates

Create a new package skeleton with:

```bash
idd-cli docs init <pkg_path>
```

The command never overwrites existing IDD or legacy metadata and does not
invent placeholder requirements or identifiers. It creates minimal identity
frontmatter and role-appropriate narrative scaffolds:

```yaml
idd:
  version: "1.0"
  package: internal/auth
  document: spec # design, contract, spec, or testing
```

The identity appears between `---` delimiters in each named Markdown file.
Semantic YAML fields such as `components`, `contracts`, `specs`, or `tests` are
invalid; migrate them into Markdown records.

**Canonical record shapes:**

| Document | Level-two record | Required content |
| --- | --- | --- |
| `design.md` | `## Component: <name>` | A human-readable component narrative |
| `contract.md` | `## Contract: <name>` | Inputs, outputs, errors, and invariants |
| `spec.md` | `## SPEC-...: <title>` | `Design`, `Contract`, and `Requirement` |
| `testing.md` | `## TEST-...: <title>` | `Kind`, `Covers`, and `Purpose` |

This table explains the format; it is not a declaration registry. Each actual
record is a Markdown section. `Kind` is `test` for behavioral tests or
`contract` for contract tests. Every SPEC must be named by at least one TEST
`Covers` field.

**`design.md`:**

```markdown
---
idd:
  version: "1.0"
  package: internal/auth
  document: design
---

# Design: internal/auth

## Component: AuthModule

Explain the component's responsibility and boundaries in concrete terms.

## Architecture

## Package Layout

## Function Composition

## Dependencies

## Testability Hooks
```

Each section requires concrete package-specific content. Empty headings produced
by `docs init` are an honest scaffold, not completed documentation.

**`contract.md`:**

````markdown
---
idd:
  version: "1.0"
  package: internal/auth
  document: contract
---

# Contracts: internal/auth

## Contract: Authenticator

State the observable inputs, outputs, errors, and invariants.

```go
type Authenticator interface {
    Authenticate(ctx context.Context, credentials Credentials) (Identity, error)
}
```
````

Do not add SPEC backlinks; the SPEC record's `Contract` field owns that
reference.

**`spec.md`:**

```markdown
---
idd:
  version: "1.0"
  package: internal/auth
  document: spec
---

# Specifications: internal/auth

## SPEC-INTERNAL_AUTH-001: Credential rejection details

- **Design:** `AuthModule`
- **Contract:** `Authenticator`

**Requirement:** Authenticate users with validated credentials.

Record acceptance examples, edge cases, and rationale as normal Markdown.
```

The level-two ID heading and fixed fields are required. Do not copy source
signatures, parameter lists, or return lists into this file.

**`testing.md`:**

```markdown
---
idd:
  version: "1.0"
  package: internal/auth
  document: testing
---

# Testing: internal/auth

## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_AUTH-001`

**Purpose:** Verify accepted and rejected credentials.

## Strategy

Explain test boundaries, fixtures, important scenarios, and why contract tests
are separate from behavioral tests.
```

Do not repeat `Spec Coverage`; this TEST record's `Covers` field owns the
relationship.

### Code Annotation Templates

Where `<pkg_path>` is the package directory relative to project root (e.g., `internal/auth`, `pkg/pattern`). `<MODULE>` is the abbreviated identifier derived from `<pkg_path>` — see **Module prefix** rules below.

**Source file** (after package declaration):

```go
// Package mymodule provides <description>.
//
// Spec: docs/<pkg_path>/spec.md
// Contract: docs/<pkg_path>/contract.md
package mymodule

// @implement SPEC-<MODULE>-001
func PublicFunction() {
    // implementation
}
```

**Test file** (`<pkg>_test.go`):

```go
// Package mymodule provides tests for <description>.
//
// Spec: docs/<pkg_path>/spec.md
// Test: docs/<pkg_path>/testing.md
package mymodule

// @test TEST-<MODULE>-001
func TestPublicFunction(t *testing.T) {
    // test implementation
}
```

**Contract test file** (`<pkg>_contract_test.go`):

Contract tests validate that the implementation satisfies the contract
interface. Each test function uses `@test-contract TEST-<MODULE>-NNN` and the
corresponding `testing.md` record has `kind: contract`.

```go
// Package mymodule provides contract tests for <description>.
//
// Spec: docs/<pkg_path>/spec.md
// Test: docs/<pkg_path>/testing.md
// Contract: docs/<pkg_path>/contract.md
package mymodule

// @test-contract TEST-<MODULE>-002
func TestContractPublicFunction(t *testing.T) {
    // Validate implementation matches contract
}
```

## Identifier Format

**Pattern:** `<TYPE>-<MODULE>-<NUMBER>` — Exactly 3 segments separated by hyphens

| Prefix | Meaning | Example |
| ------ | ------- | ------- |
| `SPEC-` | Functionality specification | `SPEC-INTERNAL_AUTH-001` |
| `TEST-` | Test case | `TEST-INTERNAL_AUTH-001` |

**Module prefix** derived from `<pkg_path>` by a deterministic rule — no lookup table required:

1. Split `<pkg_path>` by `/`.
2. Uppercase each component; replace hyphens (`-`) with underscores (`_`).
3. Join components with `_`.

**Examples:**

| `<pkg_path>`          | `<MODULE>`              |
|-----------------------|-------------------------|
| `internal/auth`       | `INTERNAL_AUTH`         |
| `internal/engine`     | `INTERNAL_ENGINE`       |
| `internal/collector`  | `INTERNAL_COLLECTOR`    |
| `internal/config`     | `INTERNAL_CONFIG`       |
| `pkg/pattern`         | `PKG_PATTERN`           |
| `pkg/walk`            | `PKG_WALK`              |
| `cmd/idd-cli`         | `CMD_IDD_CLI`           |
| `auth`                | `AUTH`                  |

**Invalid (will be flagged by idd-cli):**

- `SPEC-BE-007-007` — Triple-segment identifier
- `PATTERN-*`, `WALK-*` — Section identifiers, not IDD identifiers

## Linkage Rules

The owning document records each relationship once and idd-cli builds the
complete graph:

```text
design.md Component             contract.md Contract
       AuthModule                     Authenticator
          ↑                                  ↑
          └──── SPEC-AUTH-001 ───────────────┘
                    spec.md
                    ↑ derived reverse edge
                    │
          testing.md: TEST-AUTH-001 covers SPEC-AUTH-001

CODE:          @implement SPEC-AUTH-001
TEST:          @test TEST-AUTH-001
CONTRACT TEST: @test-contract TEST-AUTH-002
```

**Verification rules:**

- Every SPEC `Design` value names a `design.md` Component record.
- Every SPEC `Contract` value names a `contract.md` Contract record.
- Every TEST `Covers` value names an existing SPEC.
- Every SPEC has at least one derived TEST backlink.
- Document identifiers exist in matching code annotations.
- Record headings and fixed fields are reported at their Markdown source lines.
- `Kind: contract` records are used by contract test annotations.

### Legacy compatibility

Directories where none of the four fixed files has an `idd` block continue to
use frontmatter markers, SPEC `**Tests:**`, and TEST `**Spec Coverage:**`.
Migrate one package at a time:

1. Add the minimal `idd` identity to all four fixed documents.
2. Move component declarations to `design.md` `## Component:` sections.
3. Move contract declarations to `contract.md` `## Contract:` sections.
4. Move meaningful SPEC records and `Design`/`Contract` selections to
   `spec.md`.
5. Move TEST records and each relationship to `testing.md` `Covers` fields.
6. Remove legacy `markers`, `related_files`, and repeated relationship fields
   while keeping valuable narrative content.
7. Run `idd-cli docs fix docs/<pkg_path>` and then
   `idd-cli run . --format llm-markdown`.

**Annotation Placement:**

- `@implement`, `@test`, `@test-contract` are ONLY allowed on function/type declarations
- NOT permitted inside function bodies, closures, or inline expressions

## Testing Requirements

### Coverage

Minimum coverage: 80%

| Type | What | When |
| ---- | --- | ---- |
| Unit tests | Individual functions, utilities | Always |
| Integration tests | API endpoints, database operations | Always |
| E2E tests | Critical user flows (Playwright) | Critical paths |

**Contract tests:** For every `contract.md`, there MUST exist `*_contract_test.go` in the same directory as source.

## Anti-Patterns

- **Skip planning** — "Just start coding" leads to scope creep and rework
- **Write code then retrofit docs** — Intent and plan must come first
- **Skip tests** — TDD is mandatory for all new features
- **Skip coverage verification** — 80%+ is the minimum
- **Ignore blockers** — Document and escalate
- **Delete .planning/** when done — Keep for continuity
- **Annotate inside function bodies** — Annotations only on function/type declarations
- **Triple-segment identifiers** — Use `SPEC-BE-007` not `SPEC-BE-007-007`
- **One SPEC per symbol** — Produces API-mirror documents instead of behavioral specifications
- **Copied source API** — Signatures, parameters, and returns belong in source or generated API docs
- **Repeated backlinks** — Store coverage only in TEST `Covers`; let idd-cli derive the reverse edge
- **YAML semantic catalogs** — YAML identifies the document; Markdown records describe the behavior
- **Large declaration tables** — Tables are derived views, not editable record sources
- **Central marker catalog** — Creates one package-wide write hotspot
- **Mixed self-describing and legacy metadata** — Creates competing sources of truth
- **Placeholder document content** — Never use `TBD`, `auto-generated`, or an identifier copied as its own title
- **Stale installed skill** — Regenerate the agent skill after upgrading idd-cli
- **Package-targeted final gate** — `run docs/<pkg_path>` does not scope source collection; validate the project root

## IDD CLI Operational Loop

idd-cli is the executable companion to this skill. Use it at these checkpoints,
not as a semantic document generator.

### 1. Align the installed skill with the binary

```bash
idd-cli --version
idd-cli generate skill -o idd-skill.md
idd-cli skills --format json
```

`generate skill` reads the copy embedded in the binary. Regenerate it after an
idd-cli upgrade so the agent does not author an older document format than the
validator accepts. Install `idd-skill.md` as `SKILL.md` through the agent's
normal skill mechanism.

### 2. Bootstrap or adopt package documents

For a new package with no IDD metadata:

```bash
idd-cli docs init internal/auth
```

For an existing self-describing package, do not call `docs init`; edit the
owning Markdown records. For a legacy package, follow the explicit migration
steps first because initialization refuses mixed metadata.

After independently generating or moving documents, normalize only safe
structure:

```bash
idd-cli docs fix docs/internal/auth
idd-cli docs fix docs/internal/auth/testing.md
```

The directory form may create missing skeletons and repairs all four identities.
The file form has one write target. Both preserve the Markdown body
byte-for-byte.

### 3. Run the skill-guided repair loop

After the four documents, tests, implementation, and annotations form one
coherent checkpoint:

```bash
idd-cli run . --format llm-markdown
```

An invalid graph exits non-zero but still emits the complete report. Repair by
finding ownership:

- identity or missing-document findings → run `docs fix` on the reported file
  or directory;
- schema, Markdown, reference, migration, or TEST-kind findings → use this
  skill to edit the reported owning record;
- annotation or doc/code findings → edit the reported source or test
  declaration;
- never add a reverse backlink merely to silence a finding; repair the one
  canonical owner and let idd-cli derive the graph.

Repeat until the project report passes. Always use `.` for this validity gate:
the documentation argument can narrow document collection, but source
annotations are still collected from the current project working tree.

### 4. Run the final gate

Run repository tests first, then produce the stable CI-facing report:

```bash
idd-cli run . --format json
```

Use `llm-markdown` for agent repair, `markdown` for a human readout, and `json`
for automation. A commit is ready only when repository tests and the final
project-level IDD check pass.

**Configuration (`.idd.yaml`):**

This file configures the guardrail. It never owns Components, Contracts, SPECs,
TESTs, or coverage.

```yaml
version: "1.0"

docs:
  patterns:
    - "docs/**/*.md"
  identifier_patterns:
    spec: "SPEC-[A-Z0-9_]+-[0-9]+"
    test: "TEST-[A-Z0-9_]+-[0-9]+"
    test_contract: "TEST-[A-Z0-9_]+-[0-9]+"
  ignore_paths: []          # glob patterns for doc paths to skip

code:
  patterns:
    - "**/*.go"
  annotations:
    spec: "@implement"
    test: "@test"
    test_contract: "@test-contract"
  ignore_paths: []          # glob patterns for source files to skip

validation:
  # Self-described coverage or legacy bidirectional relationship consistency
  require_doc_link_consistency: true
  # Disallow identifiers with no connections in the graph
  allow_orphans: false
  # Every SPEC must have at least one TEST coverage relationship
  require_spec_test_coverage: true
  # Legacy contract identifiers must have contract-test coverage
  require_contract_test_coverage: true
  # design.md files must contain all required sections
  require_design_sections: true
  # Every identifier must appear in both docs and code annotations
  require_doc_code_correspondence: true
  # Every public function/type must have an @implement annotation
  require_public_func_annotation: true
  # Every .go file must have a package doc comment with Spec/Contract paths
  require_package_doc_comment: true
  # Legacy docs need related_files; self-describing files use fixed roles
  require_related_files: true
  # Every test function must have an @test or @test-contract annotation
  require_test_annotation: true
  # @implement/@test/@test-contract must include an identifier
  require_annotation_identifier: true
  # Multiple annotations of the same type must be comma-separated on one line
  require_annotation_on_same_line: true
  # Semantic similarity check between doc describe and code comments
  consistency_check:
    enabled: true
    threshold: 0.3          # 0.0–1.0; below this triggers a warning

output:
  file: "idd-report.json"
  include_graph: false
  verbose: false
```

**Validation Rules:**

1. **Document schema** — Every file contains concrete, role-owned Markdown records.
2. **Document identity** — `package` and `document` match `docs/<pkg_path>/<role>.md`.
3. **Document set** — All four fixed self-describing Markdown documents exist.
4. **Migration safety** — A package-local central marker catalog is rejected.
5. **Reference integrity** — Design, contract, and coverage references resolve.
6. **Coverage completeness** — Every SPEC has at least one derived TEST backlink.
7. **Identifier format** — Identifiers match `<TYPE>-<MODULE>-<NUMBER>`.
8. **Doc/code correspondence** — Document identifiers and code annotations agree.
9. **Annotation placement** — Annotations sit on function or type declarations.
10. **Narrative structure** — Required record headings and design sections contain meaningful content.
11. **Legacy compatibility** — Legacy rules apply only when no fixed file has an `idd` block.

**Output:** Finding-centered JSON, Markdown, or LLM-oriented Markdown with
rule, exact location, identifier, problem, and suggested fix.
