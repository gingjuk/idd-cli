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

**Acceptance:**

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

Generated skeletons prompt authors to explain responsibilities, boundaries,
rationale, failures, acceptance evidence, scenarios, fixtures, and oracles.
Those prompts do not become declarations and do not make an empty scaffold
complete. The collector rejects empty or obvious placeholder content but does
not use word counts as a substitute for semantic review.

## SPEC-INTERNAL_COLLECTOR-002: Document collector construction

- **Design:** `CollectorModule`
- **Contract:** `DocCollector`

**Requirement:** Construct a `DocCollector` with the supplied IDD
configuration.

**Acceptance:**

The collector retains the configuration used for document discovery,
validation, and ignore-path decisions.

Construction performs no filesystem access and does not copy or normalize the
configuration. The caller owns supplying a validated, non-nil configuration
and may create independent collectors for different repository policies.

## SPEC-INTERNAL_COLLECTOR-003: Source annotation collector

- **Design:** `CollectorModule`
- **Contract:** `CodeCollector`

**Requirement:** Represent collection of IDD annotations from supported source
files.

**Acceptance:**

Source annotations are normalized into the same identifier model as
documentation records.

The collector represents source evidence only. It does not resolve whether a
referenced document identifier exists, build graph relationships, or judge
behavioral consistency; those decisions require merged document and code
origins in the engine.

## SPEC-INTERNAL_COLLECTOR-004: Source collector construction

- **Design:** `CollectorModule`
- **Contract:** `CodeCollector`

**Requirement:** Construct a `CodeCollector` with the supplied IDD
configuration.

**Acceptance:**

The resulting collector applies the configured source patterns and ignore
rules.

Construction has no traversal side effects and retains the configuration for
the collector lifetime. Multiple collectors may operate independently, but a
single validation run should use the same policy for documentation and source
evidence.

## SPEC-INTERNAL_COLLECTOR-005: Legacy marker values

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Represent legacy marker identifiers with meaningful names and
descriptions.

**Acceptance:**

Legacy marker values remain available to packages that have not migrated to
self-describing documents.

Each value retains the declared ID, display name, and optional description
needed by the legacy parser. It is compatibility data, not a canonical record
for a package that has opted into minimal `idd` identity.

## SPEC-INTERNAL_COLLECTOR-006: Legacy related document set

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Represent the four related narrative document paths in legacy
frontmatter.

**Acceptance:**

The legacy relationship set keeps the four fixed document roles explicit.

Paths describe where legacy declarations and backlinks are expected. They are
validated as metadata but never copied into self-describing frontmatter,
because fixed filenames already establish those roles.

## SPEC-INTERNAL_COLLECTOR-007: Legacy frontmatter values

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Represent legacy marker and related-file metadata as one
frontmatter value.

**Acceptance:**

The legacy parser exposes both kinds of metadata without mixing them into the
new document model.

The value may contain markers, related files, or both. Absence of a leading
frontmatter block yields no legacy metadata rather than an invented empty
catalog; malformed YAML returns a parse error with file context.

## SPEC-INTERNAL_COLLECTOR-008: Legacy frontmatter parsing

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Parse the leading YAML frontmatter document without treating
fenced-code separators as metadata boundaries.

**Acceptance:**

Only the leading frontmatter block can define legacy document metadata.

The parser ignores `---` sequences inside the body and fenced examples, returns
`nil` when no leading block exists, and preserves YAML decode errors so callers
can report structural debt instead of silently collecting incomplete markers.

## SPEC-INTERNAL_COLLECTOR-009: Legacy marker validation

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Verify that frontmatter markers resolve to correctly formed
Markdown detail headings.

**Acceptance:**

Declared legacy markers must have matching human-readable definitions.

Validation distinguishes a heading definition from an inline reference. Every
declared marker must resolve to an appropriate heading, and every discovered
definition must agree with frontmatter ownership. Findings identify the file
and marker so migration can preserve the explanatory section.

## SPEC-INTERNAL_COLLECTOR-010: Legacy marker formatting

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Report bare IDD markers outside headings, frontmatter, and
fenced code.

**Acceptance:**

Legacy references remain backtick-delimited in narrative Markdown.

The scan excludes frontmatter, headings that define a marker, fenced code, and
quoted examples. Its purpose is to keep references visibly distinct for human
readers and deterministic for the legacy extractor, not to ban ordinary prose
containing coincidental words.

## SPEC-INTERNAL_COLLECTOR-011: Legacy heading validation

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Require IDD detail headings to contain a valid identifier and
a meaningful description after a colon.

