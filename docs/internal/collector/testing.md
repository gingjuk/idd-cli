---
idd:
  version: "1.1"
  package: internal/collector
  namespace: INTERNAL_COLLECTOR
---

# Testing: internal/collector

## TEST-INTERNAL_COLLECTOR-001: Document collection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Verify that Markdown discovery produces documentation-origin
identifiers and recoverable structural findings without requiring callers to
know which parser handled the file.

**Oracle:** The test passes only when its assertions confirm
Markdown discovery produces documentation-origin identifiers and recoverable structural
findings without requiring callers to know which parser handled the file.

## TEST-INTERNAL_COLLECTOR-002: Frontmatter-backed collection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-006`, `SPEC-INTERNAL_COLLECTOR-007`,
  `SPEC-INTERNAL_COLLECTOR-008`, `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Verify that an unmigrated package still collects legacy markers,
related-file metadata, and relationships without entering self-describing mode.

**Oracle:** The test passes only when its assertions confirm an
unmigrated package still collects legacy markers, related-file metadata, and
relationships without entering self-describing mode.

## TEST-INTERNAL_COLLECTOR-003: Missing documentation target

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Verify that a missing documentation target follows the documented
empty-result boundary instead of inventing identifiers or crashing collection.

**Oracle:** The test passes only when its assertions confirm a
missing documentation target follows the documented empty-result boundary instead of
inventing identifiers or crashing collection.

## TEST-INTERNAL_COLLECTOR-004: Non-Markdown target

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Verify that an explicitly selected non-Markdown file contributes
no documentation evidence or spurious parse findings.

**Oracle:** The test passes only when its assertions confirm an
explicitly selected non-Markdown file contributes no documentation evidence or spurious
parse findings.

## TEST-INTERNAL_COLLECTOR-005: Section title extraction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`

**Purpose:** Verify that a legacy identifier receives its human-readable title
from the matching detail heading rather than from an inline reference.

**Oracle:** The test passes only when its assertions confirm a
legacy identifier receives its human-readable title from the matching detail heading
rather than from an inline reference.

## TEST-INTERNAL_COLLECTOR-006: Missing section title

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`

**Purpose:** Verify that an absent matching heading yields no invented title,
preserving the distinction between a declaration and a reference.

**Oracle:** The test passes only when its assertions confirm an
absent matching heading yields no invented title, preserving the distinction between a
declaration and a reference.

## TEST-INTERNAL_COLLECTOR-007: Empty documentation directory

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Verify that an empty documentation directory returns a stable
empty result and does not assume a missing SPEC declaration.

**Oracle:** The test passes only when its assertions confirm an
empty documentation directory returns a stable empty result and does not assume a
missing SPEC declaration.

## TEST-INTERNAL_COLLECTOR-008: Marker description values

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-005`

**Purpose:** Verify that legacy marker names and descriptions survive parsing
for human-facing diagnostics and migration.

**Oracle:** The test passes only when its assertions confirm legacy
marker names and descriptions survive parsing for human-facing diagnostics and
migration.

## TEST-INTERNAL_COLLECTOR-009: Frontmatter value fields

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-006`, `SPEC-INTERNAL_COLLECTOR-007`

**Purpose:** Verify that legacy frontmatter represents marker values and the
four related narrative roles without conflating their ownership.

**Oracle:** The test passes only when its assertions confirm legacy
frontmatter represents marker values and the four related narrative roles without
conflating their ownership.

## TEST-INTERNAL_COLLECTOR-010: Frontmatter parsing

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-008`
- **Contracts:** `LegacyFrontmatter`

**Purpose:** Verify leading-frontmatter parsing across valid, absent, empty,
malformed, and fenced-example inputs, including preservation of the no-metadata
case.

**Oracle:** The test passes only when its assertions verify
leading-frontmatter parsing across valid, absent, empty, malformed, and fenced-example
inputs, including preservation of the no-metadata case.

## TEST-INTERNAL_COLLECTOR-011: Marker and heading validation

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-009`, `SPEC-INTERNAL_COLLECTOR-010`,
  `SPEC-INTERNAL_COLLECTOR-011`

**Purpose:** Verify that legacy declarations, definitions, and backtick
references agree, with actionable findings for missing or malformed headings.

**Oracle:** The test passes only when its assertions confirm legacy
declarations, definitions, and backtick references agree, with actionable findings for
missing or malformed headings.

## TEST-INTERNAL_COLLECTOR-012: Defined marker extraction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-009`

