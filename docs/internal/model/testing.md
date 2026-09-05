---
idd:
  version: "1.1"
  package: internal/model
  namespace: INTERNAL_MODEL
---

# Testing: internal/model

## TEST-INTERNAL_MODEL-001: Identifier set add, get, and presence

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-004`, `SPEC-INTERNAL_MODEL-005`

**Purpose:**

Prove that adding an identifier updates presence and first-observation lookup
while an unknown ID remains absent.

**Oracle:** The test passes only when its assertions confirm adding
an identifier updates presence and first-observation lookup while an unknown ID remains
absent.

## TEST-INTERNAL_MODEL-002: Unique identifier count

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-006`

**Purpose:**

Prove zero count for a new set and one count per distinct ID after additions.

**Oracle:** The test passes only when its assertions confirm zero count
for a new set and one count per distinct ID after additions.

## TEST-INTERNAL_MODEL-003: Unique sorted view cardinality

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-006`

**Purpose:**

Prove that ordinary iteration returns one entry for each distinct inserted ID.

**Oracle:** The test passes only when its assertions confirm
ordinary iteration returns one entry for each distinct inserted ID.

## TEST-INTERNAL_MODEL-004: Set merge

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-007`

**Purpose:**

Prove that merging two sets makes both source observations available in the
destination.

**Oracle:** The test passes only when its assertions confirm merging
two sets makes both source observations available in the destination.

## TEST-INTERNAL_MODEL-005: Forward link append

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-003`

**Purpose:**

Prove that two appended relationship targets remain in the identifier's link
slice.

**Oracle:** The test passes only when its assertions confirm two
appended relationship targets remain in the identifier's link slice.

## TEST-INTERNAL_MODEL-006: Basic annotation projection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-009`

**Purpose:**

Prove that annotation identity, type, source, and raw text survive conversion to
an identifier.

**Oracle:** The test passes only when its assertions confirm
annotation identity, type, source, and raw text survive conversion to an identifier.

## TEST-INTERNAL_MODEL-007: Validation error string

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-010`

**Purpose:**

Prove the stable `[rule] message` error representation independently of
structured location fields.

**Oracle:** The test passes only when its assertions confirm the stable
`[rule] message` error representation independently of structured location fields.

## TEST-INTERNAL_MODEL-008: Error accumulation invalidates

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-011`

**Purpose:**

Prove that adding an error preserves one finding and leaves validity false.

**Oracle:** The test passes only when its assertions confirm adding
an error preserves one finding and leaves validity false.

## TEST-INTERNAL_MODEL-009: Warning accumulation

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-011`

**Purpose:**

Prove that warnings append without being treated as errors or changing the
already-false initial validity state.

**Oracle:** The test passes only when its assertions confirm
warnings append without being treated as errors or changing the already-false initial
validity state.

## TEST-INTERNAL_MODEL-010: Finding sort by rule

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-011`

**Purpose:**

Prove that sorting places lower rule names first when several errors are
present.

**Oracle:** The test passes only when its assertions confirm sorting
places lower rule names first when several errors are present.

## TEST-INTERNAL_MODEL-011: Annotation comment becomes description

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-009`

**Purpose:**

Prove that a non-empty function comment is copied into the projected
identifier's description.

**Oracle:** The test passes only when its assertions confirm a
non-empty function comment is copied into the projected identifier's description.

## TEST-INTERNAL_MODEL-012: Missing annotation comment stays empty

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-009`

**Purpose:**

Prove that conversion does not invent a description when no declaration comment
was collected.

**Oracle:** The test passes only when its assertions confirm
conversion does not invent a description when no declaration comment was collected.

## TEST-INTERNAL_MODEL-013: All observations for one ID

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_MODEL-005`
- **Contracts:** `IdentifierCollection`

**Purpose:**

Prove duplicate-preserving lookup for two observations and an empty result for
an unknown ID.

**Oracle:** The test passes only when its assertions prove
duplicate-preserving lookup for two observations and an empty result for an unknown ID.

## TEST-INTERNAL_MODEL-014: Exhaustive observation view

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-006`

**Purpose:**

Prove that exhaustive iteration retains duplicate observations instead of
collapsing them by ID.

**Oracle:** The test passes only when its assertions confirm
exhaustive iteration retains duplicate observations instead of collapsing them by ID.

## TEST-INTERNAL_MODEL-015: Filtering by provenance

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-004`, `SPEC-INTERNAL_MODEL-014`

**Purpose:**

Prove that document and code observations are returned only by their matching
origin filter.

**Oracle:** The test passes only when its assertions confirm
document and code observations are returned only by their matching origin filter.

## TEST-INTERNAL_MODEL-016: Origin presence query

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-004`, `SPEC-INTERNAL_MODEL-014`

**Purpose:**

Prove positive, wrong-origin, and unknown-ID outcomes for origin presence.

**Oracle:** The test passes only when its assertions confirm positive,
wrong-origin, and unknown-ID outcomes for origin presence.

