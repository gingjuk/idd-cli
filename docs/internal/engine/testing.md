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
  - id: TEST-INTERNAL_ENGINE-030
    name: Engine Package Doc Comment Main Package Skipped
  - id: TEST-INTERNAL_ENGINE-031
    name: Engine Package Doc Files Missing Files
  - id: TEST-INTERNAL_ENGINE-032
    name: Engine Package Doc Files All Present
  - id: TEST-INTERNAL_ENGINE-033
    name: Engine Package Doc Files Root Level Skipped
  - id: TEST-INTERNAL_ENGINE-034
    name: Engine Duplicate IDs Doc Side
  - id: TEST-INTERNAL_ENGINE-035
    name: Engine Duplicate IDs Same Dir No Error
  - id: TEST-INTERNAL_ENGINE-036
    name: Engine Duplicate IDs Code Side
  - id: TEST-INTERNAL_ENGINE-037
    name: Engine Duplicate IDs Rename Suggestion
  - id: TEST-INTERNAL_ENGINE-038
    name: Engine Public Func Annotation Private Func Allowed
  - id: TEST-INTERNAL_ENGINE-039
    name: Engine Public Func Annotation Private Method Allowed
  - id: TEST-INTERNAL_ENGINE-040
    name: Engine Public Func Annotation Private Type Allowed
  - id: TEST-INTERNAL_ENGINE-041
    name: Engine Public Func Annotation Public Still Required
  - id: TEST-INTERNAL_ENGINE-042
    name: Engine Public Func Annotation Garbage After Implement Still Errors
  - id: TEST-INTERNAL_ENGINE-043
    name: Engine Private Implement Requires Doc
  - id: TEST-INTERNAL_ENGINE-044
    name: Engine Annotation Identifier Honors Ignore Scope
  - id: TEST-INTERNAL_ENGINE-045
    name: Engine Annotation Placement Honors Ignore Scope
  - id: TEST-INTERNAL_ENGINE-046
    name: Engine Consecutive Annotations Honors Ignore Scope

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

---

## TEST-INTERNAL_ENGINE-030: Engine Package Doc Comment Main Package Skipped

**Status:** Done

**Purpose:**

Test that the package-doc-comment check skips the main package (which has no
hosted package, only a `main` function).

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-031: Engine Package Doc Files Missing Files

**Status:** Done

**Purpose:**

Test that Engine reports an error when a package's `pkgDocFiles` declares files
that do not exist on disk.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-032: Engine Package Doc Files All Present

**Status:** Done

**Purpose:**

Test that Engine passes when all files declared in `pkgDocFiles` exist on disk.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-033: Engine Package Doc Files Root Level Skipped

**Status:** Done

**Purpose:**

Test that the `pkgDocFiles` check is skipped for root-level packages
(single-segment import paths).

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-034: Engine Duplicate IDs Doc Side

**Status:** Done

**Purpose:**

Test that Engine reports an error when the same identifier appears twice in
documentation headings.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-035: Engine Duplicate IDs Same Dir No Error

**Status:** Done

**Purpose:**

Test that Engine does not flag duplicate identifiers that are intentionally
co-located in the same source file (test/cont variants).

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-036: Engine Duplicate IDs Code Side

**Status:** Done

**Purpose:**

Test that Engine reports an error when the same identifier is declared in
multiple Go source files.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-037: Engine Duplicate IDs Rename Suggestion

**Status:** Done

**Purpose:**

Test that Engine's duplicate-ID error message includes a rename suggestion to
help the user resolve the conflict.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-038: Engine Public Func Annotation Private Func Allowed

**Status:** Done

**Purpose:**

Test that `validatePublicFuncAnnotations` does not require an `@implement`
annotation on private (lowercase) functions.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-039: Engine Public Func Annotation Private Method Allowed

**Status:** Done

**Purpose:**

Test that `validatePublicFuncAnnotations` does not require an `@implement`
annotation on private methods (methods of public types with lowercase names).

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-040: Engine Public Func Annotation Private Type Allowed

**Status:** Done

**Purpose:**

Test that `validatePublicFuncAnnotations` does not require `@implement` on
methods of a private type even when the method itself is exported.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-041: Engine Public Func Annotation Public Still Required

**Status:** Done

**Purpose:**

Test that `validatePublicFuncAnnotations` still flags an exported function
with no `@implement` annotation.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-042: Engine Public Func Annotation Garbage After Implement Still Errors

**Status:** Done

**Purpose:**

Test that garbage tokens following `@implement` (which make the annotation
unparseable) still surface as a validation error, not silently pass.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-043: Engine Private Implement Requires Doc

**Status:** Done

**Purpose:**

Test that an `@implement` on a private function still requires a matching doc
entry — the doc-code-correspondence check applies regardless of visibility.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-044: Engine Annotation Identifier Honors Ignore Scope

**Status:** Done

**Purpose:**

Test that `validateAnnotationIdentifiers` skips content inside an
`// idd:ignore start/end` block. Without scope honoring, the validator would
report every annotation inside a fixture that looks like a real `@implement`
statement.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-045: Engine Annotation Placement Honors Ignore Scope

**Status:** Done

**Purpose:**

Test that the placement pass of `validatePublicFuncAnnotations` skips content
inside an `// idd:ignore start/end` block — not just the first and third
passes.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`

---

## TEST-INTERNAL_ENGINE-046: Engine Consecutive Annotations Honors Ignore Scope

**Status:** Done

**Purpose:**

Test that `validateConsecutiveAnnotations` skips content inside an
`// idd:ignore start/end` block.

**Spec Coverage:** `SPEC-INTERNAL_ENGINE-001`