**Purpose:** Verify that only qualifying identifier headings establish legacy
definitions and ordinary mentions do not.

**Oracle:** The test passes only when its assertions confirm only
qualifying identifier headings establish legacy definitions and ordinary mentions do
not.

## TEST-INTERNAL_COLLECTOR-013: Referenced marker extraction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-009`

**Purpose:** Verify that quoted or backtick-delimited identifier uses are
recognized as legacy references without redefining their owner.

**Oracle:** The test passes only when its assertions confirm quoted
or backtick-delimited identifier uses are recognized as legacy references without
redefining their owner.

## TEST-INTERNAL_COLLECTOR-014: Multi-reference marker extraction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-009`

**Purpose:** Verify that one narrative location can retain several identifier
references without dropping order-independent relationship evidence.

**Oracle:** The test passes only when its assertions confirm one
narrative location can retain several identifier references without dropping
order-independent relationship evidence.

## TEST-INTERNAL_COLLECTOR-015: Expected narrative filename

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-012`

**Purpose:** Verify the stable legacy mapping from SPEC, TEST, CONTRACT, and
DESIGN identifier types to their four owning Markdown filenames.

**Oracle:** The test passes only when its assertions confirm the stable
legacy mapping from SPEC, TEST, CONTRACT, and DESIGN identifier types to their four
owning Markdown filenames.

## TEST-INTERNAL_COLLECTOR-016: Root documentation placement

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-013`

**Purpose:** Verify that root-level explanatory documentation may discuss
identifiers without being mistaken for a package-local declaration file.

**Oracle:** The test passes only when its assertions confirm
root-level explanatory documentation may discuss identifiers without being mistaken for
a package-local declaration file.

## TEST-INTERNAL_COLLECTOR-017: Package document placement

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-013`

**Purpose:** Verify that package-local legacy declarations remain in the
narrative file assigned to their identifier type.

**Oracle:** The test passes only when its assertions confirm
package-local legacy declarations remain in the narrative file assigned to their
identifier type.

## TEST-INTERNAL_COLLECTOR-019: IDD reference extraction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-009`

**Purpose:** Verify supported identifier syntax is extracted from legacy
headings and references while internal or malformed markers remain excluded.

**Oracle:** The test passes only when its assertions confirm supported
identifier syntax is extracted from legacy headings and references while internal or
malformed markers remain excluded.

## TEST-INTERNAL_COLLECTOR-021: Self-describing document workflows

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-002`,
  `SPEC-INTERNAL_COLLECTOR-017`
- **Contracts:** `IDDDocumentSet`, `DocCollector`

**Purpose:** Document tests verify filename-owned roles, minimal
version/package identity parsing, CommonMark records, exact diagnostics,
scaffold completion, typed derived links, safe initialization, single-file
writes, idempotent repair, and actionable split-document recovery.

**Oracle:** Parsing succeeds only when the exact canonical basename determines
the role and frontmatter contains no `document` field. Former role fields,
IDD metadata on arbitrary names, and `design-*`, `contract-*`, `spec-*`, or
`testing-*` splits produce their expected findings. Completion inspection of a
split file and a directory containing it must return a prompt naming exact
source and target paths, demanding preservation of unique semantics and
implementation boundaries, forbidding summary-only reduction, delaying
deletion until a no-loss review, and supplying both verification commands.
Generated documents first report deterministic incomplete slot names and exact
marker lines. Removing only a marker remains incomplete; replacing every slot
with valid role-owned content yields an empty work list. Normal collection
reports the same incomplete state, while repair preserves markers and authored
bodies byte-for-byte. Extended record fixtures also prove Purpose, Guarantees,
Acceptance, Oracle, named Contract coverage, dependencies, concerns, and
lifecycle validation.

## TEST-INTERNAL_COLLECTOR-022: Collector boundary cases

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-024`

**Purpose:** Verify legacy record extraction stops at the owning section
boundary and handles missing, short, and multi-paragraph content safely.

**Oracle:** Each table row returns only the requested legacy marker section,
returns empty for a missing marker, and never consumes the next record or reads
outside the input.

