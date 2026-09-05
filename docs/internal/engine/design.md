---
idd:
  version: "1.1"
  package: internal/engine
  namespace: INTERNAL_ENGINE
---

# Design: internal/engine

## Component: ValidationEngine

**Purpose:**

`ValidationEngine` is the stateful coordinator between collected identifier
evidence and a repository-level validation result. It detects cross-package
duplicates, builds the linkage graph, applies configured graph and filesystem
rules, enriches the result with statistics and an optional graph snapshot, and
sorts findings for deterministic reporting.

Collection is deliberately outside the component. The CLI invokes
`DocCollector` and `CodeCollector`, merges their `IdentifierSet` values, passes
collector structural findings through `AddStructuralErrors`, supplies
Tree-sitter-backed `SourceAnalysis` values through `SetSourceAnalyses`, and only
then calls `Run`. The engine can rescan configured files for rules that depend
on raw Markdown layout, but source annotation placement and declaration policy
consume the normalized AST model.

**Ownership:**

The component owns repository-wide policy execution over collected entities,
occurrences, typed relations, AST analyses, configured documentation units,
and filesystem evidence, plus deterministic validity and statistics.

**Boundary:**

Collectors own parsing and provenance creation; reporters own presentation;
the CLI owns configuration and process behavior. The engine never authors or
repairs prose and cannot establish whether an intention is semantically sound.

**Decisions:**

Validation remains an explicit ordered pipeline with one shared TraceIndex.
Lifecycle gates correspondence and coverage from canonical SPECs outward, and
scan failures are findings rather than evidence of absence or success.

### Responsibilities

- accept a resolved configuration and pre-collected identifier set;
- preserve collector structural errors in the shared result;
- detect duplicate document and code IDs across package directories;
- create graph nodes for every observation and retain origin, description,
  TEST-kind, and source metadata;
- preserve typed links for Component dependencies, named Contract evidence,
  lifecycle replacement, and SPEC/TEST coverage;
- infer directed edge kinds only for legacy untyped links;
- verify reciprocal edge types without automatically requiring every edge to
  be reciprocal;
- execute hard-coded validation rules according to configuration;
- validate language-specific public and test declarations from normalized
  syntax-tree analyses;
- join source and document observations by SPEC or TEST identifier, without
  file-level document paths;
- map every scanned source directory, including nested sub-packages, to its
  equally nested `docs/<package>/` four-file document set;
- read configured files for remaining Markdown and package-document-set
  evidence;
- set validity from accumulated errors, sort findings, and attach statistics;
  and
- create a complete internal report envelope when requested.

### Non-responsibilities

There is no production `validator.ValidationRule` or `Rule` interface and no
runtime rule plugin registry. Validation order and implementations are explicit
methods in `engine.go` and `annotation_validation.go`. A `Rule` interface in
`engine_contract_test.go` is a test-local historical fixture and is not a
production extension point.

The engine does not author or repair documents, render final output, load YAML,
or decide CLI configuration search order. It does not judge whether rich
documentation prose is semantically adequate and performs no lexical or
embedding-based comparison. That judgment belongs to an explicit human or LLM
review using evidence outside the validation result.

### Lifecycle, state, and concurrency

`New` allocates one graph and one result and retains the configuration pointer.
`Run` mutates all three pieces of state and does not reset them. An engine
instance is therefore designed for one validation run. Reusing it can retain
nodes, edges, findings, and graph snapshots from the previous run.

The engine and its graph are not concurrency-safe. Validation is synchronous.
The `context.Context` parameter is currently reserved but not observed; a
cancelled context does not interrupt graph construction, filesystem scans, or
rule execution, and `Run` currently returns a nil operational error.

### Graph construction

All identifier observations are added as nodes before edges. Multiple
observations of one ID merge into one graph node while retaining the first
available source and description per origin. Document TEST kind is stored as a
scalar; code annotation kinds are stored as flags. Component and Contract
records use internal package-scoped graph keys while retaining their
human-readable names in metadata.

Typed links are authoritative for self-describing records. They represent
TEST-to-SPEC coverage, contract TEST-to-Contract evidence,
Component-to-Component dependency, and replacement direction explicitly.
Document collection also supplies the derived SPEC-to-TEST coverage backlink.
The engine materializes exactly those typed directions; only legacy forward
strings use type-based inference.

Dependency-cycle traversal follows `depends_on` edges only, so other graph
relationships cannot create false cycles. Contract coverage queries inbound
contract-test edges on the derived Contract node. Lifecycle consistency and
replacement cycles are checked while parsing the owning document; typed edges
make the normalized result visible to reports and later graph policy.

