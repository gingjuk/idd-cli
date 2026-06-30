---
markers:
  - id: TEST-INTERNAL_REPORTER-001
    name: Reporter Test 1
  - id: TEST-INTERNAL_REPORTER-002
    name: Reporter Test 2
  - id: TEST-INTERNAL_REPORTER-003
    name: Reporter Test 3
  - id: TEST-INTERNAL_REPORTER-004
    name: Reporter Test 4
  - id: TEST-INTERNAL_REPORTER-005
    name: Reporter Test 5
  - id: TEST-INTERNAL_REPORTER-006
    name: Reporter Test 6
  - id: TEST-INTERNAL_REPORTER-007
    name: Reporter Test 7
  - id: TEST-INTERNAL_REPORTER-008
    name: Reporter Test 8
  - id: TEST-INTERNAL_REPORTER-009
    name: Reporter Test 9
  - id: TEST-INTERNAL_REPORTER-010
    name: Reporter Test 10
  - id: TEST-INTERNAL_REPORTER-011
    name: JSON LLM Output
  - id: TEST-INTERNAL_REPORTER-012
    name: LLM Markdown Output
  - id: TEST-INTERNAL_REPORTER-013
    name: LLM Finding Enrichment

related_files:
  spec: docs/internal/reporter/spec.md
  contract: docs/internal/reporter/contract.md
  design: docs/internal/reporter/design.md
  testing: docs/internal/reporter/testing.md
---

# Test Cases (reporter)

## TEST-INTERNAL_REPORTER-001: Reporter New

**Status:** Done

**Purpose:**

Test reporter creation.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-003`

---

## TEST-INTERNAL_REPORTER-002: Reporter Generate

**Status:** Done

**Purpose:**

Test report generation.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-003`

---

## TEST-INTERNAL_REPORTER-003: JSON Output

**Status:** Done

**Purpose:**

Test JSON output format.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-005`

---

## TEST-INTERNAL_REPORTER-004: Markdown Output

**Status:** Done

**Purpose:**

Test Markdown output format.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-005`

---

## TEST-INTERNAL_REPORTER-005: Write Stdout

**Status:** Done

**Purpose:**

Test writing to stdout.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-005`

---

## TEST-INTERNAL_REPORTER-006: Write File

**Status:** Done

**Purpose:**

Test writing to file.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-005`

---

## TEST-INTERNAL_REPORTER-007: Empty Result

**Status:** Done

**Purpose:**

Test with empty result.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-004`

---

## TEST-INTERNAL_REPORTER-008: Error Result

**Status:** Done

**Purpose:**

Test with error result.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-004`

---

## TEST-INTERNAL_REPORTER-009: Warning Result

**Status:** Done

**Purpose:**

Test with warning result.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-004`

---

## TEST-INTERNAL_REPORTER-010: Complete Report

**Status:** Done

**Purpose:**

Test complete report with all fields.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-001`, `SPEC-INTERNAL_REPORTER-003`, `SPEC-INTERNAL_REPORTER-004`, `SPEC-INTERNAL_REPORTER-005`

---

## TEST-INTERNAL_REPORTER-011: JSON LLM Output

**Status:** Done

**Purpose:**

Test that default JSON output is valid JSON using schema `idd.llm_report.v1` and includes summary, finding group, and finding fields.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-011`

---

## TEST-INTERNAL_REPORTER-012: LLM Markdown Output

**Status:** Done

**Purpose:**

Test that `llm-markdown` output is action-oriented and includes status, problem, fix, location, and related identifier context.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-011`

---

## TEST-INTERNAL_REPORTER-013: LLM Finding Enrichment

**Status:** Done

**Purpose:**

Test that validation errors are converted into self-contained findings with severity, rule metadata, structured locations, related identifiers, and grouped summaries for repeated rule failures.

**Spec Coverage:** `SPEC-INTERNAL_REPORTER-012`
