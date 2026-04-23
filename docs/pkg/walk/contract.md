---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Contract (walk)

**Status:** Done

**Requirement:**

The walk package provides file traversal utilities with pattern matching support.

**Key Contracts:**

- Walk must visit each file matching patterns exactly once
- Walk must avoid duplicate visits using visited map
- MatchAnyExtensions must check extension substrings
- FileVisitor is called for each matched file/directory

**Implementation:** `pkg/walk/files.go`

**Acceptance Criteria:**

- [x] Files are visited according to glob patterns
- [x] Directories are walked recursively
- [x] Duplicate visits are avoided
- [x] Visitor callback receives file path and info

**Related:** `SPEC-PKG_WALK-001`, `SPEC-PKG_WALK-002`, `SPEC-PKG_WALK-003`
