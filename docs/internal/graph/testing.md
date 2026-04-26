---
markers:
  - id: TEST-INTERNAL_GRAPH-001
    name: Graph Test 1
  - id: TEST-INTERNAL_GRAPH-002
    name: Graph Test 2
  - id: TEST-INTERNAL_GRAPH-003
    name: Graph Test 3
  - id: TEST-INTERNAL_GRAPH-004
    name: Graph Test 4
  - id: TEST-INTERNAL_GRAPH-005
    name: Graph Test 5
  - id: TEST-INTERNAL_GRAPH-006
    name: Graph Test 6
  - id: TEST-INTERNAL_GRAPH-007
    name: Graph Test 7
  - id: TEST-INTERNAL_GRAPH-008
    name: Graph Test 8
  - id: TEST-INTERNAL_GRAPH-009
    name: Graph Test 9
  - id: TEST-INTERNAL_GRAPH-010
    name: Graph Test 10
  - id: TEST-INTERNAL_GRAPH-011
    name: Graph Test 11
  - id: TEST-INTERNAL_GRAPH-012
    name: Graph Test 12
  - id: TEST-INTERNAL_GRAPH-013
    name: Graph Test 13
  - id: TEST-INTERNAL_GRAPH-014
    name: Graph Test 14
  - id: TEST-INTERNAL_GRAPH-015
    name: Graph Test 15

related_files:
  spec: docs/internal/graph/spec.md
  contract: docs/internal/graph/contract.md
  design: docs/internal/graph/design.md
  testing: docs/internal/graph/testing.md
---

# Test Cases (graph)

## TEST-INTERNAL_GRAPH-001: NewLinkageGraph

**Status:** Done

**Purpose:**

Test graph initialization.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-002: AddNode

**Status:** Done

**Purpose:**

Test node addition to graph.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-003: AddEdge

**Status:** Done

**Purpose:**

Test edge addition between nodes.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-004: GetNode

**Status:** Done

**Purpose:**

Test node lookup by ID.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-005: NodeInOutEdges

**Status:** Done

**Purpose:**

Test edge traversal from nodes.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-006: GetOutboundByType

**Status:** Done

**Purpose:**

Test filtering outbound edges by type.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-007: GetInboundByType

**Status:** Done

**Purpose:**

Test filtering inbound edges by type.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-008: GetBacklinks

**Status:** Done

**Purpose:**

Test reverse edge lookup using index.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-009: VerifyBidirectionalLinks

**Status:** Done

**Purpose:**

Test doc-link-consistency verification.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`, `SPEC-INTERNAL_GRAPH-009`

---

## TEST-INTERNAL_GRAPH-010: ToSnapshot

**Status:** Done

**Purpose:**

Test graph serialization to snapshot.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-011: Stats

**Status:** Done

**Purpose:**

Test statistics computation.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-012: ValidateCompleteness

**Status:** Done

**Purpose:**

Test completeness validation.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-013: Edge Verification

**Status:** Done

**Purpose:**

Test edge verification marking.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-014: Index Backlink Lookup

**Status:** Done

**Purpose:**

Test index-based backlink lookup performance.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`

---

## TEST-INTERNAL_GRAPH-015: Node Metadata

**Status:** Done

**Purpose:**

Test node metadata storage and retrieval.

**Spec Coverage:** `SPEC-INTERNAL_GRAPH-001`, `SPEC-INTERNAL_GRAPH-002`, `SPEC-INTERNAL_GRAPH-003`, `SPEC-INTERNAL_GRAPH-004`, `SPEC-INTERNAL_GRAPH-005`
