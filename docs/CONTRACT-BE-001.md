# Contract: CONTRACT-BE-001

## Collector Interface Contracts

**Version:** 1.0.0
**Last Updated:** 2026-04-18

## Overview

Collectors gather IDD identifiers from documentation and source code. This contract defines the interfaces and behaviors for all collectors.

## Interfaces

### Collector Interface

```go
type Collector interface {
    Collect(ctx context.Context, cfg *config.Config) (*model.IdentifierSet, error)
}
```

### DocCollector Interface

```go
type DocCollector interface {
    CollectDocs(ctx context.Context, patterns []string) ([]*model.Identifier, error)
}
```

### CodeCollector Interface

```go
type CodeCollector Interface {
    CollectCode(ctx context.Context, patterns []string) ([]*model.Annotation, error)
}
```

## Behavior

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

## Related Documents

- **SPEC-BE-001** — IDD Link Validator Overview
- **TEST-BE-001** — Collector Implementation Tests
