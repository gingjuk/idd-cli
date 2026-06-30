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

The reporter must generate LLM-oriented validation reports that include status, summary, grouped findings, detailed findings, and tool-derived validation context.

**Key Contracts:**

- Generate must include tool info, config summary, and result
- Write must output to stdout or file based on output path
- Report output must be serializable to JSON or Markdown
- JSON output must use the LLM report schema `idd.llm_report.v1`
- LLM outputs must derive from the complete validation result while grouping repeated rule failures for easier analysis

**Implementation:** `internal/reporter/reporter.go`

**Acceptance Criteria:**

- [x] JSON output includes LLM report schema, status, summary, finding groups, and findings
- [x] Markdown output is human-readable
- [x] Errors include identifier ID and source location
- [x] Statistics are accurate
- [x] LLM JSON includes schema, status, summary, finding groups, and findings
- [x] LLM Markdown includes actionable problem and fix sections

**Tests:** `TEST-INTERNAL_REPORTER-001`, `TEST-INTERNAL_REPORTER-002`, `TEST-INTERNAL_REPORTER-003`, `TEST-INTERNAL_REPORTER-004`, `TEST-INTERNAL_REPORTER-005`, `TEST-INTERNAL_REPORTER-006`, `TEST-INTERNAL_REPORTER-007`, `TEST-INTERNAL_REPORTER-008`, `TEST-INTERNAL_REPORTER-009`, `TEST-INTERNAL_REPORTER-010`, `TEST-INTERNAL_REPORTER-011`, `TEST-INTERNAL_REPORTER-012`, `TEST-INTERNAL_REPORTER-013`
