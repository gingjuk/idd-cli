---
markers:
  - id: SPEC-INTERNAL_MODEL-001
    name: Identifier Type Parsing
  - id: SPEC-INTERNAL_MODEL-002
    name: Identifier Creation
  - id: SPEC-INTERNAL_MODEL-003
    name: Identifier Link Management
  - id: SPEC-INTERNAL_MODEL-004
    name: IdentifierSet Collection
  - id: SPEC-INTERNAL_MODEL-005
    name: IdentifierSet Access Operations
  - id: SPEC-INTERNAL_MODEL-006
    name: IdentifierSet Count and All
  - id: SPEC-INTERNAL_MODEL-007
    name: IdentifierSet Merge Operation
  - id: SPEC-INTERNAL_MODEL-008
    name: Annotation Creation
  - id: SPEC-INTERNAL_MODEL-009
    name: Annotation Conversion
  - id: SPEC-INTERNAL_MODEL-010
    name: Validation Result Types
  - id: SPEC-INTERNAL_MODEL-011
    name: Validation Result Operations
  - id: SPEC-INTERNAL_MODEL-012
    name: Link Type Definition
  - id: SPEC-INTERNAL_MODEL-013
    name: Link Structure
  - id: SPEC-INTERNAL_MODEL-014
    name: Origin Type Definition
  - id: SPEC-INTERNAL_MODEL-015
    name: Origin Operations
  - id: SPEC-INTERNAL_MODEL-017
    name: Identifier Origin Methods
  - id: SPEC-INTERNAL_MODEL-018
    name: Identifier Model Core

related_files:
  spec: docs/internal/model/spec.md
  contract: docs/internal/model/contract.md
  design: docs/internal/model/design.md
  testing: docs/internal/model/testing.md
---

# Specification (model)

## SPEC-INTERNAL_MODEL-001: Identifier Type Parsing

**Contract:** `ParseIdentifierType`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Parse identifier type from string representation.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INTERNAL_MODEL-001`

**Public Functions:**

## SPEC-INTERNAL_MODEL-025: Identifier.ParseIdentifierType

**Function Signature:**
`func ParseIdentifierType(s string) (IdentifierType, error)`

---

## SPEC-INTERNAL_MODEL-002: Identifier Creation

**Contract:** `NewIdentifier`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Create a new identifier with given fields.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INTERNAL_MODEL-002`

**Public Functions:**

## SPEC-INTERNAL_MODEL-026: Identifier.NewIdentifier

**Function Signature:**
`func NewIdentifier(id string, idType IdentifierType, title, source string, line int) *Identifier`

---

## SPEC-INTERNAL_MODEL-003: Identifier Link Management

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

**Tests:** `TEST-INTERNAL_MODEL-003`

---

## SPEC-INTERNAL_MODEL-004: IdentifierSet Collection

**Contract:** `IdentifierSet`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

IdentifierSet stores identifiers by type and ID.

**Implementation:** `internal/model/identifier.go`

**Key Types:**

- `IdentifierSet` — Collection with type slices and byID map

**Tests:** `TEST-INTERNAL_MODEL-004`

---

## SPEC-INTERNAL_MODEL-005: IdentifierSet Access Operations

**Contract:** `IdentifierSet`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

IdentifierSet provides Get, Has operations.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INTERNAL_MODEL-005`

**Public Functions:**

## SPEC-INTERNAL_MODEL-027: IdentifierSet.Get

## SPEC-INTERNAL_MODEL-028: IdentifierSet.Has

---

## SPEC-INTERNAL_MODEL-006: IdentifierSet Count and All

**Contract:** `IdentifierSet`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

IdentifierSet provides Count and All operations.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INTERNAL_MODEL-006`

**Public Functions:**

## SPEC-INTERNAL_MODEL-029: IdentifierSet.Count

## SPEC-INTERNAL_MODEL-030: IdentifierSet.All

---

## SPEC-INTERNAL_MODEL-007: IdentifierSet Merge Operation

