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
approved request + .planning
        ↓
docs init <package> when new
        ↓
docs status reports generated work slots
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
- `docs init` creates a new four-document skeleton with explicit unfilled
  markers rather than fake records;
- `docs status` returns the exact file, line, role, slot, and reason for every
  generated location that still needs authored content;
- `docs fix` repairs identity and missing skeletons without rewriting prose;
- `trace <ID>` returns the canonical declaration plus typed occurrences and
  relation provenance across docs, source, and tests;
- `docs impacted --base <revision>` derives a stable semantic-review queue
  from Git changes; and
- `run .` validates the complete project graph and emits exact findings.

Regenerate the installed skill after upgrading idd-cli. `.idd.yaml` is only
validator configuration and must not contain package markers or semantic
records.

The approved request is the source of truth. An `intent.md` may be kept as
project history, but it is not required, SPEC has no Intent-source field, and
idd-cli creates no Intent-to-SPEC graph.

The final validation normally runs as `idd-cli run .` from the project root.
The optional positional path is itself a project root, so
`idd-cli run /path/to/another-worktree` discovers configuration and collects
both documents and source annotations there. Use `docs status` or
`docs review-context --docs-path` when documentation input should be narrowed.

## IDD Identifier Format

**Format:** `<TYPE>-<MODULE>-<NUMBER>`

| Prefix | Meaning | Example |
| ------ | ------- | ------- |
| `SPEC-` | Functionality specification | SPEC-XX-001 |
| `TEST-` | Test case | TEST-XX-001 |

Note: CONTRACT and DESIGN are no longer primary identifiers. contract.md describes contracts that implement SPECs. design.md describes architecture decisions.

`MODULE` is the documentation unit's stable `idd.namespace`. `docs init`
derives its first value from the package path, but later package moves preserve
the namespace and therefore preserve historical IDs. ID syntax is owned by the
versioned IDD protocol; obsolete `docs.identifier_patterns` is ignored with a
migration warning.

## Syntax-tree source binding

idd-cli binds comments to real declarations with pinned Tree-sitter grammars
for Go, TypeScript, TSX, JavaScript/JSX, C++, Java, and Python. Annotation text
inside strings is ignored automatically. A configured extension without a
pinned grammar or a syntax error in a supported file produces a
`source-parse` finding and disables binding for that file; there is no regex
fallback.

The annotation identifier is also the code-to-document join key. idd-cli
requires implementation and TEST evidence for active SPECs and resolves every
authored annotation. Planned SPECs may intentionally have neither. Ordinary
helper APIs and unrelated tests need no IDD record unless optional census rules
are enabled.
Do not repeat `Spec:`, `Contract:`, or `Test:` document paths in source-file
headers; those paths are redundant with the self-describing document set and
are not consumed by validation.

Use ignore ranges only when a real source comment intentionally demonstrates
annotation syntax and must not count as evidence.

### Nolint Directives

```go
// idd:ignore start
// @implement SPEC-XX-001
func ExampleOnly() {}
// idd:ignore end
```

`idd:ignore-start` and `idd:ignore-end` are accepted compatibility spellings.
The directive must be the complete normalized comment line; merely mentioning
`idd:ignore` in prose does not open a range.

## Document Structure

Documents are organized by documentation unit in `docs/`. Without an explicit
mapping, each scanned package retains its own matching directory. `docs.units`
may instead map one stable behavior boundary to several physical source globs;
every scanned source must resolve to exactly one configured unit.

```text
docs/
├── internal/auth/
│   ├── design.md     # Components and architecture decisions
│   ├── contract.md   # Contracts, interfaces, and invariants
│   ├── spec.md       # SPEC declarations and detail
│   └── testing.md    # TEST declarations, coverage, and strategy
└── IDD.md            # This file
```

Each document has an `idd` frontmatter block containing version, its current
package location, and a stable namespace:

```yaml
# spec.md
idd:
  version: "1.1"
  package: internal/auth
  namespace: INTERNAL_AUTH
```

The exact lowercase basename is the only role authority: `design.md`,
`contract.md`, `spec.md`, or `testing.md`. `idd.document` is invalid without a
compatibility period, and headings do not override the filename.

Each documentation unit must have all four files. They have no
line-count limit: keep the complete human-readable design, contract,
requirements, examples, boundaries, and test rationale in the owning files.
Do not create split role names such as `design-auth.md`, `spec.part.md`, or
`testing_extra.md`.

