---
idd:
  version: "1.0"
  package: internal/collector
---

# Design: internal/collector

## Component: CollectorModule

**Purpose:**

`CollectorModule` translates human-readable documentation and source
annotations into the common identifier model used by the validation engine. It
is the boundary between files written for people and normalized graph evidence
consumed by machines.

### Responsibilities

The component discovers documentation and source files, selects one
documentation mode per package, parses identity and semantic records, retains
source locations, converts declarations into identifiers and directed links,
and reports recoverable structural problems without discarding unrelated
evidence.

It also owns the only safe structural document mutations: creating a new
four-file skeleton and normalizing path-derived identity. Those operations are
kept next to document parsing so they share one definition of package path,
canonical filename, frontmatter, and body preservation.

### Boundaries

The collector does not decide whether a requirement is correct, whether a
design is adequate, or whether code satisfies behavior. It supplies evidence
and structural findings to `Engine`, which owns graph-wide policy. It does not
format final reports, and it does not generate semantic records from source
symbols.

Legacy metadata remains supported as an isolated input mode, not as a second
owner inside a self-describing package. A package-local `idd.yaml` is migration
debt and never participates as a marker source.

### Design rationale

Parsing a four-file set as a unit is necessary because Component, Contract,
SPEC, and TEST references cross files. Directory-level mode selection prevents
one sibling from using the new model while another silently contributes legacy
links. CommonMark AST parsing distinguishes authored headings and fields from
fenced examples, while YAML nodes retain the exact identity-field locations
needed for repairable diagnostics. The exact basename is the one role
authority, avoiding a second value that can disagree with the filesystem.

## Architecture

`CollectorModule` contains two concrete collectors; there is no shared
Collector interface.

```text
DocCollector
    ├── idd_document.go             identity and CommonMark record parsing
    ├── idd_document_validation.go  record, reference, and graph rules
    ├── document_schema.go          scaffold generation and completion status
    ├── document_files.go           non-overwriting init and per-file repair
    ├── doc_collector.go            self-describing/legacy mode selection
    └── frontmatter.go              legacy Markdown parsing

CodeCollector
    ├── source_ast.go       Tree-sitter registry and normalized declarations
    └── code_collector.go   source discovery and identifier conversion
```

### Documentation discovery and mode selection

`DocCollector.Collect` performs discovery before parsing:

1. walk the requested target and collect Markdown plus package-local
   `idd.yaml` candidates, applying documentation ignore rules;
2. reject IDD metadata on non-canonical names and reject split role names such
   as `design-auth.md`, recording the source, canonical basename, and
   `split-role` code needed by the reporter's path-specific repair prompt;
3. inspect the four fixed role files for a top-level `idd` block;
4. mark each matching directory as self-describing and sort all paths for
   deterministic processing;
5. parse every marked directory as one four-file set;
6. send unmarked Markdown through the legacy parser.

When the explicit target is one of the four fixed files, its existing siblings
are added to the same scope. This preserves package-level reference validation
without claiming that the overall CLI run is package-only; source collection
still uses the project working tree.

If any fixed file has self-describing metadata, the directory is committed to
that mode. Missing siblings, legacy fields, or a central catalog become
structured findings rather than fallback signals.

### Self-describing parse pipeline

Each document is processed in layers:

1. `splitLeadingFrontmatter` isolates only a real leading frontmatter block and
   preserves the remaining body bytes.
2. YAML decoding requires a top-level `idd` mapping containing only `version`
   and `package`. `idd.document`, other unknown keys, and semantic catalogs are
   rejected immediately without a compatibility branch.
3. The path-derived package establishes expected identity and module prefix;
   the exact lowercase filename establishes the role before Markdown records
   are parsed.
4. Goldmark parses the body. Level-two role headings define records, while
   subordinate prose, examples, lists, diagrams, and code blocks remain
   human-authored content.
5. Role validators require the small machine-readable fields and validate
   meaningful titles/narratives, duplicates, references, TEST kinds, and
   placeholders.
6. The shared completion schema checks generated markers and bounded effective
   content, producing the same work items for `docs status` and normal runs.
7. Valid declarations become document-origin identifiers. TEST `Covers`
   values create the canonical TEST-to-SPEC edge and its derived reverse graph
   relationship. Component dependencies, lifecycle replacements, and contract
   TEST mappings become typed links for graph-wide validation.

