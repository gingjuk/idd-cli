---
markers:
  - id: SPEC-CMD_IDD_CLI-001
    name: IDD CLI Overview
  - id: SPEC-CMD_IDD_CLI-002
    name: Graph Linkage Structure
  - id: SPEC-CMD_IDD_CLI-003
    name: Engine Validation Rules
  - id: SPEC-CMD_IDD_CLI-004
    name: Reporter Output
  - id: SPEC-CMD_IDD_CLI-005
    name: Identifier Model
  - id: SPEC-CMD_IDD_CLI-006
    name: Config Loading
  - id: SPEC-CMD_IDD_CLI-007
    name: Similarity Analysis
  - id: SPEC-CMD_IDD_CLI-008
    name: Embed Files
  - id: SPEC-CMD_IDD_CLI-009
    name: CLI Main Entry
  - id: SPEC-CMD_IDD_CLI-010
    name: Engine Contract Tests

related_files:
  spec: docs/cmd/idd-cli/spec.md
  contract: docs/cmd/idd-cli/contract.md
  design: docs/cmd/idd-cli/design.md
  testing: docs/cmd/idd-cli/testing.md

---

# Specification (backend)

## SPEC-CMD_IDD_CLI-001: IDD CLI Overview

**Design:** `IDDCLIModule`

**Contract:** `CLI`

**Requirement:**

idd-cli is a CLI tool that validates bidirectional linkage consistency between IDD (Intent-Driven Development) identifiers across documentation and source code.

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

**Status:** Done

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

