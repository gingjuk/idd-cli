---
markers:
  - id: DESIGN-BE-001
    name: Graph-First Architecture
  - id: DESIGN-BE-002
    name: Validation Rule Engine
---

# Design Examples

This directory contains example IDD design documents showing the expected format.

## DESIGN-BE-001: Graph-First Architecture

**Date:** 2026-04-18
**Status:** Accepted

### Context

The system needs to support multiple validation rules, multiple collector types, and bidirectional relationship checking. We need an architecture that allows easy extension.

### Decision

We adopted a **graph-first architecture** where:

1. **LinkageGraph** is the central data structure
2. All collectors feed identifiers into the graph
3. Validators operate on the graph
4. Reporters consume validation results

### Architecture

```
┌─────────────────────────────────────┐
│           CLI (main.go)              │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│        Validation Engine             │
└─────────────────────────────────────┘
        │           │           │
        ▼           ▼           ▼
┌────────────┐ ┌────────┐ ┌──────────┐
│ Collector   │ │ Linker │ │ Validator │
└────────────┘ └────────┘ └──────────┘
                                     │
                                     ▼
                             ┌────────────┐
                             │ Reporter   │
                             └────────────┘
```

### Consequences

**Benefits:**
- Easy to extend — Add new collectors by implementing interface
- Independent validators — Each rule is a separate struct
- Natural representation — Graphs model bidirectional relationships naturally
- Testable — Each component can be tested in isolation

**Related:** [SPEC-BE-001](../spec.md#spec-be-001-user-authentication)

---

## DESIGN-BE-002: Validation Rule Engine

**Date:** 2026-04-18
**Status:** Proposed

### Context

We need a flexible validation engine that can run different rule sets based on configuration.

### Decision

Validation rules are defined in `idd.yaml` and executed by the Engine in order:

1. **Completeness** — Check spec-test coverage
2. **Bidirectional** — Verify backlink existence
3. **Orphan** — Detect unlinked identifiers

### Rule Interface

```go
type Rule interface {
    Name() string
    Validate(*graph.LinkageGraph) []ValidationError
}
```

**Related:** [SPEC-BE-001](../spec.md#spec-be-001-user-authentication)
