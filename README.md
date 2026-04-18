# IDD Link Validator

A CLI tool that validates bidirectional linkage consistency between IDD identifiers across documentation and source code.

## Overview

IDD Link Validator scans documentation and source code to build a linkage graph, then validates that all references are bidirectional (spec→test→code consistency).

## Quick Start

```bash
# Install
go build -o idd-validator ./cmd/validator

# Run
./idd-validator --config idd.yaml

# Or use Make
make build && make run
```

## Project Structure

```
idd-link-validator/
├── cmd/
│   └── validator/
│       └── main.go           # CLI entrypoint, argument parsing
├── internal/
│   ├── config/
│   │   └── config.go         # idd.yaml loading & validation
│   ├── collector/
│   │   ├── collector.go      # Interface:Collector, orchestrates collection
│   │   ├── doc_collector.go  # Scans docs/**/*.md for IDD identifiers
│   │   └── code_collector.go # Scans code for @spec/@contract/@test annotations
│   ├── model/
│   │   ├── identifier.go     # Identifier types (Spec, Contract, Test, Design)
│   │   ├── link.go            # Link relationship between identifiers
│   │   └── report.go          # Validation report model
│   ├── graph/
│   │   ├── graph.go           # LinkageGraph structure
│   │   └── builder.go         # Builds graph from collected identifiers
│   ├── linker/
│   │   └── linker.go          # Interface:Linker, resolves identifier references
│   ├── validator/
│   │   ├── validator.go       # Interface:Validator, main validation logic
│   │   ├── bidirectional.go    # Bidirectional linkage validation
│   │   ├── orphan.go          # Orphan/dangling reference detection
│   │   └── consistency.go      # Cross-reference consistency checks
│   ├── reporter/
│   │   └── reporter.go        # Interface:Reporter, JSON report generation
│   └── engine/
│       └── engine.go           # Main validation engine, wires everything together
├── pkg/
│   ├── pattern/
│   │   └── idd.go             # IDD identifier regex patterns
│   ├── walk/
│   │   └── files.go           # File traversal utilities
│   └── annotation/
│       └── parser.go          # Code annotation parsing (@spec, etc.)
├── config/
│   └── idd.yaml.example       # Example configuration
├── docs/
│   ├── SPEC-001-example.md     # Example spec document
│   └── TEST-001-example.md    # Example test document
├── scripts/
│   └── generate-report.go      # Standalone report generator utility
├── Makefile
├── go.mod
└── README.md
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLI (main.go)                           │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Validation Engine                             │
│  - Orchestrates collection → linking → validation → reporting  │
└─────────────────────────────────────────────────────────────────┘
         │                    │                    │              │
         ▼                    ▼                    ▼              │
┌────────────────┐   ┌────────────────┐   ┌────────────────┐     │
│   Collector    │   │    Linker      │   │   Validator    │     │
│                │   │                │   │                │     │
│ - DocCollector │──▶│ - Builds graph  │──▶│ - Bidirectional│     │
│ - CodeCollector│   │ - Resolves refs │   │ - Orphan detection│  │
│                │   │                │   │ - Consistency   │     │
└────────────────┘   └────────────────┘   └────────────────┘     │
                                                   │               │
                                                   ▼               │
                                          ┌────────────────┐      │
                                          │   Reporter     │      │
                                          │   - JSON out   │      │
                                          └────────────────┘      │
                                                                   │
                                                                   ▼
                                                          ┌────────────────┐
                                                          │  LinkageGraph  │
                                                          │  - Nodes (IDs) │
                                                          │  - Edges (Links)│
                                                          └────────────────┘
```

## Core Interfaces

### Collector

```go
// collector.go
package collector

// Collector aggregates identifiers from all sources.
type Collector interface {
    Collect(ctx context.Context, cfg *config.Config) (*model.IdentifierSet, error)
}

// DocCollector collects IDD identifiers from documentation files.
type DocCollector interface {
    CollectDocs(ctx context.Context, patterns []string) ([]*model.Identifier, error)
}

// CodeCollector collects annotations from source code.
type CodeCollector interface {
    CollectCode(ctx context.Context, patterns []string) ([]*model.Annotation, error)
}
```

