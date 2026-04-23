---
markers:
  - id: TEST-PKG_WALK-001
    name: Walk Test 1
  - id: TEST-PKG_WALK-002
    name: Walk Test 2
  - id: TEST-PKG_WALK-003
    name: Walk Test 3
  - id: TEST-PKG_WALK-004
    name: Walk Test 4
  - id: TEST-PKG_WALK-005
    name: Walk Test 5
  - id: TEST-PKG_WALK-006
    name: Walk Test 6
  - id: TEST-PKG_WALK-007
    name: Walk Test 7
  - id: TEST-PKG_WALK-008
    name: Walk Test 8

related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Test Cases (walk)

## TEST-PKG_WALK-001: Walk Single Pattern

**Status:** Done

**Purpose:**

Test walking with a single glob pattern.

**Spec Coverage:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`, `SPEC-PKG_WALK-003`

---

## TEST-PKG_WALK-002: Walk Multiple Patterns

**Status:** Done

**Purpose:**

Test walking with multiple glob patterns.

**Spec Coverage:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`, `SPEC-PKG_WALK-003`

---

## TEST-PKG_WALK-003: Walk Recursive

**Status:** Done

**Purpose:**

Test recursive directory walking.

**Spec Coverage:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`, `SPEC-PKG_WALK-003`

---

## TEST-PKG_WALK-004: Walk Duplicate Avoidance

**Status:** Done

**Purpose:**

Test that duplicates are avoided.

**Spec Coverage:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`, `SPEC-PKG_WALK-003`

---

## TEST-PKG_WALK-005: Extension Matching

**Status:** Done

**Purpose:**

Test matching file extensions.

**Spec Coverage:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`, `SPEC-PKG_WALK-003`

---

## TEST-PKG_WALK-006: Visitor Callback

**Status:** Done

**Purpose:**

Test visitor callback is invoked correctly.

**Spec Coverage:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`, `SPEC-PKG_WALK-003`

---

## TEST-PKG_WALK-007: Walk NonExistent Path

**Status:** Done

**Purpose:**

Test walking non-existent path.

**Spec Coverage:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`, `SPEC-PKG_WALK-003`

---

## TEST-PKG_WALK-008: Walk With Subdirectory

**Status:** Done

**Purpose:**

Test walking directory with subdirectories.

**Spec Coverage:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`, `SPEC-PKG_WALK-003`
