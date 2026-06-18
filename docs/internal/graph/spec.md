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

**Design:** `GraphModule`

**Contract:** `Node`

**Requirement:**

Node represents an IDD identifier within the linkage graph with incoming and outgoing edges.

**Tests:** `TEST-INTERNAL_GRAPH-001`, `TEST-INTERNAL_GRAPH-002`, `TEST-INTERNAL_GRAPH-003`, `TEST-INTERNAL_GRAPH-004`, `TEST-INTERNAL_GRAPH-005`, `TEST-INTERNAL_GRAPH-006`, `TEST-INTERNAL_GRAPH-007`, `TEST-INTERNAL_GRAPH-008`, `TEST-INTERNAL_GRAPH-009`, `TEST-INTERNAL_GRAPH-010`, `TEST-INTERNAL_GRAPH-011`, `TEST-INTERNAL_GRAPH-012`, `TEST-INTERNAL_GRAPH-013`, `TEST-INTERNAL_GRAPH-014`, `TEST-INTERNAL_GRAPH-015`

**Status:** Done

**Implementation:** `internal/graph/graph.go`

**Key Types:**

- `Node` — Identifier node with edges and metadata

**Acceptance Criteria:**

- [x] Node stores ID, type, and metadata
- [x] Node tracks incoming and outgoing edges
- [x] Node provides methods to access edges

**Related:** `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`

## SPEC-INTERNAL_GRAPH-002: Graph Edge Structure

**Design:** `GraphModule`

**Contract:** `Edge`

**Requirement:**

Edge represents a directed relationship between two nodes in the linkage graph.

**Tests:** `TEST-INTERNAL_GRAPH-001`, `TEST-INTERNAL_GRAPH-002`, `TEST-INTERNAL_GRAPH-003`, `TEST-INTERNAL_GRAPH-004`, `TEST-INTERNAL_GRAPH-005`, `TEST-INTERNAL_GRAPH-006`, `TEST-INTERNAL_GRAPH-007`, `TEST-INTERNAL_GRAPH-008`, `TEST-INTERNAL_GRAPH-009`, `TEST-INTERNAL_GRAPH-010`, `TEST-INTERNAL_GRAPH-011`, `TEST-INTERNAL_GRAPH-012`, `TEST-INTERNAL_GRAPH-013`, `TEST-INTERNAL_GRAPH-014`, `TEST-INTERNAL_GRAPH-015`

**Status:** Done

**Implementation:** `internal/graph/graph.go`

**Key Types:**

- `Edge` — Directed edge with type, source, and verification status

**Acceptance Criteria:**

- [x] Edge stores from/to node IDs
- [x] Edge stores link type (tests, implements, references)
- [x] Edge tracks source file and line number
- [x] Edge tracks verification status

