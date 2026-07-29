---
idd:
  version: "1.0"
  package: internal/model
---

# Design: internal/model

## Component: IDDModel

**Purpose:**

`IDDModel` defines the values exchanged between collectors, the linkage graph,
the engine, and reporters. It preserves identifier provenance and relationship
evidence during validation, accumulates findings without formatting them, and
provides stable JSON shapes for both complete and LLM-oriented reports.

The package is intentionally behavior-light and dependency-free. It owns data
semantics and deterministic collection helpers, but not Markdown parsing,
source scanning, graph construction, validation policy, or output rendering.
That separation prevents low-level values from importing the packages that
produce or consume them.

### Responsibilities

- represent the four identifier types and document/code origin;
- retain title, description, source location, raw reference, TEST kind, and
  legacy and explicitly typed forward references for each observed identifier;
- distinguish authored public IDs from package-scoped Component and Contract
  nodes derived by the documentation collector;
- store multiple observations with the same ID rather than overwriting
  evidence;
- provide unique, duplicate-preserving, type-specific, and origin-specific
  collection views;
- detect duplicate IDs that cross source directories;
- translate source annotations into identifier observations;
- model relationship direction and reverse relationship types;
- accumulate and sort validation errors and warnings;
- carry graph statistics and optional snapshots; and
- define finding-centered JSON report structures.

### Ownership and mutability

Most values are pointers to mutable structs. Constructors allocate link and
finding slices where required but do not deep-copy caller strings, maps, or
slices. `IdentifierSet` owns its private ID index and public type slices; adding
an identifier stores the pointer itself. `Merge` adds pointers from the other
set without cloning.

Several accessors return internal pointers or slices. The engine treats model
values as mutable during collection and graph construction, then read-mostly
during validation and reporting. There is no synchronization, so shared
mutation across goroutines is outside the design.

### Duplicate evidence rather than uniqueness enforcement

The ID index maps one string to a slice of observations. This is essential
because the same identifier should normally appear once in documentation and
once in code, and invalid projects may contain additional conflicting
declarations. `Get` and `All` choose the first observation for convenience;
`GetAll` and `AllIdentifiers` preserve evidence for validation.

Duplicate-group helpers only report IDs that appear in more than one source
directory for the same origin. Multiple observations in one directory are not
a package-name conflict. The helpers retain one representative per directory
and sort each group by source, while group order itself is unspecified.

### Result and report boundaries

`ValidationResult` starts invalid because no engine run has established
success. Adding an error keeps it invalid; adding a warning does not change the
flag. The engine owns when to set a successful result valid.

The complete `Report` shape contains tool/config/result context. The
finding-centered `LLMReport` is a separate wire contract produced by the
reporter. Model types define serialization only; they do not enrich, group, or
prioritize findings.

## Architecture

```text
collectors
   └─ Identifier / Annotation
           |
           v
      IdentifierSet ── duplicate/origin views
           |
           v
graph and engine ── ValidationError / ValidationResult / GraphSnapshot
           |
           v
reporter ── Report ──> LLMReport
```

Relationship values are shared with the graph. `IdentifierLink` lets
self-describing collectors state semantics that cannot be inferred from a
SPEC/TEST prefix, including Component dependency, named Contract test
evidence, and lifecycle direction. The graph adds adjacency and verification
state in its own types while snapshots return to model-owned serializable
summaries.

## Package Layout

`identifier.go` contains identifier vocabulary, observations, typed-link
ownership, collections, annotation conversion, validation results, snapshots,
and report wire shapes. `link.go` contains relationship vocabulary, the
explicit `IdentifierLink`, and the source-located graph `Link` value. They
share no external dependencies and remain in one package to avoid conversion
layers between every pipeline phase.

## Function Composition

Collectors construct annotations or document identifiers, add legacy or typed
links, mark package-scoped derived nodes when appropriate, set origin, and add
them to an `IdentifierSet`. Annotation conversion copies lexical evidence into
a new identifier and uses normalized declaration context as its description
when available; the code collector subsequently sets code origin.

The engine adds findings to `ValidationResult`, sorts them for stable output,
and optionally attaches graph statistics/snapshot. Reporter then projects those
values into complete or finding-centered output.

## Dependencies

Only standard-library formatting, path, sorting, and string packages are used.
JSON behavior is expressed through struct tags and executed by the reporter or
tests.

## Testability Hooks

Constructors and collection operations are directly testable with in-memory
values. Tests verify first-observation lookup, duplicate-preserving views,
origin filtering, untyped and typed link append behavior, annotation
conversion, result mutation and sorting, every directional relationship
reversal, and JSON round trips.

Direct package tests do not exhaustively cover cross-directory duplicate
grouping, aliasing through returned slices, exact order among duplicate IDs,
concurrent access, every report `omitempty` combination, or invalid field
combinations. Engine and reporter tests provide additional integration evidence
for those consumers.