An `idd-document-filename` finding for a split file is also a repair protocol.
Its JSON `suggested_fix` and LLM Markdown `Fix` name the split source,
canonical target, semantic material that must survive the merge, safe deletion
order, and the exact follow-up checks. `docs status` returns equivalent guidance
when a split blocks completion inspection. The agent reads both documents,
merges without collapsing detailed prose into a summary, verifies that nothing
was lost, and only then removes the fragment. idd-cli reports and checks this
work but never performs the semantic merge itself.

The semantic declarations live in human-readable Markdown records:

This excerpt shows the parseable record boundary, not the expected depth of a
finished design or specification.

```markdown
<!-- design.md -->
## Component: AuthModule

**Purpose:** Own credential verification without importing transport policy.

**Ownership:** Own request-scoped verification orchestration and its result.

**Boundary:** Transports, persistence, and session state remain outside.

**Decisions:** Use one public rejection category for unknown users and bad secrets.

<!-- contract.md -->
## Contract: Authenticator

**Guarantees:** Return a complete identity on success and one stable public
error for expected credential rejection.

<!-- spec.md -->
## SPEC-INTERNAL_AUTH-001: User authentication

- **Components:** `AuthModule`
- **Contracts:** `Authenticator`

**Requirement:** Authenticate users with validated credentials.

**Acceptance:** Valid credentials return the expected identity; rejected
credentials expose no field-specific detail and return no partial identity.

<!-- testing.md -->
## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_AUTH-001`

**Purpose:** Verify accepted and rejected credentials.

**Oracle:** The identity and public error category exactly match the scenario.

### Evidence and scenarios

Explain the relevant cases and the observable property each case proves.
```

The small fixed-field block is a parsing boundary, not a prose limit. Each
record should preserve enough context for a human reader to understand the
reason for the behavior, its guarantees, design and implementation boundaries,
failure cases, trade-offs, and verification evidence. Rich paragraphs,
examples, diagrams, and subordinate headings are expected when the subject
needs them.

Structural validation cannot judge whether an architectural explanation is
complete or a decision is sound. Passing idd-cli confirms identity,
traceability, and record shape; the paired Skill and human review must still
reject thin summaries that force readers to reconstruct intent from source
code.

See
[`examples/self-describing-module-docs/`](../examples/self-describing-module-docs/)
for a complete record set with design rationale, observable contracts,
implementation boundaries, failure cases, verification strategy, and
identifier-derived source annotations. The
[`examples/README.md`](../examples/README.md) records its validation commands.

The TEST `Covers` field is the only authored SPEC/TEST relationship. idd-cli
derives the reverse SPEC-to-TEST edge. `design.md` owns `## Component:`
records and `contract.md` owns `## Contract:` records. Self-describing
documents must not repeat legacy marker, `related_files`, `Tests`, or
`Spec Coverage` metadata.

Component `Purpose`, `Ownership`, `Boundary`, and `Decisions`; Contract `Guarantees`; SPEC `Requirement` and
`Acceptance`, and TEST `Purpose` and `Oracle` are the required prose anchors.
A `Kind: contract` TEST must name the Contract records it proves in
`Contracts`; a behavior TEST must not. Component `Depends on` references and
record lifecycle relationships are normalized into typed graph edges and
checked for missing targets, ambiguity, self-links, and cycles.

Contract narrative distinguishes deliberate Guarantees from non-guarantees,
known limitations, and compatibility commitments. Current defects and
incidental behavior stay visible as limitations instead of becoming permanent
promises. A unit with no stable Contract may leave records empty only after
explaining that decision in prose.

Lifecycle controls delivery checks: `planned` records preserve approved intent
without implementation or evidence; `active` records require both; deprecated
and superseded records retain history without active correspondence gates.

Generated documents contain `idd:scaffold` markers. A slot is complete only
when the marker is gone and the bounded section contains non-placeholder
authored content. After a record heading is added, `docs status` continues to
report each missing required fixed field from the shared role schema. Thus a
SPEC without `Acceptance`, or a TEST without `Oracle`, remains incomplete even
after its collection marker is removed. `docs fix` never removes a marker or
fabricates that content.

Level-two headings delimit records. The prose inside each record is ordinary
Markdown, and large registry tables are rejected as an authored declaration
format. This removes duplicated registries without removing useful narrative.

A SPEC represents cohesive behavior, not one function. Keep function signatures,
parameters, and returns in source code or generated API documentation, while
keeping behavior, rationale, edge cases, non-goals, and acceptance evidence in
the SPEC.

