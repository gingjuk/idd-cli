---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Contract (reporter)

**Status:** Done

**Requirement:**

The reporter must generate validation reports that include all validation results, statistics, and tool metadata.

**Key Contracts:**

- Generate must include tool info, config summary, and result
- Write must output to stdout or file based on output path
- Report must be serializable to JSON or Markdown

**Implementation:** `internal/reporter/reporter.go`

**Acceptance Criteria:**

- [x] JSON output includes all validation results
- [x] Markdown output is human-readable
- [x] Errors include identifier ID and source location
- [x] Statistics are accurate

**Tests:** `TEST-INT_RPT-001`, `TEST-INT_RPT-002`, `TEST-INT_RPT-003`, `TEST-INT_RPT-004`, `TEST-INT_RPT-005`, `TEST-INT_RPT-006`, `TEST-INT_RPT-007`, `TEST-INT_RPT-008`, `TEST-INT_RPT-009`, `TEST-INT_RPT-010`
