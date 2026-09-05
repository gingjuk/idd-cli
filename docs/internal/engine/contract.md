---
idd:
  version: "1.1"
  package: internal/engine
  namespace: INTERNAL_ENGINE
---

# Contracts: internal/engine

## Contract: EngineLifecycle

**Guarantees:**

```go
func New(cfg *config.Config) *Engine
func (e *Engine) AddStructuralErrors(errors []*model.ValidationError)
func (e *Engine) SetSourceAnalyses(analyses []*collector.SourceAnalysis)
func (e *Engine) Run(ctx context.Context, ids *model.IdentifierSet) (*model.ValidationResult, error)
func (e *Engine) BuildReport() *model.Report
```

Construction retains a non-nil configuration pointer and initializes an empty
graph and result. Nil configuration or identifier inputs are unsupported and
may panic when dereferenced.

`AddStructuralErrors` ignores nil entries and copies every non-nil finding into
the engine result through `AddError`, preserving rule, message, source, link,
and code. These errors remain part of later validity decisions.

`SetSourceAnalyses` supplies the normalized syntax-tree declarations collected
for the same run. The engine uses those analyses for source annotation policy.
If library callers omit them, the engine may collect supported configured
sources itself; it never uses the removed line-oriented annotation validators.

`Run` mutates the engine and returns its owned result pointer. It is a
single-run lifecycle: state is not cleared before execution. The context is not
currently observed, and validation failures are represented in the result
rather than the Go error. The current implementation returns nil error.

If no errors remain after all rules, `Valid` becomes true. Warnings alone do
not make the result invalid. Findings are sorted by model ordering. Statistics
describe the built graph; the graph snapshot is attached only when configured.

`BuildReport` returns a value copy of the accumulated result with tool
`idd-cli`, version `1.0.0`, UTC RFC3339 timestamp, and collection-relevant
config summary. Nested maps, slices, and pointers are not deep-copied.

## Contract: GraphInterpretation

**Guarantees:**

Every observed identifier contributes origin and source metadata to one node
per internal graph key. First non-empty description/source values per origin
win. Document TEST kind and code annotation-kind presence remain
distinguishable. Self-describing Component and Contract names become
package-scoped derived nodes; authors never create DESIGN or CONTRACT marker
IDs for them.

The preferred relationship source is `Identifier.TypedLinks`. It preserves
explicit TEST coverage, named Contract evidence, Component dependencies, and
lifecycle replacement semantics without guessing from target prefixes. The
collector records these directed relationships:

- TEST `Covers` → SPEC `implements`, with derived SPEC `tests`;
- contract TEST `Contracts` → Contract `contract_tests`;
- Component `Depends on` → Component `depends_on`;
- replacement `Supersedes` → old record `supersedes`, while `Deprecated by`
  records the reverse direction when it is authored.

TEST coverage also contributes the derived SPEC → TEST `tests` backlink. For
other relationships, inbound graph queries expose the reverse traversal and
`ReverseLinkType` verifies reciprocal edges when both directions are present;
the engine does not manufacture a second authored relationship.

Legacy `Links` still use lexical inference for unmigrated documents:
SPEC-to-TEST is `tests`, SPEC-to-CONTRACT is `contract`, TEST-to-target is
`implements`, CONTRACT-to-target is `contract_implements`, and unknown
relationships are references. Edges retain the source observation's file and
line.

## Contract: RepositoryValidation

**Guarantees:**

The engine applies the enabled policies from `ValidationConfig` and also runs
the always-on compatibility checks present in `validate`.

The rule surface includes SPEC/TEST coverage, named and legacy contract-test
coverage, Component dependency target and cycle checks, required design
sections, legacy SPEC fields, self-describing TEST annotation kind, legacy
contract/design marker presence, documented implementation paths, legacy
relationship consistency, duplicate headings, orphan identifiers, doc/code
correspondence joined by identifier, syntax-tree public declaration and test
rules, legacy `related_files`, complete package doc sets for every scanned
source directory, annotation identifier/placement/duplication checks, source
parse failures, and duplicate IDs. Lexical or semantic prose comparison is not
part of the rule surface.
File-level `Spec`, `Contract`, and `Test` paths are not an input to this
contract.

Rules report all detected issues instead of stopping after the first. Errors and
warnings include the most precise source available. Legacy-only rules skip
self-describing package directories where required; the collector owns
self-describing Markdown schema, scaffold, lifecycle, concern, and local
reference validation.

### Filesystem and ignore behavior

Raw rule walkers use configured code and document patterns and silently skip
walk or read errors. Recursive `**` handling is implemented by walking the base
directory and matching filenames. Ignore matching combines basename glob
checks and path substring checks; some document-specific rules additionally use
`Docs.IgnorePaths`.

Source annotation rules are the exception to raw line walking. They consume
`SourceAnalysis` values built from the pinned language grammars. A configured
extension without a grammar or a syntax error in a supported source is a hard
`source-parse` finding and prevents policy decisions from partial trees; no
regex fallback is part of this contract.

The package-document-set rule also consumes those source analyses. Every
distinct containing directory maps from `<package>/` to
`docs/<package>/`, including nested sub-packages, and must contain
`design.md`, `contract.md`, `spec.md`, and `testing.md`. A parent package's
documents never satisfy a child package. Documentation directly at the docs
root remains project narrative rather than a package set.

The contract does not promise one unified ignore algorithm across every rule.
Changing scan or ignore semantics requires rule-specific regression tests.
