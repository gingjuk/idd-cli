---
module: BE
type: SPEC
description: Specification index for IDD Link Validator backend module
markers:
  - SPEC-BE-001
  - SPEC-BE-002
---

# Specification Index (BE)

| ID | Title | Status | Tests |
|----|-------|--------|-------|
| [SPEC-BE-001](#spec-be-001) | IDD Link Validator Overview | Done | TEST-BE-001 |
| [SPEC-BE-002](#spec-be-002) | Graph Linkage Structure | Done | TEST-BE-001 |

---

## SPEC-BE-001

### IDD Link Validator Overview

**Status:** Done

**Requirement:**

IDD Link Validator is a CLI tool that validates bidirectional linkage consistency between IDD (Intent-Driven Development) identifiers across documentation and source code.

**Implementation:** `cmd/validator/main.go`, `internal/engine/engine.go`

**Key Modules:**
- `internal/collector/` — Document and code collection
- `internal/graph/` — Linkage graph structure
- `internal/validator/` — Validation rules
- `internal/reporter/` — Report generation

**Acceptance Criteria:**
- [x] CLI tool accepts `--config` flag for configuration
- [x] Scans documentation files matching configured patterns
- [x] Extracts IDD identifiers using configurable regex
- [x] Builds linkage graph from collected identifiers
- [x] Validates bidirectional links exist
- [x] Detects orphan identifiers (unless `allow_orphans: true`)
- [x] Generates JSON report with validation results
- [x] Exits with non-zero code when validation fails

**Tests:** TEST-BE-001

**Related:** [SPEC-BE-002](#spec-be-002), DESIGN-BE-001

---

## SPEC-BE-002

### Graph Linkage Structure

**Status:** Done

**Requirement:**

The LinkageGraph must efficiently represent bidirectional relationships between IDD identifiers, supporting fast lookup by ID, type, and link direction.

**Implementation:** `internal/graph/graph.go`

**Key Structures:**
- `Node` — Identifier with incoming/outgoing edges
- `Edge` — Directed relationship with verification status
- `Index` — Fast lookup indexes by ID, type, backlinks

**Acceptance Criteria:**
- [x] Nodes store identifier ID, type, and edge lists
- [x] Edges store direction, type, source location, and verification status
- [x] Fast O(1) lookup by node ID
- [x] Fast lookup of backlinks (nodes linking TO a node)
- [x] Bidirectional link verification marks edges as verified/unverified

**Tests:** TEST-BE-001

**Related:** [SPEC-BE-001](#spec-be-001), DESIGN-BE-001
