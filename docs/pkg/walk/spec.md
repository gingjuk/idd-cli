---
markers:
  - id: SPEC-PKG_WALK-001
    name: Walk Function
  - id: SPEC-PKG_WALK-002
    name: Extension Matching
  - id: SPEC-PKG_WALK-003
    name: File Visitor Pattern

related_files:
  spec: docs/pkg/walk/spec.md
  contract: docs/pkg/walk/contract.md
  design: docs/pkg/walk/design.md
  testing: docs/pkg/walk/testing.md
---

# Specification (walk)

## SPEC-PKG_WALK-001: Walk Function

**Design:** `WalkModule`

**Contract:** `Walk`

**Requirement:**

Walk filesystem matching files against glob patterns.

**Tests:** `TEST-PKG_WALK-001`, `TEST-PKG_WALK-002`, `TEST-PKG_WALK-003`, `TEST-PKG_WALK-004`, `TEST-PKG_WALK-005`, `TEST-PKG_WALK-006`, `TEST-PKG_WALK-007`, `TEST-PKG_WALK-008`

**Status:** Done

**Implementation:** `pkg/walk/files.go`

**Public Functions:**

### Walk

**Function Signature:**
`func Walk(patterns []string, visitor FileVisitor) error`

**Purpose:** Walks the filesystem matching files against the given glob patterns.
---

## SPEC-PKG_WALK-002: Extension Matching

**Design:** `WalkModule`

**Contract:** `MatchAnyExtensions`

**Requirement:**

Check if file path has any of the specified extensions.

**Tests:** `TEST-PKG_WALK-001`, `TEST-PKG_WALK-002`, `TEST-PKG_WALK-003`, `TEST-PKG_WALK-004`, `TEST-PKG_WALK-005`, `TEST-PKG_WALK-006`, `TEST-PKG_WALK-007`, `TEST-PKG_WALK-008`

**Status:** Done

**Implementation:** `pkg/walk/files.go`

**Public Functions:**

### MatchAnyExtensions

**Function Signature:**
`func MatchAnyExtensions(path string, extensions []string) bool`

**Purpose:** Checks if a file path has any of the specified extensions.
---

## SPEC-PKG_WALK-003: File Visitor Pattern

**Design:** `WalkModule`

**Contract:** `FileVisitor`

**Requirement:**

Visitor pattern for file traversal callbacks.

**Tests:** `TEST-PKG_WALK-001`, `TEST-PKG_WALK-002`, `TEST-PKG_WALK-003`, `TEST-PKG_WALK-004`, `TEST-PKG_WALK-005`, `TEST-PKG_WALK-006`, `TEST-PKG_WALK-007`, `TEST-PKG_WALK-008`

**Status:** Done

**Implementation:** `pkg/walk/files.go`

**Key Types:**

- `FileVisitor` — Callback function type for file visitation
**Related:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`
