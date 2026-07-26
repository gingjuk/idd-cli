---
idd:
  version: "1.0"
  package: internal/engine
  document: spec
---

# Specifications: internal/engine

## SPEC-INTERNAL_ENGINE-001: Stateful repository validation orchestration

- **Design:** `ValidationEngine`
- **Contract:** `RepositoryValidation`

**Requirement:**

The engine must transform pre-collected identifier evidence and collector
structural errors into one complete repository validation result by detecting
cross-package duplicates, constructing the relationship graph, executing the
configured and compatibility rule set, and preserving precise evidence for
reporting.

### Required behavior

- Every document and code observation contributes to graph metadata without
  erasing the other origin.
- Code and document observations with the same SPEC or TEST identifier form the
  complete association; source files do not repeat document paths.
- A document record without source evidence and a source annotation without a
  document record produce `doc-code-correspondence`.
- Forward references become typed, source-located edges.
- Self-describing coverage, named Contract evidence, Component dependencies,
  and lifecycle replacements retain their explicit relationship types; only
  legacy links use type inference.
- Public, test, and annotation policy consumes normalized syntax-tree
  declarations for Go, TypeScript/TSX, JavaScript/JSX, C++, Java, and Python.
- A configured extension without a pinned grammar or invalid supported-source
  syntax yields a `source-parse` finding and never activates a regex fallback.
- Rule failures accumulate as structured errors or warnings.
- Self-describing collector findings survive engine execution.
- Legacy-only rules do not reinterpret self-describing records.
- Validity becomes true only when the final error collection is empty.
- Findings are sorted and graph statistics always describe the built graph.
- Optional graph output reflects post-verification state.

### Failure and implementation boundary

Validation findings are not returned as Go errors. Raw filesystem scan failures
are generally skipped, and rules report missing evidence only when their own
checks can observe it. The rule set is hard-coded; no production Rule plugin
interface, collector interface, or separate linker exists.

The engine is not a semantic prose judge. A structurally valid graph can still
contain incomplete or misleading human documentation and requires Skill-guided
review.

### Acceptance evidence

**Acceptance:** Focused tests exercise every major rule family, typed and
legacy graph interpretation, derived Contract coverage, Component dependency
cycles, seven-language public and test declarations, annotation attachment and
syntax failures, structural-error injection, duplicate diagnostics, exact
ignore scopes, identifier-derived code/document correspondence, headerless
production/test/contract-test files, warning thresholds, result sorting,
statistics, optional snapshots, and report construction.

## SPEC-INTERNAL_ENGINE-002: Initialized single-run engine

- **Design:** `ValidationEngine`
- **Contract:** `EngineLifecycle`

**Requirement:**

Construction must retain the supplied configuration and allocate a non-nil
empty graph and validation result ready to receive collector errors and one
validation run.

### Lifecycle boundary

Construction does not clone or validate configuration. The initial result is
invalid with empty finding slices. Engine state is not reset by `Run`; callers
must create a fresh engine for an independent run.

### Acceptance evidence

**Acceptance:**

Contract tests inspect non-nil fields, initial validity, and empty findings.

## SPEC-INTERNAL_ENGINE-004: Deterministic run and result lifecycle

- **Design:** `ValidationEngine`
- **Contract:** `EngineLifecycle`

**Requirement:**

`Run` must perform duplicate detection, typed graph construction, statistics,
validation, optional snapshotting, and return the engine-owned result with all
collector and rule findings preserved. When the caller supplies source analyses,
source rules must reuse that exact normalized evidence instead of reparsing or
applying language-blind line rules.

**Acceptance:** Lifecycle tests prove collector findings survive `Run`, typed
edges and metadata reach the graph, source analyses drive declaration rules,
validity reflects the final error set, output ordering is deterministic, and an
optional snapshot describes the same completed graph.

### Inputs, output, and context

The identifier set and configuration must be non-nil. The result pointer remains
owned by the engine and may be used by `BuildReport`. The accepted context is
currently not checked; cancellation does not abort work. Operational failure is
not currently returned, so callers inspect `result.Valid`, errors, and warnings.

### Ordering and side effects

Execution reads repository files according to configured patterns and mutates
the graph/result. Duplicate errors are added before graph validation. Statistics
are captured after graph construction; validation then sorts findings. An
optional snapshot is created after validation has calculated relationship
verification state.

### Non-goals

This specification does not promise safe engine reuse, concurrent execution,
fail-fast rules, file repair, output rendering, or cancellation. Adding those
capabilities requires explicit lifecycle and error contracts.
