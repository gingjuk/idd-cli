---
name: intent-driven-development
description: Apply Intent-Driven Development with idd-cli to author canonical package documents, bind SPEC and TEST records to declarations, validate traceability, review focused SPEC evidence, and finish approved code changes.
license: MIT
metadata:
  audience: agents
  workflow: development
  protected: true
---

# Intent-Driven Development (IDD)

> Protected Skill: change this file only with explicit user approval.

Use this workflow when a repository uses `idd-cli` and the canonical
`design.md`, `contract.md`, `spec.md`, and `testing.md` package documents.
Treat the approved request and confirmed decisions as the source of truth.

## Responsibility boundary

The agent owns meaning:

- explain behavior, rationale, boundaries, failures, exclusions, acceptance
  evidence, and test oracles;
- reconcile obsolete or conflicting prose with code and tests; and
- preserve useful semantics during repair or migration.

idd-cli owns deterministic mechanics:

- scaffold and inspect canonical documents;
- validate identity, fields, references, lifecycle-aware coverage, annotations, and documentation-unit
  correspondence;
- perform safe structural repair; and
- build one ID-centered Entity/Occurrence and Relation/Provenance index for
  trace, focused review, and change-aware impact queries.

A passing `idd-cli run` proves structural traceability, not semantic truth or
implementation correctness. `.idd.yaml` configures validation; it does not own
Component, Contract, SPEC, or TEST records.

## Workflow

### 1. Establish scope and baseline

For non-trivial work, record the approved behavior, success criteria,
constraints, current step, and verification risks in `.planning/`.
A dedicated `intent.md` is optional; do not add an Intent field to SPEC
records.

From the project root, inspect the current state:

```bash
idd-cli docs status docs --format json
idd-cli run . --format json
```

Also inspect the real code, tests, configuration, documents, and working-tree
changes. Separate existing findings from findings introduced by the work.

### 2. Initialize only new documentation units

```bash
idd-cli docs init internal/auth internal/config --format json
```

Without `docs.units`, each scanned source package owns the matching
`docs/<package>/` four-file set. An explicit unit can cover several source
package globs; every scanned source must match exactly one configured unit.

Do not initialize over an existing self-describing or legacy package.
Generated guidance is incomplete until replaced with package-specific content.
Use `docs fix` only for structural work identified by a finding; it must not
invent prose or coverage.

Batch related packages or paths to reduce tool calls, but keep unrelated
semantic reviews separate.

### 3. Author documents in dependency order

1. `design.md`: each Component's Purpose, Ownership, Boundary, and Decisions,
   plus concern-driven dependencies, failure containment, and testability.
2. `contract.md`: externally meaningful guarantees, inputs, outputs, errors,
   side effects, invariants, non-guarantees, known limitations, and compatibility.
3. `spec.md`: required behavior and observable acceptance evidence, linked to
   the owning Component and Contract.
4. `testing.md`: evidence strategy, TEST records, scenarios, fixtures,
   isolation, exclusions, and exact oracles.

Use ordinary Markdown prose. Record fields are traceability anchors, not a
substitute for explanation. Do not create large semantic YAML catalogs or add
values merely to silence validation.

### 4. Bind declarations

Use the annotation spellings configured in `.idd.yaml`. With the defaults:

```go
// @implement SPEC-INTERNAL_AUTH-001
func Authenticate(input LoginRequest) (LoginResponse, error) { ... }

// @test TEST-INTERNAL_AUTH-001
func TestAuthenticate(t *testing.T) { ... }

// @test-contract TEST-INTERNAL_AUTH-002
func TestAuthenticatorContract(t *testing.T) { ... }
```

Place annotations on the declarations they describe. One annotation may list
several same-kind identifiers separated by commas.

Validation proceeds from canonical records outward: every active SPEC needs an
implementation and TEST evidence, and every authored annotation must resolve.
Ordinary helper APIs and unrelated test functions do not require invented
records unless the optional public/test census rules are enabled.

Do not add repeated `Spec:`, `Contract:`, or `Test:` paths to source headers.
idd-cli associates declarations with documents through identifiers. An
annotation is evidence, not a replacement for a useful declaration comment.

### 5. Validate and review

Use focused status while authoring, then the project graph:

```bash
idd-cli docs status docs/internal/auth docs/internal/config --format json
idd-cli run . --format llm-markdown
idd-cli run . --format json
```

Treat `incomplete_slots` and `suggested_fix` values as the mechanical work
queue. Fix the earliest causal finding first. Read all affected content before
merging or deleting anything, preserve unique meaning, then rerun the focused
command and the full gate.

When meaning or evidence changed, request focused review input:

```bash
idd-cli docs review-context \
  SPEC-INTERNAL_AUTH-001 \
  SPEC-INTERNAL_CONFIG-001 \
  --format llm-markdown
```

Start from Git changes when the affected IDs are not known, or trace any ID
directly across docs, source, tests, and relation provenance:

