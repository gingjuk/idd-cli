---
related_files:
  design: docs/architecture/design.md
---

# IDD CLI architecture

## Design goals

The architecture is shaped by four goals that must hold together:

1. **Human comprehension:** IDD documents explain intent, behavior, decisions,
   boundaries, failure modes, and evidence; they are not merely indexes into
   source code.
2. **Single ownership:** each declaration and relationship has one authored
   home, so adding a SPEC or changing coverage does not rewrite several files.
3. **Deterministic enforcement:** idd-cli can locate every structural fact at a
   stable Markdown or source line and produce a complete repair report.
4. **Safe adoption:** existing legacy packages continue to validate until they
   are deliberately migrated, and structural repair never guesses semantics.

These goals distinguish structural economy from textual brevity. Minimal
identity YAML and derived backlinks remove coordination overhead; detailed
Markdown remains the durable knowledge layer.

## Architecture

The system has an authoring plane, a structural-maintenance boundary, and a
validation/reporting plane.

### Authoring plane

The embedded IDD Skill tells an agent how to turn approved intent into design,
contracts, cohesive behavioral SPECs, TEST evidence, and code annotations. It
also owns the semantic review: whether a design is understandable, a boundary
is correct, or a requirement is complete cannot be inferred from field
presence alone.

The approved request itself is sufficient intent. A project may preserve an
`intent.md` as discussion history, but neither the document schema nor the
graph requires an Intent source or Intent-to-SPEC relationship.

The Skill is embedded in the binary so the instructions can evolve with the
parser and validator. `generate skill` exports that exact copy; regenerating it
after an idd-cli upgrade prevents an agent from writing an older format against
a newer validator.

### Structural-maintenance boundary

`docs init` and `docs fix` are intentionally narrow mutators, while
`docs status` is their read-only completion observer:

- initialization creates minimal identity plus stable scaffold markers and
  honest authoring guidance for a new package, or preserves an existing plain
  narrative while adopting it;
- status accepts one or more targets, applies the same shared role schema, and
  emits deduplicated file, line, role, slot, and reason work items;
- repair normalizes document identity and can create missing skeletons for a
  directory target;
- initialization and repair preflight all requested targets before writes;
- neither command invents a component, contract, requirement, purpose,
  coverage relationship, architectural decision, or explanatory paragraph;
- existing Markdown bodies remain byte-preserved during repair.

This boundary makes automation safe: a tool may repair facts derived from the
path, but only the Skill and human intent may decide what the system means.
Removing a marker is not sufficient completion: its bounded Markdown must also
contain effective non-placeholder authored content.

### Validation and reporting plane

idd-cli follows a graph-first validation architecture with two inputs:

1. `DocCollector` reads complete self-describing document sets or isolated
   legacy Markdown metadata.
2. `CodeCollector` uses pinned Tree-sitter grammars to bind annotations to real
   Go, TypeScript/TSX, JavaScript/JSX, C++, Java, and Python declarations.
3. `Engine.buildGraph` merges the two origins into a `LinkageGraph`, preserves
   explicit typed relationships, including the SPEC/TEST coverage backlink
   already derived by document collection.
4. Concrete validation methods check structural findings, graph integrity,
   document/code correspondence, and package conventions.
5. `Reporter` converts one validation result into human Markdown, repair-focused
   LLM Markdown, or stable JSON.

There is no separate Linker or Validator interface in the current
implementation. Graph construction and rules live in `internal/engine`; the
architecture documentation names the concrete ownership rather than the
interfaces once proposed by an older README.

```text
Embedded IDD skill ── generate skill ──→ Agent authoring workflow
                                             │
docs init/status/fix ── work list + safe structure ─┤
                                             ▼
                         Markdown records + code annotations
                                  │                  │
                                  ▼                  ▼
                            DocCollector       CodeCollector
                                  └─────────┬────────┘
                                            ▼
                                        Engine.Run
                                            │
                                            ▼
                                         Reporter
                                  llm-markdown │ JSON
                                               │
                           Agent repair loop ←──┘──→ CI gate
```

The final validity check normally runs as `idd-cli run .` from project root.
The optional positional path denotes another complete project root:
configuration discovery, `DocCollector`, `CodeCollector`, and engine filesystem
checks all execute against that same tree. Focused documentation inspection
belongs to `docs status` or `docs review-context --docs-path`, not `run`.

