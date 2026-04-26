---
related_files:
  spec: docs/internal/model/spec.md
  contract: docs/internal/model/contract.md
  design: docs/internal/model/design.md
  testing: docs/internal/model/testing.md
---

# Design (model)

**Status:** Done

## Architecture

The model module defines core data types:

1. **Value Types**: Identifier and Link are designed as value types with value semantics
2. **Slice-Based Collection**: IdentifierSet uses slices for ordered iteration
3. **Map-Based Lookup**: byID map enables O(1) identifier lookup
4. **Origin Tracking**: Doc vs Code origin enables validation of doc-link-consistency

## Package Layout

```go
internal/model/
├── identifier.go     # Identifier and IdentifierSet types
├── link.go           # Link type and LinkType enum
└── origin.go          # Origin enum (Doc/Code)
```

## Function Composition

1. **NewIdentifier()** - Creates an identifier with type, module, number
2. **NewLink()** - Creates a directed link between two identifiers
3. **IdentifierSet.Add()** - Adds identifier to set
4. **IdentifierSet.Has()** - Checks if identifier exists in set
5. **LinkType.String()** - Converts link type to string representation

## Testability Hooks

- All types are plain Go structs with no I/O
- Constructor functions allow controlled test data creation
- Value semantics simplify equality checking

## Dependencies

- No external dependencies (self-contained)
