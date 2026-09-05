---
idd:
  version: "1.1"
  package: internal/model
  namespace: INTERNAL_MODEL
---

# Specifications: internal/model

## SPEC-INTERNAL_MODEL-001: Identifier type vocabulary and parsing

- **Components:** `IDDModel`
- **Contracts:** `IdentifierVocabulary`

**Requirement:**

The model must represent SPEC, CONTRACT, TEST, and DESIGN kinds and parse their
names case-insensitively while rejecting unknown or empty values.

### Acceptance evidence

**Acceptance:**

Table-driven cases cover lowercase, uppercase, and mixed-case SPEC plus every
other supported kind, an unknown name, and empty input.

## SPEC-INTERNAL_MODEL-002: Initialized identifier construction

- **Components:** `IDDModel`
- **Contracts:** `IdentifierRecord`

**Requirement:**

Identifier constructors must preserve supplied identity, title, optional
description, source, and line while initializing raw reference, empty links, and
document origin.

**Acceptance:** Constructor tests compare the supplied identity, title,
description, source, and line with the resulting record. They also require the
raw reference to be initialized, both link collections to be non-nil, and the
initial origin to be documentation.

### Boundary

Construction does not validate grammar or source paths. The describe-aware and
basic forms differ only in whether `Describe` is populated.

## SPEC-INTERNAL_MODEL-003: Append-only forward references

- **Components:** `IDDModel`
- **Contracts:** `IdentifierRecord`

**Requirement:** An identifier must retain both legacy untyped targets and
explicitly typed relationship targets in insertion order, allowing repeated
targets without validation or deduplication.

### Acceptance evidence

**Acceptance:** Tests append untyped and typed targets, then observe the exact
target values, relationship semantics, and insertion order without implicit
normalization.

## SPEC-INTERNAL_MODEL-004: Duplicate-preserving identifier collection

- **Components:** `IDDModel`
- **Contracts:** `IdentifierCollection`

**Requirement:**

The collection must preserve every observation, maintain type and ID indexes,
filter by origin, and distinguish same-directory repetition from IDs duplicated
across packages.

**Acceptance:** After adding document and code observations, lookup by ID and
origin returns only the matching values; wrong-origin and unknown-ID queries
are empty. Duplicate analysis retains separate observations, ignores repeated
files from one directory as a package conflict, and groups occurrences only
when the same ID spans directories.

### Implementation boundary

Adding values performs no uniqueness or nil check. Cross-directory duplicate
groups retain one source representative per directory; callers decide whether
that evidence is an error.

## SPEC-INTERNAL_MODEL-005: First and complete ID lookup

- **Components:** `IDDModel`
- **Contracts:** `IdentifierCollection`

**Requirement:**

Callers must be able to test presence, retrieve the first observation for an
ID, and retrieve all observations without losing doc/code duplicates.

**Acceptance:** Lookup tests distinguish an inserted ID from an unknown ID,
return the first inserted observation through the compatibility lookup, and
return both observations through the exhaustive lookup when one ID has
document and code evidence.

### Ownership

Returned pointers and the all-observations slice reference collection-owned
values and are not defensive copies.

## SPEC-INTERNAL_MODEL-006: Deterministic unique and exhaustive views

- **Components:** `IDDModel`
- **Contracts:** `IdentifierCollection`

**Requirement:**

The collection must report unique-ID count, return one first observation per ID
for ordinary iteration, and return every observation for duplicate analysis.
Both iteration forms are sorted by identifier string.

**Acceptance:** A new set reports zero unique IDs, and repeated observations of
one ID do not increase that count. Ordinary iteration returns one entry per ID,
exhaustive iteration returns every duplicate observation, and both results are
ordered by identifier.

### Edge cases

Relative order among observations sharing an ID is not a stable secondary sort
contract.

## SPEC-INTERNAL_MODEL-007: Lossless set merge

- **Components:** `IDDModel`
- **Contracts:** `IdentifierCollection`

**Requirement:**

Merging must add every observation from the source set through the destination's
normal indexes without cloning or collapsing shared IDs.

**Acceptance:** After merging independently populated sets, destination lookup
retrieves every source observation. When both sets contain the same ID,
exhaustive lookup returns both observations rather than an overwritten value.

## SPEC-INTERNAL_MODEL-008: Source annotation evidence

- **Components:** `IDDModel`
- **Contracts:** `SourceAnnotation`

**Requirement:**

Parsed source annotations must retain kind, reference, raw text, source
location, lexical context, and optional declaration comment so later stages can
diagnose syntax and compare descriptions.

**Acceptance:** A constructed annotation exposes the supplied kind, reference,
raw text, file, line, and lexical context unchanged. Supplying a declaration
comment preserves it; omitting the comment leaves that evidence empty.

### Boundary

Constructors accept lexical evidence as supplied and do not determine placement
or declaration validity.

## SPEC-INTERNAL_MODEL-009: Annotation-to-identifier projection

- **Components:** `IDDModel`
- **Contracts:** `SourceAnnotation`

**Requirement:**

Annotation conversion must create an identifier with the annotation identity,
type, source, and line, preserve full raw annotation text, and use a non-empty
declaration comment as the identifier description.

