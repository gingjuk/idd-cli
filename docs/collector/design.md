---
markers:
  - id: DESIGN-COLLECTOR-001
    name: Collector Architecture
---

# Design (collector)

## DESIGN-COLLECTOR-001: Collector Architecture

**Date:** 2026-04-19
**Status:** Accepted

### Context

idd-cli needs to collect IDD identifiers from two distinct sources: markdown documentation and source code. Each source has different formats, parsing requirements, and validation rules.

### Decision

We adopted a **dual-collector architecture** with shared utilities:

1. **DocCollector** — Specialized for markdown documentation
2. **CodeCollector** — Specialized for source code annotations
3. **Shared frontmatter utilities** — Common validation and parsing

### Architecture

```text
┌─────────────────────────────────────────────────────────┐
│                    Collector Package                      │
├─────────────────────────────────────────────────────────┤
│  ┌─────────────────┐     ┌─────────────────────────┐   │
│  │  DocCollector   │     │     CodeCollector       │   │
│  │                 │     │                         │   │
│  │ - Markdown      │     │ - Go (.go)              │   │
│  │ - Frontmatter   │     │ - TypeScript (.ts/.tsx) │   │
│  │ - IDD refs      │     │ - JavaScript (.js)      │   │
│  └─────────────────┘     └─────────────────────────┘   │
│              │                       │                  │
│              └───────────┬───────────┘                  │
│                          ▼                               │
│              ┌─────────────────────┐                   │
│              │   frontmatter.go     │                   │
│              │                     │                   │
│              │ - ParseFrontmatter  │                   │
│              │ - ValidateMarkers   │                   │
│              │ - ModulePrefix      │                   │
│              └─────────────────────┘                   │
└─────────────────────────────────────────────────────────┘
```

### Key Design Decisions

#### 1. Separate Collectors by Source Type

**Rationale:** Documentation and code have fundamentally different structures:

- Docs: hierarchical with frontmatter metadata
- Code: linear with annotation prefixes

**Tradeoff:** Duplication of directory walking logic → Mitigated by shared `shouldIgnore()` and `matchGlob()` patterns.

#### 2. Frontmatter-First Document Processing

**Rationale:** Frontmatter defines the "owned" identifiers, content references are "links" to those.

**Flow:**

1. Parse frontmatter → extract defined markers
2. Scan content → extract referenced IDs
3. Create nodes for defined markers
4. Add links from nodes to content references

#### 3. Function Context Extraction for Code

**Rationale:** Code annotations need more context than docs because:

- Annotations are often on preceding lines
- Function names clarify intent
- Comments before annotation provide documentation

**Implementation:** `extractFunctionComment()` searches lines after annotation for comments and function declaration.

### Consequences

**Benefits:**

- Clear separation of concerns
- Easy to add new source types (e.g., Python, Rust)
- Shared validation logic ensures consistency
- Testable in isolation

**Drawbacks:**

- Some code duplication in file walking
- Frontmatter parsing is YAML-dependent

**Related Specs:** `SPEC-COLLECTOR-001`, `SPEC-COLLECTOR-002`, `SPEC-BE-001`, `DESIGN-BE-001`