Issues retain a rule, message, path, line, related identifier, and field code.
YAML line offsets and Markdown AST segments are normalized before the finding
leaves the package.

The parser deliberately does not impose a prose or document length.
Structural checks can reject empty or obvious placeholder content, but
semantic sufficiency belongs to the Skill and human review. Long documents
remain in the canonical role file rather than producing `design-*`,
`contract-*`, `spec-*`, or `testing-*` fragments.

`InspectDocumentCompletion` treats those split names as a blocking structural
state rather than pretending their content belongs to a completion slot. Its
error prompt instructs the agent to read both source and target, structurally
create a missing canonical set if needed, preserve detailed semantics during
the merge, compare before deleting, and rerun the package and project gates.
The collector formats guidance but retains a read-only boundary.

### Legacy isolation

Legacy parsing remains in `frontmatter.go` and the legacy portions of
`doc_collector.go`. It recognizes `markers`, `related_files`, identifier
headings, and repeated relationship fields exactly for packages where no fixed
file opts into the new model.

Keeping the old path operational enables package-by-package migration.
Preventing both paths from running in one package avoids duplicate identifiers
and conflicting relationship ownership.

### Source annotation collection

`CodeCollector` walks configured files only when `source_ast.go` has a pinned
grammar for their extension. The official Tree-sitter grammars cover Go,
TypeScript, TSX, JavaScript/JSX, C++, Java, and Python. One normalized
`SourceAnalysis` retains syntax errors, real comment nodes, declaration
boundaries, visibility, test classification, and declaration-attached
annotations.

Binding permits intervening declaration comments and limited whitespace but no
executable tokens. Comments inside declaration bodies and detached comments
remain visible for placement diagnostics but cannot establish traceability.
Annotation-looking strings never enter the comment stream. A syntax error
invalidates evidence for the complete file and becomes a `source-parse`
finding; there is intentionally no regex fallback.

Each attached code identifier retains origin, file, line, annotation context,
and normalized declaration context. TEST identifiers also retain kind:
`@test` becomes `test`, and `@test-contract` becomes `contract`. The engine
receives the full analyses to enforce language-specific public and test rules,
then compares that evidence with the `Kind` declared by `testing.md`.

The collector requires an annotation prefix rather than treating every
identifier-looking string as evidence. Examples and fixtures can use ignore
directives when they intentionally contain annotation syntax.

## Document mutation model

### Initialization

`InitDocuments` validates that the package path is project-relative and the
source package exists. Every nested source package is initialized separately
at its matching nested `docs/<package>/` path. Before writing, initialization
checks all four target files for self-describing metadata, legacy
markers/relationships, malformed generic frontmatter, and a package-local
central catalog.

Only after complete preflight does it create the documentation directory and
write each role. Missing files receive minimal identity plus detailed
authoring guidance. Existing plain narratives or generic frontmatter are
adopted without discarding their body or unrelated metadata.

Initialization is intentionally not idempotent after adoption: a second call
refuses to overwrite the now self-describing set. Ongoing edits belong to
authors and `docs fix`, not repeated initialization.

### Repair

`RepairDocuments` derives the package from `docs/<package>/` and the role from
the canonical basename. A file target parses and potentially replaces only
that file. A directory target prepares all four outputs before the write phase
and may add missing skeletons.

Existing bodies are passed through unchanged. Updated files are written to a
same-directory temporary file, synchronized, closed, and renamed over the
target while preserving permissions. Malformed frontmatter, semantic YAML,
legacy metadata, or a central catalog aborts repair rather than being
discarded.

Repair can normalize version and package identity proven by the path; it
cannot decide how to rewrite a requirement, contract, design decision,
purpose, or coverage relationship. It does not accept or emit
`idd.document`.

## Package Layout

```text
internal/collector/
├── idd_document.go              # identity, CommonMark parser, source indexes
├── idd_document_validation.go   # role, reference, lifecycle, graph rules
├── document_schema.go           # templates, slots, completion inspection
├── document_files.go            # init, repair, and atomic writes
├── doc_collector.go             # discovery and legacy routing
├── frontmatter.go               # legacy metadata parsing and validation
├── source_ast.go                # Tree-sitter registry and declarations
├── code_collector.go            # source evidence conversion
└── *_test.go                    # focused and compatibility evidence
```

