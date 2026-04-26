---
related_files:
  spec: docs/pkg/pattern/spec.md
  contract: docs/pkg/pattern/contract.md
  design: docs/pkg/pattern/design.md
  testing: docs/pkg/pattern/testing.md
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

**Related:** `SPEC-PKG_PATTERN-001` through `SPEC-PKG_PATTERN-009`
