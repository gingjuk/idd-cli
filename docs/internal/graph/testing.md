---
markers:
  - id: TEST-INT_GRPH-001
    name: Graph Test 1
  - id: TEST-INT_GRPH-002
    name: Graph Test 2
  - id: TEST-INT_GRPH-003
    name: Graph Test 3
  - id: TEST-INT_GRPH-004
    name: Graph Test 4
  - id: TEST-INT_GRPH-005
    name: Graph Test 5
  - id: TEST-INT_GRPH-006
    name: Graph Test 6
  - id: TEST-INT_GRPH-007
    name: Graph Test 7
  - id: TEST-INT_GRPH-008
    name: Graph Test 8
  - id: TEST-INT_GRPH-009
    name: Graph Test 9
  - id: TEST-INT_GRPH-010
    name: Graph Test 10
  - id: TEST-INT_GRPH-011
    name: Graph Test 11
  - id: TEST-INT_GRPH-012
    name: Graph Test 12
  - id: TEST-INT_GRPH-013
    name: Graph Test 13
  - id: TEST-INT_GRPH-014
    name: Graph Test 14
  - id: TEST-INT_GRPH-015
    name: Graph Test 15

related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Test Cases (graph)

## TEST-INT_GRPH-001: NewLinkageGraph

**Status:** Done

**Purpose:**

Test graph initialization.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-002: AddNode

**Status:** Done

**Purpose:**

Test node addition to graph.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-003: AddEdge

**Status:** Done

**Purpose:**

Test edge addition between nodes.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-004: GetNode

**Status:** Done

**Purpose:**

Test node lookup by ID.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-005: NodeInOutEdges

**Status:** Done

**Purpose:**

Test edge traversal from nodes.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-006: GetOutboundByType

**Status:** Done

**Purpose:**

Test filtering outbound edges by type.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-007: GetInboundByType

**Status:** Done

**Purpose:**

Test filtering inbound edges by type.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-008: GetBacklinks

**Status:** Done

**Purpose:**

Test reverse edge lookup using index.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-009: VerifyBidirectionalLinks

**Status:** Done

**Purpose:**

Test doc-link-consistency verification.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`, `SPEC-INT_GRPH-009`

---

## TEST-INT_GRPH-010: ToSnapshot

**Status:** Done

**Purpose:**

Test graph serialization to snapshot.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-011: Stats

**Status:** Done

**Purpose:**

Test statistics computation.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-012: ValidateCompleteness

**Status:** Done

**Purpose:**

Test completeness validation.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-013: Edge Verification

**Status:** Done

**Purpose:**

Test edge verification marking.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-014: Index Backlink Lookup

**Status:** Done

**Purpose:**

Test index-based backlink lookup performance.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`

---

## TEST-INT_GRPH-015: Node Metadata

**Status:** Done

**Purpose:**

Test node metadata storage and retrieval.

**Spec Coverage:** `SPEC-INT_GRPH-001`, `SPEC-INT_GRPH-002`, `SPEC-INT_GRPH-003`, `SPEC-INT_GRPH-004`, `SPEC-INT_GRPH-005`