## TEST-INTERNAL_COLLECTOR-023: Coverage and receiver parsing

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-024`

**Purpose:** Verify legacy `Spec Coverage` fields retain one or several SPEC
references without manufacturing coverage when the field is absent.

**Oracle:** Table rows produce the exact ordered SPEC lists for single,
multiple, and whitespace-separated values and nil for a section without the
legacy field.

## TEST-INTERNAL_COLLECTOR-024: Source collection

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-003`,
  `SPEC-INTERNAL_COLLECTOR-004`, `SPEC-INTERNAL_COLLECTOR-024`
- **Contracts:** `CodeCollector`

**Purpose:** Verify source discovery binds supported language comments to real
declarations and retains normalized source and semantic TEST-kind evidence.

**Oracle:** Table-driven Go, TypeScript, TSX, JavaScript/JSX, C++, Java, and
Python fixtures each produce the expected declaration name, kind, visibility,
source line, and attached annotation. Strings and body or detached comments
produce no binding; standalone ignore ranges suppress evidence; malformed
syntax and configured unsupported extensions each produce one `source-parse`
finding and no regex-fallback identifiers. Configured prefix fixtures prove
custom spellings bind while canonical spellings no longer match that
configured collector. Test-classification rows also cover Go tests,
benchmarks, fuzz targets, and examples; JavaScript `test`/`it` modifiers; Java
ordinary and parameterized test annotations; C++ test macros; and
ecosystem-specific filename conventions.

## TEST-INTERNAL_COLLECTOR-025: Go annotation collection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-024`

**Purpose:** Verify Go annotation parsing across implementation and test
declarations, including valid source locations and legacy relationship context.

**Oracle:** The test passes only when its assertions confirm Go
annotation parsing across implementation and test declarations, including valid source
locations and legacy relationship context.

## TEST-INTERNAL_COLLECTOR-026: Multiple source annotations

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-024`,
  `SPEC-INTERNAL_COLLECTOR-025`

**Purpose:** Verify comma-separated annotations produce independent identifiers
at one source location and preserve behavioral versus contract TEST kind.

**Oracle:** The test passes only when its assertions verify
comma-separated annotations produce independent identifiers at one source location and
preserve behavioral versus contract TEST kind.

## TEST-INTERNAL_COLLECTOR-027: Document links and source context

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-017`,
  `SPEC-INTERNAL_COLLECTOR-024`

**Purpose:** Verify documentation relationships and source declaration context
survive collection into the common identifier model used by graph validation.

**Oracle:** The test passes only when its assertions verify
documentation relationships and source declaration context survive collection into the
common identifier model used by graph validation.

## TEST-INTERNAL_COLLECTOR-028: Legacy contract links

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Verify an unmigrated contract-to-SPEC relationship remains
available until that package is explicitly moved to canonical SPEC ownership.

**Oracle:** The test passes only when its assertions confirm an
unmigrated contract-to-SPEC relationship remains available until that package is
explicitly moved to canonical SPEC ownership.

## TEST-INTERNAL_COLLECTOR-029: Documentation ignore paths

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Verify configured documentation ignore patterns suppress matching
files and directories without hiding unrelated package evidence.

**Oracle:** The test passes only when its assertions confirm configured
documentation ignore patterns suppress matching files and directories without hiding
unrelated package evidence.

## TEST-INTERNAL_COLLECTOR-030: Collector public boundary contracts

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-002`,
  `SPEC-INTERNAL_COLLECTOR-003`, `SPEC-INTERNAL_COLLECTOR-004`,
  `SPEC-INTERNAL_COLLECTOR-017`, `SPEC-INTERNAL_COLLECTOR-024`
- **Contracts:** `IDDDocumentSet`, `DocCollector`, `CodeCollector`,
  `LegacyFrontmatter`

**Purpose:** Exercise the collector constructors, collection boundaries,
frontmatter behavior, source-language inputs, and recoverable failure results
through the exported package surface.

**Oracle:** Every table-driven contract case returns the documented identifier
set, origin, finding, or error boundary without relying on an undocumented
collector interface or accepting detached source annotations.

## TEST-INTERNAL_COLLECTOR-031: Focused SPEC review evidence

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-026`
- **Contracts:** `SpecReviewContext`

**Purpose:** Verify deterministic collection of one SPEC, its Contract,
covering TEST records, and annotated implementation and test declaration
excerpts, plus shared-scan multi-SPEC collection, ordered deduplication, the
ten-SPEC limit, cross-package evidence, peer-H2 record boundaries, and missing,
duplicate, ignored, and bounded-excerpt cases.