### Linker

```go
// linker.go
package linker

// Linker builds the linkage graph by resolving identifier references.
type Linker interface {
    // BuildGraph constructs the full linkage graph from collected identifiers.
    BuildGraph(ctx context.Context, ids *model.IdentifierSet) (*graph.LinkageGraph, error)

    // ResolveLink resolves a reference string to a concrete identifier.
    ResolveLink(ref string, g *graph.LinkageGraph) (*model.Identifier, error)
}
```

### Validator

```go
// validator.go
package validator

// Validator validates the linkage graph for consistency.
type Validator interface {
    // Validate runs all validation rules and returns a validation report.
    Validate(ctx context.Context, g *graph.LinkageGraph) (*model.ValidationResult, error)
}

// ValidationRule is a single validation rule.
type ValidationRule interface {
    Name() string
    Validate(ctx context.Context, g *graph.LinkageGraph) []model.ValidationError
}
```

### Reporter

```go
// reporter.go
package reporter

// Reporter generates validation reports.
type Reporter interface {
    // Generate creates a report from the validation result.
    Generate(ctx context.Context, result *model.ValidationResult) (*model.Report, error)

    // Write writes the report to the configured output.
    Write(ctx context.Context, report *model.Report, output string) error
}
```

## Data Models

### Identifier

```go
// identifier.go
package model

// IdentifierType represents the type of IDD identifier.
type IdentifierType string

const (
    TypeSpec      IdentifierType = "SPEC"
    TypeContract  IdentifierType = "CONTRACT"
    TypeTest      IdentifierType = "TEST"
    TypeDesign    IdentifierType = "DESIGN"
)

// Identifier represents a single IDD identifier found in docs or code.
type Identifier struct {
    ID       string
    Type     IdentifierType
    Title    string
    Source   string    // File path or "inline"
    Line     int       // Line number
    RawRef   string    // Raw reference text (e.g., "SPEC-001")
    Links    []string  // Forward references from this identifier
}

// IdentifierSet is a collection of all collected identifiers.
type IdentifierSet struct {
    Specs     []*Identifier
    Contracts []*Identifier
    Tests     []*Identifier
    Designs   []*Identifier

    // Index for fast lookup
    byID map[string]*Identifier
}
```

### Link

```go
// link.go
package model

// LinkType represents the type of relationship between identifiers.
type LinkType string

const (
    LinkImplements   LinkType = "implements"
    LinkTests        LinkType = "tests"
    LinkReferences   LinkType = "references"
    LinkAnnotates    LinkType = "annotates"
)

// Link represents a directed relationship between two identifiers.
type Link struct {
    From        string    // Source identifier ID
    To          string    // Target identifier ID
    Type        LinkType
    Source      string    // File where the link was found
    Line        int
}
```

### Report

```go
// report.go
package model

// ValidationError represents a single validation error.
type ValidationError struct {
    Rule    string `json:"rule"`
    Message string `json:"message"`
    Source  string `json:"source,omitempty"`
    Link    string `json:"link,omitempty"`
    Code    string `json:"code,omitempty"`
}

// ValidationResult contains the result of validation.
type ValidationResult struct {
    Valid          bool              `json:"valid"`
    Errors         []ValidationError `json:"errors,omitempty"`
    Warnings       []ValidationError `json:"warnings,omitempty"`
    Stats          ValidationStats   `json:"stats"`
    Graph          *GraphSnapshot    `json:"graph,omitempty"`
}

// ValidationStats contains statistics about the validation run.
type ValidationStats struct {
    TotalIdentifiers int `json:"total_identifiers"`
    TotalLinks       int `json:"total_links"`
    SpecsAnalyzed    int `json:"specs_analyzed"`
    TestsAnalyzed    int `json:"tests_analyzed"`
    ContractsAnalyzed int `json:"contracts_analyzed"`
    DesignsAnalyzed  int `json:"designs_analyzed"`
}

// GraphSnapshot is a summary of the linkage graph for the report.
type GraphSnapshot struct {
    Nodes []NodeSummary `json:"nodes"`
    Edges []EdgeSummary `json:"edges"`
}

// Report is the final output report.
type Report struct {
    Tool      string            `json:"tool"`
    Version   string            `json:"version"`
    Timestamp string            `json:"timestamp"`
    Config    ConfigSummary     `json:"config"`
    Result    ValidationResult  `json:"result"`
}

// ConfigSummary summarizes the config used for the run.
type ConfigSummary struct {
    DocPatterns  []string `json:"doc_patterns"`
    CodePatterns []string `json:"code_patterns"`
    Annotations  []string `json:"annotations"`
}
```

