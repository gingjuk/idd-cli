# Design: DESIGN-BE-001

## Graph-First Architecture

**Date:** 2026-04-18
**Status:** Accepted

## Context

The IDD Link Validator needs to support multiple validation rules, multiple collector types, and bidirectional relationship checking. We needed an architecture that allows easy extension.

## Decision

We adopted a **graph-first architecture** where:

1. **LinkageGraph** is the central data structure
2. All collectors feed identifiers into the graph
3. Validators operate on the graph
4. Reporters consume validation results

## Architecture

```
┌─────────────────────────────────────┐
│           CLI (main.go)              │
└─────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────┐
│        Validation Engine              │
│  Orchestrates:                       │
│  collection → linking → validation  │
└─────────────────────────────────────┘
        │           │           │
        ▼           ▼           ▼
┌────────────┐ ┌────────┐ ┌──────────┐
│ Collector   │ │ Linker │ │ Validator │
│             │ │        │ │           │
│ DocCollector│ │BuildGr │ │Bidirect.  │
│ CodeCollector│ │ph      │ │Orphan     │
│ (extensible)│ │        │ │Consist.   │
└────────────┘ └────────┘ └──────────┘
                                     │
                                     ▼
                             ┌────────────┐
                             │ Reporter  │
                             │ JSON output│
                             └────────────┘
```

## Consequences

### Benefits

1. **Easy to extend** — Add new collectors by implementing interface
2. **Independent validators** — Each rule is a separate struct implementing `ValidationRule`
3. **Natural representation** — Graphs naturally model bidirectional relationships
4. **Testable** — Each component can be tested in isolation

### Drawbacks

1. **Memory overhead** — Graph structure requires more memory than flat representation
2. **Index maintenance** — Must keep indexes updated as graph changes

## Related Documents

- **SPEC-BE-001** — IDD Link Validator Overview
- **SPEC-BE-002** — Graph Linkage Structure