### Rule execution and failure containment

Most rule failures are accumulated into `ValidationResult`; they are not
returned as Go errors. File walkers skip unreadable paths and invalid matches,
so missing evidence typically becomes a validation finding only when a specific
rule detects the consequence.

Some rules are gated by configuration, while self-describing TEST-kind checks,
graph target/cycle checks, legacy contract/design marker checks, and documented
implementation-path checks run unconditionally. Self-describing package
parsing, scaffold, lifecycle, concern, and local reference errors are produced
by the collector and injected before `Run`; legacy-only engine rules explicitly
skip directories with `idd` metadata.

### Decisions and trade-offs

Hard-coded rule order makes interactions reviewable and gives one complete
result instead of fail-fast behavior. It also means adding a rule requires
editing the engine and tests rather than registering a plugin.

Several legacy and package-layout rules still rescan files after collectors
have read them. The source declaration rules no longer do: a shared normalized
analysis prevents Go-shaped regular expressions from being applied to other
languages. Remaining repeated I/O and ignore-path logic can be consolidated
later without weakening the AST binding contract.

## Architecture

```text
DocCollector ----\
                  +--> IdentifierSet + structural errors
CodeCollector ---/                 +--> SourceAnalysis[]
                                    |          |
                                    v          v
                         AddStructuralErrors
                                    |
                                    v
                               Engine.Run
                    / duplicate ID validation
                   /  graph build and verification
                  /   graph + filesystem rule suite
                 /    validity, sorting, stats, snapshot
                v
         ValidationResult ──> reporter
```

The graph is the relationship query layer. Normalized source analyses support
annotation and declaration rules and contribute the complete source-package
inventory. Raw document walkers remain for legacy field order, headings, and
existing four-file evidence, so a package with no documentation directory is
still reportable.

## Package Layout

`engine.go` owns lifecycle, graph construction, and the main validation
passes.
`source_validation.go` owns syntax-tree analysis loading and language-neutral
public, test, annotation, and placement policy.
`annotation_validation.go` owns duplicate-ID reporting and package-derived
rename suggestions.
`filesystem.go` owns configured source/document traversal and the optional
internal report envelope.

The split is organizational only; all files implement methods on the same
stateful `Engine`. There are no `validator`, `linker`, or separate rule
packages.

## Function Composition

`Run` performs duplicate detection, graph construction, statistics capture,
validation, and optional graph projection. `validate` invokes rule methods in a
fixed order, sets valid only when no errors remain, and sorts errors and
warnings.

`buildGraph` first merges observations into nodes, then adds every supplied
typed or legacy forward relationship. It verifies reciprocal types when both
directions exist but does not manufacture missing edges. Source rules consume
the supplied normalized analyses; remaining raw file rules call
`walkCodeFiles` or `walkDocFiles`, which dispatch configured patterns to
recursive or non-recursive standard-library walkers and skip ignored/unreadable
files.

`validatePkgDocFiles` derives the configured `docs` root, groups each source
analysis by its containing project-relative directory, and maps that directory
without collapsing sub-packages. It unions this inventory with discovered
documentation directories, then reports each missing canonical basename at
the path where the file belongs. Parent and child packages are exact path
matches rather than prefix matches.

## Dependencies

- `internal/config` supplies rule, scan, and graph-output policy.
- `internal/model` supplies evidence, findings, stats, and reports.
- `internal/graph` supplies relationship storage and traversal.
- `pkg/pattern` supplies identifier recognition and type parsing.
- standard-library filesystem, regex, context, string, path, and time packages
  support repository scans and report metadata.

## Testability Hooks

The engine accepts an in-memory `IdentifierSet`; tests can enable only relevant
flags and use temporary repositories for raw-file rules. Same-package tests
inspect graph/result state and invoke focused validators. Structural findings
can be injected without reproducing collector parsing.

The suite exercises typed and legacy graph interpretation, correspondence,
required fields, design content, TEST kind, seven-language declaration
binding, public/test policy, detached/body/string annotations, parse failures,
headerless production/test/contract-test sources, complete doc sets, wholly
absent nested-package doc sets, ignore scopes, duplicates, consistency
warnings, and report construction.
Important exclusions include actual context cancellation, engine reuse,
concurrent access, broad I/O failure propagation, every interaction among rule
flags, and complete semantic quality of documentation prose.
