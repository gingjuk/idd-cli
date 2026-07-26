---
name: intent-driven-development
description: Use Intent-Driven Development with idd-cli to turn an approved user request into human-readable design, contract, specification, and testing documents; bind those records to source declarations in supported languages; validate and repair traceability; and carry the work through tests, review, and commit
license: MIT
metadata:
  audience: agents
  workflow: development
  protected: true
---

# Intent-Driven Development (IDD)

> Protected Skill: change this file only with explicit user approval.

Use this workflow when a repository uses the four IDD documents and
`idd-cli`. Treat the approved user request and confirmed decisions as the
source of truth. A dedicated `intent.md` is optional: idd-cli neither requires
one nor builds Intent-to-SPEC relationships.

## Skill and CLI responsibilities

The agent authors and reviews meaning. idd-cli checks deterministic evidence.

- The agent explains requirements, behavior, rationale, boundaries,
  guarantees, failures, acceptance evidence, test oracles, and exclusions.
- idd-cli creates visible scaffolds, parses fixed Markdown records, validates
  references and lifecycle rules, binds source annotations through syntax
  trees, builds the linkage graph, and reports exact findings.
- `.idd.yaml` configures scanning and rules. It never owns Component,
  Contract, SPEC, TEST, or completion-marker records.
- `docs fix` repairs only path-proven identity and missing structure. It must
  not invent prose, coverage, or lifecycle decisions.

A passing command proves structural consistency, not that prose is true,
complete for the domain, or implemented correctly. Review those questions
explicitly.

## End-to-end workflow

### 1. Capture and track the approved work

For non-trivial work, create or update `.planning/` with:

- the user request and confirmed decisions;
- observable success criteria;
- constraints and out-of-scope items;
- implementation phases and one current next action;
- verification and migration risks.

Do not add an optional Intent source field to SPEC records. Keep plan status
current until no required item remains.

### 2. Establish the repository baseline

From the project root:

```bash
idd-cli docs status docs --format json
idd-cli run . --format json
```

Also inspect the real code, tests, configuration, existing documents, and
working-tree changes. Preserve unrelated user changes. Separate existing
findings from findings introduced by the requested work.

### 3. Bootstrap only a new package

```bash
idd-cli docs init internal/auth --format json
```

This creates `docs/internal/auth/{design,contract,spec,testing}.md` and returns
the initial `incomplete_slots` work list. Generated files are intentionally
invalid until authored.

Create one independent four-file set for every scanned source package
directory, including nested sub-packages. For example,
`internal/auth/token` owns `docs/internal/auth/token/`; it does not inherit the
documents of `internal/auth`.

Do not run `docs init` over existing self-describing or legacy records.
For an existing package, migrate deliberately or use:

```bash
idd-cli docs fix docs/internal/auth
```

`docs fix` preserves the Markdown body and every scaffold marker. It may
normalize minimal identity or create a missing skeleton, but it cannot complete
semantic work.

### 4. Author the four documents in dependency order

Use ordinary Markdown prose, lists, diagrams, small comparison tables, and
examples. Do not turn the documents into large registry tables or semantic
YAML catalogs.

1. `design.md`: ownership, architecture, boundaries, dependencies, decisions,
   failure containment, and testability.
2. `contract.md`: externally meaningful guarantees, inputs, outputs, errors,
   side effects, invariants, and compatibility.
3. `spec.md`: required behavior and acceptance evidence, referring to named
   Components and Contracts.
4. `testing.md`: evidence strategy, TEST records, concrete oracles, fixtures,
   isolation, and exclusions.

When implementation reveals a missing decision, update the owning document
before relying on the code change.

### 5. Bind code and tests

Add declaration comments:

```go
// @implement SPEC-INTERNAL_AUTH-001
func Authenticate(input LoginRequest) (LoginResponse, error) { ... }

// @test TEST-INTERNAL_AUTH-001
func TestAuthenticate(t *testing.T) { ... }

// @test-contract TEST-INTERNAL_AUTH-002
func TestAuthenticatorContract(t *testing.T) { ... }
```

One annotation may reference several same-kind IDs separated by commas.
Annotations are evidence, not substitutes for declaration comments or document
prose.

