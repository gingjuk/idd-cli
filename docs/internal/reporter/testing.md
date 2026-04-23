---
markers:
  - id: TEST-INT_RPT-001
    name: Reporter Test 1
  - id: TEST-INT_RPT-002
    name: Reporter Test 2
  - id: TEST-INT_RPT-003
    name: Reporter Test 3
  - id: TEST-INT_RPT-004
    name: Reporter Test 4
  - id: TEST-INT_RPT-005
    name: Reporter Test 5
  - id: TEST-INT_RPT-006
    name: Reporter Test 6
  - id: TEST-INT_RPT-007
    name: Reporter Test 7
  - id: TEST-INT_RPT-008
    name: Reporter Test 8
  - id: TEST-INT_RPT-009
    name: Reporter Test 9
  - id: TEST-INT_RPT-010
    name: Reporter Test 10

related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Test Cases (reporter)

## TEST-INT_RPT-001: Reporter New

**Status:** Done

**Purpose:**

Test reporter creation.

**Spec Coverage:** `SPEC-INT_RPT-003`

---

## TEST-INT_RPT-002: Reporter Generate

**Status:** Done

**Purpose:**

Test report generation.

**Spec Coverage:** `SPEC-INT_RPT-003`

---

## TEST-INT_RPT-003: JSON Output

**Status:** Done

**Purpose:**

Test JSON output format.

**Spec Coverage:** `SPEC-INT_RPT-005`

---

## TEST-INT_RPT-004: Markdown Output

**Status:** Done

**Purpose:**

Test Markdown output format.

**Spec Coverage:** `SPEC-INT_RPT-005`

---

## TEST-INT_RPT-005: Write Stdout

**Status:** Done

**Purpose:**

Test writing to stdout.

**Spec Coverage:** `SPEC-INT_RPT-005`

---

## TEST-INT_RPT-006: Write File

**Status:** Done

**Purpose:**

Test writing to file.

**Spec Coverage:** `SPEC-INT_RPT-005`

---

## TEST-INT_RPT-007: Empty Result

**Status:** Done

**Purpose:**

Test with empty result.

**Spec Coverage:** `SPEC-INT_RPT-004`

---

## TEST-INT_RPT-008: Error Result

**Status:** Done

**Purpose:**

Test with error result.

**Spec Coverage:** `SPEC-INT_RPT-004`

---

## TEST-INT_RPT-009: Warning Result

**Status:** Done

**Purpose:**

Test with warning result.

**Spec Coverage:** `SPEC-INT_RPT-004`

---

## TEST-INT_RPT-010: Complete Report

**Status:** Done

**Purpose:**

Test complete report with all fields.

**Spec Coverage:** `SPEC-INT_RPT-001`, `SPEC-INT_RPT-003`, `SPEC-INT_RPT-004`, `SPEC-INT_RPT-005`
