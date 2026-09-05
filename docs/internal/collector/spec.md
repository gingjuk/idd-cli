---
idd:
  version: "1.1"
  package: internal/collector
  namespace: INTERNAL_COLLECTOR
---

# Specifications: internal/collector

## SPEC-INTERNAL_COLLECTOR-001: Distributed ownership and safety

- **Components:** `CollectorModule`
- **Contracts:** `IDDDocumentSet`

**Requirement:** Parse and validate minimal IDD version/package frontmatter,
derive the document role only from its canonical filename, and preserve
human-readable role-specific Markdown records while retaining legacy Markdown
collection and safe document initialization and repair.

**Acceptance:**

Each source package and nested sub-package maps independently to
`docs/<package>/` and owns `design.md`, `contract.md`, `spec.md`, and
`testing.md`. The exact lowercase basename is the only role authority. YAML
frontmatter contains only `version` and `package`; the former `idd.document`
field is rejected without a compatibility mode. Level-two Markdown records own
design components, contracts, SPECs, or TESTs and coverage. The collector
derives reverse SPEC-to-TEST links from `testing.md` `Covers` fields, so a
relationship is written once without creating a shared write target.

CommonMark AST parsing keeps normal prose, examples, and subordinate headings
free-form while providing exact record and field source locations. Syntax
errors, semantic YAML catalogs, unknown fields, path/package disagreement,
non-canonical or split role filenames, package-derived module disagreement,
missing files, placeholders, duplicates, large registry tables, and unresolved
references become structured findings. A bad document set does not prevent
unrelated packages from being collected.

Initialization is non-overwriting and refuses legacy metadata before its first
write. Repair normalizes only version/package identity and atomically replaces
each selected Markdown file; it never emits role metadata, reformats body
prose, or invents semantic records.

Generated skeletons prompt authors to explain responsibilities, boundaries,
rationale, failures, acceptance evidence, scenarios, fixtures, and oracles.
Those prompts do not become declarations and do not make an empty scaffold
complete. The collector rejects empty or obvious placeholder content but does
not use word, line, or file-size counts as a substitute for semantic review.
Long role documents stay in the canonical file instead of being split into
feature- or size-suffixed filenames.

A split-role validation finding must retain enough structured evidence for the
reporter to generate a path-specific agent prompt. Direct completion
inspection of a split file or a tree containing one must return equivalent
instructions: read the split and canonical documents, preserve all unique
still-valid semantics rather than summarize them away, remove the split only
after a no-loss review, and rerun package status plus full-project validation.
Neither collection nor status inspection performs the merge or deletion.

## SPEC-INTERNAL_COLLECTOR-002: Document collector construction

- **Components:** `CollectorModule`
- **Contracts:** `DocCollector`

**Requirement:** Construct a `DocCollector` with the supplied IDD
configuration.

**Acceptance:**

Constructor contract tests require a non-nil collector and compare its retained
configuration pointer with the supplied value. Collection tests then observe
that its document patterns and ignore paths select the expected Markdown
evidence.

Construction performs no filesystem access and does not copy or normalize the
configuration. The caller owns supplying a validated, non-nil configuration
and may create independent collectors for different repository policies.

## SPEC-INTERNAL_COLLECTOR-003: Source annotation collector

- **Components:** `CollectorModule`
- **Contracts:** `CodeCollector`

**Requirement:** Represent collection of IDD annotations from supported source
files.

**Acceptance:**

Seven-language table-driven fixtures compare declaration name, kind,
visibility, source line, attached annotation kind, and resulting code-origin
identifier. String literals, body comments, detached comments, ignored ranges,
malformed syntax, and configured-but-unsupported extensions must not produce a
fallback identifier.

The collector represents source evidence only. It does not resolve whether a
referenced document identifier exists, build graph relationships, or judge
behavioral consistency; those decisions require merged document and code
origins in the engine.

## SPEC-INTERNAL_COLLECTOR-004: Source collector construction

- **Components:** `CollectorModule`
- **Contracts:** `CodeCollector`

