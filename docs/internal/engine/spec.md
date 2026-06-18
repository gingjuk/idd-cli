---
markers:
  - id: SPEC-INTERNAL_ENGINE-001
    name: Engine Core Struct
  - id: SPEC-INTERNAL_ENGINE-002
    name: Engine.New Method
  - id: SPEC-INTERNAL_ENGINE-004
    name: Engine.Run Method

related_files:
  spec: docs/internal/engine/spec.md
  contract: docs/internal/engine/contract.md
  design: docs/internal/engine/design.md
  testing: docs/internal/engine/testing.md
---

# Specification (engine)

## SPEC-INTERNAL_ENGINE-001: Engine Core Struct

**Design:** `EngineModule`

**Contract:** `Engine`

**Requirement:**

The validation engine is the core orchestration component that coordinates collection, graph building, and validation rules.

**Tests:** `TEST-INTERNAL_ENGINE-001`, `TEST-INTERNAL_ENGINE-002`, `TEST-INTERNAL_ENGINE-003`, `TEST-INTERNAL_ENGINE-004`, `TEST-INTERNAL_ENGINE-005`, `TEST-INTERNAL_ENGINE-006`, `TEST-INTERNAL_ENGINE-007`, `TEST-INTERNAL_ENGINE-008`, `TEST-INTERNAL_ENGINE-009`, `TEST-INTERNAL_ENGINE-010`, `TEST-INTERNAL_ENGINE-011`, `TEST-INTERNAL_ENGINE-012`

**Status:** Done

**Implementation:** `internal/engine/engine.go`

**Key Types:**

- `Engine` — Core validation engine with config, graph, and result

**Acceptance Criteria:**

- [x] Engine struct holds config, graph, and result
- [x] Engine coordinates validation pipeline
- [x] Engine supports context for cancellation

**Related:** `SPEC-INTERNAL_ENGINE-001`, `SPEC-INTERNAL_ENGINE-002`

## SPEC-INTERNAL_ENGINE-002: Engine.New

**Design:** `EngineModule`

**Contract:** `New`

**Requirement:**

Factory function to create a new Engine instance with configuration.

**Tests:** `TEST-INTERNAL_ENGINE-002`, `TEST-INTERNAL_ENGINE-003`

---

**Status:** Done

**Implementation:** `internal/engine/engine.go`

**Function Signature:**
`func New(cfg *config.Config) *Engine`

**Purpose:** Creates a new validation engine with the given configuration.

**Parameters:**

- `cfg`: Configuration pointer

**Returns:** A new Engine ready to run validation

**Acceptance Criteria:**

- [x] Engine created with config reference
- [x] Empty graph initialized
- [x] Empty result initialized

---

## SPEC-INTERNAL_ENGINE-004: Engine.Run

**Design:** `EngineModule`

**Contract:** `Run`

**Requirement:**

Execute the validation pipeline for the given identifiers and return the accumulated validation result.

**Tests:** `TEST-INTERNAL_ENGINE-004`, `TEST-INTERNAL_ENGINE-005`
**Function Signature:**
`func (e *Engine) Run(ctx context.Context, ids *model.IdentifierSet) (*model.ValidationResult, error)`

**Purpose:** Executes the validation pipeline for the given identifiers.

**Parameters:**

- `ctx`: Context for cancellation
- `ids`: Identifier set to validate

**Returns:** Validation result with stats and errors

**Acceptance Criteria:**

- [x] Builds linkage graph from identifiers
- [x] Runs all enabled validation rules
- [x] Returns accumulated validation result

