---
related_files:
  spec: docs/internal/reporter/spec.md
  contract: docs/internal/reporter/contract.md
  design: docs/internal/reporter/design.md
  testing: docs/internal/reporter/testing.md
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

**Tests:** `TEST-INTERNAL_REPORTER-001`, `TEST-INTERNAL_REPORTER-002`, `TEST-INTERNAL_REPORTER-003`, `TEST-INTERNAL_REPORTER-004`, `TEST-INTERNAL_REPORTER-005`, `TEST-INTERNAL_REPORTER-006`, `TEST-INTERNAL_REPORTER-007`, `TEST-INTERNAL_REPORTER-008`, `TEST-INTERNAL_REPORTER-009`, `TEST-INTERNAL_REPORTER-010`
