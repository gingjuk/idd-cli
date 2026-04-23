---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Design (graph)

**Status:** Done

## Architecture

The graph module provides the central data structure for tracking identifier relationships:

1. **Map-Based Node Storage**: Nodes stored in map for O(1) lookup by ID
2. **Slice-Based Edge Storage**: Edges stored in slice for ordered traversal
3. **Pre-built Index**: Backlinks index maintained for fast reverse edge lookup
4. **Bidirectional Verification**: Edge verification marks reverse link existence

## Package Layout

```text
internal/graph/
├── graph.go          # LinkageGraph implementation
├── node.go           # Node structure
└── edge.go           # Edge structure and types
```

## Function Composition

1. **AddNode()** - Adds an identifier node to the graph
2. **AddLink()** - Creates a directed link between nodes
3. **GetNode()** - Retrieves node by ID (O(1) lookup)
4. **Nodes()** - Returns all nodes for iteration
5. **Links()** - Returns all links for the graph
6. **Verify()** - Checks doc-link-consistency

## Testability Hooks

- Graph can be constructed with test data without I/O
- Node and Link structures are plain Go types
- Iterator methods allow inspection without modification

## Dependencies

- `internal/model` - For Identifier and Link types
- No external dependencies (self-contained)
