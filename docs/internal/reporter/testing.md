---
idd:
  version: "1.1"
  package: internal/reporter
  namespace: INTERNAL_REPORTER
---

# Testing: internal/reporter

## TEST-INTERNAL_REPORTER-001: Complete report generation

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_REPORTER-001`, `SPEC-INTERNAL_REPORTER-004`
- **Contracts:** `ReporterLifecycle`, `cmd/idd-cli#Reporter`

**Purpose:**

Prove tool name, hard-coded version, and non-empty timestamp in a generated
report from a valid result and default configuration.

**Oracle:** The test passes only when its assertions confirm tool name,
hard-coded version, and non-empty timestamp in a generated report from a valid result
and default configuration.

## TEST-INTERNAL_REPORTER-002: JSON file output

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-005`, `SPEC-INTERNAL_REPORTER-011`

**Purpose:**

Prove that writing JSON creates a decodable finding report with the expected
schema and pass status in an isolated temporary directory.

**Oracle:** The test passes only when its assertions confirm writing
JSON creates a decodable finding report with the expected schema and pass status in an
isolated temporary directory.

## TEST-INTERNAL_REPORTER-003: Stdout destination

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-005`

**Purpose:**

Prove that `-` selects stdout and returns no error for a valid JSON report. The
test checks the call result rather than capturing process output.

**Oracle:** The test passes only when its assertions confirm `-`
selects stdout and returns no error for a valid JSON report. The test checks the call
result rather than capturing process output.

## TEST-INTERNAL_REPORTER-004: Human Markdown file

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-005`

**Purpose:**

Prove that human Markdown contains the linkage-report heading and a supplied
warning after a file round trip.

**Oracle:** The test passes only when its assertions confirm human
Markdown contains the linkage-report heading and a supplied warning after a file round
trip.

## TEST-INTERNAL_REPORTER-005: Unsupported format rejection

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_REPORTER-005`
- **Contracts:** `ReportDestination`, `HumanMarkdownReport`

**Purpose:**

Prove that an unsupported format returns an error. The `/dev/null` fixture
isolates format dispatch from ordinary file-content assertions.

**Oracle:** The test passes only when its assertions confirm an
unsupported format returns an error. The `/dev/null` fixture isolates format dispatch
from ordinary file-content assertions.

## TEST-INTERNAL_REPORTER-006: Human Markdown error details

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-005`

**Purpose:**

Prove that an error's rule and message appear in human Markdown rendered to an
in-memory writer.

**Oracle:** The test passes only when its assertions confirm an
error's rule and message appear in human Markdown rendered to an in-memory writer.

## TEST-INTERNAL_REPORTER-007: Human Markdown statistics

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-005`

**Purpose:**

Prove that populated validation statistics produce the expected metrics
section.

**Oracle:** The test passes only when its assertions confirm
populated validation statistics produce the expected metrics section.

## TEST-INTERNAL_REPORTER-008: Human Markdown graph evidence

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-005`

**Purpose:**

Prove that a supplied graph snapshot renders identifier nodes and relationship
details.

**Oracle:** The test passes only when its assertions confirm a
supplied graph snapshot renders identifier nodes and relationship details.

## TEST-INTERNAL_REPORTER-009: Reporter construction and default format

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-003`

**Purpose:**

Prove explicit formats are retained, empty format becomes JSON, and original
format casing is preserved until dispatch.

**Oracle:** The test passes only when its assertions confirm explicit
formats are retained, empty format becomes JSON, and original format casing is preserved
until dispatch.

## TEST-INTERNAL_REPORTER-010: Human status symbols

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-005`

**Purpose:**

Prove exact PASS and FAIL status strings for valid and invalid results.

**Oracle:** The test passes only when its assertions confirm exact PASS
and FAIL status strings for valid and invalid results.

## TEST-INTERNAL_REPORTER-011: Default and explicit JSON finding schema

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_REPORTER-003`, `SPEC-INTERNAL_REPORTER-011`
- **Contracts:** `FindingReport`

**Purpose:**

Prove empty-format and explicit JSON outputs both use
`idd.llm_report.v1`, and that a failing sample preserves status, error count,
structured location, primary identifier, and repair guidance.

**Oracle:** The test passes only when its assertions confirm empty-format
and explicit JSON outputs both use `idd.llm_report.v1`, and that a failing sample
preserves status, error count, structured location, primary identifier, and repair
guidance.

## TEST-INTERNAL_REPORTER-012: Agent-readable Markdown

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-011`

**Purpose:**

Prove that LLM Markdown includes failure status, grouping, exact file/line,
rule, problem, fix, and related identifier context.

**Oracle:** The test passes only when its assertions confirm LLM
Markdown includes failure status, grouping, exact file/line, rule, problem, fix, and
related identifier context.

## TEST-INTERNAL_REPORTER-013: Finding metadata, severity, and grouping

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-012`

**Purpose:**

Prove top-rule selection, one-based group indexes, expected/actual enrichment,
graph-related identifiers, warning-severity preservation, repeated-finding
aggregation, sorted file/identifier summaries, and curated guidance for every
self-describing document rule plus deprecated configuration. For split-role filenames, prove the detailed
finding contains a complete agent merge prompt while aggregate rule guidance
does not select only the first file.

**Oracle:** The test passes only when its assertions confirm top-rule
selection, one-based group indexes, expected/actual enrichment, graph-related
identifiers, warning-severity preservation, repeated-finding aggregation, sorted
file/identifier summaries, and curated guidance for every self-describing
document rule. A split filename finding must name its exact source and
canonical target, demand preservation of unique semantic content and
implementation boundaries, prohibit summary-only reduction, delay deletion
until a no-loss review, and include package `docs status` plus project `run`
commands. The group-level fix must contain neither one split source path nor
another.

## TEST-INTERNAL_REPORTER-014: Review-context formats

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_REPORTER-013`
- **Contracts:** `ReviewContextReport`

**Purpose:** Verify JSON and Markdown review-context output preserve the same
single-SPEC evidence, preserve ordered multi-SPEC batch evidence, and never add
semantic verdict language.

**Oracle:** Table-driven format cases preserve schema, SPEC, TEST, Contract,
declaration, and record-or-declaration truncation fields. One context retains
the single schema; multiple contexts use the batch schema and independent
Markdown sections; unsupported formats return an error.

## TEST-INTERNAL_REPORTER-016: Trace and review projection formats

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_REPORTER-013`

**Purpose:** Verify trace and review-context v2 JSON/LLM Markdown preserve the
same ID-centered evidence without changing lifecycle, resolution, provenance,
truncation, or semantic-review boundaries.

**Oracle:** Table-driven formats pass only when trace emits schema
`idd.trace.v1`, review-context emits its v2 single or batch schema, Components
and multiple Contracts remain visible, lists are stable, Markdown contains no
approval or semantic score, and unsupported formats return explicit errors.

## Strategy

Tests use immutable model fixtures, `strings.Builder`, and temporary output
files. Exact schema fields and selected human text are the oracles; no snapshots
of entire reports make harmless prose changes artificially expensive.

The suite excludes file permission failures, parent-directory creation,
unsupported-format truncation, close/write failures, nil inputs, unknown-rule
fallbacks, Windows source paths, top-rule truncation beyond five, and group tie
ordering. It also does not validate whether a suggested fix is semantically
correct for a real project.
