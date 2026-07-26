---
idd:
  version: "1.0"
  package: internal/graph
  document: spec
---

# Specifications: internal/graph

## SPEC-INTERNAL_GRAPH-001: Identifier node state

- **Design:** `LinkageGraphStore`
- **Contract:** `GraphQuery`

**Requirement:**

Each graph identifier must have one node containing its ID, identifier type,
mutable metadata, and private incoming and outgoing edge collections.

### Ownership and boundary

Node identity is keyed only by the ID string. The first insertion determines
the retained type. Metadata is not interpreted by the graph, and adjacency
records only edges whose endpoints existed when those edges were added.

### Acceptance evidence

**Acceptance:**

Node insertion and adjacency tests verify initialized metadata, stable pointer
reuse, type indexing, and incoming/outgoing edge attachment.

## SPEC-INTERNAL_GRAPH-002: Directed relationship evidence

- **Design:** `LinkageGraphStore`
- **Contract:** `RelationshipVerification`

**Requirement:**

Every relationship edge must retain source ID, target ID, link type, source file
and line, plus mutable verification state so diagnostics can explain both the
relationship and its evidence location.

**Acceptance:** Every relationship edge
must retain source ID, target ID, link type, source file and line, plus mutable
verification state so diagnostics can explain both the relationship and its evidence
location.

### Edge cases

Edges are not unique and do not require existing endpoints. Verification is
derived later and can be reset on each verification pass.

## SPEC-INTERNAL_GRAPH-003: Mutable graph and typed traversal

- **Design:** `LinkageGraphStore`
- **Contract:** `GraphQuery`

**Requirement:**

The graph must support direct node lookup, complete node and edge views,
per-node adjacency, typed inbound and outbound filtering, counts, reciprocal
verification, completeness checks, statistics, and snapshot projection over one
coherent in-memory state.

### Implementation boundary

The graph does not infer relationships or choose validation rules. Returned
collection views are not immutable, and concurrent reads during mutation are
outside the contract.

### Acceptance evidence

**Acceptance:**

Tests cover lookup success and failure, typed filtering, reciprocal and
one-direction relationships, counts, statistics, snapshots, and completeness
outcomes.

## SPEC-INTERNAL_GRAPH-004: Reverse-reference index

- **Design:** `LinkageGraphStore`
- **Contract:** `GraphQuery`

**Requirement:**

Every added edge must append its source ID to an index keyed by target ID so
reverse references can be retrieved without scanning the global edge list.

**Acceptance:** Every added edge must
append its source ID to an index keyed by target ID so reverse references can be
retrieved without scanning the global edge list.

### Boundaries

The index preserves insertion and duplicates. It contains an entry even when
the target node does not exist. Successful lookup exposes the stored slice
without copying.

## SPEC-INTERNAL_GRAPH-005: Fully initialized empty graph

- **Design:** `LinkageGraphStore`
- **Contract:** `GraphMutation`

**Requirement:**

Graph construction must initialize all maps and slices so callers can add nodes
and edges immediately without nil checks.

### Acceptance evidence

**Acceptance:**

The constructor test asserts a non-nil graph, zero counts, and initialized
internal collections through same-package access.

## SPEC-INTERNAL_GRAPH-007: Idempotent node insertion

- **Design:** `LinkageGraphStore`
- **Contract:** `GraphMutation`

**Requirement:**

Adding a previously unseen ID must create and index one node; adding the same ID
again must return the existing node without changing its original type or
duplicating its type-index entry.

**Acceptance:** Adding a previously
unseen ID must create and index one node; adding the same ID again must return the
existing node without changing its original type or duplicating its type-index entry.

### Non-goals

Insertion does not merge metadata, diagnose a conflicting type, or repair edges
that were added before the node.

## SPEC-INTERNAL_GRAPH-008: Edge insertion across graph views

- **Design:** `LinkageGraphStore`
- **Contract:** `GraphMutation`

**Requirement:**

Adding an edge must append one global edge, update available endpoint adjacency,
and append the source to the target backlink index using the supplied type and
location evidence.

**Acceptance:** Adding an edge must
append one global edge, update available endpoint adjacency, and append the source to
the target backlink index using the supplied type and location evidence.

### Failure and compatibility boundary

The operation returns no error for missing endpoints and performs no duplicate
check. Engine construction order is responsible for normal endpoint integrity;
changing to rejection or deduplication would be observable behavior.
