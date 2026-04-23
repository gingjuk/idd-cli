---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Design (walk)

**Status:** Done

## Architecture

The walk package provides file system traversal:

1. **Visited Set**: Uses map to track visited directories and avoid duplicates
2. **Glob Pattern**: Uses filepath.Match for pattern matching
3. **Visitor Callback**: Uses function type for flexible file processing

## Package Layout

```text
pkg/walk/
└── files.go          # File walking implementation
```

## Function Composition

1. **Walk()** - Recursively walks directory tree matching patterns
2. **matchGlob()** - Checks if path matches glob pattern
3. **shouldIgnore()** - Determines if path matches ignore patterns
4. **visit()** - Internal callback for processing matched files

## Testability Hooks

- WalkFunc callback pattern enables testing without real files
- Mock file systems can be passed for testing
- Pattern matching is isolated in pure functions

## Dependencies

- No external dependencies (uses stdlib only)