**Acceptance:**

Malformed or description-free legacy headings produce structural findings.

The accepted shape is a level-two identifier followed by a colon and meaningful
title. Identifier syntax remains strict and package-oriented; internal section
markers or a copied identifier used as its own title do not become valid
requirements.

## SPEC-INTERNAL_COLLECTOR-012: Narrative filename mapping

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Map each IDD identifier type to its fixed narrative Markdown
filename.

**Acceptance:**

The mapping keeps legacy declarations in the same four conceptual roles as the
self-describing format.

SPEC maps to `spec.md`, TEST to `testing.md`, CONTRACT to `contract.md`, and
DESIGN to `design.md`. Unknown and internal marker types return no narrative
filename rather than being assigned to an arbitrary document.

## SPEC-INTERNAL_COLLECTOR-013: Narrative document placement

- **Design:** `CollectorModule`
- **Contract:** `LegacyFrontmatter`

**Requirement:** Verify that legacy identifiers are declared in the narrative
file assigned to their type while allowing root documentation.

**Acceptance:**

Package documentation is checked for role placement without applying the rule
to root-level explanatory documents.

The check uses the file basename and identifier type. Root documentation may
quote identifiers for onboarding or architecture discussion, while a
package-local declaration must live in its owning narrative file so collectors
and readers share one location convention.

## SPEC-INTERNAL_COLLECTOR-017: Deterministic document collection

- **Design:** `CollectorModule`
- **Contract:** `DocCollector`

**Requirement:** Discover self-describing document sets before Markdown and
select self-describing or legacy parsing deterministically for each package.

**Acceptance:**

Directory collection has two ordered passes:

1. discover directories where a fixed document has an `idd` frontmatter block;
2. parse each four-document set and then parse all other Markdown through the
   legacy frontmatter path.

When the target is one fixed self-describing document, its three siblings are
included in the same validation scope. The body may use a declared ID as an
optional detail heading, but it cannot reintroduce legacy `markers` or
`related_files`.

Traversal and parse order are stable so repeated runs produce stable findings.
Recoverable document defects are returned as structured validation errors while
unavailable targets or traversal failures follow the documented collector error
boundary.

## SPEC-INTERNAL_COLLECTOR-024: Source annotation collection

- **Design:** `CollectorModule`
- **Contract:** `CodeCollector`

**Requirement:** Bind `@implement`, `@test`, and `@test-contract` comments to
real declarations in Go, TypeScript, TSX, JavaScript/JSX, C++, Java, and Python
source files.

**Acceptance:**

The source collector walks configured `.go`, `.ts`, `.tsx`, `.js`, `.jsx`,
`.cpp`, `.cc`, `.cxx`, `.hpp`, `.hh`, `.hxx`, `.java`, and `.py` files. A
pinned Tree-sitter grammar creates a normalized declaration and comment model
for each file. The collector honors configured ignore paths and standalone
`idd:ignore` ranges and creates one code-origin identifier per valid
declaration-attached annotation reference.

Configuration may rename the three annotation prefixes. Parsing maps those
lexical values back to the stable semantic kinds `implement`, `test`, and
`test-contract`, so custom spelling does not change identifier type, TEST kind,
or engine policy.

Declaration context is retained for consistency diagnostics. `@test` produces
TEST kind `test`; `@test-contract` produces TEST kind `contract`, allowing the
engine to compare source annotations with `testing.md` records.

Annotations are evidence attached to declarations, not free-form identifier
mentions. Strings, function-body comments, and detached comments are never
bound. A configured file with no pinned grammar or a supported file with
invalid syntax produces an exact `source-parse` finding and no partial
evidence; line-oriented or regex fallback is forbidden because it would make
placement and visibility results language-dependent and untrustworthy.

The normalized model retains declaration name, kind, source range, visibility,
test classification, annotations, and parse errors. Collection does not
require a corresponding document ID; the engine reports that mismatch after
both origins are available.

## SPEC-INTERNAL_COLLECTOR-025: Multiple references

- **Design:** `CollectorModule`
- **Contract:** `CodeCollector`

**Requirement:** Split comma-separated annotation references into trimmed
identifiers without losing their annotation kind.

**Acceptance:**

One annotation may contain comma-separated identifiers. Each value is trimmed
and collected independently while retaining the shared annotation kind and
source location.

Empty segments are ignored. Splitting does not validate cross-reference
existence or change identifier case; syntax validation and graph correspondence
remain with their owning stages.

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
