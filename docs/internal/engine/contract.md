---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Contract (engine)

**Status:** Done

**Requirement:**

The engine must provide a complete validation pipeline that processes identifiers through collection, graph building, and validation rules.

**Key Contracts:**

- Engine.Run must build graph and run validation rules
- Engine.Run must return ValidationResult with errors and stats
- Engine.AddStructuralErrors must append errors to result
- Engine.BuildReport must return complete Report

**Implementation:** `internal/engine/engine.go`

**Acceptance Criteria:**

- [x] Run processes identifiers through validation pipeline
- [x] Run returns result with errors and statistics
- [x] Run supports context for cancellation
- [x] AddStructuralErrors appends errors correctly
- [x] BuildReport returns complete report structure

**Related:** `SPEC-INT_ENG-001`, `SPEC-INT_ENG-002`
