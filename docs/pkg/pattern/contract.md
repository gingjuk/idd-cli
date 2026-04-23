---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Contract (pattern)

**Status:** Done

**Requirement:**

The pattern package provides IDD identifier pattern matching and extraction utilities.

**Key Contracts:**

- ExtractIDDReferences must find all valid identifiers
- ExtractAnnotations must filter out quoted identifiers
- SplitAnnotationRefs must handle comma-separated values
- GetIdentifierType must return correct type for known patterns

**Implementation:** `pkg/pattern/idd.go`

**Acceptance Criteria:**

- [x] Valid identifiers are extracted correctly
- [x] Quoted identifiers are filtered out
- [x] Multiple references in one annotation are split
- [x] Unknown identifier types return empty string

**Related:** `SPEC-PKG_PAT-001` through `SPEC-PKG_PAT-009`
