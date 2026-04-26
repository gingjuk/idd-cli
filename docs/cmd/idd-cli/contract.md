---
related_files:
  spec: docs/cmd/idd-cli/spec.md
  contract: docs/cmd/idd-cli/contract.md
  design: docs/cmd/idd-cli/design.md
  testing: docs/cmd/idd-cli/testing.md
---

# Contracts (backend)

## Collector Interface Contracts

**Status:** Done

**Overview:**

Collectors gather IDD identifiers from documentation and source code. This contract defines the interfaces and behaviors for all collectors.

### Interfaces

```go
// Collector is the main interface all collectors implement
type Collector interface {
    Collect(ctx context.Context, cfg *config.Config) (*model.IdentifierSet, error)
}

// DocCollector collects identifiers from documentation
type DocCollector interface {
    CollectDocs(ctx context.Context, patterns []string) ([]*model.Identifier, error)
}

// CodeCollector collects annotations from source code
type CodeCollector interface {
    CollectCode(ctx context.Context, patterns []string) ([]*model.Annotation, error)
}
```

### Collection Process

1. **Initialization** — Collector receives configuration on creation
2. **Target Resolution** — Accepts path (file or directory) as target
3. **File Discovery** — Recursively finds files matching configured patterns
4. **Identifier Extraction** — Parses files to extract IDD identifiers
5. **Link Detection** — Identifies references between identifiers within same file
6. **Result Assembly** — Returns IdentifierSet containing all found identifiers

### Identifier Extraction Rules

| Source | Pattern |
| ------ | ------- |
| Documentation | `SPEC-[A-Z]+-[0-9]+`, `TEST-[A-Z]+-[0-9]+`, etc. |
| Code | `@implement SPEC-XXX`, `@test TEST-XXX`, `@test-contract TEST-XXX`, etc. |

### Error Handling

- **File not found** — Return error, do not continue
- **Permission denied** — Return error, do not continue
- **Parse error** — Log warning, skip file, continue processing
- **Empty file** — No identifiers found, return empty result

**Related Specs:** `SPEC-CMD_IDD_CLI-001`

---

## Engine Validation Contracts

**Status:** Done

**Overview:**

The validation engine orchestrates the collection, linking, and validation process. This contract defines the engine's responsibilities and expected behaviors.

### Engine Interface

```go
// Engine is the main validation orchestrator
type Engine struct {
    cfg        *config.Config
    collector  *collector.Collector
    graph      *graph.LinkageGraph
    rules      []Rule
}

// Rule defines a validation rule
type Rule interface {
    Name() string
    Validate(g *graph.LinkageGraph) []model.ValidationError
}
```

### Validation Process

1. **Collect** — Gather identifiers from docs and code
2. **Link** — Build graph with edges from forward links
3. **Verify** — Check bidirectional links are reciprocated
4. **Report** — Aggregate errors and generate output

### Validation Rules

| Rule | Description |
| ---- | ----------- |
| BidirectionalLinkRule | Verifies links are bidirectional |
| OrphanRule | Detects unreferenced identifiers |

**Related Specs:** `SPEC-CMD_IDD_CLI-004`, `SPEC-CMD_IDD_CLI-002`