**Requirement:** Construct a `CodeCollector` with the supplied IDD
configuration.

**Acceptance:**

Constructor tests require a non-nil collector retaining the supplied
configuration. Temporary source trees then demonstrate that configured
extensions and annotation prefixes are collected, configured ignore paths are
excluded, and unrelated Markdown files contribute no source evidence.

Construction has no traversal side effects and retains the configuration for
the collector lifetime. Multiple collectors may operate independently, but a
single validation run should use the same policy for documentation and source
evidence.

## SPEC-INTERNAL_COLLECTOR-005: Legacy marker values

- **Components:** `CollectorModule`
- **Contracts:** `LegacyFrontmatter`

**Requirement:** Represent legacy marker identifiers with meaningful names and
descriptions.

**Acceptance:**

Parsing a legacy marker with an ID, display name, and description preserves
all three values. The marker-description test compares the exact authored text
used by human-facing diagnostics.

Each value retains the declared ID, display name, and optional description
needed by the legacy parser. It is compatibility data, not a canonical record
for a package that has opted into minimal `idd` identity.

## SPEC-INTERNAL_COLLECTOR-006: Legacy related document set

- **Components:** `CollectorModule`
- **Contracts:** `LegacyFrontmatter`

**Requirement:** Represent the four related narrative document paths in legacy
frontmatter.

**Acceptance:**

Decoding frontmatter with explicit design, contract, specification, and
testing paths must return each authored value in its matching field.
Collecting that fixture must retain those related-file paths without entering
self-describing mode.

Paths describe where legacy declarations and backlinks are expected. They are
validated as metadata but never copied into self-describing frontmatter,
because fixed filenames already establish those roles.

## SPEC-INTERNAL_COLLECTOR-007: Legacy frontmatter values

- **Components:** `CollectorModule`
- **Contracts:** `LegacyFrontmatter`

**Requirement:** Represent legacy marker and related-file metadata as one
frontmatter value.

**Acceptance:**

Given valid legacy YAML containing markers and related narrative paths, the
parsed value exposes both groups without losing or reclassifying either. The
same package remains on the legacy collection path; minimal
`idd.version`/`idd.package` documents are parsed through the separate
self-describing representation.

The value may contain markers, related files, or both. Absence of a leading
frontmatter block yields no legacy metadata rather than an invented empty
catalog; malformed YAML returns a parse error with file context.

## SPEC-INTERNAL_COLLECTOR-008: Legacy frontmatter parsing

- **Components:** `CollectorModule`
- **Contracts:** `LegacyFrontmatter`

**Requirement:** Parse the leading YAML frontmatter document without treating
fenced-code separators as metadata boundaries.

**Acceptance:**

Table-driven parsing distinguishes valid leading YAML, no frontmatter, an
empty marker list, malformed YAML, and `---` delimiters inside fenced examples.
Only the leading block produces metadata; malformed leading YAML returns an
error rather than a partial value.

The parser ignores `---` sequences inside the body and fenced examples, returns
`nil` when no leading block exists, and preserves YAML decode errors so callers
can report structural debt instead of silently collecting incomplete markers.

## SPEC-INTERNAL_COLLECTOR-009: Legacy marker validation

- **Components:** `CollectorModule`
- **Contracts:** `LegacyFrontmatter`

**Requirement:** Verify that frontmatter markers resolve to correctly formed
Markdown detail headings.

**Acceptance:**

Validation fixtures distinguish a marker with its matching detail heading from
a missing heading, a malformed or description-free heading, and an ordinary
inline mention. The failing cases report the affected identifier instead of
inventing a definition.

Validation distinguishes a heading definition from an inline reference. Every
declared marker must resolve to an appropriate heading, and every discovered
definition must agree with frontmatter ownership. Findings identify the file
and marker so migration can preserve the explanatory section.

## SPEC-INTERNAL_COLLECTOR-010: Legacy marker formatting

- **Components:** `CollectorModule`
- **Contracts:** `LegacyFrontmatter`