## TEST-INTERNAL_MODEL-017: First-observation lookup

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-005`

**Purpose:**

Prove successful lookup of an inserted value and the false result for an
unknown ID.

**Oracle:** The test passes only when its assertions confirm successful
lookup of an inserted value and the false result for an unknown ID.

## TEST-INTERNAL_MODEL-018: Multi-rule result ordering

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_MODEL-011`
- **Contracts:** `ValidationResult`

**Purpose:**

Prove deterministic ascending rule order across three distinct validation
rules.

**Oracle:** The test passes only when its assertions prove
deterministic ascending rule order across three distinct validation rules.

## TEST-INTERNAL_MODEL-019: Identifier type parsing

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_MODEL-001`
- **Contracts:** `IdentifierVocabulary`

**Purpose:**

Prove case-insensitive parsing of all supported types and errors for invalid and
empty input through table-driven cases.

**Oracle:** The test passes only when its assertions prove
case-insensitive parsing of all supported types and errors for invalid and empty input
through table-driven cases.

## TEST-INTERNAL_MODEL-020: Basic identifier construction

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_MODEL-002`, `SPEC-INTERNAL_MODEL-018`
- **Contracts:** `IdentifierRecord`, `cmd/idd-cli#Identifier`

**Purpose:**

Prove identity, title, source, line, raw reference, and non-nil link
initialization for the basic constructor.

**Oracle:** The test passes only when its assertions confirm identity,
title, source, line, raw reference, and non-nil link initialization for the basic
constructor.

## TEST-INTERNAL_MODEL-021: New validation result state

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-010`

**Purpose:**

Prove the initial invalid state and allocated empty error and warning slices.

**Oracle:** The test passes only when its assertions confirm the initial
invalid state and allocated empty error and warning slices.

## TEST-INTERNAL_MODEL-022: Annotation construction

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_MODEL-008`
- **Contracts:** `SourceAnnotation`

**Purpose:**

Prove that annotation construction retains type, target reference, source, and
line evidence.

**Oracle:** The test passes only when its assertions confirm
annotation construction retains type, target reference, source, and line evidence.

## TEST-INTERNAL_MODEL-023: Description-aware identifier construction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-002`, `SPEC-INTERNAL_MODEL-018`

**Purpose:**

Prove preservation of the supplied description and default document origin in
addition to the basic identifier fields.

**Oracle:** The test passes only when its assertions confirm preservation
of the supplied description and default document origin in addition to the basic
identifier fields.

## TEST-INTERNAL_MODEL-024: Annotation construction with comment

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-008`

**Purpose:**

Prove that the optional declaration comment is retained alongside ordinary
annotation evidence.

**Oracle:** The test passes only when its assertions confirm the
optional declaration comment is retained alongside ordinary annotation evidence.

## TEST-INTERNAL_MODEL-028: Relationship constants

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_MODEL-012`
- **Contracts:** `Relationship`

**Purpose:**

Prove stable serialized strings for the primary relationship constants consumed
by graph and reporter logic.

**Oracle:** The test passes only when its assertions confirm stable
serialized strings for the primary relationship constants consumed by graph and reporter
logic.

## TEST-INTERNAL_MODEL-029: Reverse relationship mapping

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-012`

**Purpose:**

Prove directional pairs, symmetric reference/annotation types, and pass-through
behavior for an unknown type.

**Oracle:** The test passes only when its assertions confirm directional
pairs, symmetric reference/annotation types, and pass-through behavior for an unknown
type.

## TEST-INTERNAL_MODEL-030: Directed link construction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-013`, `SPEC-INTERNAL_MODEL-015`

**Purpose:**

Prove exact preservation of endpoints, type, source, and line in a newly
constructed link.

**Oracle:** The test passes only when its assertions confirm exact
preservation of endpoints, type, source, and line in a newly constructed link.

## TEST-INTERNAL_MODEL-031: LLM report JSON round trip

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_MODEL-037`
- **Contracts:** `LLMFindingReport`

**Purpose:**

Prove that a populated schema, summary group, finding, location, repair hint,
and related identifier survive JSON marshal and unmarshal.

**Oracle:** The test passes only when its assertions confirm a
populated schema, summary group, finding, location, repair hint, and related identifier
survive JSON marshal and unmarshal.

## TEST-INTERNAL_MODEL-032: Explicit origin reassignment

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_MODEL-017`

**Purpose:**

Prove that `SetOrigin` replaces the constructor's document origin with code
origin on the same identifier.

**Oracle:** The test passes only when its assertions confirm
`SetOrigin` replaces the constructor's document origin with code origin on the same
identifier.

## Strategy

Tests are deterministic in-memory unit tests. Table-driven cases cover enum and
reverse-link mappings; exact field assertions cover constructors and mutation;
JSON round trips cover the wire model. No test relies on filesystem state,
time, randomness, or network access.

The suite directly covers canonical-owner conflicts across and within
directories and the non-conflicting source-evidence case. It does not claim
coverage of nil insertion, mutation aliasing, every duplicate-ID secondary
ordering, concurrency, every link constant, every JSON omission rule, or
invalid report field combinations. Engine and reporter tests cover some
consumer behavior, but these model-level boundaries remain explicit.
