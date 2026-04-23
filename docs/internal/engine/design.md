---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Design (engine)

**Status:** Done

## Architecture

The engine orchestrates the validation pipeline:

1. **Pipeline Pattern**: Engine uses a three-phase pipeline: Collect → Build Graph → Validate
2. **Graph-Centric**: All identifiers are stored in a linkage graph for efficient traversal
3. **Result Aggregation**: Errors from all validation rules are aggregated in ValidationResult

## Package Layout

```text
internal/engine/
├── engine.go         # Main engine implementation
├── collect.go        # Collection orchestration
├── graph.go          # Graph building logic
└── validate.go       # Validation rule execution
```

## Function Composition

1. **Run()** - Main entry point that executes full validation pipeline
2. **Collect()** - Orchestrates doc and code collection
3. **BuildGraph()** - Constructs the linkage graph from collected identifiers
4. **Validate()** - Runs all validation rules on the graph
5. **Report()** - Formats and outputs validation results

## Testability Hooks

- Engine accepts mock collectors for testing
- Graph can be constructed directly for unit testing validators
- ValidationResult can be inspected without I/O

## Dependencies

- `internal/collector` - For doc and code collection
- `internal/graph` - For LinkageGraph structure
- `internal/model` - For Identifier types
- `internal/reporter` - For output formatting
