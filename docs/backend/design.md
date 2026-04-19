---
markers:
  - id: DESIGN-BE-001
    name: Graph-First Architecture
  - id: DESIGN-BE-002
    name: Config-Driven Design
---

# Design (backend)

## DESIGN-BE-001: Graph-First Architecture

**Date:** 2026-04-18
**Status:** Accepted

### Context

idd-cli needs to support multiple validation rules, multiple collector types, and bidirectional relationship checking. We needed an architecture that allows easy extension.

### Decision

We adopted a **graph-first architecture** where:

1. **LinkageGraph** is the central data structure
2. All collectors feed identifiers into the graph
3. Validators operate on the graph
4. Reporters consume validation results

### Architecture

```text
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

**Related Specs:** `SPEC-BE-001`, `SPEC-BE-002`

---

## DESIGN-BE-002: Config-Driven Design

**Date:** 2026-04-19
**Status:** Accepted

### Config Context

idd-cli needs to support flexible configuration for different project structures, multiple languages, and various validation rules. We needed a design that allows configuration without code changes.

### Config Decision

We adopted a **config-driven design** where:

1. **YAML Configuration** — All settings in `.idd.yaml`
2. **CLI Flag Override** — Command-line flags override config file
3. **Sensible Defaults** — Works without any config file
4. **Pattern-Based Discovery** — Glob patterns for file discovery

### Configuration Structure

```yaml
# .idd.yaml
ignore_paths:
  - "examples/**"
  - "node_modules/**"

identifier_patterns:
  marker: "[A-Z]+-[A-Z]+-[0-9]+"
  annotation: "@(spec|contract|test|design)\\s+([A-Z]+-[A-Z]+-[0-9]+)"

reporter:
  format: "json"  # or "markdown"
  output: "stdout"  # or file path
```

### Config Consequences

**Benefits:**

- No recompilation needed for configuration changes
- Teams can customize validation per project
- Easy to share configuration via version control
- Supports monorepos with different rules per module

**Drawbacks:**

- Configuration errors only detected at runtime
- Additional file to maintain in project
- Potential drift between config and actual behavior

**Related Specs:** `SPEC-BE-003`
