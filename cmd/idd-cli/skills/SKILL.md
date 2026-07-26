---
name: intent-driven-development
description: Use Intent-Driven Development with idd-cli to turn an approved user request into human-readable design, contract, specification, and testing documents; bind those records to source declarations; validate and repair traceability; and carry the work through tests, review, and commit
license: MIT
metadata:
  audience: agents
  workflow: development
  protected: true
---

# Intent-Driven Development (IDD)

> Protected Skill: change this file only with explicit user approval.

Use this workflow when a repository uses the four IDD documents and
`idd-cli`. Treat the approved request and confirmed decisions as the source of
truth. A dedicated `intent.md` is optional, and SPEC records do not contain an
Intent-source field.

## Responsibility boundary

The agent authors and reviews meaning. idd-cli owns deterministic mechanics.

The agent must:

- explain requirements, behavior, rationale, guarantees, boundaries, failure
  behavior, exclusions, acceptance evidence, and test oracles;
- keep the documents understandable to humans instead of reducing them to
  registries, placeholders, or summaries;
- reconcile conflicting or obsolete prose and decide whether implementation
  and tests actually satisfy the documented intent; and
- preserve useful information while migrating or repairing documents.

Use idd-cli to:

- create the canonical four-file scaffold and report unfinished slots;
- validate filenames, identity, record fields, references, lifecycle rules,
  package coverage, source annotations, and doc/code correspondence;
- produce path- and finding-specific repair instructions; and
- repair only safe structure with `docs fix`.

`.idd.yaml` configures scanning and validation. It does not own Component,
Contract, SPEC, TEST, or completion-marker records. A passing validation proves
structural consistency, not semantic truth or implementation correctness.

## Workflow

### 1. Record the approved work

For non-trivial work, create or update `.planning/` with:

- the approved request and decisions;
- observable success criteria;
- constraints and out-of-scope items;
- implementation phases and the current next action; and
- verification and migration risks.

Keep the plan current until every accepted item is complete. Do not add an
Intent-source field to SPEC records.

### 2. Establish a baseline

From the project root, inspect the current machine-readable state:

```bash
idd-cli docs status docs --format json
idd-cli run . --format json
```

Also inspect the real code, tests, configuration, documents, and working-tree
changes. Separate existing findings from findings introduced by the requested
work, and preserve unrelated user changes.

### 3. Initialize only new packages

For a new source package, run:

```bash
idd-cli docs init internal/auth --format json
```

The command creates `docs/internal/auth/{design,contract,spec,testing}.md` and
returns the initial `incomplete_slots` work list. Generated files are
intentionally incomplete until their guidance is replaced with
package-specific content.

Every scanned source package directory, including a nested sub-package, owns
an equally nested four-file document set. A child package never inherits its
parent package's records.

Do not initialize over an existing self-describing or legacy package. Use
`docs fix` only when a CLI finding explicitly calls for missing structure or
path-proven identity repair. It must not invent prose, coverage, or lifecycle
decisions.

### 4. Author documents in dependency order

Write the four files in this order:

1. `design.md`: define ownership, architecture, boundaries, dependencies,
   decisions, failure containment, and testability.
2. `contract.md`: define externally meaningful guarantees, inputs, outputs,
   errors, side effects, invariants, and compatibility.
3. `spec.md`: define required behavior and observable acceptance evidence,
   referring to the owning Components and Contracts.
4. `testing.md`: define evidence strategy, TEST records, concrete oracles,
   fixtures, isolation, and intentional exclusions.

Use ordinary Markdown prose, diagrams, examples, and small local tables where
they improve understanding. Do not turn the documents into large tables or
semantic YAML catalogs. When implementation exposes a missing decision, update
the owning document before relying on the code change.

### 5. Bind declarations to records

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

Do not add repeated `Spec:`, `Contract:`, or `Test:` document paths to source
headers. idd-cli derives code-to-document ownership from identifiers and the
owning package documents. Keep ordinary language package or module comments
when they help human readers.

### 6. Validate each coherent slice

Run focused completion inspection while authoring, followed by the project
gate:

```bash
idd-cli docs status docs/internal/auth --format json
idd-cli run . --format llm-markdown
idd-cli run . --format json
```

Treat `incomplete_slots` and each finding's `suggested_fix` as the
machine-owned work queue. Fix the earliest causal finding first; missing
records or declarations can produce several downstream findings.

Then run the repository's relevant focused tests and its complete test, race,
lint, build, and IDD gates.

### 7. Review and finish

Before committing, confirm that:

- every acceptance statement is observable;
- each Contract guarantee has named contract-test evidence;
- every TEST Oracle distinguishes pass from fail;
- relevant failure, concurrency, persistence, security, performance, and
  compatibility boundaries are explained;
- annotations bind to the intended declarations;
- prose agrees with the resulting code and tests;
- `.planning` contains no unfinished accepted item; and
- final validation is clean.

## Authoring contract

### Canonical package documents

The exact lowercase basename is the sole document-role authority:
`design.md`, `contract.md`, `spec.md`, and `testing.md`.

Each file begins with version and package identity:

```yaml
---
idd:
  version: "1.0"
  package: internal/auth
---
```