**Contract:** `IdentifierSet`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Combine two identifier sets.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INTERNAL_MODEL-007`

**Public Functions:**

## SPEC-INTERNAL_MODEL-031: IdentifierSet.Merge

---

## SPEC-INTERNAL_MODEL-008: Annotation Creation

**Contract:** `NewAnnotation`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Create annotation from code.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INTERNAL_MODEL-008`

**Public Functions:**

## SPEC-INTERNAL_MODEL-032: Annotation.NewAnnotation

---

## SPEC-INTERNAL_MODEL-009: Annotation Conversion

**Contract:** `Annotation`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Convert annotation to identifier.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INTERNAL_MODEL-009`

**Public Functions:**

## SPEC-INTERNAL_MODEL-033: Annotation.ToIdentifier

---

## SPEC-INTERNAL_MODEL-010: Validation Result Types

**Contract:** `ValidationResult`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Validation result stores errors and warnings.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INTERNAL_MODEL-010`

**Key Types:**

- `ValidationResult` — Result with errors and warnings
- `ValidationError` — Error with rule, message, source, link

---

## SPEC-INTERNAL_MODEL-011: Validation Result Operations

**Contract:** `ValidationResult`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Add errors/warnings and sort results.

**Implementation:** `internal/model/identifier.go`

**Tests:** `TEST-INTERNAL_MODEL-011`, `TEST-INTERNAL_MODEL-016`

**Public Functions:**

## SPEC-INTERNAL_MODEL-034: ValidationResult.AddError

## SPEC-INTERNAL_MODEL-035: ValidationResult.AddWarning

## SPEC-INTERNAL_MODEL-036: ValidationResult.Sort

**Tests:** `TEST-INTERNAL_MODEL-011`, `TEST-INTERNAL_MODEL-016`

---

## SPEC-INTERNAL_MODEL-012: Link Type Definition

**Contract:** `LinkType`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

LinkType enum for edge types.

**Implementation:** `internal/model/link.go`

**Key Types:**

- `LinkType` — Enum for LinkTests, LinkImplements, LinkReferences

**Tests:** `TEST-INTERNAL_MODEL-012`, `TEST-INTERNAL_MODEL-028`

---

## SPEC-INTERNAL_MODEL-013: Link Structure

**Contract:** `Link`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Link represents a reference from one identifier to another.

**Implementation:** `internal/model/link.go`

**Key Types:**

- `Link` — Reference with type and target ID

**Tests:** `TEST-INTERNAL_MODEL-013`, `TEST-INTERNAL_MODEL-029`, `TEST-INTERNAL_MODEL-030`

---

## SPEC-INTERNAL_MODEL-014: Origin Type Definition

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

**Tests:** `TEST-INTERNAL_MODEL-014`

---

## SPEC-INTERNAL_MODEL-015: Origin Operations

**Contract:** `Origin`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Set and check identifier origin.

**Implementation:** `internal/model/identifier.go`

**Methods:**

- `SetOrigin(o Origin)`
- `Origin() Origin`

**Tests:** `TEST-INTERNAL_MODEL-015`

---

## SPEC-INTERNAL_MODEL-017: Identifier Origin Methods

**Contract:** `SetOrigin`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Methods to get/set identifier origin.

**Implementation:** `internal/model/identifier.go`

**Methods:**

- `SetOrigin(origin Origin)` — implemented as part of `SPEC-INTERNAL_MODEL-017`
- `GetOrigin() Origin` — not present in code (only SetOrigin exists)

**Tests:** `TEST-INTERNAL_MODEL-017`

**Note:** `GetOrigin()` method does not exist in code. Only `SetOrigin()` is implemented.

---

## SPEC-INTERNAL_MODEL-018: Identifier Model Core

**Contract:** `Identifier`

**Design:** `ModelModule`

**Status:** Done

**Requirement:**

Core identifier type with all fields.

**Implementation:** `internal/model/identifier.go`

**Key Types:**

- `Identifier` — IDD identifier with type, module, number, links

**Tests:** `TEST-INTERNAL_MODEL-018`

**Tests:** ``