## Linkage Graph

```go
// graph.go
package graph

// LinkageGraph represents the bidirectional linkage graph.
type LinkageGraph struct {
    nodes map[string]*Node
    edges []*Edge
    idx   *Index
}

// Node represents an identifier node in the graph.
type Node struct {
    ID       string
    Type     model.IdentifierType
    Metadata map[string]interface{}
    inEdges  []*Edge  // Incoming edges (backlinks)
    outEdges []*Edge  // Outgoing edges (forward links)
}

// Edge represents a directed edge in the graph.
type Edge struct {
    From      string
    To        string
    Type      model.LinkType
    Source    string
    Line      int
    Verified  bool  // Bidirectional link confirmed
}

// Index provides fast lookups.
type Index struct {
    byID       map[string]*Node
    byType     map[model.IdentifierType][]*Node
    backlinks  map[string][]string  // nodeID -> list of nodes linking TO it
}
```

## Bidirectional Linkage Validation

The core validation logic checks:

1. **Completeness**: Every SPEC has at least one TEST link (and vice versa)
2. **Consistency**: If SPEC→TEST exists, TEST→SPEC backlink must also exist
3. **Transitivity**: CODE annotations → SPEC → TEST chain consistency
4. **Orphan Detection**: No identifiers with zero connections

```go
// bidirectional.go
package validator

// BidirectionalRule validates bidirectional linkages.
type BidirectionalRule struct{}

func (r *BidirectionalRule) Name() string {
    return "bidirectional-linkage"
}

func (r *BidirectionalRule) Validate(ctx context.Context, g *graph.LinkageGraph) []model.ValidationError {
    var errors []model.ValidationError

    for _, node := range g.Nodes() {
        if node.Type != model.TypeSpec && node.Type != model.TypeTest {
            continue
        }

        // Check: specs must have test links
        if node.Type == model.TypeSpec {
            testLinks := g.GetOutboundByType(node.ID, model.LinkTests)
            if len(testLinks) == 0 {
                errors = append(errors, model.ValidationError{
                    Rule:    "bidirectional-linkage",
                    Message: fmt.Sprintf("SPEC %s has no test links", node.ID),
                    Link:    node.ID,
                })
            }
        }

        // Check: tests must have spec backlinks
        if node.Type == model.TypeTest {
            specLinks := g.GetInboundByType(node.ID, model.LinkTests)
            if len(specLinks) == 0 {
                errors = append(errors, model.ValidationError{
                    Rule:    "bidirectional-linkage",
                    Message: fmt.Sprintf("TEST %s is not linked from any SPEC", node.ID),
                    Link:    node.ID,
                })
            }
        }

        // Check: forward link must have corresponding backlink
        for _, edge := range node.OutEdges() {
            reverseType := reverseLinkType(edge.Type)
            reverseEdges := g.GetInboundByType(edge.To, reverseType)
            if len(reverseEdges) == 0 {
                errors = append(errors, model.ValidationError{
                    Rule:    "bidirectional-linkage",
                    Message: fmt.Sprintf("Link %s → %s has no backlink", edge.From, edge.To),
                    Link:    fmt.Sprintf("%s→%s", edge.From, edge.To),
                    Source:  fmt.Sprintf("%s:%d", edge.Source, edge.Line),
                })
            }
        }
    }

    return errors
}
```