Passing this plane proves that the declared graph is structurally coherent. It
does not prove that the prose is sufficient or the architectural decision is
sound; that remains an explicit authoring/review responsibility.

## Document ownership model

Each scanned source package directory, including every nested sub-package,
maps to an equally nested `docs/<package>/` directory and uses four fixed
Markdown files because their concerns evolve at different rates:

- `design.md` owns named components and the reasons responsibilities are
  divided as they are;
- `contract.md` owns observable boundaries, errors, invariants, side effects,
  ownership rules, and compatibility promises;
- `spec.md` owns cohesive behavioral requirements and their selected design
  and contract;
- `testing.md` owns TEST evidence and the TEST-to-SPEC `Covers` relationship.

The exact lowercase basename is the sole role authority and eliminates both
repeated `related_files` and a duplicate `idd.document` value. Minimal
frontmatter identifies only version and package; the former role field is
invalid without a compatibility period. TEST coverage is authored once in
`testing.md`; the engine derives the reverse edge rather than requiring a
second backlink in `spec.md`.

Small required prose anchors make generated omissions mechanically visible:
Component `Purpose`, Contract `Guarantees`, SPEC `Requirement` and
`Acceptance`, and TEST `Purpose` and `Oracle`. Contract TEST records own
`Contracts`, the evidence relationship to named Contract guarantees.
Components may own `Depends on`; optional lifecycle and concern fields remain
inside the record they describe. The collector normalizes these facts into
typed graph links and checks local structure before the engine performs
cross-package target, coverage, and dependency-cycle policy.

The design deliberately permits unrestricted Markdown within each record.
Subordinate headings, diagrams, examples, rationale, edge cases, and operational
notes remain readable to people while the level-two record boundary and small
fixed-field block remain deterministic for tools.

There is no document line-count limit. A long role remains one canonical file;
the collector rejects size- or feature-derived names such as
`design-auth.md`, `contract_part.md`, `spec.api.md`, and
`testing-extra.md`. Package granularity, rather than file splitting, bounds the
documentation set.

## Validation lifecycle

One `run` invocation has the following lifecycle:

1. Resolve configuration from the explicit flag, project defaults, or
   in-memory defaults, then apply invocation-only output overrides.
2. Discover documentation and lock each package into either self-describing or
   legacy mode. A partial self-describing set is an error rather than an
   invitation to mix models.
3. Parse document identity and role-owned CommonMark records, retaining exact
   source lines for fields and findings.
4. Parse configured supported sources with Tree-sitter and bind `@implement`,
   `@test`, and `@test-contract` only from real comments adjacent to normalized
   declarations. Syntax errors become findings and never trigger regex
   fallback.
5. Merge identifiers by ID and origin. Matching document and source
   observations establish code-to-document correspondence directly; source
   files do not repeat `Spec`, `Contract`, or `Test` paths.
6. Construct the explicitly typed or legacy directed links supplied by
   collection, and verify reciprocal types where both directions are present.
7. Add collector findings, run graph and repository validation rules, sort the
   result deterministically, and build optional graph statistics.
8. Write the complete selected report. Validation failure produces a non-zero
   exit only after the report is available; traversal or I/O failure returns an
   operational error.

This ordering separates malformed evidence from unavailable evidence. A bad
record becomes an actionable finding while unrelated packages continue to be
collected; a filesystem failure that prevents trustworthy collection aborts
the command.

## Mutation and failure boundaries

Document initialization performs package, path-traversal, legacy-metadata, and
central-catalog preflight checks before its first write. Repair parses every
selected target before replacement and uses same-directory temporary files plus
rename for atomic file updates.

A file-targeted repair has one possible write target. A directory-targeted
repair owns the four-file structural set and may create missing skeletons.
Neither path silently migrates legacy relationships, because doing so would
require semantic choices about canonical ownership. These document commands
never rewrite source files; code association is derived during validation from
declaration identifiers rather than inserted document-path headers.

Report formats share one validation result. Verbose diagnostics use stderr so
JSON stdout stays machine-readable. LLM Markdown may recommend an owner and
repair action, but the reporter never mutates the reported file.

## Package Layout