**Acceptance:** Projection tests compare the resulting identifier's ID, type,
source, line, and raw reference with its annotation. A non-empty declaration
comment becomes the description, while absent comment evidence leaves the
description empty.

### Failure and ownership boundary

The projected identifier initially has document origin because it uses the
general constructor. The code collector must set code origin explicitly.
Lexical context is not copied.

## SPEC-INTERNAL_MODEL-010: Validation and complete-report value shapes

- **Components:** `IDDModel`
- **Contracts:** `ValidationResult`

**Requirement:**

The model must carry structured findings, aggregate validity and statistics,
optional graph snapshots, tool/config metadata, and stable JSON field names
without embedding reporter behavior.

**Acceptance:** A newly constructed result is invalid with allocated empty
error and warning slices. Populated values preserve structured rule, message,
and location evidence, and JSON round trips retain the documented report
schema, summary, statistics, optional graph data, and nested findings.

### Invariants

New results have allocated empty finding slices and begin invalid until the
engine establishes success. Error string formatting exposes rule and message
while retaining structured evidence separately.

## SPEC-INTERNAL_MODEL-011: Finding accumulation and stable rule ordering

- **Components:** `IDDModel`
- **Contracts:** `ValidationResult`

**Requirement:**

Errors and warnings must append without losing evidence; errors must invalidate
the result; and both collections must be sortable by rule then message for
stable reports.

**Acceptance:** Adding an error increases only the error collection and leaves
the result invalid; adding a warning increases only the warning collection.
Multi-finding tests sort distinct rules into deterministic ascending order
without dropping their messages or locations.

### Non-goals

The model does not deduplicate findings, assign severity from a rule, or define
a source-based tie breaker.

## SPEC-INTERNAL_MODEL-012: Relationship vocabulary and reversal

- **Components:** `IDDModel`
- **Contracts:** `Relationship`

**Requirement:** The model must expose the relationship strings consumed by
graph validation, including Component dependency and lifecycle replacement,
and map every directional type to its correct reverse while leaving symmetric
and unknown values unchanged.

### Acceptance evidence

**Acceptance:** Constant and table-driven reversal tests cover
tests/implements, contract directions, dependency directions, lifecycle
directions, self-reversing contract-test/reference/annotation links, and an
unknown value.

## SPEC-INTERNAL_MODEL-013: Source-located directed link

- **Components:** `IDDModel`
- **Contracts:** `Relationship`

**Requirement:**

A relationship value must retain directed endpoints, relationship type, and
the source file and line that provided the evidence.

**Acceptance:** The directed-link test constructs one relationship and compares
both endpoint IDs, relationship type, source path, and line with the supplied
values.

## SPEC-INTERNAL_MODEL-014: Document and code provenance

- **Components:** `IDDModel`
- **Contracts:** `IdentifierVocabulary`

**Requirement:**

Every identifier observation must be classifiable as documentation or code so
the engine can require correspondence while retaining both copies of one ID.

### Acceptance evidence

**Acceptance:**

Origin filtering and presence tests distinguish document and code observations.

## SPEC-INTERNAL_MODEL-015: Link construction without graph policy

- **Components:** `IDDModel`
- **Contracts:** `Relationship`

**Requirement:**

Link construction must copy endpoints, type, and location into a new non-nil
value without checking endpoint existence or reverse-link completeness.

**Acceptance:** Construction returns a non-nil link whose endpoints, type,
source, and line exactly match the inputs even though no graph nodes or reverse
relationship were supplied.

## SPEC-INTERNAL_MODEL-017: Explicit origin reassignment

- **Components:** `IDDModel`
- **Contracts:** `IdentifierRecord`

**Requirement:**

Collectors must be able to replace an identifier's default origin after
construction so annotation-derived observations can be marked as code.

### Acceptance evidence

**Acceptance:**

A direct regression test constructs a default document-origin identifier,
assigns code origin through the method, and observes the new value.

## SPEC-INTERNAL_MODEL-018: Complete identifier evidence record

- **Components:** `IDDModel`
- **Contracts:** `IdentifierRecord`

**Requirement:**

The core identifier value must retain all evidence needed for graph creation,
diagnostics, TEST-kind validation, derived Component and
Contract nodes, typed relationships, and source-location reporting without
depending on collector-specific types.

**Acceptance:** Constructors allocate both link collections and preserve
ordinary identity evidence. Tests can then add typed relationships and mark a
derived node without importing collector or graph types, and graph tests can
materialize those exact semantics.

### Non-goals

The record does not enforce field combinations or own reverse backlinks. Its
mutable fields represent facts accumulated by pipeline stages.

## SPEC-INTERNAL_MODEL-037: Finding-centered LLM wire model

- **Components:** `IDDModel`
- **Contracts:** `LLMFindingReport`

**Requirement:**

The model must serialize an LLM-oriented report containing schema, status,
summary counts, grouped repeated rules, indexed findings, repair guidance,
structured locations, and related identifiers.

### Compatibility boundary

Field names and nesting are a machine-consumed wire surface. Reporter decides
content and ordering; changing JSON tags or required nesting requires output
contract tests and consumer review.

### Acceptance evidence

**Acceptance:**

A populated report round-trips through JSON while preserving schema, findings,
rule groups, and related identifier data.