Do not add `idd.document`; the filename already supplies the role. Every source
package maps to the same project-relative path below `docs/`.

There is no document line limit. Keep rich explanation in the canonical file
instead of creating `design-*`, `contract-*`, `spec-*`, `testing-*`, or similar
split-role files. When such a file is reported, follow its `suggested_fix`,
merge all still-valid information into the canonical target, compare both
files, and only then remove the split source.

### Semantic responsibilities

Use record fields as traceability anchors, not as replacements for narrative:

- A Component explains why an implementation boundary exists, what it owns,
  what it delegates, how failures are contained, and how it can be tested.
- A Contract explains stable behavior visible across a boundary, including
  success, failure, invariants, side effects, and compatibility expectations.
- A SPEC explains required behavior, its design and contract context, and
  concrete evidence by which a reviewer can accept or reject it.
- A TEST explains the evidence it supplies, the scenario or boundary exercised,
  and the exact result that forms its oracle.

Address only concerns relevant to the package, but describe them fully enough
that a maintainer can understand the intended implementation boundary without
reconstructing it from code.

### Minimal record syntax

Use `docs init` as the template authority and `docs status` as the completion
authority. When adding a record to an existing document, preserve these
minimum shapes:

```markdown
## Component: AuthModule

**Purpose:** Explain the component's ownership and boundary.

## Contract: Authenticator

**Guarantees:** Explain externally visible success and failure behavior.

## SPEC-INTERNAL_AUTH-001: Authenticate credentials

- **Design:** `AuthModule`
- **Contract:** `Authenticator`

**Requirement:** State the required behavior.

**Acceptance:** State observable acceptance evidence.

## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_AUTH-001`

**Purpose:** State what evidence this test supplies.

**Oracle:** State the exact result that distinguishes pass from fail.
```

A contract TEST uses `Kind: contract` and adds `Contracts` naming the
guarantees it actually exercises. Do not assign every package Contract by
default. Let the current CLI schema and its findings govern optional fields,
allowed values, and conditional requirements.

### Relationship rules the author must understand

- Component and Contract names are package-local. Use the local name inside
  the package and `<package>#<name>` for a cross-package reference.
- `Covers` is the single authored TEST-to-SPEC relationship; idd-cli derives
  the reverse relationship.
- A contract TEST's `Contracts` field is the authored evidence link to Contract
  records.
- Dependencies, lifecycle links, and concern declarations must express real
  design decisions. Do not add values merely to silence validation.
- An annotation is implementation or test evidence, not a substitute for
  document prose or a declaration's normal explanatory comment.

## Work from CLI findings

Do not reproduce the validator's field catalogs or infer completion from
document length. Follow the current binary's output:

- `docs status` identifies each unfinished file, line, role, slot, and reason.
- A scaffold location is complete only after concrete content is written and
  its temporary marker is removed. Removing a marker alone does not pass.
- `run` reports reference, lifecycle, package-set, source-parse, annotation,
  coverage, and correspondence failures with a repair hint.
- `docs fix` may normalize safe identity or create missing skeletons. It never
  performs semantic merges, deletes split documents, or authors meaning.

Read an affected document and its canonical target completely before applying
a destructive repair. Preserve requirements, rationale, examples, diagrams,
boundaries, identifiers, and test evidence; do not replace detailed content
with a thin summary. Re-run the focused command from the finding and then the
full gate.

## Human review boundary

idd-cli can establish that expected files and fields exist, references resolve,
annotations attach to supported-language declarations, and the configured
traceability graph is internally consistent.

The agent must still decide whether:

- the requirement represents the approved product decision;
- the prose covers the relevant domain and operational risks;
- the design and contract describe coherent boundaries;
- implementation behavior fulfills each guarantee;
- acceptance evidence is meaningful; and
- tests, fixtures, and assertions provide a sufficient oracle.

Treat optional lexical-consistency warnings only as drift hints, never as
approval or a replacement for review.

## Existing-project upgrade

Upgrade package by package on a dedicated branch:

1. Build the intended idd-cli version, use its paired Skill, and record the old
   test and validation baseline.
2. Inventory legacy records, central catalogs, split-role files, repeated
   source-header paths, annotations, and current findings before editing.
3. Preserve rich prose while moving records into the owning canonical files.
   Add filename-derived identity, not role metadata or large YAML registries.
4. Give every nested source package its own nested four-file set. Merge split
   files with the finding-specific prompt and compare before deleting them.
5. Remove repeated source-header document paths only after confirming that
   declaration annotations resolve to the intended records.
6. Add concrete Purpose, Guarantees, Requirement, Acceptance, and Oracle
   content. Never insert placeholders merely to make a migration pass.
7. Resolve lifecycle, dependency, concern, contract-evidence, and cross-package
   relationships from their actual meaning rather than bulk-filling values.
8. Move annotations onto real declarations and resolve parse failures before
   trusting doc/code correspondence results.
9. Run `docs status`, `run`, and relevant tests after each package, then run
   the repository's complete gates before committing.

Do not mix legacy marker records with self-describing `idd` frontmatter inside
one package. Use the current CLI findings to determine the remaining mechanical
work and the documents themselves to preserve semantic depth.