**Requirement:** Report bare IDD markers outside headings, frontmatter, and
fenced code.

**Acceptance:**

Formatting-validation fixtures place the same declared marker in a heading, a
backtick-delimited reference, and bare body prose. Definitions and explicit
references are accepted; a bare body occurrence produces a source-located
formatting finding, while frontmatter and fenced examples remain excluded.

The scan excludes frontmatter, headings that define a marker, fenced code, and
quoted examples. Its purpose is to keep references visibly distinct for human
readers and deterministic for the legacy extractor, not to ban ordinary prose
containing coincidental words.

## SPEC-INTERNAL_COLLECTOR-011: Legacy heading validation

- **Components:** `CollectorModule`
- **Contracts:** `LegacyFrontmatter`

**Requirement:** Require IDD detail headings to contain a valid identifier and
a meaningful description after a colon.

**Acceptance:**

Heading-validation cases accept an identifier heading only when its ID matches
the declaration and a colon is followed by meaningful text. Missing, bare,
mismatched, or malformed headings produce a source-located finding naming the
declared marker.

The accepted legacy shape is a Markdown identifier heading followed by a colon
and meaningful title. Identifier syntax remains strict and package-oriented;
internal section markers or a copied identifier used as its own title do not
become valid requirements.

## SPEC-INTERNAL_COLLECTOR-012: Narrative filename mapping

- **Components:** `CollectorModule`
- **Contracts:** `LegacyFrontmatter`

**Requirement:** Map each IDD identifier type to its fixed narrative Markdown
filename.

**Acceptance:**

The mapping test compares SPEC, TEST, CONTRACT, and DESIGN identifiers with
`spec.md`, `testing.md`, `contract.md`, and `design.md` respectively, and
requires an unknown identifier kind to return an empty filename.

SPEC maps to `spec.md`, TEST to `testing.md`, CONTRACT to `contract.md`, and
DESIGN to `design.md`. Unknown and internal marker types return no narrative
filename rather than being assigned to an arbitrary document.

## SPEC-INTERNAL_COLLECTOR-013: Narrative document placement

- **Components:** `CollectorModule`
- **Contracts:** `LegacyFrontmatter`

**Requirement:** Verify that legacy identifiers are declared in the narrative
file assigned to their type while allowing root documentation.

**Acceptance:**

For a package-local legacy declaration, a matching role filename returns no
placement error and a mismatched role filename reports the expected canonical
file. Root-document cases return no placement error, allowing explanatory
Markdown to mention identifiers without becoming a package-local declaration
owner.

The check uses the file basename and identifier type. Root documentation may
quote identifiers for onboarding or architecture discussion, while a
package-local declaration must live in its owning narrative file so collectors
and readers share one location convention.

## SPEC-INTERNAL_COLLECTOR-017: Deterministic document collection

- **Components:** `CollectorModule`
- **Contracts:** `DocCollector`

**Requirement:** Discover self-describing document sets before Markdown and
select self-describing or legacy parsing deterministically for each package.

**Acceptance:**

Directory collection has ordered phases:

1. reject IDD metadata on arbitrary filenames and role-derived split files;
2. discover directories where a canonical document has an `idd` frontmatter
   block;
3. parse each four-document set and then parse other permitted Markdown through
   the legacy frontmatter path.

When the target is one fixed self-describing document, its three siblings are
included in the same validation scope. The body may use a declared ID as an
optional detail heading, but it cannot reintroduce legacy `markers` or
`related_files`.

Traversal and parse order are stable so repeated runs produce stable findings.
Recoverable document defects are returned as structured validation errors while
unavailable targets or traversal failures follow the documented collector error
boundary.

## SPEC-INTERNAL_COLLECTOR-024: Source annotation collection

- **Components:** `CollectorModule`
- **Contracts:** `CodeCollector`

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

- **Components:** `CollectorModule`
- **Contracts:** `CodeCollector`

**Requirement:** Split comma-separated annotation references into trimmed
identifiers without losing their annotation kind.

**Acceptance:**

