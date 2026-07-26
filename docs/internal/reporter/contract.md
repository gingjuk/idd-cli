---
idd:
  version: "1.0"
  package: internal/reporter
  document: contract
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
include the first finding's title and repair guidance.

Rule presentation metadata supplies known explanations and hints, including
syntax-tree parse/binding failures, scaffold completion, named Contract
coverage, and Component dependency cycles. Unknown rules receive a humanized
title and generic guidance. Enrichment does not change the engine's underlying
result.
