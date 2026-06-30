---
markers:
  - id: SPEC-INTERNAL_REPORTER-001
    name: Reporter Structure
  - id: SPEC-INTERNAL_REPORTER-003
    name: Reporter.New
  - id: SPEC-INTERNAL_REPORTER-004
    name: Reporter.Generate
  - id: SPEC-INTERNAL_REPORTER-005
    name: Reporter.Write
  - id: SPEC-INTERNAL_REPORTER-011
    name: LLM Report Output
  - id: SPEC-INTERNAL_REPORTER-012
    name: LLM Finding Enrichment

related_files:
  spec: docs/internal/reporter/spec.md
  contract: docs/internal/reporter/contract.md
  design: docs/internal/reporter/design.md
  testing: docs/internal/reporter/testing.md

---

# Specification (reporter)

## SPEC-INTERNAL_REPORTER-001: Reporter Structure

**Design:** `ReporterModule`

**Contract:** `Reporter`

**Requirement:**

Reporter generates validation reports in LLM-oriented JSON and Markdown formats.

**Tests:** `TEST-INTERNAL_REPORTER-001`, `TEST-INTERNAL_REPORTER-002`, `TEST-INTERNAL_REPORTER-003`, `TEST-INTERNAL_REPORTER-004`, `TEST-INTERNAL_REPORTER-005`, `TEST-INTERNAL_REPORTER-006`, `TEST-INTERNAL_REPORTER-007`, `TEST-INTERNAL_REPORTER-008`, `TEST-INTERNAL_REPORTER-009`, `TEST-INTERNAL_REPORTER-010`, `TEST-INTERNAL_REPORTER-011`, `TEST-INTERNAL_REPORTER-012`, `TEST-INTERNAL_REPORTER-013`

**Status:** Done

**Implementation:** `internal/reporter/reporter.go`

**Key Types:**

- `Reporter` — Report generator with config and format

**Acceptance Criteria:**

- [x] Reporter stores config and format
- [x] Reporter.Generate creates complete report
- [x] Reporter.Write outputs report to stdout or file
- [x] Reporter JSON output uses a finding-centered LLM schema

## SPEC-INTERNAL_REPORTER-003: Reporter.New

**Design:** `ReporterModule`

**Contract:** `New`

**Requirement:**

Create a new Reporter with the given configuration and output format, defaulting to JSON when format is empty.

**Tests:** `TEST-INTERNAL_REPORTER-001`

**Function Signature:**
`func New(cfg *config.Config, format string) *Reporter`

**Purpose:** Creates a new Reporter with the given configuration and output format. Defaults to JSON if format is empty. JSON uses the LLM report schema.

**Parameters:**

- `cfg`: Configuration pointer
- `format`: Output format ("json", "markdown", or "llm-markdown")

**Returns:** A new Reporter instance

## SPEC-INTERNAL_REPORTER-004: Reporter.Generate

**Design:** `ReporterModule`

**Contract:** `Generate`

**Requirement:**

Generate a complete report from a validation result, including tool metadata, config summary, and the validation result.

**Tests:** `TEST-INTERNAL_REPORTER-002`

**Function Signature:**
`func (r *Reporter) Generate(result *model.ValidationResult) (*model.Report, error)`

**Purpose:** Generates a complete report from a validation result, including tool metadata, config summary, and the validation result.

**Parameters:**

- `result`: The validation result to include in the report

**Returns:** Complete Report structure or error

## SPEC-INTERNAL_REPORTER-005: Reporter.Write

**Design:** `ReporterModule`

**Contract:** `Write`

**Requirement:**

Write the report to stdout or the specified output file path.

**Tests:** `TEST-INTERNAL_REPORTER-001`, `TEST-INTERNAL_REPORTER-002`, `TEST-INTERNAL_REPORTER-003`, `TEST-INTERNAL_REPORTER-004`, `TEST-INTERNAL_REPORTER-005`, `TEST-INTERNAL_REPORTER-006`, `TEST-INTERNAL_REPORTER-007`, `TEST-INTERNAL_REPORTER-008`, `TEST-INTERNAL_REPORTER-009`, `TEST-INTERNAL_REPORTER-010`, `TEST-INTERNAL_REPORTER-011`, `TEST-INTERNAL_REPORTER-012`, `TEST-INTERNAL_REPORTER-013`

**Function Signature:**
`func (r *Reporter) Write(report *model.Report, output string) error`

**Purpose:** Writes the report to the specified output destination. If output is empty or "-", writes to stdout. Otherwise creates a file at the given path.

**Parameters:**

- `report`: The report to write
- `output`: File path or "-" for stdout

**Returns:** Error if writing fails

**Acceptance Criteria:**

- [x] Generate creates report with tool info, config, and result
- [x] Write outputs to stdout for "-" or empty output
- [x] Write creates file for non-empty non-dash output path
- [x] Write supports JSON LLM reports and LLM Markdown reports

**Related:** `SPEC-INTERNAL_REPORTER-001`

## SPEC-INTERNAL_REPORTER-011: LLM Report Output

**Design:** `ReporterModule`

**Contract:** `Reporter`

**Requirement:**

Reporter must support JSON and LLM Markdown formats that transform the validation report into finding-centered output for LLM analysis and automated repair.

**Tests:** `TEST-INTERNAL_REPORTER-011`, `TEST-INTERNAL_REPORTER-012`

**Acceptance Criteria:**

- [x] `json` emits schema `idd.llm_report.v1`
- [x] `json` includes status, summary, finding groups, and findings
- [x] `llm-markdown` includes status, prioritized findings, locations, problems, fixes, and related identifiers

**Related:** `SPEC-INTERNAL_REPORTER-001`, `SPEC-INTERNAL_REPORTER-005`

## SPEC-INTERNAL_REPORTER-012: LLM Finding Enrichment

**Design:** `ReporterModule`

**Contract:** `Reporter`

**Requirement:**

Each LLM finding must be self-contained enough for an LLM to diagnose and repair the issue without manually joining global report sections.

**Tests:** `TEST-INTERNAL_REPORTER-013`

**Acceptance Criteria:**

- [x] Each finding contains severity, rule, title, problem, suggested fix, and structured location when available
- [x] Repeated findings are grouped by severity and rule in the summary
- [x] Rule metadata supplies default severity, explanation, and fix hints
- [x] Source strings in `path` or `path:line` form are normalized into structured locations
- [x] Link identifiers and graph edges are exposed as related identifiers when available