**Oracle:** Table-driven fixtures return complete contexts in first-request
order, collect documentation/source once per batch, collapse repeated IDs, and
return atomic operational errors for malformed, absent, multiply owned, empty,
or oversized requests. No case produces a semantic score or verdict.

## TEST-INTERNAL_COLLECTOR-032: Batch document operations

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-027`
- **Contracts:** `IDDDocumentSet`

**Purpose:** Verify multi-target completion aggregation and whole-batch
initialization and repair preflight.

**Oracle:** Table-driven fixtures prove normalized first-occurrence targets,
relative-versus-absolute path identity, stable incomplete-slot ordering,
overlap and duplicate deduplication, and single-target compatibility.
Initialization and repair apply each planned path once, while a later invalid
batch member leaves every earlier valid member unchanged.

## TEST-INTERNAL_COLLECTOR-033: Shared trace-project queries

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-026`, `SPEC-INTERNAL_COLLECTOR-028`

**Purpose:** Verify one collected TraceProject supplies ID dossiers and the v2
review-context projection without losing canonical ownership, design context,
multiple Contracts, declaration evidence, findings, or provenance.

**Oracle:** Table-driven fixtures pass only when resolved, unresolved,
ambiguous, and planned IDs retain the documented status; depth and mentions
options bound traversal; external project workdirs keep paths project-relative;
Components, Contracts, implementations, and TEST declarations are stable; and
review-context v2 is a projection of the same index. Invalid IDs, depth, scan
errors, and malformed canonical records return explicit operational errors.

## Strategy

This document owns TEST titles, purposes, kinds, and coverage. Test source uses
`@test` annotations for behavior records; the body keeps suite-level strategy
and fixture boundaries.

## Self-describing document scenarios

`TEST-INTERNAL_COLLECTOR-021` is table-driven around the new document model. It
covers:

- canonical filename role derivation, minimal version/package YAML, and
  CommonMark record parsing;
- rejection of the former `idd.document` key, arbitrary IDD filenames, and
  hyphenated, underscored, or dotted role splits;
- semantic-YAML migration findings, unknown fields, malformed YAML,
  module/path mismatches, placeholders, and unresolved references;
- exact record/field source lines;
- Component Purpose, Contract Guarantees, SPEC Acceptance, TEST Oracle, and
  contract TEST `Contracts` enforcement;
- concern-section, dependency-target/cycle, and lifecycle replacement
  validation;
- typed reverse coverage, named Contract, dependency, and replacement links;
- shared scaffold slots, deterministic `docs status` ordering, missing sibling
  work items, marker/content dual completion, split-document recovery prompts,
  and repair preservation;
- multi-target status aggregation, overlap deduplication, and whole-batch
  initialization and repair preflight;
- narrative metadata and heading validation;
- minimal-identity round-trip serialization with byte-preserved bodies;
- large registry-table rejection and fenced-example isolation;
- non-overwriting initialization, generic-frontmatter merging, legacy and
  central-catalog refusal, and path traversal;
- single-file write scope, body preservation, and idempotent identity repair.

Temporary directory trees provide real package and `docs/<package>` layouts.
No test mutates repository documentation.

## Legacy scenarios

`TEST-INTERNAL_COLLECTOR-001` through
`TEST-INTERNAL_COLLECTOR-017`, plus `TEST-INTERNAL_COLLECTOR-019`, preserve
frontmatter, heading, filename, and legacy relationship behavior. These tests
prevent self-describing mode from changing packages that have not migrated.

## Source annotation scenarios

`TEST-INTERNAL_COLLECTOR-024` through
`TEST-INTERNAL_COLLECTOR-026` cover source discovery, normalized declaration
context, multiple references, and semantic test-kind preservation. The
`TEST-INTERNAL_COLLECTOR-024` contract fixtures cover every
supported grammar, visibility and test classification, attachment boundaries,
annotation-looking strings, exact ignore directives, syntax errors, and the
absence of regex fallback. `TEST-INTERNAL_COLLECTOR-027` separately proves
document-side relationship retention. Fixture comment ranges are ignored only
where the test source itself intentionally contains real annotation comments.

## Integration boundary

`TEST-INTERNAL_COLLECTOR-028` and `TEST-INTERNAL_COLLECTOR-029` exercise legacy
contract links and ignore-path behavior. `collector_contract_test.go` separately
checks concrete public collector behavior and failure handling; there is no
mock collector interface.
