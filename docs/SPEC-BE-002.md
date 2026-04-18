# Specification: SPEC-BE-002

## Graph Linkage Structure

**Version:** 1.0.0
**Last Updated:** 2026-04-18

## Requirement

The LinkageGraph must efficiently represent bidirectional relationships between IDD identifiers, supporting fast lookup by ID, type, and link direction.

## Data Structures

### Node

Represents an IDD identifier in the graph:

```go
type Node struct {
    ID       string                 // Unique identifier (e.g., "SPEC-BE-001")
    Type     model.IdentifierType   // SPEC, TEST, CONTRACT, or DESIGN
    Metadata map[string]interface{} // Additional metadata
    inEdges  []*Edge               // Incoming edges (backlinks)
    outEdges []*Edge               // Outgoing edges (forward links)
}
```

### Edge

Represents a directed relationship between two identifiers:

```go
type Edge struct {
    From     string           // Source identifier
    To       string           // Target identifier
    Type     model.LinkType  // LinkTests, LinkImplements, etc.
    Source   string          // File where link was found
    Line     int             // Line number
    Verified bool            // Bidirectional link confirmed
}
```

### LinkType

| Type | Description | Reverse |
|------|-------------|---------|
| `LinkTests` | SPEC → TEST | `LinkImplements` |
| `LinkImplements` | TEST → SPEC | `LinkTests` |
| `LinkReferences` | DESIGN/CONTRACT → other | self |
| `LinkAnnotates` | Code annotation → identifier | self |

## API

### LinkageGraph Methods

| Method | Description |
|--------|-------------|
| `AddNode(id, type)` | Add identifier node |
| `AddEdge(from, to, type, source, line)` | Add directed edge |
| `GetNode(id)` | Get node by ID |
| `GetOutboundByType(nodeID, linkType)` | Get outgoing edges by type |
| `GetInboundByType(nodeID, linkType)` | Get incoming edges by type |
| `GetBacklinks(nodeID)` | Get all nodes linking TO this node |
| `Nodes()` | Get all nodes |
| `Edges()` | Get all edges |
| `VerifyBidirectionalLinks()` | Mark edges with verified=false |
| `ValidateCompleteness()` | Return errors for missing links |
| `ToSnapshot()` | Convert to report model |
| `Stats()` | Get validation statistics |

### Index

For fast lookups:

```go
type Index struct {
    byID      map[string]*Node                     // nodeID -> Node
    byType    map[model.IdentifierType][]*Node    // type -> nodes
    backlinks map[string][]string                   // nodeID -> sourceIDs
}
```

## Acceptance Criteria

- [x] Nodes store identifier ID, type, and edge lists
- [x] Edges store direction, type, source location, and verification status
- [x] Fast O(1) lookup by node ID
- [x] Fast lookup by node type
- [x] Fast lookup of backlinks (nodes linking TO a node)
- [x] Bidirectional link verification marks edges as verified/unverified

## Related Documents

- **SPEC-BE-001** — IDD Link Validator Overview
- **DESIGN-BE-001** — Graph-First Architecture Design
- **TEST-BE-001** — Core Validation Tests (edge validation tests)