Parsing, deterministic validation, completion inspection, and file mutation
share package-private types but remain in separate cohesive files. This keeps
the schema a single source of truth without concentrating every concern in one
oversized source file; collection remains read-only while init and repair own
the bounded write paths.

## Function Composition

1. `DocCollector.Collect → filepath.Walk → collectIDDDocumentSet/collectFile`
   discovers inputs and locks package mode.
2. `parseIDDDocument → splitLeadingFrontmatter → decodeIDDMetadata →
   parseIDDMarkdownRecords` builds one parsed document while retaining body and
   source indexes.
3. `collectIDDDocumentSet → validateIDDDocument →
   validateIDDDocumentReferences` evaluates role-local and cross-file
   invariants.
4. `addIDDDocumentIdentifiers` converts valid records and derives reverse SPEC
   backlinks from TEST coverage.
5. `validateIDDDocumentMarkdown` checks undeclared or legacy-style references
   in narrative Markdown without treating fenced examples as declarations.
6. `collectFile` preserves the legacy frontmatter and section parser for
   unmigrated packages.
7. `CodeCollector.CollectWithErrors → AnalyzeSource → bindSourceAnnotations`
   gathers trustworthy declaration-attached evidence and retains normalized
   analyses for engine policy.
8. Configured files without a pinned grammar and supported files with invalid
   syntax both return actionable `source-parse` findings; neither path falls
   back to regex binding.

Initialization and repair are separate from collection:
`InitDocuments → MarshalIDDDocument → writeExclusiveFile`, while
`RepairDocuments → normalizeIDDDocument → replaceFileAtomically`.

## Dependencies

- `internal/config` supplies ignore patterns and collection settings.
- `internal/model` supplies identifiers and validation errors.
- `pkg/pattern` validates identifier syntax.
- `gopkg.in/yaml.v3` parses and emits identity and legacy frontmatter.
- `github.com/yuin/goldmark` parses CommonMark records with source positions.
- `github.com/tree-sitter/go-tree-sitter` owns parser lifecycle and normalized
  syntax-tree access.
- Official Tree-sitter Go, JavaScript, TypeScript/TSX, C++, Java, and Python
  grammar modules define language syntax. Versions are pinned together at an
  ABI-compatible generation to preserve the repository's Go 1.22 baseline.

## Testability Hooks

- Pure IDD document parsing accepts a filename and byte slice, while marshaling
  accepts identity plus body bytes. Tests assert filename-owned roles, exact
  body preservation, minimal frontmatter shape, fenced-example isolation,
  wrapped fields, and source lines.
- Identity diagnostics use YAML node locations; semantic diagnostics use
  Markdown AST source positions. Focused cases verify offsets instead of only
  checking error text.
- File operations use temporary package trees to test traversal refusal,
  preflight-before-write, generic-frontmatter merge, permissions, atomic
  replacement, file-vs-directory scope, and idempotence.
- Table-driven document-set cases cover every role, reference direction,
  placeholder class, mode conflict, and TEST kind without constructing the
  graph manually.
- Self-describing and legacy modes share the same `IdentifierSet` output
  contract, allowing engine tests to remain independent from documentation
  syntax.
- Source tests use table-driven real-language fixtures to assert declaration
  kind and visibility, test classification, source locations, multiple
  references, ignore behavior, detached/body/string isolation, malformed
  syntax handling, and semantic annotation kind.

## Trade-offs and evolution

Supporting two documentation modes increases collector complexity, but explicit
package-level selection contains that cost and gives users a safe migration
path. The legacy path can be removed only through a separate compatibility
decision.

CommonMark AST parsing costs more than line matching but avoids false document
records inside examples. Tree-sitter adds CGO and binary-size cost, but makes
source binding trustworthy across seven language syntaxes and avoids silently
different regex heuristics. Both parsers extract only fields required for
traceability; adding a new semantic field should be justified by a graph or
validation need rather than a desire to structure every paragraph.

Future language support belongs in source collection. Future document roles or
record shapes affect the canonical filename set, discovery, reference
validation, templates, report metadata, the embedded Skill, and migration
rules together.
