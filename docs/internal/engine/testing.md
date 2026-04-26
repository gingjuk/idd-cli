---
markers:
  - id: TEST-INTERNAL_ENGINE-001
    name: Engine Run Basic
  - id: TEST-INTERNAL_ENGINE-002
    name: Engine BuildGraph
  - id: TEST-INTERNAL_ENGINE-003
    name: Engine Validation
  - id: TEST-INTERNAL_ENGINE-004
    name: Engine Structural Errors
  - id: TEST-INTERNAL_ENGINE-005
    name: Engine Report Generation
  - id: TEST-INTERNAL_ENGINE-006
    name: Engine Context Cancellation
  - id: TEST-INTERNAL_ENGINE-007
    name: Engine Consistency Check
  - id: TEST-INTERNAL_ENGINE-008
    name: Engine Validation Rules
  - id: TEST-INTERNAL_ENGINE-009
    name: Engine Config Validation
  - id: TEST-INTERNAL_ENGINE-010
    name: Engine Pattern Matching
  - id: TEST-INTERNAL_ENGINE-011
    name: Engine Annotation Validation
  - id: TEST-INTERNAL_ENGINE-012
    name: Engine File Collection
  - id: TEST-INTERNAL_ENGINE-020
    name: Engine Duplicate Heading Validation
  - id: TEST-INTERNAL_ENGINE-021
    name: Engine Duplicate Heading No Duplicates
  - id: TEST-INTERNAL_ENGINE-022
    name: Engine Contract Design Markers With Markers
  - id: TEST-INTERNAL_ENGINE-023
    name: Engine Contract Design Markers Without Markers
  - id: TEST-INTERNAL_ENGINE-024
    name: Engine Contract Design Markers Design File
  - id: TEST-INTERNAL_ENGINE-025
    name: Engine Package Exists
  - id: TEST-INTERNAL_ENGINE-026
    name: Engine Ignored Doc Path
  - id: TEST-INTERNAL_ENGINE-027
    name: Engine Validate Doc Path Exists
  - id: TEST-INTERNAL_ENGINE-028
    name: Engine Validate Related Files
  - id: TEST-INTERNAL_ENGINE-029
    name: Engine Validate Related Files Valid

related_files:
  spec: docs/internal/engine/spec.md
  contract: docs/internal/engine/contract.md
  design: docs/internal/engine/design.md
  testing: docs/internal/engine/testing.md
---

# Test Cases (engine)

## Test

## TEST-INTERNAL_ENGINE-001: Engine Run Basic

**Status:** Done

**Purpose:**

Test basic Engine.Run functionality.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-002: Engine BuildGraph

**Status:** Done

**Purpose:**

Test that Engine correctly builds the linkage graph from identifiers.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-003: Engine Validation

**Status:** Done

**Purpose:**

Test that Engine runs all validation rules correctly.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-004: Engine Structural Errors

**Status:** Done

**Purpose:**

Test that Engine.AddStructuralErrors correctly appends errors.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-005: Engine Report Generation

**Status:** Done

**Purpose:**

Test that Engine.BuildReport returns correct report structure.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-006: Engine Context Cancellation

**Status:** Done

**Purpose:**

Test that Engine properly handles context cancellation.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-007: Engine Consistency Check

**Status:** Done

**Purpose:**

Test that Engine performs consistency checks between doc and code identifiers.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## Contract Test

## TEST-INTERNAL_ENGINE-013: Engine Struct Contract Test

**Status:** Done

**Purpose:**

Test Engine struct initialization and field validation.

---

## TEST-INTERNAL_ENGINE-014: Engine Run Contract Test

**Status:** Done

**Purpose:**

Test Engine.Run execution and result generation.

---

## TEST-INTERNAL_ENGINE-015: Engine Validation Contract Test

**Status:** Done

**Purpose:**

Test Engine validation logic.

---

## TEST-INTERNAL_ENGINE-016: Engine AddStructuralErrors Contract Test

**Status:** Done

**Purpose:**

Test Engine.AddStructuralErrors functionality.

---

## TEST-INTERNAL_ENGINE-017: Engine BuildReport Contract Test

**Status:** Done

**Purpose:**

Test Engine.BuildReport generation.

---

## TEST-INTERNAL_ENGINE-018: Engine Context Contract Test

**Status:** Done

**Purpose:**

Test Engine context handling.

---

## TEST-INTERNAL_ENGINE-019: Engine Consistency Contract Test

**Status:** Done

**Purpose:**

Test Engine consistency check.

---

## TEST-INTERNAL_ENGINE-008: Engine Validation Rules

**Status:** Done

**Purpose:**

Test Engine validation rules including orphan detection.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-009: Engine Config Validation

**Status:** Done

**Purpose:**

Test Engine configuration validation.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-010: Engine Pattern Matching

**Status:** Done

**Purpose:**

Test Engine pattern matching for identifiers.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-011: Engine Annotation Validation

**Status:** Done

**Purpose:**

Test Engine annotation validation.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-012: Engine File Collection

**Status:** Done

**Purpose:**

Test Engine file collection from configured paths.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

---

## TEST-INTERNAL_ENGINE-020: Engine Duplicate Heading Validation

**Status:** Done

**Purpose:**

Test that Engine validates duplicate heading identifiers in documentation.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-021: Engine Duplicate Heading No Duplicates

**Status:** Done

**Purpose:**

Test that Engine passes when headings have unique identifiers.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-022: Engine Contract Design Markers With Markers

**Status:** Done

**Purpose:**

Test that Engine reports error when contract.md has markers.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-023: Engine Contract Design Markers Without Markers

**Status:** Done

**Purpose:**

Test that Engine passes when contract.md has no markers.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-024: Engine Contract Design Markers Design File

**Status:** Done

**Purpose:**

Test that Engine reports error when design.md has markers.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-025: Engine Package Exists

**Status:** Done

**Purpose:**

Test Engine.packageExists method.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-026: Engine Ignored Doc Path

**Status:** Done

**Purpose:**

Test Engine.isIgnoredDocPath method.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-027: Engine Validate Doc Path Exists

**Status:** Done

**Purpose:**

Test Engine.validateDocPathExists method.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-028: Engine Validate Related Files

**Status:** Done

**Purpose:**

Test Engine.validateRelatedFiles when related_files is missing.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-029: Engine Validate Related Files Valid

**Status:** Done

**Purpose:**

Test Engine.validateRelatedFiles when related_files is present.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`