Packages where none of the four files has an `idd` block continue to use the
legacy frontmatter format, so migration can happen one package at a time.
`docs init` refuses legacy marker metadata and central catalogs before writing;
migrate those relationships explicitly before enabling self-describing mode.
Legacy marker compatibility does not accept the removed `idd.document` field:
an `idd` block immediately uses filename-owned roles and version/package/namespace
identity.
The migration must preserve useful rationale, boundary discussion, failure
behavior, examples, and test strategy; a smaller file is not automatically a
better document.

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

Ordinary language package/module comments remain useful human documentation,
but they do not carry IDD linkage. `docs init` and `docs fix` never edit source;
the explicit declaration annotation and matching document record are the
complete association.

## Validation Rules

1. **Document Schema** — Each role owns correctly shaped Markdown records.
2. **Filename Authority** — Only the four canonical basenames define IDD
   roles; split role documents receive a path-specific agent merge prompt.
3. **Stable Identity** — Version, current package path, and namespace agree
   across the four-file unit.
4. **Document Set** — Every configured documentation unit has all four
   self-describing documents and each source maps unambiguously.
5. **Migration Safety** — A package-local central marker catalog is rejected.
6. **Reference Integrity** — Component, Contract, coverage, dependency, and lifecycle references resolve.
7. **Lifecycle Completeness** — Every active SPEC has implementation and TEST evidence.
8. **TEST Kind** — `Kind: test` uses `@test`; `Kind: contract` uses
   `@test-contract`.
9. **Doc/Code Correspondence** — Active records have declarations and every annotation resolves.
10. **Legacy Compatibility** — Legacy marker rules apply when no `idd` block
    is present; this does not accept the removed self-describing
    `idd.document` field.

## Configuration Example

This project configuration controls collection, validation, and reporting. It
does not own Components, Contracts, SPECs, TESTs, or coverage.

```yaml
version: "1.0"

docs:
  patterns:
    - "docs/**/*.md"
  units:
    - package: auth
      sources: ["internal/auth/**", "internal/session/**"]

code:
  patterns:
    - "**/*.go"
    - "**/*.ts"
    - "**/*.tsx"
    - "**/*.js"
    - "**/*.jsx"
    - "**/*.cpp"
    - "**/*.cc"
    - "**/*.cxx"
    - "**/*.hpp"
    - "**/*.hh"
    - "**/*.hxx"
    - "**/*.java"
    - "**/*.py"
  annotations:
    spec: "@implement"
    test: "@test"
    test_contract: "@test-contract"

validation:
  require_doc_link_consistency: true
  require_spec_fields: true
  allow_orphans: false
  require_spec_test_coverage: true
  # This joins document records to declaration annotations by identifier.
  require_doc_code_correspondence: true

output:
  file: "idd-report.json"
  include_graph: true
  verbose: false
```

## Usage

```bash
# Build
CGO_ENABLED=1 go build -o idd-cli ./cmd/idd-cli

# Export and install this binary's paired skill
./idd-cli generate skill -o idd-skill.md

# Create documents for new packages only
./idd-cli docs init internal/auth internal/config

# Inspect several packages as one authoring work list
./idd-cli docs status docs/internal/auth docs/internal/config --format json

# Normalize safe structural identity
./idd-cli docs fix docs/internal/auth docs/internal/config

# Normalize selected independently generated documents
./idd-cli docs fix docs/internal/auth/testing.md docs/internal/config/testing.md

# Let the skill repair the complete finding report
./idd-cli run . --format llm-markdown

# Inspect one ID or derive the queue from Git changes
./idd-cli trace SPEC-INTERNAL_AUTH-001 --format llm-markdown
./idd-cli docs impacted --base main --format llm-markdown

# Final automation-facing gate
./idd-cli run . --config .idd.yaml --format json
```

An invalid graph exits non-zero after writing its report. Use the skill to edit
semantic records or source annotations; use `docs fix` only for identity and
missing-document findings.

The three structural document commands accept one or more targets, deduplicate
repeated or overlapping paths, and return deterministic aggregate output.
`docs init` and `docs fix` preflight the whole batch before expected writes.
`run` and `lint` intentionally retain one project root because they build one
cross-package graph. Relative configuration paths resolve inside that root;
relative report output paths remain anchored to the invoking directory.

Deprecated configuration is migration-only. Explicitly configuring
`validation.consistency_check` produces a non-failing `deprecated-config`
warning even when its values are false or zero. Remove the complete mapping;
semantic review uses `docs review-context`. Remove `docs.identifier_patterns`
as well; its values no longer override the protocol grammar.