**Related:** [`SPEC-CMD_IDD_CLI-002`](#spec-cmd_idd_cli-002-graph-linkage-structure)

## SPEC-CMD_IDD_CLI-002: Graph Linkage Structure

**Design:** `IDDCLIModule`

**Contract:** `LinkageGraph`

**Requirement:**

The LinkageGraph must efficiently represent bidirectional relationships between IDD identifiers, supporting fast lookup by ID, type, and link direction.

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

**Status:** Done

**Implementation:** `internal/graph/graph.go`

**Key Structures:**

- `Node` — Identifier with incoming/outgoing edges
- `Edge` — Directed relationship with verification status
- `Index` — Fast lookup indexes by ID, type, backlinks

**Public Functions:**

### LinkageGraph.NewLinkageGraph

**Function Signature:**
`func NewLinkageGraph() *LinkageGraph`

**Purpose:** Creates a new empty linkage graph with initialized node map, edge slice, and index structures.

**Returns:** A new LinkageGraph pointer ready to accept nodes and edges

---

### LinkageGraph.AddNode

**Function Signature:**
`func (g *LinkageGraph) AddNode(id string, idType model.IdentifierType) *Node`

**Purpose:** Adds a node to the graph if it doesn't already exist. If the node already exists, returns the existing node. Creates the node with empty edge lists and initializes metadata map.

**Parameters:**

- `id`: The identifier ID for the node
- `idType`: The type of identifier (SPEC, CONTRACT, TEST, DESIGN)

**Returns:** The newly created or existing Node

---

### LinkageGraph.AddEdge

### LinkageGraph.AddEdge

**Function Signature:**
`func (g *LinkageGraph) AddEdge(from, to string, edgeType model.LinkType, source string, line int)`

**Purpose:** Adds a directed edge between two nodes. Creates the edge with source location information and appends it to the graph's edge list. Also updates the from node's outEdges and to node's inEdges.

**Parameters:**

- `from`: Source node ID
- `to`: Target node ID
- `edgeType`: Type of link (LinkTests, LinkImplements, LinkReferences)
- `source`: File path where the link was found
- `line`: Line number where the link was found

---

### LinkageGraph.GetNode

### LinkageGraph.GetNode

**Function Signature:**
`func (g *LinkageGraph) GetNode(id string) (*Node, bool)`

**Purpose:** Retrieves a node by its ID using O(1) map lookup.

**Parameters:**

- `id`: The node ID to look up

**Returns:** The Node and true if found, or nil and false if not found

---

### LinkageGraph.Nodes

**Function Signature:**
`func (g *LinkageGraph) Nodes() map[string]*Node`

**Purpose:** Returns the underlying node map for iteration. Use for traversing all nodes in the graph.

**Returns:** Map of node ID to Node pointer

---

### LinkageGraph.Edges

**Function Signature:**
`func (g *LinkageGraph) Edges() []*Edge`

**Purpose:** Returns all edges in the graph.

**Returns:** Slice of all Edge pointers

---

### LinkageGraph.GetOutboundByType

### LinkageGraph.GetOutboundByType

**Function Signature:**
`func (g *LinkageGraph) GetOutboundByType(nodeID string, linkType model.LinkType) []*Edge`

**Purpose:** Returns all outbound edges from a node that match the specified link type.

**Parameters:**

- `nodeID`: The source node ID
- `linkType`: The type of links to retrieve

**Returns:** Slice of matching edges (empty if node not found)

---

### LinkageGraph.GetInboundByType

### LinkageGraph.GetInboundByType

**Function Signature:**
`func (g *LinkageGraph) GetInboundByType(nodeID string, linkType model.LinkType) []*Edge`

**Purpose:** Returns all inbound edges to a node that match the specified link type.

**Parameters:**

- `nodeID`: The target node ID
- `linkType`: The type of links to retrieve

**Returns:** Slice of matching edges (empty if node not found)

---

### LinkageGraph.GetBacklinks

### LinkageGraph.GetBacklinks

**Function Signature:**
`func (g *LinkageGraph) GetBacklinks(nodeID string) []string`

**Purpose:** Returns all node IDs that have edges pointing TO the specified node. Uses the pre-built index for fast lookup.

**Parameters:**

- `nodeID`: The node ID to get backlinks for

**Returns:** Slice of source node IDs that link to this node

---

### LinkageGraph.NodeCount

### LinkageGraph.NodeCount

**Function Signature:**
`func (g *LinkageGraph) NodeCount() int`

**Purpose:** Returns the total number of nodes in the graph.

**Returns:** Node count

---

### LinkageGraph.EdgeCount

### LinkageGraph.EdgeCount

**Function Signature:**
`func (g *LinkageGraph) EdgeCount() int`

**Purpose:** Returns the total number of edges in the graph.

**Returns:** Edge count

---

### LinkageGraph.VerifyBidirectionalLinks

**Function Signature:**
`func (g *LinkageGraph) VerifyBidirectionalLinks()`

**Purpose:** Verifies that for every edge, there exists a corresponding reverse edge. Sets the Verified flag on each edge based on whether a matching reverse edge exists. For example, if A→B exists with LinkTests, then B→A should exist with reverse link type.

---

### LinkageGraph.ToSnapshot

**Function Signature:**
`func (g *LinkageGraph) ToSnapshot() *model.GraphSnapshot`

**Purpose:** Creates a serializable snapshot of the graph for inclusion in validation reports. Converts all nodes and edges to summary structures.

**Returns:** GraphSnapshot containing node and edge summaries

---

### LinkageGraph.Stats

**Function Signature:**
`func (g *LinkageGraph) Stats() model.ValidationStats`

**Purpose:** Computes validation statistics from the graph, including counts of total identifiers, links, and breakdowns by type (SPEC, TEST, CONTRACT, DESIGN).

**Returns:** ValidationStats with computed counts

---

### LinkageGraph.ValidateCompleteness

**Function Signature:**
`func (g *LinkageGraph) ValidateCompleteness() []model.ValidationError`

**Purpose:** Validates that SPEC nodes have at least one test link and TEST nodes are linked from at least one SPEC. Returns validation errors for any completeness violations.

**Returns:** Slice of ValidationError for completeness violations

---

### Node.InEdges

**Function Signature:**
`func (n *Node) InEdges() []*Edge`

**Purpose:** Returns all edges pointing TO this node (inbound edges).

**Returns:** Slice of inbound edges

---

### Node.OutEdges

**Function Signature:**
`func (n *Node) OutEdges() []*Edge`

**Purpose:** Returns all edges originating from this node (outbound edges).

**Returns:** Slice of outbound edges

---

**Acceptance Criteria:**

- [x] Nodes store identifier ID, type, and edge lists
- [x] Edges store direction, type, source location, and verification status
- [x] Fast O(1) lookup by node ID
- [x] Fast lookup of backlinks (nodes linking TO a node)
- [x] Bidirectional link verification marks edges as verified/unverified

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

**Related:** [`SPEC-CMD_IDD_CLI-001`](#spec-cmd_idd_cli-001-idd-cli-overview)

---

---

## SPEC-CMD_IDD_CLI-003: Configuration Module

**Design:** `IDDCLIModule`

**Contract:** `Config`

**Requirement:**

The configuration module must load IDD settings from YAML configuration files, supporting CLI flag overrides, sensible defaults, and pattern-based file discovery.

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

**Status:** Done

**Implementation:** `internal/config/config.go`

**Key Functionality:**

- Load config from `.idd.yaml` file (or path specified via `--config` flag)
- Support `ignore_paths` with glob patterns (including `**` for recursive)
- Define identifier patterns for documentation markers and code annotations
- Configure reporter output format (JSON/Markdown)

**Public Functions:**

### Config.Load

**Function Signature:**
`func Load(path string) (*Config, error)`

**Purpose:** Loads configuration from a YAML file at the specified path. Validates the configuration after loading and returns an error if validation fails.

**Parameters:**

- `path`: Path to the YAML configuration file

**Returns:** Parsed Config pointer or error if file cannot be read or validation fails

### Config.Default

**Function Signature:**
`func Default() *Config`

**Purpose:** Returns a Config with sensible default values. Used when no config file is provided or as a baseline to override.

**Returns:** Config with default values for all settings

---

### Config.Validate

**Function Signature:**
`func (c *Config) Validate() error`

**Purpose:** Validates and normalizes the configuration. Sets default values for empty fields and ensures consistency (e.g., threshold bounds).

**Returns:** Error if validation fails, nil otherwise

---

**Acceptance Criteria:**

- [x] Config file is optional; defaults are applied if not found
- [x] CLI `--config` flag overrides default config paths
- [x] `ignore_paths` correctly excludes files/directories from validation
- [x] Identifier patterns are configurable via config file
- [x] Reporter format can be set to JSON or Markdown

**Related:** [`SPEC-CMD_IDD_CLI-001`](#spec-cmd_idd_cli-001-idd-cli-overview)

---

---

## SPEC-CMD_IDD_CLI-004: Validation Engine

**Design:** `IDDCLIModule`

**Contract:** `Engine`

**Requirement:**

The validation engine orchestrates the collection of identifiers, building the linkage graph, and running validation rules to detect orphaned or improperly linked identifiers.

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

**Status:** Done

**Implementation:** `internal/engine/engine.go`

**Key Functionality:**

- Coordinate doc and code collectors
- Build linkage graph from collected identifiers
- Run validation rules (bidirectional links, orphan detection)
- Aggregate results and errors

**Public Functions:**

### Engine.New

**Function Signature:**
`func New(cfg *config.Config) *Engine`

**Purpose:** Creates a new validation engine with the given configuration. Initializes an empty linkage graph and validation result.

**Parameters:**

- `cfg`: Configuration pointer with validation rules and settings

**Returns:** A new Engine ready to run validation

### Engine.Run

**Function Signature:**
`func (e *Engine) Run(ctx context.Context, ids *model.IdentifierSet) (*model.ValidationResult, error)`

**Purpose:** Runs the complete validation pipeline: builds the graph from identifiers, runs all validation rules, and returns the result. If IncludeGraph is enabled in config, includes graph snapshot in result.

**Parameters:**

- `ctx`: Context for cancellation
- `ids`: Collected identifier set to validate

**Returns:** ValidationResult with errors, warnings, and stats, or error

---

### Engine.AddStructuralErrors

**Function Signature:**
`func (e *Engine) AddStructuralErrors(errors []*model.ValidationError)`

**Purpose:** Adds pre-collected structural errors (e.g., from collectors) to the engine's result. These are errors that were found during identifier collection phase.

**Parameters:**

- `errors`: Slice of validation errors to add

---

### Engine.BuildReport

**Function Signature:**
`func (e *Engine) BuildReport() *model.Report`

**Purpose:** Builds a complete report structure from the engine's current state, including tool info, config summary, and validation result.

**Returns:** Complete Report ready for output

---

**Acceptance Criteria:**

- [x] Engine collects identifiers from all configured sources
- [x] Graph is built with all nodes and edges
- [x] Bidirectional link validation detects unverified links
- [x] Orphan validation detects unreferenced identifiers
- [x] Validation results include errors and warnings

**Related:** [`SPEC-CMD_IDD_CLI-001`](#spec-cmd_idd_cli-001-idd-cli-overview), [`SPEC-CMD_IDD_CLI-002`](#spec-cmd_idd_cli-002-graph-linkage-structure)

---

---

## SPEC-CMD_IDD_CLI-005: Identifier Model

**Design:** `IDDCLIModule`

**Contract:** `Identifier`

**Requirement:**

The identifier model defines data structures for representing IDD identifiers, annotations, and the identifier set collection with support for links and merging.

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

**Status:** Done

**Implementation:** `internal/model/identifier.go`

**Key Structures:**

- `Identifier` — IDD identifier with type, module, number, and links
- `Annotation` — Code annotation linking to spec/contract/test/design
- `IdentifierSet` — Collection of identifiers with add/get/has operations
- `IdentifierType` — Enum for SPEC, CONTRACT, TEST, DESIGN, BACKLINK

**Public Functions:**

### Identifier.ParseIdentifierType

**Function Signature:**
`func ParseIdentifierType(s string) (IdentifierType, error)`

**Purpose:** Converts a string to IdentifierType. Used to parse identifier type from string representations (e.g., "SPEC", "TEST").

**Parameters:**

- `s`: The string to parse

**Returns:** The corresponding IdentifierType or an error if the string is not a valid type.

### Identifier.NewIdentifier

**Function Signature:**
`func NewIdentifier(id string, idType IdentifierType, title, source string, line int) *Identifier`

**Purpose:** Creates a new identifier with the given fields. Initializes an empty Links slice and sets Origin to OriginDoc.

**Parameters:**

- `id`: The identifier ID (e.g., `SPEC-CMD_IDD_CLI-001`)
- `idType`: The type of identifier (TypeSpec, TypeContract, etc.)
- `title`: The title/name of the identifier
- `source`: The file path where the identifier was found
- `line`: The line number where the identifier was found

**Returns:** A new Identifier pointer

---

### Identifier.NewIdentifierWithDescribe

**Function Signature:**
`func NewIdentifierWithDescribe(id string, idType IdentifierType, title, describe, source string, line int) *Identifier`

**Purpose:** Creates a new identifier with an additional describe field for extended description text.

**Parameters:**

- `id`: The identifier ID
- `idType`: The type of identifier
- `title`: The title/name of the identifier
- `describe`: Extended description text
- `source`: The file path where the identifier was found
- `line`: The line number

**Returns:** A new Identifier pointer

---

### IdentifierSet.NewIdentifierSet

**Function Signature:**
`func NewIdentifierSet() *IdentifierSet`

**Purpose:** Creates a new empty identifier set with initialized slices for each identifier type and an empty byID map.

**Returns:** A new IdentifierSet pointer ready to accept identifiers

---

### IdentifierSet.Add

**Function Signature:**
`func (s *IdentifierSet) Add(id *Identifier)`

**Purpose:** Adds an identifier to the set. Registers the identifier in the byID map for fast lookup and appends to the appropriate type slice (Specs, Contracts, Tests, or Designs).

**Parameters:**

- `id`: The identifier to add

---

### IdentifierSet.Get

**Function Signature:**
`func (s *IdentifierSet) Get(id string) (*Identifier, bool)`

**Purpose:** Returns the first identifier with the given ID for backward compatibility. Use GetAll to retrieve all identifiers with the same ID from different origins.

**Parameters:**

- `id`: The identifier ID to look up

**Returns:** The first matching identifier and true, or nil and false if not found

---

### IdentifierSet.GetAll

**Function Signature:**
`func (s *IdentifierSet) GetAll(id string) []*Identifier`

**Purpose:** Returns all identifiers with the given ID, including those from different origins (doc vs code).

**Parameters:**

- `id`: The identifier ID to look up

**Returns:** A slice of all identifiers with the ID (may be empty)

---

### IdentifierSet.Has

**Function Signature:**
`func (s *IdentifierSet) Has(id string) bool`

**Purpose:** Returns true if at least one identifier with the given ID exists in the set.

**Parameters:**

- `id`: The identifier ID to check

**Returns:** True if identifier exists, false otherwise

---

### IdentifierSet.All

**Function Signature:**
`func (s *IdentifierSet) All() []*Identifier`

**Purpose:** Returns all unique identifiers as a slice (one per ID). When multiple identifiers exist with the same ID (from different origins), only the first one is returned.

**Returns:** A slice of unique identifiers

---

### IdentifierSet.AllIdentifiers

**Function Signature:**
`func (s *IdentifierSet) AllIdentifiers() []*Identifier`

**Purpose:** Returns all identifiers including duplicates (multiple origins). Unlike All, this includes every identifier even if they share the same ID.

**Returns:** A slice of all identifiers

---

### IdentifierSet.ByOrigin

**Function Signature:**
`func (s *IdentifierSet) ByOrigin(origin Origin) []*Identifier`

**Purpose:** Returns all identifiers that have the specified origin (OriginDoc or OriginCode).

**Parameters:**

- `origin`: The origin to filter by

**Returns:** A slice of identifiers with the specified origin

---

### IdentifierSet.HasOrigin

**Function Signature:**
`func (s *IdentifierSet) HasOrigin(id string, origin Origin) bool`

**Purpose:** Returns true if at least one identifier with the given ID has the specified origin.

**Parameters:**

- `id`: The identifier ID to check
- `origin`: The origin to check for

**Returns:** True if at least one matching identifier has the origin

---

### IdentifierSet.Count

**Function Signature:**
`func (s *IdentifierSet) Count() int`

**Purpose:** Returns the total number of unique identifier IDs in the set.

**Returns:** The count of unique IDs

---

### IdentifierSet.Merge

**Function Signature:**
`func (s *IdentifierSet) Merge(other *IdentifierSet)`

**Purpose:** Combines another identifier set into this one by adding all identifiers from the other set.

**Parameters:**

- `other`: The identifier set to merge in

---

### Annotation.NewAnnotation

**Function Signature:**
`func NewAnnotation(typ IdentifierType, ref, source, raw, context string, line int) *Annotation`

**Purpose:** Creates a new annotation with the given fields. Used for code annotations extracted from source files.

**Parameters:**

- `typ`: The annotation type (SPEC, CONTRACT, TEST, DESIGN)
- `ref`: The identifier reference
- `source`: The file path
- `raw`: The raw annotation text
- `context`: Surrounding code context
- `line`: Line number

**Returns:** A new Annotation pointer

---

### Annotation.NewAnnotationWithComment

**Function Signature:**
`func NewAnnotationWithComment(typ IdentifierType, ref, source, raw, context, funcComment string, line int) *Annotation`

**Purpose:** Creates a new annotation that includes the associated function comment. Use when the annotation is attached to a function with documentation.

**Parameters:**

- `typ`: The annotation type
- `ref`: The identifier reference
- `source`: The file path
- `raw`: The raw annotation text
- `context`: Surrounding code context
- `funcComment`: The function's documentation comment
- `line`: Line number

**Returns:** A new Annotation pointer

---

### Annotation.ToIdentifier

**Function Signature:**
`func (a *Annotation) ToIdentifier() *Identifier`

**Purpose:** Converts an annotation to an identifier. The describe field is populated from FunctionComment if present.

**Returns:** A new Identifier representing the annotation

---

### ValidationResult.NewValidationResult

**Function Signature:**
`func NewValidationResult() *ValidationResult`

**Purpose:** Creates a new validation result with initialized empty slices for errors and warnings.

**Returns:** A new ValidationResult pointer

---

### ValidationResult.AddError

**Function Signature:**
`func (r *ValidationResult) AddError(rule, msg, source, link, code string)`

**Purpose:** Adds a validation error to the result and sets Valid to false.

**Parameters:**

- `rule`: The validation rule that failed
- `msg`: Human-readable error message
- `source`: File path where error was found
- `link`: The identifier/link involved
- `code`: Specific code or line involved

---

### ValidationResult.AddWarning

**Function Signature:**
`func (r *ValidationResult) AddWarning(rule, msg, source, link, code string)`

**Purpose:** Adds a validation warning to the result. Does not affect the Valid flag.

**Parameters:**

- `rule`: The validation rule that triggered the warning
- `msg`: Human-readable warning message
- `source`: File path where warning was found
- `link`: The identifier/link involved
- `code`: Specific code or line involved

---

### ValidationResult.Sort

**Function Signature:**
`func (r *ValidationResult) Sort()`

**Purpose:** Sorts errors and warnings by rule name, then by message. Ensures consistent output ordering.

---

**Acceptance Criteria:**

- [x] Identifiers store type, module, number, and local ID
- [x] Forward links connect identifiers to their dependencies
- [x] Backlinks are computed from forward links
- [x] IdentifierSet supports add, get, has, count, all, merge operations
- [x] Annotations can be converted to identifiers

**Related:** [`SPEC-CMD_IDD_CLI-002`](#spec-cmd_idd_cli-002-graph-linkage-structure)

---

---

## SPEC-CMD_IDD_CLI-006: Reporter Module

**Design:** `IDDCLIModule`

**Contract:** `Reporter`

**Requirement:**

The reporter module generates validation reports in multiple formats (JSON, Markdown), presenting errors, warnings, and statistics clearly.

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

### Reporter.Generate

**Status:** Done

**Implementation:** `internal/reporter/reporter.go`

**Key Functionality:**

- Generate JSON report with validation summary
- Generate Markdown report with formatted output
- Include statistics: total identifiers, total links, error count
- List validation errors with file locations

**Public Functions:**

### New

**Function Signature:**
`func New(cfg *config.Config, format string) *Reporter`

**Purpose:** Creates a new Reporter with the given configuration and output format. Defaults to JSON if format is empty.

**Parameters:**

- `cfg`: Configuration pointer
- `format`: Output format ("json" or "markdown")

**Returns:** A new Reporter instance

**Function Signature:**
`func (r *Reporter) Generate(result *model.ValidationResult) (*model.Report, error)`

**Purpose:** Generates a complete report from a validation result, including tool metadata, config summary, and the validation result.

**Parameters:**

- `result`: The validation result to include in the report

**Returns:** Complete Report structure or error

---

### Reporter.Write

**Function Signature:**
`func (r *Reporter) Write(report *model.Report, output string) error`

**Purpose:** Writes the report to the specified output destination. If output is empty or "-", writes to stdout. Otherwise creates a file at the given path.

**Parameters:**

- `report`: The report to write
- `output`: File path or "-" for stdout

**Returns:** Error if writing fails

---

**Acceptance Criteria:**

- [x] JSON output includes all validation results
- [x] Markdown output is human-readable
- [x] Errors include identifier ID and source location
- [x] Statistics are accurate

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

**Related:** [`SPEC-CMD_IDD_CLI-001`](#spec-cmd_idd_cli-001-idd-cli-overview)

---

---

## SPEC-CMD_IDD_CLI-007: Similarity Analysis

**Design:** `IDDCLIModule`

**Contract:** `TFIDF`

**Requirement:**

The similarity module provides TF-IDF based document similarity analysis to help detect duplicate or very similar documentation files.

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

### TFIDF.Tokenize

**Status:** Done

**Implementation:** `internal/similarity/tfidf.go`

**Key Functionality:**

- TF-IDF vectorization of document content
- Cosine similarity computation between documents
- Threshold-based duplicate detection

**Public Functions:**

### TFIDF.NewTFIDF

**Function Signature:**
`func NewTFIDF() *TFIDF`

**Purpose:** Creates a new TFIDF indexer with an empty IDF cache.

**Returns:** A new TFIDF pointer

**Function Signature:**
`func (t *TFIDF) Tokenize(text string) []string`

**Purpose:** Tokenizes text into lowercase alphanumeric tokens, filtering out stop words and single-character tokens.

**Parameters:**

- `text`: The text to tokenize

**Returns:** Slice of filtered tokens

---

### TFIDF.ComputeTF

**Function Signature:**
`func (t *TFIDF) ComputeTF(tokens []string) map[string]float64`

**Purpose:** Computes Term Frequency (TF) for each token in the document. TF = (count of token) / (total tokens).

**Parameters:**

- `tokens`: Slice of tokens from Tokenize

**Returns:** Map of token to TF value

---

### TFIDF.ComputeIDF

**Function Signature:**
`func (t *TFIDF) ComputeIDF(documents [][]string)`

**Purpose:** Computes Inverse Document Frequency (IDF) across a corpus of documents. IDF = log((N - df + 0.5) / (df + 0.5)) where N is total docs and df is document frequency.

**Parameters:**

- `documents`: Slice of token slices representing documents

---

### TFIDF.ComputeTFIDF

**Function Signature:**
`func (t *TFIDF) ComputeTFIDF(tf map[string]float64) map[string]float64`

**Purpose:** Computes TF-IDF vector by multiplying TF values with pre-computed IDF values.

**Parameters:**

- `tf`: Term frequency map from ComputeTF

**Returns:** TF-IDF vector map

---

### CosineSimilarity

**Function Signature:**
`func CosineSimilarity(vec1, vec2 map[string]float64) float64`

**Purpose:** Computes cosine similarity between two TF-IDF vectors. Returns value between 0.0 and 1.0.

**Parameters:**

- `vec1`: First TF-IDF vector
- `vec2`: Second TF-IDF vector

**Returns:** Cosine similarity score (0.0 to 1.0)

---

### TFIDF.Score

**Function Signature:**
`func (t *TFIDF) Score(docText, codeText string) float64`

**Purpose:** Computes similarity score between document text and code text using TF-IDF and cosine similarity.

**Parameters:**

- `docText`: Documentation text
- `codeText`: Code comment text

**Returns:** Similarity score (0.0 to 1.0)

---

### Score

**Function Signature:**
`func Score(docText, codeText string) float64`

**Purpose:** Convenience function that creates a temporary TFIDF instance and computes similarity in one call.

**Parameters:**

- `docText`: Documentation text
- `codeText`: Code comment text

**Returns:** Similarity score (0.0 to 1.0)

---

### NormalizeText

**Function Signature:**
`func NormalizeText(text string) string`

**Purpose:** Normalizes text by converting to lowercase, removing non-alphanumeric characters (except spaces), and collapsing whitespace.

**Parameters:**

- `text`: Text to normalize

**Returns:** Normalized text

---

**Acceptance Criteria:**

- [x] Documents are vectorized using TF-IDF
- [x] Similarity scores range from 0.0 to 1.0
- [x] Configurable similarity threshold
- [x] Similar files are flagged for review

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

**Related:** [`SPEC-CMD_IDD_CLI-001`](#spec-cmd_idd_cli-001-idd-cli-overview)

---

## Pattern Package

**Implementation:** `pkg/pattern/idd.go`

**Purpose:**

The pattern package provides IDD identifier pattern matching and extraction utilities. It defines regex patterns for IDD identifiers and code annotations.

**Public Functions:**

### ExtractIDDReferences

**Function Signature:**
`func ExtractIDDReferences(content string) []string`

**Purpose:** Extracts all IDD identifier references from content using configured patterns. Filters out identifiers that are quoted (wrapped in backticks or double quotes).

**Parameters:**

- `content`: Text content to search for IDD references

**Returns:** Slice of unique IDD identifier strings found

---

### ExtractAnnotations

**Function Signature:**
`func ExtractAnnotations(content string) []string`

**Purpose:** Extracts all IDD references from code annotations (e.g., `@implement`, `@test`, `@test-contract`) in source code.

**Parameters:**

- `content`: Source code content to search

**Returns:** Slice of identifier references found in annotations

---

### SplitAnnotationRefs

**Function Signature:**
`func SplitAnnotationRefs(s string) []string`

**Purpose:** Splits comma-separated IDD references and trims whitespace. Used to handle multiple references in a single annotation.

**Parameters:**

- `s`: Comma-separated string of references

**Returns:** Slice of individual trimmed references

---

### GetIdentifierType

**Function Signature:**
`func GetIdentifierType(ref string) string`

**Purpose:** Determines the type of an IDD identifier by matching against known patterns.

**Parameters:**
**Parameters:**

- `ref`: The identifier reference string

**Returns:** Type name ("SPEC", "CONTRACT", "TEST", "DESIGN") or empty string if no match

---

### GetAnnotationType

**Function Signature:**
`func GetAnnotationType(prefix string) string`

**Purpose:** Maps an annotation prefix to its corresponding identifier type.

**Parameters:**

- `prefix`: Annotation prefix (e.g., "@implement", "@test", "@test-contract")

**Returns:** Corresponding identifier type or empty string if unknown

---

### ValidateIDPattern

**Function Signature:**
`func ValidateIDPattern(id string) error`

**Purpose:** Validates that an identifier string matches one of the known IDD patterns.

**Parameters:**

- `id`: The identifier to validate

**Returns:** Error if identifier doesn't match any known pattern

---

## Walk Package

**Implementation:** `pkg/walk/files.go`

**Purpose:**

The walk package provides file traversal utilities with pattern matching support.

**Public Functions:**

### Walk

**Function Signature:**
`func Walk(patterns []string, visitor FileVisitor) error`

**Purpose:** Walks the filesystem matching files against the given glob patterns. Avoids duplicate visits using a visited map. Directories are walked recursively.

**Parameters:**

- `patterns`: Slice of glob patterns to match
- `visitor`: Callback function called for each matched file/directory

**Returns:** Error if visitation fails, nil otherwise

---

### MatchAnyExtensions

**Function Signature:**
`func MatchAnyExtensions(path string, extensions []string) bool`

**Purpose:** Checks if a file path has any of the specified extensions.

**Parameters:**

**Parameters:**

- `path`: File path to check
- `extensions`: Slice of extensions to match (e.g., ".go", ".ts")

**Returns:** True if path has one of the extensions

---

---

## SPEC-CMD_IDD_CLI-008: Embed Files

**Design:** `IDDCLIModule`

**Contract:** `Embed`

**Requirement:**

The CLI must be able to embed skill files for distribution as a single binary.

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

**Status:** Done

**Implementation:** `cmd/idd-cli/embed.go`

**Public Functions:**

### ListEmbeddedSkills

**Function Signature:**
`func ListEmbeddedSkills() []string`

**Purpose:** Lists all embedded skill files available in the binary.

**Returns:** Slice of skill file paths

---

### ReadEmbeddedSkill

**Function Signature:**
`func ReadEmbeddedSkill(path string) ([]byte, error)`

**Purpose:** Reads an embedded skill file by path.

**Parameters:**

- `path`: Path to the skill file within the embedded filesystem

**Returns:** File contents and error if not found

---

## SPEC-CMD_IDD_CLI-009: CLI Main Entry

**Design:** `IDDCLIModule`

**Contract:** `SkillInfo`

**Requirement:**

The CLI must provide a main entry point that parses flags and runs the validation engine.

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`

**Status:** Done

**Implementation:** `cmd/idd-cli/main.go`

**Public Types:**

### SkillInfo

**Type Definition:**
`type SkillInfo struct { ... }`

**Purpose:** Represents metadata about an available skill.

**Fields:**

- `Name` — Skill name
- `Description` — Skill description
- `License` — License identifier
- `Compatibility` — Compatibility tag
- `Audience` — Target audience
- `Workflow` — Workflow type
- `Protected` — Whether skill is protected
- `Module` — Module path
- `Path` — File path

---

## SPEC-CMD_IDD_CLI-010: Engine Contract Tests

**Design:** `IDDCLIModule`

**Contract:** `Rule`

**Requirement:**

The engine package must define a Rule interface that all validation rules implement.

**Tests:** `TEST-CMD_IDD_CLI-001`, `TEST-CMD_IDD_CLI-002`
**Status:** Done

**Implementation:** `internal/engine/engine_contract_test.go`

**Public Interfaces:**

### Rule

**Interface Definition:**

```go
type Rule interface {
    Name() string
    Validate(g *graph.LinkageGraph) []model.ValidationError
}
```

**Purpose:** Interface for all validation rules.

**Methods:**

- `Name()` — Returns the rule name
- `Validate(g *LinkageGraph)` — Validates the graph and returns errors
