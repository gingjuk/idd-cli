---
idd:
  version: "1.0"
  package: internal/collector
  document: spec
---

# Specifications: internal/collector

## SPEC-INTERNAL_COLLECTOR-001: Distributed ownership and safety

- **Design:** `CollectorModule`
- **Contract:** `IDDDocumentSet`

**Requirement:** Parse and validate minimal IDD identity frontmatter plus
human-readable role-specific Markdown records while preserving legacy Markdown
collection and safe document initialization and repair.

Each of the four Markdown documents is self-describing. Its YAML frontmatter
contains only identity; its level-two Markdown records own design components,
contracts, SPECs, or TESTs and coverage. The collector derives reverse
SPEC-to-TEST links from `testing.md` `Covers` fields, so a relationship is
written once without creating a shared write target.

CommonMark AST parsing keeps normal prose, examples, and subordinate headings
free-form while providing exact record and field source locations. Syntax
errors, semantic YAML catalogs, unknown fields, path/package/role disagreement,
package-derived module disagreement, missing files, placeholders, duplicates,
large registry tables, and unresolved references become structured findings. A
bad document set does not prevent unrelated packages from being collected.

Initialization is non-overwriting and refuses legacy metadata before its first
write. Repair normalizes only identity and atomically replaces each selected
Markdown file; it never reformats body prose or invents semantic records.

## SPEC-INTERNAL_COLLECTOR-002: Document collector construction

- **Design:** `CollectorModule`
- **Contract:** `DocCollector`

**Requirement:** Construct a `DocCollector` with the supplied IDD
configuration.

The collector retains the configuration used for document discovery,
validation, and ignore-path decisions.

## SPEC-INTERNAL_COLLECTOR-003: Source annotation collector

- **Design:** `CollectorModule`
- **Contract:** `CodeCollector`

**Requirement:** Represent collection of IDD annotations from supported source
files.

Source annotations are normalized into the same identifier model as
documentation records.

## SPEC-INTERNAL_COLLECTOR-004: Source collector construction

- **Design:** `CollectorModule`
- **Contract:** `CodeCollector`

**Requirement:** Construct a `CodeCollector` with the supplied IDD
configuration.

The resulting collector applies the configured source patterns and ignore
rules.

## SPEC-INTERNAL_COLLECTOR-005: Legacy marker values

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Represent legacy marker identifiers with meaningful names and
descriptions.

Legacy marker values remain available to packages that have not migrated to
self-describing documents.

## SPEC-INTERNAL_COLLECTOR-006: Legacy related document set

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Represent the four related narrative document paths in legacy
frontmatter.

The legacy relationship set keeps the four fixed document roles explicit.

## SPEC-INTERNAL_COLLECTOR-007: Legacy frontmatter values

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Represent legacy marker and related-file metadata as one
frontmatter value.

The legacy parser exposes both kinds of metadata without mixing them into the
new document model.

## SPEC-INTERNAL_COLLECTOR-008: Legacy frontmatter parsing

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Parse the leading YAML frontmatter document without treating
fenced-code separators as metadata boundaries.

Only the leading frontmatter block can define legacy document metadata.

## SPEC-INTERNAL_COLLECTOR-009: Legacy marker validation

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Verify that frontmatter markers resolve to correctly formed
Markdown detail headings.

Declared legacy markers must have matching human-readable definitions.

## SPEC-INTERNAL_COLLECTOR-010: Legacy marker formatting

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Report bare IDD markers outside headings, frontmatter, and
fenced code.

Legacy references remain backtick-delimited in narrative Markdown.

## SPEC-INTERNAL_COLLECTOR-011: Legacy heading validation

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Require IDD detail headings to contain a valid identifier and
a meaningful description after a colon.

Malformed or description-free legacy headings produce structural findings.

## SPEC-INTERNAL_COLLECTOR-012: Narrative filename mapping

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Map each IDD identifier type to its fixed narrative Markdown
filename.

The mapping keeps legacy declarations in the same four conceptual roles as the
self-describing format.

## SPEC-INTERNAL_COLLECTOR-013: Narrative document placement

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Verify that legacy identifiers are declared in the narrative
file assigned to their type while allowing root documentation.

Package documentation is checked for role placement without applying the rule
to root-level explanatory documents.

## SPEC-INTERNAL_COLLECTOR-017: Deterministic document collection

- **Design:** `CollectorModule`
- **Contract:** `DocCollector`

**Requirement:** Discover self-describing document sets before Markdown and
select self-describing or legacy parsing deterministically for each package.

Directory collection has two ordered passes:

1. discover directories where a fixed document has an `idd` frontmatter block;
2. parse each four-document set and then parse all other Markdown through the
   legacy frontmatter path.

When the target is one fixed self-describing document, its three siblings are
included in the same validation scope. The body may use a declared ID as an
optional detail heading, but it cannot reintroduce legacy `markers` or
`related_files`.

## SPEC-INTERNAL_COLLECTOR-024: Source annotation collection

- **Design:** `CollectorModule`
- **Contract:** `CodeCollector`

**Requirement:** Collect `@implement`, `@test`, and `@test-contract`
annotations with source context from Go, TypeScript, and JavaScript files.

The source collector walks supported Go, TypeScript, TSX, and JavaScript files,
honors configured ignore paths and `idd:ignore` scopes, and creates one
code-origin identifier per annotation reference.

Declaration context is retained for consistency diagnostics. `@test` produces
TEST kind `test`; `@test-contract` produces TEST kind `contract`, allowing the
engine to compare source annotations with `testing.md` records.

## SPEC-INTERNAL_COLLECTOR-025: Multiple references

- **Design:** `CollectorModule`
- **Contract:** `CodeCollector`

**Requirement:** Split comma-separated annotation references into trimmed
identifiers without losing their annotation kind.

One annotation may contain comma-separated identifiers. Each value is trimmed
and collected independently while retaining the shared annotation kind and
source location.

## Legacy compatibility

The legacy path remains isolated in `frontmatter.go`:

- `SPEC-INTERNAL_COLLECTOR-005` through
  `SPEC-INTERNAL_COLLECTOR-007` define marker and frontmatter values;
- `SPEC-INTERNAL_COLLECTOR-008` through
  `SPEC-INTERNAL_COLLECTOR-011` parse and validate legacy metadata and headings;
- `SPEC-INTERNAL_COLLECTOR-012` and
  `SPEC-INTERNAL_COLLECTOR-013` enforce narrative filename placement.

Self-describing packages do not execute these legacy relationship-parsing
rules.
