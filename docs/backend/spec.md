---
markers:
  - id: SPEC-BE-001
    name: IDD CLI Overview
  - id: SPEC-BE-002
    name: Graph Linkage Structure
  - id: SPEC-BE-003
    name: Configuration Module
  - id: SPEC-BE-004
    name: Validation Engine
  - id: SPEC-BE-005
    name: Identifier Model
  - id: SPEC-BE-006
    name: Reporter Module
  - id: SPEC-BE-007
    name: Similarity Analysis
---

# Specification (backend)

## SPEC-BE-001: IDD CLI Overview

**Status:** Done

**Requirement:**

idd-cli is a CLI tool that validates bidirectional linkage consistency between IDD (Intent-Driven Development) identifiers across documentation and source code.

**Implementation:** `cmd/idd-cli/main.go`, `internal/engine/engine.go`

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

**Tests:** `TEST-BE-001`

**Related:** [`SPEC-BE-002`](#spec-be-002-graph-linkage-structure), `DESIGN-BE-001`

---

## SPEC-BE-002: Graph Linkage Structure

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

**Tests:** `TEST-BE-001`

**Related:** [`SPEC-BE-001`](#spec-be-001-idd-cli-overview), `DESIGN-BE-001`

---

## SPEC-BE-003: Configuration Module

**Status:** Done

**Requirement:**

The configuration module must load IDD settings from YAML configuration files, supporting CLI flag overrides, sensible defaults, and pattern-based file discovery.

**Implementation:** `internal/config/config.go`

**Key Functionality:**

- Load config from `.idd.yaml` file (or path specified via `--config` flag)
- Support `ignore_paths` with glob patterns (including `**` for recursive)
- Define identifier patterns for documentation markers and code annotations
- Configure reporter output format (JSON/Markdown)

**Acceptance Criteria:**

- [x] Config file is optional; defaults are applied if not found
- [x] CLI `--config` flag overrides default config paths
- [x] `ignore_paths` correctly excludes files/directories from validation
- [x] Identifier patterns are configurable via config file
- [x] Reporter format can be set to JSON or Markdown

**Tests:** `TEST-BE-001`

**Related:** [`SPEC-BE-001`](#spec-be-001-idd-cli-overview)

---

## SPEC-BE-004: Validation Engine

**Status:** Done

**Requirement:**

The validation engine orchestrates the collection of identifiers, building the linkage graph, and running validation rules to detect orphaned or improperly linked identifiers.

**Implementation:** `internal/engine/engine.go`

**Key Functionality:**

- Coordinate doc and code collectors
- Build linkage graph from collected identifiers
- Run validation rules (bidirectional links, orphan detection)
- Aggregate results and errors

**Acceptance Criteria:**

- [x] Engine collects identifiers from all configured sources
- [x] Graph is built with all nodes and edges
- [x] Bidirectional link validation detects unverified links
- [x] Orphan validation detects unreferenced identifiers
- [x] Validation results include errors and warnings

**Tests:** `TEST-BE-001`

**Related:** [`SPEC-BE-001`](#spec-be-001-idd-cli-overview), [`SPEC-BE-002`](#spec-be-002-graph-linkage-structure)

---

## SPEC-BE-005: Identifier Model

**Status:** Done

**Requirement:**

The identifier model defines data structures for representing IDD identifiers, annotations, and the identifier set collection with support for links and merging.

**Implementation:** `internal/model/identifier.go`

**Key Structures:**

- `Identifier` — IDD identifier with type, module, number, and links
- `Annotation` — Code annotation linking to spec/contract/test/design
- `IdentifierSet` — Collection of identifiers with add/get/has operations
- `IdentifierType` — Enum for SPEC, CONTRACT, TEST, DESIGN, BACKLINK

**Acceptance Criteria:**

- [x] Identifiers store type, module, number, and local ID
- [x] Forward links connect identifiers to their dependencies
- [x] Backlinks are computed from forward links
- [x] IdentifierSet supports add, get, has, count, all, merge operations
- [x] Annotations can be converted to identifiers

**Tests:** `TEST-BE-001`

**Related:** [`SPEC-BE-002`](#spec-be-002-graph-linkage-structure)

---

## SPEC-BE-006: Reporter Module

**Status:** Done

**Requirement:**

The reporter module generates validation reports in multiple formats (JSON, Markdown), presenting errors, warnings, and statistics clearly.

**Implementation:** `internal/reporter/reporter.go`

**Key Functionality:**

- Generate JSON report with validation summary
- Generate Markdown report with formatted output
- Include statistics: total identifiers, total links, error count
- List validation errors with file locations

**Acceptance Criteria:**

- [x] JSON output includes all validation results
- [x] Markdown output is human-readable
- [x] Errors include identifier ID and source location
- [x] Statistics are accurate

**Tests:** `TEST-BE-001`

**Related:** [`SPEC-BE-001`](#spec-be-001-idd-cli-overview)

---

## SPEC-BE-007: Similarity Analysis

**Status:** Done

**Requirement:**

The similarity module provides TF-IDF based document similarity analysis to help detect duplicate or very similar documentation files.

**Implementation:** `internal/similarity/tfidf.go`

**Key Functionality:**

- TF-IDF vectorization of document content
- Cosine similarity computation between documents
- Threshold-based duplicate detection

**Acceptance Criteria:**

- [x] Documents are vectorized using TF-IDF
- [x] Similarity scores range from 0.0 to 1.0
- [x] Configurable similarity threshold
- [x] Similar files are flagged for review

**Tests:** `TEST-BE-001`

**Related:** [`SPEC-BE-001`](#spec-be-001-idd-cli-overview)