```bash
idd-cli docs impacted --base main --format llm-markdown
idd-cli trace SPEC-INTERNAL_AUTH-001 --format llm-markdown
```

An impacted `related-spec-unchanged` warning requests review; it does not claim
that the canonical prose is stale. Trace and review-context reuse the same
project index and must not reinterpret mentions as authoritative relationships.

One request accepts up to ten SPEC IDs. Review each returned SPEC independently
against its Requirement, Acceptance, Contract, covering TEST records, and
implementation/test declarations. One SPEC may have several TESTs, and one
TEST may cover several SPECs.

Use review-context for changed, ambiguous, high-risk, or weakly evidenced
SPECs. Do not review every function merely to prove association: annotations
already establish that graph. The command supplies evidence; it does not invoke
an LLM, approve a SPEC, or participate in `run` validity.

### 6. Finish

Before committing, confirm that:

- prose matches the approved behavior and resulting implementation;
- every Acceptance is observable;
- Contract guarantees have contract-test evidence;
- every TEST Oracle distinguishes pass from fail;
- relevant failure, security, concurrency, persistence, performance, and
  compatibility boundaries are addressed;
- annotations bind to the intended declarations;
- `.planning/` has no unfinished accepted item; and
- focused tests plus the repository's build, full test/race, lint, and IDD
  gates pass.

## Authoring contract

### Canonical package documents

The exact lowercase basename is the document-role authority:
`design.md`, `contract.md`, `spec.md`, and `testing.md`.

Each file begins with:

```yaml
---
idd:
  version: "1.1"
  package: internal/auth
  namespace: INTERNAL_AUTH
---
```

Do not add `idd.document`. Package records the current document location;
namespace is the stable ID prefix and survives package moves. All four files in
a unit must agree on both values.
There is no line limit: keep rich content in the canonical file instead of
creating `design-*`, `contract-*`, `spec-*`, or `testing-*` fragments.

Use `docs init` as the template authority and `docs status` as the completion
authority. The minimal record shapes are:

```markdown
## Component: AuthModule

**Purpose:** Explain why this component exists.

**Ownership:** Explain the responsibilities and state it owns.

**Boundary:** Explain what remains outside and how collaboration crosses it.

**Decisions:** Explain the choices and trade-offs behind this shape.

## Contract: Authenticator

**Guarantees:** Explain observable success, failure, and invariants.

## SPEC-INTERNAL_AUTH-001: Authenticate credentials

- **Components:** `AuthModule`
- **Contracts:** `Authenticator`

**Requirement:** State the required behavior.

**Acceptance:** State evidence that can accept or reject the behavior.

## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_AUTH-001`

**Purpose:** State what evidence this test supplies.

**Oracle:** State the exact result that distinguishes pass from fail.
```

A contract TEST uses `Kind: contract` and lists only the Contract guarantees it
actually exercises in `Contracts`.

`Components` is required and may name several owners. `Contracts` is optional
and may name several stable boundaries; do not invent a Contract for internal
behavior. A Contract file with no records must explain that choice explicitly.

### Relationship rules

- Component and Contract names are package-local. Use `<package>#<name>` for a
  cross-package reference.
- TEST `Covers` is the single authored TEST-to-SPEC relationship; idd-cli
  derives the reverse edge. Do not add a TEST list to SPEC records or force
  one-to-one numbering.
- A contract TEST's `Contracts` field links evidence to Contract records.
- Dependencies, lifecycle links, and concerns must represent real decisions.
- `planned` SPECs preserve approved intent without implementation or evidence;
  active SPECs require both. Deprecated and superseded records retain history
  without delivery correspondence requirements.
- Required fields and allowed values come from the current `docs init`,
  `docs status`, and finding output; do not copy validator catalogs into prose.

## Repair and migration rules

- `docs fix` may normalize identity or create missing skeletons. It does not
  merge prose, delete fragments, or author meaning.
- For a split role file, read both files, preserve all still-valid unique
  content in the canonical target, compare them, and only then remove the
  fragment.
- Do not mix legacy markers or `related_files` metadata with self-describing
  `idd` frontmatter.
- Migrate package by package using the intended idd-cli version and its paired
  Skill. Preserve rich prose, then repair references, annotations, coverage,
  lifecycle, and package ownership from their actual meaning.
- Remove repeated source-header document paths only after identifier
  annotations resolve to the intended records.
- If `validation.consistency_check` is configured, remove the complete mapping.
  It is a deprecated, ignored compatibility key; use `docs review-context` for
  semantic review instead of carrying forward TF-IDF thresholds.
- Remove obsolete `docs.identifier_patterns`; the versioned IDD protocol owns
  ID syntax. Moving a package changes `idd.package`, not its stable namespace.

Run `docs status`, `run`, and relevant tests after each migrated package, then
run the repository's complete gates before committing.
