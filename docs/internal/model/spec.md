---
markers:
  - id: SPEC-INT_MOD-001
    name: Identifier Type Parsing
  - id: SPEC-INT_MOD-002
    name: Identifier Creation
  - id: SPEC-INT_MOD-003
    name: Identifier Link Management
  - id: SPEC-INT_MOD-004
    name: IdentifierSet Collection
  - id: SPEC-INT_MOD-005
    name: IdentifierSet Access Operations
  - id: SPEC-INT_MOD-006
    name: IdentifierSet Count and All
  - id: SPEC-INT_MOD-007
    name: IdentifierSet Merge Operation
  - id: SPEC-INT_MOD-008
    name: Annotation Creation
  - id: SPEC-INT_MOD-009
    name: Annotation Conversion
  - id: SPEC-INT_MOD-010
    name: Validation Result Types
  - id: SPEC-INT_MOD-011
    name: Validation Result Operations
  - id: SPEC-INT_MOD-012
    name: Link Type Definition
  - id: SPEC-INT_MOD-013
    name: Link Structure
  - id: SPEC-INT_MOD-014
    name: Origin Type Definition
  - id: SPEC-INT_MOD-015
    name: Origin Operations
  - id: SPEC-INT_MOD-017
    name: Identifier Origin Methods
  - id: SPEC-INT_MOD-018
    name: Identifier Model Core

related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Specification (model)

## SPEC-INT_MOD-001: Identifier Type Parsing

**Contract:** `ParseIdentifierType`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Parse identifier type from string representation.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INT_MOD-001`

**Public Functions:**

## SPEC-INT_MOD-025: Identifier.ParseIdentifierType

**Function Signature:**
`func ParseIdentifierType(s string) (IdentifierType, error)`

---

## SPEC-INT_MOD-002: Identifier Creation

**Contract:** `NewIdentifier`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Create a new identifier with given fields.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INT_MOD-002`

**Public Functions:**

## SPEC-INT_MOD-026: Identifier.NewIdentifier

**Function Signature:**
`func NewIdentifier(id string, idType IdentifierType, title, source string, line int) *Identifier`

---

## SPEC-INT_MOD-003: Identifier Link Management

**Contract:** `AddLink`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Identifiers store forward links to dependencies.

**Implementation:** `internal/model/identifier.go`

**Acceptance Criteria:**

- [x] Identifier has Links slice
- [x] AddLink adds forward reference
- [x] Backlinks are computed from forward links

**Tests:** `TEST-INT_MOD-003`

---

## SPEC-INT_MOD-004: IdentifierSet Collection

**Contract:** `IdentifierSet`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

IdentifierSet stores identifiers by type and ID.

**Implementation:** `internal/model/identifier.go`

**Key Types:**

- `IdentifierSet` — Collection with type slices and byID map

**Tests:** `TEST-INT_MOD-004`

---

## SPEC-INT_MOD-005: IdentifierSet Access Operations

**Contract:** `IdentifierSet`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

IdentifierSet provides Get, Has operations.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INT_MOD-005`

**Public Functions:**

## SPEC-INT_MOD-027: IdentifierSet.Get

## SPEC-INT_MOD-028: IdentifierSet.Has

---

## SPEC-INT_MOD-006: IdentifierSet Count and All

**Contract:** `IdentifierSet`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

IdentifierSet provides Count and All operations.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INT_MOD-006`

**Public Functions:**

## SPEC-INT_MOD-029: IdentifierSet.Count

## SPEC-INT_MOD-030: IdentifierSet.All

---

## SPEC-INT_MOD-007: IdentifierSet Merge Operation

**Contract:** `IdentifierSet`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Combine two identifier sets.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INT_MOD-007`

**Public Functions:**

## SPEC-INT_MOD-031: IdentifierSet.Merge

---

## SPEC-INT_MOD-008: Annotation Creation

**Contract:** `NewAnnotation`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Create annotation from code.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INT_MOD-008`

**Public Functions:**

## SPEC-INT_MOD-032: Annotation.NewAnnotation

---

## SPEC-INT_MOD-009: Annotation Conversion

**Contract:** `Annotation`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Convert annotation to identifier.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INT_MOD-009`

**Public Functions:**

## SPEC-INT_MOD-033: Annotation.ToIdentifier

---

## SPEC-INT_MOD-010: Validation Result Types

**Contract:** `ValidationResult`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Validation result stores errors and warnings.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INT_MOD-010`

**Key Types:**

- `ValidationResult` — Result with errors and warnings
- `ValidationError` — Error with rule, message, source, link

---

## SPEC-INT_MOD-011: Validation Result Operations

**Contract:** `ValidationResult`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Add errors/warnings and sort results.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INT_MOD-011`, `TEST-INT_MOD-016`

**Public Functions:**

## SPEC-INT_MOD-034: ValidationResult.AddError

## SPEC-INT_MOD-035: ValidationResult.AddWarning

## SPEC-INT_MOD-036: ValidationResult.Sort

**Tests:** `TEST-INT_MOD-011`, `TEST-INT_MOD-016`

---

## SPEC-INT_MOD-012: Link Type Definition

**Contract:** `LinkType`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

LinkType enum for edge types.

**Implementation:** `internal/model/link.go`

**Key Types:**

- `LinkType` — Enum for LinkTests, LinkImplements, LinkReferences

**Tests:** `TEST-INT_MOD-012`, `TEST-INT_MOD-028`

---

## SPEC-INT_MOD-013: Link Structure

**Contract:** `Link`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Link represents a reference from one identifier to another.

**Implementation:** `internal/model/link.go`

**Key Types:**

- `Link` — Reference with type and target ID

**Tests:** `TEST-INT_MOD-013`, `TEST-INT_MOD-029`, `TEST-INT_MOD-030`

---

## SPEC-INT_MOD-014: Origin Type Definition

**Contract:** `Origin`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Origin represents the source location type of an identifier.

**Implementation:** `internal/model/identifier.go`

**Key Types:**

- `Origin` — Enum type with values `Doc` and `Code`
- `OriginDoc` — Indicates identifier came from documentation
- `OriginCode` — Indicates identifier came from source code

**Tests:** `TEST-INT_MOD-014`

---

## SPEC-INT_MOD-015: Origin Operations

**Contract:** `Origin`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Set and check identifier origin.

**Implementation:** `internal/model/identifier.go`

**Methods:**

- `SetOrigin(o Origin)`
- `Origin() Origin`

**Tests:** `TEST-INT_MOD-015`

---

## SPEC-INT_MOD-017: Identifier Origin Methods

**Contract:** `SetOrigin`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Methods to get/set identifier origin.

**Implementation:** `internal/model/identifier.go`

**Methods:**

- `SetOrigin(origin Origin)` — implemented as part of `SPEC-INT_MOD-017`
- `GetOrigin() Origin` — not present in code (only SetOrigin exists)

**Tests:** `TEST-INT_MOD-017`

**Note:** `GetOrigin()` method does not exist in code. Only `SetOrigin()` is implemented.

---

## SPEC-INT_MOD-018: Identifier Model Core

**Contract:** `Identifier`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Core identifier type with all fields.

**Implementation:** `internal/model/identifier.go`

**Key Types:**

- `Identifier` — IDD identifier with type, module, number, links

**Tests:** `TEST-INT_MOD-018`

**Tests:** ``