## Configuration (idd.yaml)

```yaml
# idd.yaml - IDD Link Validator Configuration

version: "1.0"

# Documentation scan settings
docs:
  # Glob patterns for documentation files
  patterns:
    - "docs/**/*.md"
    - "SPEC-*.md"
    - "TEST-*.md"
  # Regex patterns to extract identifiers from docs
  identifier_patterns:
    spec: "SPEC-[0-9]+"
    contract: "CONTRACT-[0-9]+"
    test: "TEST-[0-9]+"
    design: "DESIGN-[0-9]+"

# Source code scan settings
code:
  # Glob patterns for source files
  patterns:
    - "**/*.go"
    - "**/*.ts"
    - "**/*.js"
  # Annotation markers to look for
  annotations:
    - "@spec"
    - "@contract"
    - "@test"
    - "@design"

# Validation rules
validation:
  # Require bidirectional links
  require_bidirectional: true
  # Allow orphan identifiers (not linked from anywhere)
  allow_orphans: false
  # Require all specs to have tests
  require_spec_test_coverage: true

# Output settings
output:
  # Output file path (default: stdout)
  file: "idd-report.json"
  # Include graph in report
  include_graph: false
  # Verbose output
  verbose: false
```

## JSON Output Format

```json
{
  "tool": "idd-link-validator",
  "version": "1.0.0",
  "timestamp": "2026-04-18T12:00:00Z",
  "config": {
    "doc_patterns": ["docs/**/*.md", "SPEC-*.md"],
    "code_patterns": ["**/*.go", "**/*.ts"],
    "annotations": ["@spec", "@contract", "@test"]
  },
  "result": {
    "valid": false,
    "errors": [
      {
        "rule": "bidirectional-linkage",
        "message": "SPEC-001 has no test links",
        "link": "SPEC-001",
        "code": "docs/SPEC-001.md:15"
      },
      {
        "rule": "orphan-detection",
        "message": "CONTRACT-042 is not referenced by any identifier",
        "link": "CONTRACT-042",
        "source": "docs/CONTRACT-042.md"
      }
    ],
    "warnings": [
      {
        "rule": "consistency",
        "message": "TEST-003 references SPEC-001 but test file name suggests SPEC-002",
        "link": "TEST-003→SPEC-001"
      }
    ],
    "stats": {
      "total_identifiers": 24,
      "total_links": 18,
      "specs_analyzed": 8,
      "tests_analyzed": 6,
      "contracts_analyzed": 5,
      "designs_analyzed": 5
    },
    "graph": {
      "nodes": [
        {
          "id": "SPEC-001",
          "type": "SPEC",
          "outbound": 2,
          "inbound": 1
        },
        {
          "id": "TEST-001",
          "type": "TEST",
          "outbound": 1,
          "inbound": 2
        }
      ],
      "edges": [
        {
          "from": "SPEC-001",
          "to": "TEST-001",
          "type": "tests",
          "verified": true
        }
      ]
    }
  }
}
```

## Key Design Decisions

### 1. Graph-First Architecture

The linkage graph is the central data structure. All collection feeds into it, and validation operates on it. This makes the system:
- Easy to extend with new validation rules
- Simple to add new collectors (just feed into the graph)
- Natural representation of bidirectional relationships

### 2. Extensible Collectors

New collectors can be added by implementing the `Collector` interface:
- Add a `JSONCollector` for JSON schema specs
- Add a `YAMLCollector` for config specs
- Add a `GitCollector` for git commit messages

### 3. Validation Rules as Plugins

Each `ValidationRule` is independent and testable in isolation. Rules are composed in the `Validator`:
```go
rules := []validator.ValidationRule{
    &validator.BidirectionalRule{},
    &validator.OrphanDetectionRule{},
    &validator.ConsistencyRule{},
}
```

### 4. Clear Separation

- **Collector**: Finds identifiers, doesn't know about links
- **Linker**: Builds the graph, resolves references
- **Validator**: Checks graph properties
- **Reporter**: Formats output

This separation allows testing each component independently and swapping implementations.