---
idd:
  version: "1.0"
  package: internal/reporter
  document: spec
---

# Specifications: internal/reporter

## SPEC-INTERNAL_REPORTER-001: Multi-audience report rendering

- **Design:** `ReportRenderer`
- **Contract:** `ReporterLifecycle`

**Requirement:**

One reporter instance must render the same validation evidence for automation,
human inspection, and agent repair without rerunning or changing validation.

### Boundaries

The reporter owns presentation and destination side effects only. It must not
mutate validation findings, decide validity, or claim that repair guidance is
itself a fix.

### Acceptance evidence

**Acceptance:**

The suite generates a report once and exercises JSON, human Markdown, LLM
Markdown, stdout, file output, statistics, graph output, and finding enrichment.

## SPEC-INTERNAL_REPORTER-003: Format-aware reporter construction

- **Design:** `ReportRenderer`
- **Contract:** `ReporterLifecycle`

**Requirement:**

Construction must retain the supplied configuration and format, defaulting only
an empty format to JSON. Case normalization must be deferred to output
dispatch.

**Acceptance:** Construction must retain
the supplied configuration and format, defaulting only an empty format to JSON. Case
normalization must be deferred to output dispatch.

### Edge cases

A nil configuration is accepted by construction but will panic when generation
dereferences it; callers are responsible for supplying the resolved config.

## SPEC-INTERNAL_REPORTER-004: Complete internal report envelope

- **Design:** `ReportRenderer`
- **Contract:** `ReporterLifecycle`

**Requirement:**

Generation must combine tool identity, hard-coded reporter version, an RFC3339
timestamp, collection-relevant configuration summary, and the supplied
validation result into one report.

**Acceptance:** Generation must combine
tool identity, hard-coded reporter version, an RFC3339 timestamp, collection-relevant
configuration summary, and the supplied validation result into one report.

### Ownership

The result is copied by value, but nested reference values are shared. The
operation does not serialize or write and currently has no expected error path.

## SPEC-INTERNAL_REPORTER-005: Destination and human output behavior

- **Design:** `ReportRenderer`
- **Contract:** `ReportDestination`

**Requirement:**

Writing must select stdout or create the requested file, dispatch
case-insensitively among supported formats, return creation/rendering errors,
and produce a readable human Markdown report with status, findings, statistics,
and optional graph evidence.

### Failure and side-effect boundary

Unsupported formats return an error. With a file destination, creation and
truncation happen before that error is detected. Parent directories and atomic
writes are outside the specification.

### Acceptance evidence

**Acceptance:**

Tests cover stdout success, valid JSON files, human Markdown warnings/errors,
stats and graph content, pass/fail icons, and unsupported format rejection.

## SPEC-INTERNAL_REPORTER-011: Finding-centered JSON and LLM Markdown

- **Design:** `ReportRenderer`
- **Contract:** `FindingReport`

**Requirement:**

JSON and LLM Markdown must expose schema status, summary counts and groups, and
actionable findings from one shared projection. Empty-format output must use
the same JSON schema.

**Acceptance:** JSON and LLM Markdown
must expose schema status, summary counts and groups, and actionable findings from one
shared projection. Empty-format output must use the same JSON schema.

### Required behavior

JSON must be indented and machine-decodable as `model.LLMReport`. LLM Markdown
must show status, finding groups, numbered locations, rule, severity, problem,
optional expected/actual data, fix guidance, and related identifiers. A result
with no findings must explicitly say so.

## SPEC-INTERNAL_REPORTER-012: Self-contained finding enrichment

- **Design:** `ReportRenderer`
- **Contract:** `FindingReport`

**Requirement:**

Raw validation errors must be enriched with enough local context for a repair
agent to act without joining unrelated report sections.

**Acceptance:** Raw validation errors
must be enriched with enough local context for a repair agent to act without joining
unrelated report sections.

### Required behavior

- known rules receive curated title, explanation, and fix guidance;
- unknown rules receive deterministic fallback text;
- `path:line` evidence becomes a structured location;
- primary and related identifiers are extracted without duplicating the
  primary;
- graph-adjacent identifiers are included when a snapshot is available;
- error and warning severity is preserved;
- repeated severity/rule findings are grouped with sorted file and identifier
  summaries; and
- top rule ranking is deterministic and capped.

### Boundary

Enrichment is lexical and graph-adjacent, not causal analysis. Suggested fixes
remain guidance and may require Skill-guided semantic judgment.
