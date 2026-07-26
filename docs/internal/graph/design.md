---
idd:
  version: "1.0"
  package: internal/graph
  document: design
---

# Design: internal/graph

## Component: LinkageGraphStore

**Purpose:**

`LinkageGraphStore` is the engine's in-memory projection of IDD identifiers and
their directed relationships. It makes forward and reverse traversal cheap
enough for validation rules, attaches source evidence to edges, marks
relationships that have a reverse counterpart, and exports statistics or a
serializable snapshot for reports.

The graph is a mutable working structure, not the canonical source of
documentation. Collectors create identifiers and relationship references;
`Engine.buildGraph` decides which nodes and edges to add. This package does not
parse Markdown, infer link types, apply configuration, or format findings.

### Responsibilities

- store one node per identifier string;
- retain identifier type and mutable metadata on the first node creation;
- append directed edges in insertion order with source and line evidence;
- maintain per-node incoming and outgoing adjacency when endpoints exist;
- maintain a reverse lookup from target ID to source IDs;
- query edges by direction and relationship type;
- mark edges that have a correctly typed reverse edge;
- calculate graph-level validation statistics and basic SPEC/TEST completeness;
  and
- create report snapshots without exposing private adjacency fields.

### Mutation and ownership boundaries

All maps and slices are package-owned mutable state, but several query methods
return the underlying collections directly. Callers inside the repository can
therefore mutate nodes, metadata, edge slices, and backlink slices without a
copy. The design relies on disciplined single-threaded engine construction and
read-mostly validation after the graph is built.

There is no locking. A graph must not be mutated concurrently with queries,
verification, snapshotting, or statistics.

### Edge and endpoint behavior

`AddEdge` always appends a global edge and backlink, even if one or both endpoint
nodes do not yet exist. It attaches adjacency only to nodes present at insertion
time; adding a missing node later does not retroactively attach that edge.
Duplicate edges are allowed. The engine therefore creates all known nodes
before adding relationships and owns endpoint integrity.

### Decisions and trade-offs

A map provides direct node lookup, adjacency slices preserve insertion order,
and the backlink index avoids scanning all edges for reverse references. This
duplicates relationship state across the global edge list, node adjacency, and
index. The cost is acceptable for validation-sized graphs, but mutation must go
through `AddEdge` to keep those views coherent.

`VerifyBidirectionalLinks` records verification on each edge instead of
rejecting immediately. That lets the engine report complete evidence and decide
which link types require reciprocity.

## Architecture

```text
IdentifierSet
     |
     v
Engine.buildGraph
     |
     +--> AddNode(id, type) ---------------------+
     |                                           |
     +--> AddEdge(from, to, type, source, line)  |
             |                 |                 |
             v                 v                 v
        global edges       node adjacency    backlink index
             |
             +--> VerifyBidirectionalLinks
             +--> validation queries
             +--> Stats / ToSnapshot
```

The graph does not own validation policy beyond its narrow
`ValidateCompleteness` helper. Most rule selection and source-specific
diagnostics remain in the engine.

## Package Layout

`graph.go` contains nodes, edges, indexes, mutation, traversal, verification,
statistics, and snapshot projection. These operations remain together because
they must preserve the invariants of one mutable representation.

Relationship vocabulary and serializable report types live in
`internal/model`; this package depends on those value types but does not define
their wire format.

## Function Composition

Construction initializes every internal collection. Node insertion updates the
node map plus ID and type indexes. Edge insertion updates the global list,
available endpoint adjacency, and target backlink index.

Verification maps each edge type through `model.ReverseLinkType`, inspects
incoming edges to the original source, and marks the edge when the opposite
endpoint and reverse type are present. Snapshot and statistics then read the
resulting state without changing it.

## Dependencies

The package uses `internal/model` for identifier types, link types, validation
errors, statistics, and graph snapshot structures. Its only other dependency is
`fmt` for completeness messages.

## Testability Hooks

Tests build graphs entirely in memory and inspect public queries, node
adjacency, verification flags, statistics, snapshots, and completeness errors.
No collectors, filesystem, configuration, or reporter is needed.

The suite does not cover edges added before nodes, duplicate edges,
first-type-wins node insertion, mutation through returned maps/slices,
deterministic snapshot node ordering, or concurrent access. These are important
boundaries for future hardening.
