---
module: BE
type: CONTRACT
description: Contract interface definitions for IDD Link Validator
markers:
  - CONTRACT-BE-001
---

# Contract Index (BE)

| ID | Title | Status |
|----|-------|--------|
| [CONTRACT-BE-001](#contract-be-001) | Collector Interface Contracts | Done |

---

## CONTRACT-BE-001

### Collector Interface Contracts

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
|--------|---------|
| Documentation | `SPEC-[A-Z]+-[0-9]+`, `TEST-[A-Z]+-[0-9]+`, etc. |
| Code | `@spec SPEC-XXX`, `@test TEST-XXX`, etc. |

### Error Handling

- **File not found** — Return error, do not continue
- **Permission denied** — Return error, do not continue
- **Parse error** — Log warning, skip file, continue processing
- **Empty file** — No identifiers found, return empty result

**Related Specs:** SPEC-BE-001
