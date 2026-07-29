---
idd:
  version: "1.0"
  package: internal/graph
---

# Testing: internal/graph

## TEST-INTERNAL_GRAPH-001: Node insertion and reuse

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-007`
- **Contracts:** `GraphMutation`, `GraphQuery`, `cmd/idd-cli#LinkageGraph`

**Purpose:**

Prove creation, initialized metadata, node count, ID lookup, and pointer reuse
when the same ID is added twice.

**Oracle:** The test passes only when its assertions confirm creation,
initialized metadata, node count, ID lookup, and pointer reuse when the same ID is added
twice.

## TEST-INTERNAL_GRAPH-002: Edge insertion

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-008`
- **Contracts:** `RelationshipVerification`

**Purpose:**

Prove that adding an edge between existing nodes records global edge data and
updates source/outbound and target/inbound adjacency.

**Oracle:** The test passes only when its assertions confirm adding
an edge between existing nodes records global edge data and updates source/outbound and
target/inbound adjacency.

## TEST-INTERNAL_GRAPH-003: Node lookup presence

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_GRAPH-003`
- **Contracts:** `GraphProjection`

**Purpose:**

Prove the pointer-and-boolean lookup contract for present and absent IDs.

**Oracle:** The test passes only when its assertions confirm the
pointer-and-boolean lookup contract for present and absent IDs.

## TEST-INTERNAL_GRAPH-004: Node adjacency views

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`

**Purpose:**

Prove that a directed relationship is visible in the source's outgoing and
target's incoming collections with the expected endpoint IDs.

**Oracle:** The test passes only when its assertions confirm a
directed relationship is visible in the source's outgoing and target's incoming
collections with the expected endpoint IDs.

## TEST-INTERNAL_GRAPH-005: Outbound relationship type filtering

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-003`

**Purpose:**

Prove that outbound filtering separates relationship kinds and returns nil for
a missing node.

**Oracle:** The test passes only when its assertions confirm
outbound filtering separates relationship kinds and returns nil for a missing node.

## TEST-INTERNAL_GRAPH-006: Inbound relationship type filtering

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-003`

**Purpose:**

Prove that inbound filtering returns all edges of the selected type targeting
one node.

**Oracle:** The test passes only when its assertions confirm inbound
filtering returns all edges of the selected type targeting one node.

## TEST-INTERNAL_GRAPH-007: Backlink lookup

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-004`

**Purpose:**

Prove that adding an edge makes its source directly retrievable from the
target-keyed backlink index.

**Oracle:** The test passes only when its assertions confirm adding
an edge makes its source directly retrievable from the target-keyed backlink index.

## TEST-INTERNAL_GRAPH-008: Reciprocal relationships verify

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`

**Purpose:**

Prove that opposite-direction `tests` and `implements` edges are both marked
verified after recalculation.

**Oracle:** The test passes only when its assertions confirm
opposite-direction `tests` and `implements` edges are both marked verified after
recalculation.

## TEST-INTERNAL_GRAPH-009: One-way relationships remain unverified

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`

**Purpose:**

Prove that a relationship without its typed reverse edge has `Verified` false.

**Oracle:** The test passes only when its assertions confirm a
relationship without its typed reverse edge has `Verified` false.

## TEST-INTERNAL_GRAPH-010: Snapshot projection

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-003`

**Purpose:**

Prove that snapshot projection includes the expected number of node and edge
summaries for a small graph.

**Oracle:** The test passes only when its assertions confirm
snapshot projection includes the expected number of node and edge summaries for a small
graph.

## TEST-INTERNAL_GRAPH-011: Statistics by identifier type

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-003`

**Purpose:**

Prove total and per-type counts for a graph containing SPEC, TEST, and CONTRACT
nodes.

**Oracle:** The test passes only when its assertions confirm total and
per-type counts for a graph containing SPEC, TEST, and CONTRACT nodes.

## TEST-INTERNAL_GRAPH-012: Complete SPEC and TEST links

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-003`

**Purpose:**

Prove that reciprocal SPEC/TEST coverage relationships yield no basic
completeness errors.

**Oracle:** The test passes only when its assertions confirm
reciprocal SPEC/TEST coverage relationships yield no basic completeness errors.

## TEST-INTERNAL_GRAPH-013: Missing SPEC and TEST links

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-003`

**Purpose:**

Prove that an unlinked SPEC and unlinked TEST produce two completeness errors,
one for each missing outbound relationship.

**Oracle:** The test passes only when its assertions confirm an
unlinked SPEC and unlinked TEST produce two completeness errors, one for each missing
outbound relationship.

## TEST-INTERNAL_GRAPH-014: Complete collection views

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-003`

**Purpose:**

Prove that the node-map and edge-slice accessors expose all inserted entries.

**Oracle:** The test passes only when its assertions confirm the
node-map and edge-slice accessors expose all inserted entries.

## TEST-INTERNAL_GRAPH-015: Empty graph construction

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_GRAPH-005`

**Purpose:**

Prove that construction returns an initialized graph with non-nil internal
collections and zero nodes and edges.

**Oracle:** The test passes only when its assertions confirm
construction returns an initialized graph with non-nil internal collections and zero
nodes and edges.

## Strategy

Tests use small in-memory graphs and exact counts, IDs, types, and flags as
oracles. Because tests are in the same package, constructor checks can inspect
private collection initialization while behavioral cases use the public API.

The suite does not cover missing endpoints at edge insertion, edges inserted
before nodes, duplicate relationships, conflicting node types, mutation through
borrowed collection views, node-order stability in snapshots, or concurrent
access. Those exclusions define where future invariant hardening needs new
tests.