Do not repeat package document paths in source-file headers. Lines such as
`Spec: docs/<package>/spec.md`, `Contract: ...`, and `Test: ...` are neither
required nor consumed. idd-cli owns the association by joining each
syntax-tree-bound identifier to its self-describing SPEC or TEST record.
Ordinary language package/module comments remain human-owned and may be kept
when they are useful.

### 6. Validate after each coherent slice

```bash
idd-cli docs status docs/internal/auth --format json
idd-cli run . --format llm-markdown
idd-cli run . --format json
```

Fix the earliest causal finding first. A missing declaration may cause several
coverage and orphan findings; do not patch each symptom independently.
Re-run focused tests, then the repository's complete test, race, lint, build,
and IDD gates.

### 7. Review and finish

Confirm:

- each acceptance statement is observable;
- every Contract guarantee has named contract-test evidence;
- every TEST Oracle says what result distinguishes pass from fail;
- boundaries, failures, concurrency, persistence, security, performance, and
  compatibility are addressed where relevant;
- code annotations bind to the intended declarations;
- `.planning` has no unfinished accepted item;
- final validation is clean before commit.

## Canonical document model

The lowercase basename is the sole authority for document role:

| Filename | Role |
| --- | --- |
| `design.md` | design |
| `contract.md` | contract |
| `spec.md` | specification |
| `testing.md` | testing |

Each file begins with minimal version and package identity only:

```yaml
---
idd:
  version: "1.0"
  package: internal/auth
---
```

Do not add an `idd.document` field: it is invalid and idd-cli does not provide
a compatibility mode for it. Do not infer the role from the H1 or prose.
Semantic records live in Markdown.

Every source package directory maps to the same project-relative path below
`docs/` and owns exactly these four role files. Nested source packages own
nested four-file sets. Documents have no line-count limit. Keep rich prose,
diagrams, examples, subordinate headings, and small local tables in the owning
file; never create `design-*`, `contract-*`, `spec-*`, or `testing-*` files to
split a role by size or feature.

### Repair a split-role finding

The default JSON report exposes the executable instruction in each
`idd-document-filename` finding's `suggested_fix`; LLM Markdown prints the same
instruction under `Fix`. Follow the per-finding prompt rather than applying
the group summary to only one file.

For each split file:

1. read the split source completely and read its canonical target if present;
2. if the target is missing, use `idd-cli docs fix <docs-package>` only to
   create the structural four-file set;
3. merge every unique, still-valid requirement, behavior description, design
   rationale, contract, implementation boundary, identifier relationship,
   example, diagram, and test-evidence note into the canonical target;
4. reconcile duplicate headings while preserving human-readable Markdown and
   semantic depth—do not replace detailed text with a summary and do not create
   another split file;
5. remove the split source only after comparing both files and verifying that
   no information was lost; and
6. run the exact `docs status` command from the prompt, then
   `idd-cli run . --format json`, until no related finding or incomplete slot
   remains.

`docs status` returns this repair prompt as an operational error when the
requested file or tree contains a split role document. The CLI never performs
the merge or deletion itself because both operations require semantic
judgment.

### design.md

```markdown
## Component: AuthModule

- **Status:** `active`
- **Depends on:** `CredentialStore`, `internal/audit#AuditLog`
- **Concerns:** `security`, `concurrency`

**Purpose:** Own credential verification and isolate it from transport code.

### Security

Plain-text credentials never leave the call boundary.

### Concurrency

The component has no mutable process-global authentication state.
```

`Purpose` is required. `Depends on`, lifecycle fields, and `Concerns` are
optional. The document also requires non-empty `## Architecture`,
`## Package Layout`, `## Function Composition`, `## Dependencies`, and
`## Testability Hooks` sections.

### contract.md

```markdown
## Contract: Authenticator

- **Status:** `active`
- **Concerns:** `compatibility`

**Guarantees:** Valid credentials return a token; rejected credentials return
one stable public error without revealing which field failed.

### Compatibility

Existing callers retain the same input and public error categories.
```

