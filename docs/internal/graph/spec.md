---
markers:
  - id: SPEC-INTERNAL_GRAPH-001
    name: Graph Node Structure
  - id: SPEC-INTERNAL_GRAPH-002
    name: Graph Edge Structure
  - id: SPEC-INTERNAL_GRAPH-003
    name: Linkage Graph Structure
  - id: SPEC-INTERNAL_GRAPH-004
    name: Graph Index Structure
  - id: SPEC-INTERNAL_GRAPH-005
    name: Graph Factory Function
  - id: SPEC-INTERNAL_GRAPH-006
    name: LinkageGraph.NewLinkageGraph
  - id: SPEC-INTERNAL_GRAPH-007
    name: LinkageGraph.AddNode
  - id: SPEC-INTERNAL_GRAPH-008
    name: LinkageGraph.AddEdge
  - id: SPEC-INTERNAL_GRAPH-009
    name: NewLinkageGraph

related_files:
  spec: docs/internal/graph/spec.md
  contract: docs/internal/graph/contract.md
  design: docs/internal/graph/design.md
  testing: docs/internal/graph/testing.md
---

# Specification (graph)

## SPEC-INTERNAL_GRAPH-001: Graph Node Structure

**Contract:** `Node`

**Design:** `GraphModule`

**Status:** Done

**Requirement:**

Node represents an IDD identifier within the linkage graph with incoming and outgoing edges.

**Implementation:** `internal/graph/graph.go`

**Key Types:**

- `Node` — Identifier node with edges and metadata

**Acceptance Criteria:**

- [x] Node stores ID, type, and metadata
- [x] Node tracks incoming and outgoing edges
- [x] Node provides methods to access edges

**Tests:** `TEST-INTERNAL_GRAPH-001`, `TEST-INTERNAL_GRAPH-002`, `TEST-INTERNAL_GRAPH-003`, `TEST-INTERNAL_GRAPH-004`, `TEST-INTERNAL_GRAPH-005`, `TEST-INTERNAL_GRAPH-006`, `TEST-INTERNAL_GRAPH-007`, `TEST-INTERNAL_GRAPH-008`, `TEST-INTERNAL_GRAPH-009`, `TEST-INTERNAL_GRAPH-010`, `TEST-INTERNAL_GRAPH-011`, `TEST-INTERNAL_GRAPH-012`, `TEST-INTERNAL_GRAPH-013`, `TEST-INTERNAL_GRAPH-014`, `TEST-INTERNAL_GRAPH-015`

**Related:** `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`

---

## SPEC-INTERNAL_GRAPH-002: Graph Edge Structure

**Contract:** `Edge`

**Design:** `GraphModule`

**Status:** Done

**Requirement:**

Edge represents a directed relationship between two nodes in the linkage graph.

**Implementation:** `internal/graph/graph.go`

**Key Types:**

- `Edge` — Directed edge with type, source, and verification status

**Acceptance Criteria:**

- [x] Edge stores from/to node IDs
- [x] Edge stores link type (tests, implements, references)
- [x] Edge tracks source file and line number
- [x] Edge tracks verification status

**Tests:** `TEST-INTERNAL_GRAPH-001`, `TEST-INTERNAL_GRAPH-002`, `TEST-INTERNAL_GRAPH-003`, `TEST-INTERNAL_GRAPH-004`, `TEST-INTERNAL_GRAPH-005`, `TEST-INTERNAL_GRAPH-006`, `TEST-INTERNAL_GRAPH-007`, `TEST-INTERNAL_GRAPH-008`, `TEST-INTERNAL_GRAPH-009`, `TEST-INTERNAL_GRAPH-010`, `TEST-INTERNAL_GRAPH-011`, `TEST-INTERNAL_GRAPH-012`, `TEST-INTERNAL_GRAPH-013`, `TEST-INTERNAL_GRAPH-014`, `TEST-INTERNAL_GRAPH-015`

