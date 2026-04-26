---
related_files:
  spec: docs/internal/model/spec.md
  contract: docs/internal/model/contract.md
  design: docs/internal/model/design.md
  testing: docs/internal/model/testing.md
---

# Contract (model)

**Status:** Done

**Requirement:**

The model package provides data structures for representing IDD identifiers, annotations, and the identifier set collection.

**Key Contracts:**

- Identifiers must have unique IDs within the set
- Links must have valid target IDs
- Annotations must convert to identifiers correctly
- Validation results must accumulate errors and warnings

**Implementation:** `internal/model/identifier.go`, `internal/model/link.go`

**Acceptance Criteria:**

- [x] Identifiers store type, module, number, and local ID
- [x] Forward links connect identifiers to their dependencies
- [x] Backlinks are computed from forward links
- [x] IdentifierSet supports add, get, has, count, all, merge operations
- [x] Annotations can be converted to identifiers

**Related:** `SPEC-INTERNAL_MODEL-001` through `SPEC-INTERNAL_MODEL-018`
