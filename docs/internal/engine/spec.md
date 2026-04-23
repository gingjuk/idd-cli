---
markers:
  - id: SPEC-INT_ENG-001
    name: Engine Core Struct
  - id: SPEC-INT_ENG-002
    name: Engine.New Method
  - id: SPEC-INT_ENG-004
    name: Engine.Run Method

related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Specification (engine)

## SPEC-INT_ENG-001: Engine Core Struct

**Contract:** `Engine`

**Design:** `EngineModule`

**Status:** Done

**Requirement:**

The validation engine is the core orchestration component that coordinates collection, graph building, and validation rules.

**Implementation:** `internal/engine/engine.go`

**Key Types:**

- `Engine` — Core validation engine with config, graph, and result

**Acceptance Criteria:**

- [x] Engine struct holds config, graph, and result
- [x] Engine coordinates validation pipeline
- [x] Engine supports context for cancellation

**Tests:** `TEST-INT_ENG-001`, `TEST-INT_ENG-002`, `TEST-INT_ENG-003`, `TEST-INT_ENG-004`, `TEST-INT_ENG-005`, `TEST-INT_ENG-006`, `TEST-INT_ENG-007`, `TEST-INT_ENG-008`, `TEST-INT_ENG-009`, `TEST-INT_ENG-010`, `TEST-INT_ENG-011`, `TEST-INT_ENG-012`

**Related:** `SPEC-INT_ENG-001`, `SPEC-INT_ENG-002`

---

## SPEC-INT_ENG-002: Engine.New

**Contract:** `New`

**Design:** `EngineModule`

**Status:** Done

**Requirement:**

Factory function to create a new Engine instance with configuration.

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

**Tests:** `TEST-INT_ENG-002`, `TEST-INT_ENG-003`

---

## SPEC-INT_ENG-004: Engine.Run

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

**Tests:** `TEST-INT_ENG-004`, `TEST-INT_ENG-005`