One annotation may contain comma-separated identifiers. Each value is trimmed
and collected independently while retaining the shared annotation kind and
source location.

Empty segments are ignored. Splitting does not validate cross-reference
existence or change identifier case; syntax validation and graph correspondence
remain with their owning stages.

## SPEC-INTERNAL_COLLECTOR-026: Focused SPEC review evidence

- **Components:** `CollectorModule`
- **Contracts:** `SpecReviewContext`

**Requirement:** Assemble bounded canonical documentation, test, contract, and
source declaration evidence for one or a bounded batch of requested SPECs
without judging semantic quality.

**Acceptance:**

Each context uses schema `idd.spec_review_context.v2`, resolves one canonical
SPEC owner through TraceProject, and includes the complete bounded SPEC,
Component design context, zero or more Contracts, and covering TEST record
Markdown plus matching non-ignored implementation and test declarations,
preserves subordinate authored `Details`, and sorts evidence deterministically.
Record Markdown ends at the next peer H2 heading, so package-wide strategy or
guidance sections are not attributed to the preceding record. Qualified
`<package>#<name>` Contract references and TEST records in another collected
package remain eligible evidence.
Test declarations are selected through the covering TEST IDs, including both
`@test` and `@test-contract`. Records and declaration excerpts expose truncation.
Malformed, absent, and duplicate SPEC owners return operational errors.

Batch collection accepts at most ten unique IDs, deduplicates repeated inputs
in first-request order, performs one documentation scan and one source scan,
and returns schema `idd.spec_review_context_batch.v2`. Every relationship comes
from the same ID-centered index used by trace, not an independent join. Any invalid member fails
the whole request so a reviewer cannot mistake partial evidence for a complete
batch.

## SPEC-INTERNAL_COLLECTOR-027: Batch document operations

- **Components:** `CollectorModule`
- **Contracts:** `IDDDocumentSet`

**Requirement:** Inspect, initialize, or structurally repair one or more
document targets with deterministic deduplication and whole-batch preflight.

**Acceptance:**

Completion inspection accepts one or more role files, package directories, or
documentation trees. It cleans repeated inputs in first-request order,
aggregates the same role-schema work items used by single-target inspection,
deduplicates findings produced by overlapping targets, and sorts the resulting
work list deterministically. Relative and absolute spellings of the same
filesystem target share one identity. Schema `idd.document_status.v1` includes
the normalized first-occurrence `targets` that were inspected.

Initialization accepts one or more project-relative source packages. Repair
accepts one or more canonical role files or package directories. Both mutators
prepare every selected document result before the first write, merge repeated
or overlapping file plans, and reject conflicting plans. An invalid package,
legacy or central-catalog input, malformed document, or unsupported path
aborts the batch without changing a valid sibling target.

New files still use exclusive creation and existing files still use atomic
replacement. Whole-batch preflight does not claim a cross-filesystem
transaction: an unexpected I/O failure while applying an already validated
plan may leave earlier individually atomic writes in place. Batch operations
do not author semantic records.

## SPEC-INTERNAL_COLLECTOR-028: Shared ID-centered trace project

- **Components:** `CollectorModule`
- **Contracts:** `SpecReviewContext`

**Requirement:** Collect documentation and source once into a reusable project
query boundary that exposes the shared TraceIndex, AST declaration analyses,
and bounded dossiers for SPEC, TEST, Contract, and Component IDs.

**Acceptance:**

`BuildTraceProject` uses the configured external workdir for both document and
source collection, merges their observations once, preserves collector
findings, and builds the canonical Entity/Occurrence and
Relation/Provenance index. Read-only accessors expose the index and copies of
collected analyses so review-context and impacted never trigger a second scan.

`Trace` accepts a valid public or scoped ID, depth zero through five, and an
explicit mentions option. It preserves unresolved, ambiguous, planned, and
active-but-incomplete state; returns canonical and related semantic records,
declaration excerpts, provenance, and findings in stable order; and bounds
content without suppressing evidence. Plain Markdown mentions are discoverable
context only and never satisfy correspondence or coverage.

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