**Related:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-003`

---

## SPEC-INTERNAL_GRAPH-003: Linkage Graph Structure

**Contract:** `LinkageGraph`

**Design:** `GraphModule`

**Status:** Done

**Requirement:**

LinkageGraph manages nodes and edges for IDD identifier validation.

**Implementation:** `internal/graph/graph.go`

**Key Types:**

- `LinkageGraph` — Main graph structure with nodes, edges, and index

**Tests:** `TEST-INTERNAL_GRAPH-001`, `TEST-INTERNAL_GRAPH-002`, `TEST-INTERNAL_GRAPH-003`, `TEST-INTERNAL_GRAPH-004`, `TEST-INTERNAL_GRAPH-005`, `TEST-INTERNAL_GRAPH-006`, `TEST-INTERNAL_GRAPH-007`, `TEST-INTERNAL_GRAPH-008`, `TEST-INTERNAL_GRAPH-009`, `TEST-INTERNAL_GRAPH-010`, `TEST-INTERNAL_GRAPH-011`, `TEST-INTERNAL_GRAPH-012`, `TEST-INTERNAL_GRAPH-013`, `TEST-INTERNAL_GRAPH-014`, `TEST-INTERNAL_GRAPH-015`

**Public Functions:**

## SPEC-INTERNAL_GRAPH-006: NewLinkageGraph

**Function Signature:**
`func NewLinkageGraph() *LinkageGraph`

**Purpose:** Creates a new empty linkage graph.

**Tests:** `TEST-INTERNAL_GRAPH-006`

---

## SPEC-INTERNAL_GRAPH-007: LinkageGraph.AddNode

**Function Signature:**
`func (g *LinkageGraph) AddNode(id string, idType model.IdentifierType) *Node`

**Purpose:** Adds a node to the graph if it doesn't exist.

**Tests:** `TEST-INTERNAL_GRAPH-007`

---

## SPEC-INTERNAL_GRAPH-008: LinkageGraph.AddEdge

**Function Signature:**
`func (g *LinkageGraph) AddEdge(from, to string, edgeType model.LinkType, source string, line int)`

**Purpose:** Adds a directed edge between two nodes.

**Tests:** `TEST-INTERNAL_GRAPH-008`

---

## SPEC-INTERNAL_GRAPH-004: Graph Index Structure

**Contract:** `Index`

**Design:** `GraphModule`

**Status:** Done

**Requirement:**

Index provides fast lookup structures for nodes by ID, type, and backlinks.

**Implementation:** `internal/graph/graph.go`

**Key Types:**

- `Index` — Fast lookup indexes for graph traversal

**Acceptance Criteria:**

- [x] Index stores nodes by ID for O(1) lookup
- [x] Index stores nodes by type for filtering
- [x] Index stores backlinks for reverse edge lookup

**Tests:** `TEST-INTERNAL_GRAPH-001`, `TEST-INTERNAL_GRAPH-002`, `TEST-INTERNAL_GRAPH-003`, `TEST-INTERNAL_GRAPH-004`, `TEST-INTERNAL_GRAPH-005`, `TEST-INTERNAL_GRAPH-006`, `TEST-INTERNAL_GRAPH-007`, `TEST-INTERNAL_GRAPH-008`, `TEST-INTERNAL_GRAPH-009`, `TEST-INTERNAL_GRAPH-010`, `TEST-INTERNAL_GRAPH-011`, `TEST-INTERNAL_GRAPH-012`, `TEST-INTERNAL_GRAPH-013`, `TEST-INTERNAL_GRAPH-014`, `TEST-INTERNAL_GRAPH-015`

**Related:** `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-005`

---

## SPEC-INTERNAL_GRAPH-005: Graph Factory Function

**Contract:** `NewLinkageGraph`

**Design:** `GraphModule`

**Status:** Done

**Requirement:**

Factory function to create a new LinkageGraph instance.

**Implementation:** `internal/graph/graph.go`

**Tests:** `TEST-INTERNAL_GRAPH-009`, `TEST-INTERNAL_GRAPH-010`, `TEST-INTERNAL_GRAPH-011`, `TEST-INTERNAL_GRAPH-012`, `TEST-INTERNAL_GRAPH-013`, `TEST-INTERNAL_GRAPH-014`, `TEST-INTERNAL_GRAPH-015`

**Public Functions:**

## SPEC-INTERNAL_GRAPH-009: NewLinkageGraph

**Contract:** `LinkageGraph`

**Design:** `GraphModule`

**Function Signature:**
`func NewLinkageGraph() *LinkageGraph`

**Purpose:** Creates a new empty linkage graph with initialized maps.

**Returns:** A new LinkageGraph pointer ready to accept nodes and edges

**Acceptance Criteria:**

- [x] Creates empty node map
- [x] Creates empty edge slice
- [x] Initializes index with empty lookup maps

**Tests:** `TEST-INTERNAL_GRAPH-009`

---

**Related:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`