**Related:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-003`

---

---

## SPEC-INTERNAL_GRAPH-003: Linkage Graph Structure

**Design:** `GraphModule`

**Contract:** `LinkageGraph`

**Requirement:**

LinkageGraph manages nodes and edges for IDD identifier validation.

**Tests:** `TEST-INTERNAL_GRAPH-001`, `TEST-INTERNAL_GRAPH-002`, `TEST-INTERNAL_GRAPH-003`, `TEST-INTERNAL_GRAPH-004`, `TEST-INTERNAL_GRAPH-005`, `TEST-INTERNAL_GRAPH-006`, `TEST-INTERNAL_GRAPH-007`, `TEST-INTERNAL_GRAPH-008`, `TEST-INTERNAL_GRAPH-009`, `TEST-INTERNAL_GRAPH-010`, `TEST-INTERNAL_GRAPH-011`, `TEST-INTERNAL_GRAPH-012`, `TEST-INTERNAL_GRAPH-013`, `TEST-INTERNAL_GRAPH-014`, `TEST-INTERNAL_GRAPH-015`

**Status:** Done

**Implementation:** `internal/graph/graph.go`

**Key Types:**

- `LinkageGraph` — Main graph structure with nodes, edges, and index

**Public Functions:**

## SPEC-INTERNAL_GRAPH-006: NewLinkageGraph

**Design:** `GraphModule`

**Contract:** `NewLinkageGraph`

**Requirement:**

Create a new empty linkage graph.

**Tests:** `TEST-INTERNAL_GRAPH-006`

**Function Signature:**
`func NewLinkageGraph() *LinkageGraph`

**Purpose:** Creates a new empty linkage graph.

---

## SPEC-INTERNAL_GRAPH-007: LinkageGraph.AddNode

**Design:** `GraphModule`

**Contract:** `AddNode`

**Requirement:**

Add a node to the graph if it does not already exist.

**Tests:** `TEST-INTERNAL_GRAPH-007`

**Function Signature:**
`func (g *LinkageGraph) AddNode(id string, idType model.IdentifierType) *Node`

**Purpose:** Adds a node to the graph if it doesn't exist.

---

## SPEC-INTERNAL_GRAPH-008: LinkageGraph.AddEdge

**Design:** `GraphModule`

**Contract:** `AddEdge`

**Requirement:**

Add a directed edge between two graph nodes.

**Tests:** `TEST-INTERNAL_GRAPH-008`

**Function Signature:**
`func (g *LinkageGraph) AddEdge(from, to string, edgeType model.LinkType, source string, line int)`

**Purpose:** Adds a directed edge between two nodes.

---

## SPEC-INTERNAL_GRAPH-004: Graph Index Structure

**Design:** `GraphModule`

**Contract:** `Index`

**Requirement:**

Index provides fast lookup structures for nodes by ID, type, and backlinks.

**Tests:** `TEST-INTERNAL_GRAPH-001`, `TEST-INTERNAL_GRAPH-002`, `TEST-INTERNAL_GRAPH-003`, `TEST-INTERNAL_GRAPH-004`, `TEST-INTERNAL_GRAPH-005`, `TEST-INTERNAL_GRAPH-006`, `TEST-INTERNAL_GRAPH-007`, `TEST-INTERNAL_GRAPH-008`, `TEST-INTERNAL_GRAPH-009`, `TEST-INTERNAL_GRAPH-010`, `TEST-INTERNAL_GRAPH-011`, `TEST-INTERNAL_GRAPH-012`, `TEST-INTERNAL_GRAPH-013`, `TEST-INTERNAL_GRAPH-014`, `TEST-INTERNAL_GRAPH-015`

**Status:** Done

**Implementation:** `internal/graph/graph.go`

**Key Types:**

- `Index` — Fast lookup indexes for graph traversal

**Acceptance Criteria:**

- [x] Index stores nodes by ID for O(1) lookup
- [x] Index stores nodes by type for filtering
- [x] Index stores backlinks for reverse edge lookup

**Related:** `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-005`

## SPEC-INTERNAL_GRAPH-005: Graph Factory Function

**Design:** `GraphModule`

**Contract:** `NewLinkageGraph`

**Requirement:**

Factory function to create a new LinkageGraph instance.

**Tests:** `TEST-INTERNAL_GRAPH-009`, `TEST-INTERNAL_GRAPH-010`, `TEST-INTERNAL_GRAPH-011`, `TEST-INTERNAL_GRAPH-012`, `TEST-INTERNAL_GRAPH-013`, `TEST-INTERNAL_GRAPH-014`, `TEST-INTERNAL_GRAPH-015`

**Status:** Done

**Implementation:** `internal/graph/graph.go`

**Public Functions:**

---

## SPEC-INTERNAL_GRAPH-009: NewLinkageGraph

**Design:** `GraphModule`

**Contract:** `LinkageGraph`

**Requirement:**

Create a new empty linkage graph with initialized maps.

**Tests:** `TEST-INTERNAL_GRAPH-009`

---

**Function Signature:**
`func NewLinkageGraph() *LinkageGraph`

**Purpose:** Creates a new empty linkage graph with initialized maps.

**Returns:** A new LinkageGraph pointer ready to accept nodes and edges

**Acceptance Criteria:**

- [x] Creates empty node map
- [x] Creates empty edge slice
- [x] Initializes index with empty lookup maps

**Related:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`
