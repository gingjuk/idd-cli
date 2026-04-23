---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Design (architecture)

**Status:** Done

## Architecture

The backend module follows a graph-first architecture:

1. **LinkageGraph** is the central data structure
2. All collectors feed identifiers into the graph
3. Validators operate on the graph
4. Reporters consume validation results

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

## Package Layout

```text
idd-cli/
├── cmd/idd-cli/       # CLI entry point
├── internal/
│   ├── collector/        # Doc/code identifier collection
│   ├── engine/           # Validation engine
│   ├── graph/            # Linkage graph
│   ├── model/            # Data models
│   └── reporter/         # Output formatting
└── pkg/
    ├── pattern/          # IDD regex patterns
    └── walk/             # File system traversal
```

## Function Composition

1. **CLI Run** - Entry point via cobra command
2. **Engine.Run()** - Orchestrates Collect → Build Graph → Validate
3. **Collectors** - DocCollector and CodeCollector gather identifiers
4. **Graph.Build()** - Constructs linkage graph from collected data
5. **Validators** - Run rules against the graph
6. **Reporter.Write()** - Formats and outputs results

## Testability Hooks

- Each component can be tested in isolation with interfaces
- Mock collectors allow testing without real files
- Graph can be pre-populated for validator unit tests
- Reporter tests use bytes.Buffer to capture output

## Dependencies

- `internal/collector` - For identifier collection
- `internal/engine` - For orchestration
- `internal/graph` - For graph structure
- `internal/model` - For data types
- `internal/reporter` - For output
- `internal/similarity` - For consistency checking
- `pkg/pattern` - For IDD patterns
- `pkg/walk` - For file traversal
- `gopkg.in/yaml.v3` - For config parsing
- `github.com/spf13/cobra` - For CLI