`Guarantees` is required. Keep signatures and concise comparison tables near
the prose they clarify; the prose remains authoritative.

### spec.md

```markdown
## SPEC-INTERNAL_AUTH-001: Authenticate credentials

- **Design:** `AuthModule`
- **Contract:** `Authenticator`
- **Status:** `active`

**Requirement:** Authenticate valid credentials without exposing rejection
details.

**Acceptance:** A valid fixture yields a non-empty token; every invalid fixture
yields the documented public rejection and no credential detail.
```

`Design`, `Contract`, `Requirement`, and `Acceptance` are required. There is no
required or optional Intent-source field in the idd-cli schema.

### testing.md

```markdown
## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_AUTH-001`

**Purpose:** Exercise accepted and rejected credential scenarios.

**Oracle:** The returned token or public rejection exactly matches each
scenario row.

## TEST-INTERNAL_AUTH-002: Authenticator boundary

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_AUTH-001`
- **Contracts:** `Authenticator`

**Purpose:** Protect the externally visible authentication boundary.

**Oracle:** All documented success, rejection, and redaction guarantees hold
for every implementation under test.
```

`Kind`, `Covers`, `Purpose`, and `Oracle` are required. `Kind: contract` also
requires one or more named `Contracts`; a behavior TEST must not declare that
field.

## Relationship rules

- Component and Contract names are package-local human names. idd-cli derives
  internal graph keys; authors do not create `DESIGN-*` or `CONTRACT-*`
  markers for them. Names must not contain `#` or `,`, because those characters
  delimit package-qualified and list references.
- A local reference uses the name. A cross-package Component or Contract
  reference uses `<package>#<name>`.
- TEST `Covers` is the single authored TEST-to-SPEC relationship. The reverse
  graph relationship is derived.
- Contract TEST `Contracts` is the single authored evidence link to named
  Contract records.
- `Depends on` targets must resolve and the Component dependency graph must be
  acyclic.

Optional lifecycle values are:

- `Status`: `active`, `deprecated`, or `superseded`;
- `Supersedes`: same-kind records replaced by this record;
- `Deprecated by`: the same-kind replacement for this record.

Either side may establish a replacement, avoiding mandatory edits to both old
and new records. If both sides are authored they must agree. A `superseded`
record must have exactly one derived replacement, and replacement relationships
must not self-link or form cycles.

`Concerns` accepts `security`, `concurrency`, `persistence`, `performance`, and
`compatibility`. Every selected concern requires a non-empty same-named
level-three section inside that record. Do not infer concern tags from prose.

## Scaffold completion protocol

Generated fill locations contain stable temporary markers:

```markdown
<!-- idd:scaffold slot="design.architecture" -->
```

An agent knows what remains by reading `docs status`, not by guessing from
length or keywords. Each `incomplete_slots` entry contains file, line, role,
slot, and reason.

A slot is complete only when:

1. its marker is removed; and
2. its bounded canonical record or section contains effective,
   non-placeholder authored content.

Deleting a marker alone does not pass. Headings, HTML comments, generated
guidance blockquotes, empty containers, `TBD`, `TODO`, and `auto-generated` do
not count as authored content. Once a record heading exists, `docs status`
also returns a record-field slot for every missing required field, including
conditional contract TEST `Contracts`. Remove a marker only after filling its
location, and continue until both `docs status` and `run` pass.

## Syntax-tree source binding

idd-cli parses configured source files with pinned Tree-sitter grammars:

- Go: `.go`
- TypeScript and TSX: `.ts`, `.tsx`
- JavaScript and JSX: `.js`, `.jsx`
- C++: `.cpp`, `.cc`, `.cxx`, `.hpp`, `.hh`, `.hxx`
- Java: `.java`
- Python: `.py`

The parser recognizes real comment nodes, declarations, visibility, and test
conventions. It rejects annotation-looking strings, executable-body comments,
and detached annotations. Configured extensions without a pinned grammar and
syntax errors in supported files produce `source-parse`; idd-cli does not
silently fall back to line matching.

The spellings below are defaults. Read `code.annotations` in `.idd.yaml`; when
a repository configures different `spec`, `test`, or `test_contract` prefixes,
use those exact spellings. idd-cli normalizes them back to the same three
semantic kinds.