```text
idd-cli/
├── cmd/idd-cli/       # CLI entry point
├── internal/
│   ├── collector/      # document modes, source annotations, safe file operations
│   ├── config/         # configuration loading and normalization
│   ├── engine/         # graph construction and concrete validation rules
│   ├── graph/          # directed traceability data structure
│   ├── model/          # cross-stage identifiers, findings, and report values
│   └── reporter/       # validation reports and review-context projections
└── pkg/
    ├── pattern/        # identifier and annotation syntax
    └── walk/           # reusable filesystem traversal
```

`cmd/idd-cli` remains a composition root. Parsing and validation decisions stay
in internal packages so command handlers do not become a second source of
domain rules.

## Function Composition

1. **Skill commands:** `ListEmbeddedSkills` and `ReadEmbeddedSkill` expose the
   binary-owned authoring instructions; `generateSkill` copies them without
   transformation.
2. **Document commands:** `InitDocumentPackages` creates one or more structural
   sets; `InspectDocumentCompletions` exposes deduplicated remaining scaffold
   work; `RepairDocumentTargets` normalizes derived identity while preserving
   bodies and markers; `BuildSpecReviewContexts` gathers ordered evidence-only
   semantic review bundles through one shared scan.
3. **Collection:** `DocCollector.Collect` selects self-describing or legacy
   mode per package; `CodeCollector.CollectWithErrors` independently produces
   source evidence and normalized language analyses.
4. **Graph and validation:** `Engine.SetSourceAnalyses` and `Engine.Run` merge
   origins, derive typed links, apply concrete checks, and retain structural
   findings collected earlier.
5. **Reporting:** `Reporter.Write` projects the result for people, agents, or
   automation and selects stdout or an explicit output file.

The composition is intentionally one-way. Reporters do not collect, the engine
does not author documents, and document repair does not infer requirements.

## Testability Hooks

- Pure document parsing and marshaling accept byte slices, allowing exact
  frontmatter, CommonMark, and body-preservation cases without repository I/O.
- File operations use temporary package trees to verify traversal refusal,
  preflight atomicity, one-file write scope, deterministic scaffold status,
  marker preservation, and idempotence.
- Graph tests pre-populate identifiers and edges directly so validation
  semantics do not depend on collector fixtures.
- Source and engine tests use table-driven syntax fixtures for every supported
  language, including visibility, test classification, detached/body/string
  comments, exact ignore directives, and syntax failures.
- Engine graph tests combine focused typed relationships with temporary
  repository files for rules that genuinely depend on package layout.
- Reporter tests write to buffers or temporary outputs and assert stable
  schemas, grouped findings, locations, and stderr/stdout isolation.
- The repository-level acceptance gate builds the real binary and validates
  the repository itself, closing the gap between package tests and the shipped
  workflow.

## Dependencies

- `internal/collector` - For identifier collection
- `internal/engine` - For orchestration
- `internal/graph` - For graph structure
- `internal/model` - For data types
- `internal/reporter` - For validation and review-context output
- `pkg/pattern` - For IDD patterns
- `pkg/walk` - For file traversal
- `gopkg.in/yaml.v3` - For config and document identity parsing
- `github.com/yuin/goldmark` - For CommonMark AST parsing and source positions
- `github.com/tree-sitter/go-tree-sitter` plus pinned official grammars - For
  declaration-aware Go, TypeScript/TSX, JavaScript/JSX, C++, Java, and Python
  binding
- `github.com/spf13/cobra` - For CLI

Goldmark is used because line-oriented regular expressions cannot reliably
distinguish real headings and fields from fenced examples. `yaml.v3` nodes are
retained where exact identity-field locations or merge-preserving writes are
required. Tree-sitter provides the corresponding source-language boundary:
strings and body comments are distinguishable from declaration comments, and a
malformed file can be rejected rather than guessed through a Go-shaped regex.

## Trade-offs and non-goals

- The system accepts richer prose without trying to score completeness by word
  count. Mechanical length thresholds would reward filler and still miss
  absent decisions.
- The CLI is not an API-documentation generator. Exact symbol inventories and
  signatures belong in source or generated references; IDD records explain the
  behavior and design those symbols realize.
- Legacy compatibility adds two parsing paths, but package-level mode
  isolation prevents a half-migration from producing competing truth.
- Tree-sitter increases binary size and requires CGO, but one pinned,
  language-aware binding model is safer than seven drifting regular-expression
  heuristics.
- Concrete validation methods are less extensible than a plugin rule
  interface, but they match the current implementation and keep rule ordering,
  configuration, and result semantics explicit.
