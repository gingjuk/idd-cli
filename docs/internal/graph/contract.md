---
idd:
  version: "1.0"
  package: internal/graph
  document: contract
---

# Contracts: internal/graph

## Contract: GraphMutation

**Guarantees:**

`NewLinkageGraph` returns a non-nil empty graph with initialized node, edge, ID,
type, and backlink collections.

`AddNode(id, type)` creates a node only when the ID is absent. Repeating the ID
returns the existing pointer; the original type and metadata remain unchanged.
New nodes start with non-nil metadata and empty adjacency.

`AddEdge(from, to, type, source, line)` appends a directed edge and target
backlink without deduplication. Existing source and target nodes receive
outgoing and incoming adjacency respectively. Missing endpoints do not produce
an error and later node insertion does not repair adjacency.

All mutation is synchronous and in-memory. No method is safe for concurrent
writers.

## Contract: GraphQuery

**Guarantees:**

Node lookup returns a node pointer and presence boolean. Type-filtered inbound
and outbound queries return nil for a missing node and a newly built slice for
an existing node's matches. Backlink lookup returns the stored slice when
present and a non-nil empty slice otherwise.

`Nodes`, `Edges`, `Node.InEdges`, `Node.OutEdges`, and successful backlink
queries expose underlying mutable collections. Callers must treat them as
borrowed read-only views unless they intentionally accept invariant risk.

No query sorts its result. Node-map and snapshot-node order are unspecified;
global and adjacency edge order follows insertion.

## Contract: RelationshipVerification

**Guarantees:**

`VerifyBidirectionalLinks` recalculates every edge's `Verified` flag. An edge is
verified when an opposite-direction edge exists whose type is
`model.ReverseLinkType` of the original. The method does not add missing reverse
edges or emit validation errors.

`ValidateCompleteness` reports a `spec-missing-tests` error for each SPEC node
without an outbound `tests` edge and a `test-missing-coverage` error for each
TEST node without an outbound `implements` edge. It does not require
verification, inspect other node types, or validate endpoint existence.

## Contract: GraphProjection

**Guarantees:**

`Stats` reports total nodes and edges plus counts by the four identifier types.
Unknown identifier types contribute only to the total.

`ToSnapshot` allocates serializable node and edge summary slices. It copies
counts and scalar edge evidence, not metadata or adjacency. Node order is
unspecified; edge order matches insertion. Snapshot mutation does not mutate the
graph.