Public production declarations require `@implement`. Test files are excluded
from the public-API rule, while real test declarations require `@test` or
`@test-contract`. Contract-test-named files require the contract form; a normal
test file may still use `@test-contract` when the documented TEST kind is
contract.

Code-to-document association is identifier-derived for every supported
language. A documented identifier without a source annotation and a source
annotation without a document record both produce `doc-code-correspondence`.
Do not add a second file-level path mechanism.

Ignore directives must be standalone comment directives:

```text
idd:ignore start
idd:ignore end
```

A prose mention of `idd:ignore` does not activate an ignore range.

The official grammar bindings require CGO. Use the repository's CGO-enabled
build and test commands rather than treating a CGO-disabled compile failure as
a document defect.

## Deterministic checks and review boundary

idd-cli can check:

- document identity, required records and fields, placeholder state, and exact
  locations;
- canonical role filenames and one complete four-file set for every scanned
  source package, including nested packages;
- reference resolution, duplicate names/IDs, lifecycle agreement and cycles,
  conditional concern sections, dependency cycles, coverage, and orphans;
- source parseability, annotation syntax, declaration attachment, visibility,
  test kind, and doc/code correspondence;
- package document sets and configured legacy compatibility rules.

idd-cli cannot prove:

- that a requirement is the right product decision;
- that prose covers every domain risk;
- that a guarantee is implemented correctly;
- that assertions or test data form a sufficient oracle;
- that a lexical similarity score represents semantic consistency.

The optional TF-IDF consistency warning is disabled by default. Enable it only
as a drift hint; never use it as approval or as a replacement for review.

## Existing-project upgrade

Upgrade package by package on a dedicated branch:

1. Build the intended idd-cli version and export its paired Skill.
2. Inventory legacy markers, central catalogs, duplicate relationships,
   repeated source-file document paths, source annotations, and current
   validation findings. Remove `Spec:`, `Contract:`, and `Test:` header paths
   after confirming the declaration annotations resolve; preserve useful
   language package/module descriptions. Remove the obsolete
   `require_package_doc_comment` configuration key; code/document pairing is
   controlled by `require_doc_code_correspondence`.
3. Preserve rich prose. Add minimal identity and move declarations into the
   owning H2 records; delete any obsolete `idd.document` field immediately and
   do not replace prose with YAML or a large table.
   Map every source sub-package to its own equally nested `docs/<package>/`
   directory. For every `design-*`, `contract-*`, `spec-*`, or `testing-*`
   fragment, follow the finding's path-specific repair prompt and compare the
   merged result before deleting the fragment; document length is not a
   validation constraint.
4. Add concrete `Purpose`, `Guarantees`, `Acceptance`, and `Oracle` labels to
   existing evidence. Do not insert placeholders.
5. Derive each contract TEST's `Contracts` from the guarantees it actually
   exercises; do not assign every package Contract indiscriminately.
6. Normalize lifecycle and dependency names, then resolve cross-package
   targets with package-qualified references.
7. Move annotations onto real declarations and fix syntax errors before
   trusting correspondence results.
8. Use `docs status` and `run` after each package. Use `docs fix` only for
   structural identity; verify it is idempotent.
9. Run the old test baseline plus the full new gates, review the final graph,
   and commit only after the plan and findings are closed.

Legacy packages may remain on the legacy parser until explicitly migrated.
Never mix legacy `markers` or `related_files` with self-describing `idd`
frontmatter in one package.

## Operational commands

```bash
# Export the workflow embedded in this binary
idd-cli generate skill --output idd-skill.md

# Create new-package scaffolds and return the work list
idd-cli docs init internal/auth --format json

# Inspect completion without modifying files
idd-cli docs status docs/internal/auth --format json

# Normalize safe version/package identity while preserving body and markers
idd-cli docs fix docs/internal/auth

# Run the authoritative project gate
idd-cli run . --format llm-markdown
idd-cli run . --format json --output idd-report.json
```

Run from project root unless the repository documents another scope. An
invalid graph still writes the requested report and exits non-zero.
