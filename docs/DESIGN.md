---
module: BE
type: DESIGN
description: Architecture design decisions for IDD Link Validator
markers:
  - DESIGN-BE-001
---

# Design Index (BE)

| ID | Title | Status |
|----|-------|--------|
| [DESIGN-BE-001](#design-be-001) | Graph-First Architecture | Accepted |

---

## DESIGN-BE-001

### Graph-First Architecture

**Date:** 2026-04-18
**Status:** Accepted

### Context

The IDD Link Validator needs to support multiple validation rules, multiple collector types, and bidirectional relationship checking. We needed an architecture that allows easy extension.

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
│  Orchestrates:                       │
│  collection → linking → validation   │
└─────────────────────────────────────┘
        │           │           │
        ▼           ▼           ▼
┌────────────┐ ┌────────┐ ┌──────────┐
│ Collector   │ │ Linker │ │ Validator │
│             │ │        │ │           │
│ DocCollector│ │BuildGr │ │Bidirect.   │
│ CodeCollector│ │ph      │ │Orphan      │
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

**Drawbacks:**
- Memory overhead — Graph structure requires more memory
- Index maintenance — Must keep indexes updated

**Related Specs:** SPEC-BE-001, SPEC-BE-002
