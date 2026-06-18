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

**Design:** `ModelModule`

**Contract:** `ParseIdentifierType`

**Requirement:**

Parse identifier type from string representation.

**Tests:** `TEST-INTERNAL_MODEL-001`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Public Functions:**

## SPEC-INTERNAL_MODEL-025: Identifier.ParseIdentifierType

**Function Signature:**
`func ParseIdentifierType(s string) (IdentifierType, error)`

## SPEC-INTERNAL_MODEL-002: Identifier Creation

**Design:** `ModelModule`

**Contract:** `NewIdentifier`

**Requirement:**

Create a new identifier with given fields.

**Tests:** `TEST-INTERNAL_MODEL-002`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Public Functions:**

---

## SPEC-INTERNAL_MODEL-026: Identifier.NewIdentifier

**Function Signature:**
`func NewIdentifier(id string, idType IdentifierType, title, source string, line int) *Identifier`

## SPEC-INTERNAL_MODEL-003: Identifier Link Management

**Design:** `ModelModule`

**Contract:** `AddLink`

**Requirement:**

Identifiers store forward links to dependencies.

**Tests:** `TEST-INTERNAL_MODEL-003`

---

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Acceptance Criteria:**

- [x] Identifier has Links slice
- [x] AddLink adds forward reference
- [x] Backlinks are computed from forward links

---

## SPEC-INTERNAL_MODEL-004: IdentifierSet Collection

**Design:** `ModelModule`

**Contract:** `IdentifierSet`

**Requirement:**

IdentifierSet stores identifiers by type and ID.

**Tests:** `TEST-INTERNAL_MODEL-004`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Key Types:**

- `IdentifierSet` — Collection with type slices and byID map

---

## SPEC-INTERNAL_MODEL-005: IdentifierSet Access Operations

**Design:** `ModelModule`

**Contract:** `IdentifierSet`

**Requirement:**

IdentifierSet provides Get, Has operations.

**Tests:** `TEST-INTERNAL_MODEL-005`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Public Functions:**

## SPEC-INTERNAL_MODEL-027: IdentifierSet.Get

## SPEC-INTERNAL_MODEL-028: IdentifierSet.Has

## SPEC-INTERNAL_MODEL-006: IdentifierSet Count and All

**Design:** `ModelModule`

**Contract:** `IdentifierSet`

**Requirement:**

IdentifierSet provides Count and All operations.

**Tests:** `TEST-INTERNAL_MODEL-006`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Public Functions:**

---

## SPEC-INTERNAL_MODEL-029: IdentifierSet.Count

## SPEC-INTERNAL_MODEL-030: IdentifierSet.All

## SPEC-INTERNAL_MODEL-007: IdentifierSet Merge Operation

**Design:** `ModelModule`

**Contract:** `IdentifierSet`

**Requirement:**

Combine two identifier sets.

**Tests:** `TEST-INTERNAL_MODEL-007`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Public Functions:**

---

## SPEC-INTERNAL_MODEL-031: IdentifierSet.Merge

## SPEC-INTERNAL_MODEL-008: Annotation Creation

**Design:** `ModelModule`

**Contract:** `NewAnnotation`

**Requirement:**

Create annotation from code.

**Tests:** `TEST-INTERNAL_MODEL-008`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Public Functions:**

---

## SPEC-INTERNAL_MODEL-032: Annotation.NewAnnotation

## SPEC-INTERNAL_MODEL-009: Annotation Conversion

**Design:** `ModelModule`

**Contract:** `Annotation`

**Requirement:**

Convert annotation to identifier.

**Tests:** `TEST-INTERNAL_MODEL-009`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Public Functions:**

---

## SPEC-INTERNAL_MODEL-033: Annotation.ToIdentifier

## SPEC-INTERNAL_MODEL-010: Validation Result Types

**Design:** `ModelModule`

**Contract:** `ValidationResult`

**Requirement:**

Validation result stores errors and warnings.

**Tests:** `TEST-INTERNAL_MODEL-010`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Key Types:**

- `ValidationResult` — Result with errors and warnings
- `ValidationError` — Error with rule, message, source, link

---

---

## SPEC-INTERNAL_MODEL-011: Validation Result Operations

**Design:** `ModelModule`

**Contract:** `ValidationResult`

**Requirement:**

Add errors/warnings and sort results.

**Tests:** `TEST-INTERNAL_MODEL-011`, `TEST-INTERNAL_MODEL-016`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Public Functions:**

## SPEC-INTERNAL_MODEL-034: ValidationResult.AddError

## SPEC-INTERNAL_MODEL-035: ValidationResult.AddWarning

## SPEC-INTERNAL_MODEL-036: ValidationResult.Sort

**Tests:** `TEST-INTERNAL_MODEL-011`, `TEST-INTERNAL_MODEL-016`

## SPEC-INTERNAL_MODEL-012: Link Type Definition

**Design:** `ModelModule`

**Contract:** `LinkType`

**Requirement:**

LinkType enum for edge types.

**Tests:** `TEST-INTERNAL_MODEL-012`, `TEST-INTERNAL_MODEL-028`

---

**Status:** Done

**Implementation:** `internal/model/link.go`

**Key Types:**

- `LinkType` — Enum for LinkTests, LinkImplements, LinkReferences

---

## SPEC-INTERNAL_MODEL-013: Link Structure

**Design:** `ModelModule`

**Contract:** `Link`

**Requirement:**

Link represents a reference from one identifier to another.

**Tests:** `TEST-INTERNAL_MODEL-013`, `TEST-INTERNAL_MODEL-029`, `TEST-INTERNAL_MODEL-030`

**Status:** Done

**Implementation:** `internal/model/link.go`

**Key Types:**

- `Link` — Reference with type and target ID

---

## SPEC-INTERNAL_MODEL-014: Origin Type Definition

**Design:** `ModelModule`

**Contract:** `Origin`

**Requirement:**

Origin represents the source location type of an identifier.

**Tests:** `TEST-INTERNAL_MODEL-014`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Key Types:**

- `Origin` — Enum type with values `Doc` and `Code`
- `OriginDoc` — Indicates identifier came from documentation
- `OriginCode` — Indicates identifier came from source code

---

## SPEC-INTERNAL_MODEL-015: Origin Operations

**Design:** `ModelModule`

**Contract:** `Origin`

**Requirement:**

Set and check identifier origin.

**Tests:** `TEST-INTERNAL_MODEL-015`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Methods:**

- `SetOrigin(o Origin)`
- `Origin() Origin`

---

## SPEC-INTERNAL_MODEL-017: Identifier Origin Methods

**Design:** `ModelModule`

**Contract:** `SetOrigin`

**Requirement:**

Methods to get/set identifier origin.

**Tests:** `TEST-INTERNAL_MODEL-017`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Methods:**

- `SetOrigin(origin Origin)` — implemented as part of `SPEC-INTERNAL_MODEL-017`
- `GetOrigin() Origin` — not present in code (only SetOrigin exists)

**Note:** `GetOrigin()` method does not exist in code. Only `SetOrigin()` is implemented.

---

## SPEC-INTERNAL_MODEL-018: Identifier Model Core

**Design:** `ModelModule`

**Contract:** `Identifier`

**Requirement:**

Core identifier type with all fields.

**Tests:** `TEST-INTERNAL_MODEL-018`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Key Types:**

- `Identifier` — IDD identifier with type, module, number, links

**Tests:** ``
