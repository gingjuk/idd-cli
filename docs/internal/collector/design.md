---
related_files:
  spec: docs/internal/collector/spec.md
  contract: docs/internal/collector/contract.md
  design: docs/internal/collector/design.md
  testing: docs/internal/collector/testing.md
---

# Design (collector)

**Status:** Done

## Architecture

The collector module uses a dual-collector architecture:

1. **DocCollector** - Specialized for markdown documentation
2. **CodeCollector** - Specialized for source code annotations
3. **Shared frontmatter utilities** - Common validation and parsing

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
│              ┌─────────────────────┐                     │
│              │   frontmatter.go    │                     │
│              │                     │                     │
│              │ - ParseFrontmatter  │                     │
│              │ - ValidateMarkers   │                     │
│              │ - ModulePrefix      │                     │
│              └─────────────────────┘                     │
└─────────────────────────────────────────────────────────┘
```

## Package Layout

```text
internal/collector/
├── collector.go      # Shared collector interface
├── doc.go            # DocCollector implementation
├── code.go           # CodeCollector implementation
└── frontmatter.go    # Frontmatter parsing utilities
```

## Function Composition

1. **Collect()** - Entry point coordinating doc and code collection
2. **ParseFrontmatter()** - Extracts markers from YAML frontmatter
3. **extractFunctionComment()** - Gets function context for code annotations
4. **shouldIgnore()** - Determines if path matches ignore patterns
5. **matchGlob()** - Checks if path matches glob pattern

## Testability Hooks

- Mock collectors implementing Collector interface for testing
- Frontmatter parsing isolated in `frontmatter.go` for unit testing
- File walking abstracted via `WalkFunc` callback pattern
- No external I/O dependencies in parsing logic

## Dependencies

- `internal/model` - For Identifier and Link types
- `internal/graph` - For building the linkage graph
- `pkg/walk` - For file system traversal
- `pkg/pattern` - For IDD regex patterns
- `gopkg.in/yaml.v3` - For YAML frontmatter parsing
