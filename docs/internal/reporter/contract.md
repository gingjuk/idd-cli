---
idd:
  version: "1.1"
  package: internal/reporter
  namespace: INTERNAL_REPORTER
---

# Contracts: internal/reporter

## Contract: ReporterLifecycle

**Guarantees:**

```go
func New(cfg *config.Config, format string) *Reporter
func (r *Reporter) Generate(result *model.ValidationResult) (*model.Report, error)
```

Construction retains the configuration pointer and requested format. Only an
empty format is replaced with `json`; casing is otherwise preserved and
interpreted case-insensitively during writing.

Generation requires non-nil configuration and result values. It returns a
complete report with tool `idd-cli`, version `1.0.0`, current RFC3339
timestamp, selected config patterns and annotations, and a value copy of the
validation result. It currently returns no operational error. Nested maps,
slices, graph pointers, and identifier pointers are not deep-copied.

## Contract: ReportDestination

**Guarantees:**

```go
func (r *Reporter) Write(report *model.Report, output string) error
```

Empty output and `-` select stdout. Any other value is passed to `os.Create`,
which creates or truncates the file; parent directories must already exist.
File creation errors are wrapped. The file is closed on return and close errors
are ignored.

Supported formats, case-insensitively, are `json`, `markdown`/`md`, and
`llm-markdown`/`llm-md`. An unsupported format returns an error after
destination selection, so a file destination may already have been truncated.
Write does not use the config's output file field and does not perform atomic
replacement.

## Contract: HumanMarkdownReport

**Guarantees:**

Human Markdown includes tool/version/timestamp, pass/fail status, errors,
warnings, statistics, and—when present—graph nodes and edges. Status is derived
only from `ValidationResult.Valid`; warnings may coexist with PASS.

Finding order follows the input slices and graph edge order follows the
snapshot. The renderer returns the writer's error and does not escape arbitrary
finding text as a separate security boundary.

## Contract: FindingReport

**Guarantees:**

JSON and LLM Markdown share a finding-centered projection with schema
`idd.llm_report.v1`. The status is `pass` when the result is valid and `fail`
otherwise. Errors precede warnings and retain their respective severity.

Each finding contains rule, human title, problem, suggested fix, and structured
location when available. Source values ending in `:<integer>` are split at the
last colon; otherwise the full source is the file. Identifier detection prefers
the raw link field and then scans message, code, and source.

Top rules combine errors and warnings, rank by descending frequency and
ascending name, and retain at most five. Rule groups share severity and rule,
use one-based finding indexes, deduplicate and sort files/identifiers, and
include the first finding's title plus package-neutral rule guidance. A group
never embeds the first finding's path-specific instruction because one group
may span several packages.

Rule presentation metadata supplies known explanations and hints, including
syntax-tree parse/binding failures, scaffold completion, named Contract
coverage, Component dependency cycles, and non-canonical or split IDD role
filenames. `deprecated-config` guidance removes the named obsolete YAML key and
points semantic review to review-context. Filename guidance names the four
canonical files and directs content back to the owning file rather than
recommending another fragment.

For a `split-role` filename finding, the individual `suggested_fix` is a
self-contained agent repair prompt. It derives the canonical target beside the
reported source, requires both documents to be read, preserves every unique
still-valid requirement, behavior, rationale, contract, implementation
boundary, identifier relationship, example, diagram, and test-evidence note,
and forbids summary-only replacement. The prompt permits removal of the split
source only after a no-loss comparison and ends with package `docs status` and
project `run` commands. A missing canonical file may be structurally created
with `docs fix`, but that command is never presented as the semantic merge.

Unknown rules receive a humanized title and generic guidance. Enrichment does
not change the engine's underlying result.

## Contract: ReviewContextReport

**Guarantees:**

Review-context rendering serializes the collector-owned
`idd.spec_review_context.v2` value without inventing findings or verdicts.
JSON preserves the schema exactly. Markdown and LLM Markdown present the same
SPEC, Contract, TEST, declaration, issue, and truncation evidence with neutral
questions for a human or model reviewer.

When more than one context is requested, JSON uses
`idd.spec_review_context_batch.v2` with ordered `requested_spec_ids` and
`contexts`. Markdown uses one batch heading and an independent nested review
section for each SPEC. A single context retains the existing single-SPEC JSON
and Markdown shape.
