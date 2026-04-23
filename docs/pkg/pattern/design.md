---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Design (pattern)

**Status:** Done

## Architecture

The pattern package provides IDD regex patterns:

1. **Compiled Regex**: Patterns are compiled once for efficiency
2. **Filter Strategy**: Quoted identifiers filtered by checking surrounding characters
3. **Type Mapping**: Simple prefix-to-type mapping for annotation conversion

## Package Layout

```text
pkg/pattern/
└── idd.go            # IDD pattern definitions and utilities
```

## Function Composition

1. **MatchSpecifier()** - Checks if string matches spec pattern
2. **MatchTest()** - Checks if string matches test pattern
3. **MatchContract()** - Checks if string matches contract pattern
4. **ExtractIdentifiers()** - Extracts all IDD identifiers from text
5. **IsQuoted()** - Checks if identifier is in quotes (not an annotation)

## Testability Hooks

- Pure regex matching functions
- No I/O or external state
- Deterministic output for given inputs

## Dependencies

- No external dependencies (self-contained)
