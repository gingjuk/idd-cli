---
markers:
  - id: SPEC-PKG_WALK-001
    name: Walk Function
  - id: SPEC-PKG_WALK-002
    name: Extension Matching
  - id: SPEC-PKG_WALK-003
    name: File Visitor Pattern

related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Specification (walk)

## SPEC-PKG_WALK-001: Walk Function

**Contract:** `Walk`

**Design:** `WalkModule`

**Status:** Done

**Requirement:**

Walk filesystem matching files against glob patterns.

**Implementation:** `pkg/walk/files.go`

**Public Functions:**

### Walk

**Function Signature:**
`func Walk(patterns []string, visitor FileVisitor) error`

**Purpose:** Walks the filesystem matching files against the given glob patterns.
**Tests:** `TEST-PKG_WALK-001`, `TEST-PKG_WALK-002`, `TEST-PKG_WALK-003`, `TEST-PKG_WALK-004`, `TEST-PKG_WALK-005`, `TEST-PKG_WALK-006`, `TEST-PKG_WALK-007`, `TEST-PKG_WALK-008`

---

## SPEC-PKG_WALK-002: Extension Matching

**Contract:** `MatchAnyExtensions`

**Design:** `WalkModule`

**Status:** Done

**Requirement:**

Check if file path has any of the specified extensions.

**Implementation:** `pkg/walk/files.go`

**Public Functions:**

### MatchAnyExtensions

**Function Signature:**
`func MatchAnyExtensions(path string, extensions []string) bool`

**Purpose:** Checks if a file path has any of the specified extensions.
**Tests:** `TEST-PKG_WALK-001`, `TEST-PKG_WALK-002`, `TEST-PKG_WALK-003`, `TEST-PKG_WALK-004`, `TEST-PKG_WALK-005`, `TEST-PKG_WALK-006`, `TEST-PKG_WALK-007`, `TEST-PKG_WALK-008`

---

## SPEC-PKG_WALK-003: File Visitor Pattern

**Contract:** `FileVisitor`

**Design:** `WalkModule`

**Status:** Done

**Requirement:**

Visitor pattern for file traversal callbacks.

**Implementation:** `pkg/walk/files.go`

**Key Types:**

- `FileVisitor` — Callback function type for file visitation
**Tests:** `TEST-PKG_WALK-001`, `TEST-PKG_WALK-002`, `TEST-PKG_WALK-003`, `TEST-PKG_WALK-004`, `TEST-PKG_WALK-005`, `TEST-PKG_WALK-006`, `TEST-PKG_WALK-007`, `TEST-PKG_WALK-008`

**Related:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`
