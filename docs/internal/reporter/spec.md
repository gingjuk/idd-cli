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

related_files:
  spec: docs/internal/reporter/spec.md
  contract: docs/internal/reporter/contract.md
  design: docs/internal/reporter/design.md
  testing: docs/internal/reporter/testing.md
---

# Specification (reporter)

## SPEC-INTERNAL_REPORTER-001: Reporter Structure

**Contract:** `Reporter`

**Design:** `ReporterModule`

**Status:** Done

**Requirement:**

Reporter generates validation reports in JSON and Markdown formats.

**Implementation:** `internal/reporter/reporter.go`

**Key Types:**

- `Reporter` — Report generator with config and format

**Acceptance Criteria:**

- [x] Reporter stores config and format
- [x] Reporter.Generate creates complete report
- [x] Reporter.Write outputs report to stdout or file

**Tests:** `TEST-INTERNAL_REPORTER-001`, `TEST-INTERNAL_REPORTER-002`, `TEST-INTERNAL_REPORTER-003`, `TEST-INTERNAL_REPORTER-004`, `TEST-INTERNAL_REPORTER-005`, `TEST-INTERNAL_REPORTER-006`, `TEST-INTERNAL_REPORTER-007`, `TEST-INTERNAL_REPORTER-008`, `TEST-INTERNAL_REPORTER-009`, `TEST-INTERNAL_REPORTER-010`

## SPEC-INTERNAL_REPORTER-003: Reporter.New

**Function Signature:**
`func New(cfg *config.Config, format string) *Reporter`

**Purpose:** Creates a new Reporter with the given configuration and output format. Defaults to JSON if format is empty.

**Parameters:**

- `cfg`: Configuration pointer
- `format`: Output format ("json" or "markdown")

**Returns:** A new Reporter instance

**Tests:** `TEST-INTERNAL_REPORTER-001`

## SPEC-INTERNAL_REPORTER-004: Reporter.Generate

**Function Signature:**
`func (r *Reporter) Generate(result *model.ValidationResult) (*model.Report, error)`

**Purpose:** Generates a complete report from a validation result, including tool metadata, config summary, and the validation result.

**Parameters:**

- `result`: The validation result to include in the report

**Returns:** Complete Report structure or error

**Tests:** `TEST-INTERNAL_REPORTER-002`

## SPEC-INTERNAL_REPORTER-005: Reporter.Write

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

**Tests:** `TEST-INTERNAL_REPORTER-001`, `TEST-INTERNAL_REPORTER-002`, `TEST-INTERNAL_REPORTER-003`, `TEST-INTERNAL_REPORTER-004`, `TEST-INTERNAL_REPORTER-005`, `TEST-INTERNAL_REPORTER-006`, `TEST-INTERNAL_REPORTER-007`, `TEST-INTERNAL_REPORTER-008`, `TEST-INTERNAL_REPORTER-009`, `TEST-INTERNAL_REPORTER-010`

**Related:** `SPEC-INTERNAL_REPORTER-001`
