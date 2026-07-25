---
idd:
  version: "1.0"
  package: internal/collector
  document: testing
---

# Testing: internal/collector

## TEST-INTERNAL_COLLECTOR-001: Document collection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Function `TestDocCollector_Collect` verifies basic documentation
collection.

## TEST-INTERNAL_COLLECTOR-002: Frontmatter-backed collection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-006`, `SPEC-INTERNAL_COLLECTOR-007`,
  `SPEC-INTERNAL_COLLECTOR-008`, `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Function `TestDocCollector_Collect_WithFrontmatter` verifies legacy
frontmatter collection.

## TEST-INTERNAL_COLLECTOR-003: Missing documentation target

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Function `TestDocCollector_Collect_FileNotFound` verifies a missing
target returns an empty set.

## TEST-INTERNAL_COLLECTOR-004: Non-Markdown target

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Function `TestDocCollector_Collect_NonMarkdownFile` verifies
unsupported document files are ignored.

## TEST-INTERNAL_COLLECTOR-005: Section title extraction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`

**Purpose:** Function `TestDocCollector_extractTitle` verifies titles are read
from identifier headings.

## TEST-INTERNAL_COLLECTOR-006: Missing section title

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`

**Purpose:** Function `TestDocCollector_extractTitle_NotFound` verifies absent
headings return no title.

## TEST-INTERNAL_COLLECTOR-007: Empty documentation directory

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Function `TestDocCollector_Collect_DirectoryWithNoSpec` verifies an
empty directory is handled safely.

## TEST-INTERNAL_COLLECTOR-008: Marker description values

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-005`

**Purpose:** Function `TestMarker_Describe` verifies legacy marker descriptions
are retained.

## TEST-INTERNAL_COLLECTOR-009: Frontmatter value fields

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-006`, `SPEC-INTERNAL_COLLECTOR-007`

**Purpose:** Function `TestFrontmatter_Markers` verifies marker and related-file
values.

## TEST-INTERNAL_COLLECTOR-010: Frontmatter parsing

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-008`

**Purpose:** Function `TestParseFrontmatter` verifies valid, absent, empty, and
fenced-code inputs.

## TEST-INTERNAL_COLLECTOR-011: Marker and heading validation

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-009`, `SPEC-INTERNAL_COLLECTOR-010`,
  `SPEC-INTERNAL_COLLECTOR-011`

**Purpose:** Function `TestValidateFrontmatterMarkers` verifies declared
markers and heading formatting.

## TEST-INTERNAL_COLLECTOR-012: Defined marker extraction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-009`

**Purpose:** Function `TestExtractDefinedMarkers` verifies identifier headings
are recognized as definitions.

## TEST-INTERNAL_COLLECTOR-013: Referenced marker extraction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-009`

**Purpose:** Function `TestExtractReferencedMarkers` verifies referenced
identifiers are recognized.

## TEST-INTERNAL_COLLECTOR-014: Multi-reference marker extraction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-009`

**Purpose:** Function `TestExtractReferencedMarkers` verifies multiple
references are retained.

## TEST-INTERNAL_COLLECTOR-015: Expected narrative filename

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-012`

**Purpose:** Function `TestGetExpectedFilename` verifies identifier types map to
fixed Markdown filenames.

## TEST-INTERNAL_COLLECTOR-016: Root documentation placement

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-013`

**Purpose:** Function `TestIsRootDocFile` verifies root documentation is exempt
from package placement.

## TEST-INTERNAL_COLLECTOR-017: Package document placement

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-013`

**Purpose:** Function `TestValidateDocumentStructure` verifies identifier types
use the correct narrative file.

## TEST-INTERNAL_COLLECTOR-019: IDD reference extraction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-009`

**Purpose:** Function `TestExtractIDDRefs` verifies supported identifiers are
extracted from headings and references.

## TEST-INTERNAL_COLLECTOR-021: Self-describing document workflows

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-002`,
  `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Document tests verify minimal identity parsing, CommonMark
records, exact diagnostics, derived links, safe initialization, single-file
writes, and idempotent repair.

## TEST-INTERNAL_COLLECTOR-022: Collector boundary cases

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-024`

**Purpose:** Collector tests verify section parsing and function-context
bounds.

## TEST-INTERNAL_COLLECTOR-023: Coverage and receiver parsing

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-024`

**Purpose:** Collector tests verify SPEC coverage parsing and pointer-receiver
context.

## TEST-INTERNAL_COLLECTOR-024: Source collection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-003`,
  `SPEC-INTERNAL_COLLECTOR-004`, `SPEC-INTERNAL_COLLECTOR-024`

**Purpose:** Function `TestCodeCollector_Collect` verifies construction,
collection, and legacy TEST-field parsing.

## TEST-INTERNAL_COLLECTOR-025: Go annotation collection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-024`

**Purpose:** Function `TestCodeCollector_CollectGoFile` verifies Go annotations
and legacy contract fields.

## TEST-INTERNAL_COLLECTOR-026: Multiple source annotations

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-024`,
  `SPEC-INTERNAL_COLLECTOR-025`

**Purpose:** Collector tests verify comma-separated references, annotation
kinds, and legacy implements fields.

## TEST-INTERNAL_COLLECTOR-027: Document links and source context

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-017`,
  `SPEC-INTERNAL_COLLECTOR-024`

**Purpose:** Collector tests verify derived legacy links and source function
context.

## TEST-INTERNAL_COLLECTOR-028: Legacy contract links

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-001`, `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Function `TestDocCollector_Collect_WithContractLink` verifies
legacy contract references.

## TEST-INTERNAL_COLLECTOR-029: Documentation ignore paths

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_COLLECTOR-017`

**Purpose:** Function `TestDocCollector_shouldIgnore` verifies configured
documentation paths are ignored during discovery.

## Strategy

This document owns TEST titles, purposes, kinds, and coverage. Test source uses
`@test` annotations for behavior records; the body keeps suite-level strategy
and fixture boundaries.

## Self-describing document scenarios

`TEST-INTERNAL_COLLECTOR-021` is table-driven around the new document model. It
covers:

- minimal identity YAML and CommonMark record parsing;
- semantic-YAML migration findings, unknown fields, malformed YAML,
  module/path mismatches, placeholders, and unresolved references;
- exact record/field source lines;
- derived reverse coverage links;
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

`TEST-INTERNAL_COLLECTOR-022` through
`TEST-INTERNAL_COLLECTOR-027` cover source discovery, declaration context,
multiple references, and semantic test-kind preservation. Go source fixtures
are wrapped in `idd:ignore` scopes where their sample annotations must not be
collected from the test file itself.

## Integration boundary

`TEST-INTERNAL_COLLECTOR-028` and `TEST-INTERNAL_COLLECTOR-029` exercise legacy
contract links and ignore-path behavior. `collector_contract_test.go` separately
checks concrete public collector behavior and failure handling; there is no
mock collector interface.
