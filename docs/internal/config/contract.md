---
related_files:
  spec: docs/internal/config/spec.md
  contract: docs/internal/config/contract.md
  design: docs/internal/config/design.md
  testing: docs/internal/config/testing.md
---

# Contracts (config)

**Status:** Done

**Overview:**

Contracts for the configuration loading and validation module.

## Config Interface

```go
type Config struct {
    Version    string
    Docs       DocsConfig
    Code       CodeConfig
    Validation ValidationConfig
    Output     OutputConfig
}
```

**Invariants:**

- `Version` defaults to `"1.0"` when empty
- `Docs.Patterns` defaults to `["docs/**/*.md"]` when empty
- `Code.Patterns` defaults to `["**/*.go"]` when empty
- `ConsistencyCheck.Threshold` is clamped to `[0.0, 1.0]`

## Load Contract

`Load(path string) (*Config, error)` reads and parses a YAML config file, validates it, and returns the result. Returns an error if the file cannot be read, parsed, or fails validation.

## Default Contract

`Default() *Config` returns a fully populated Config with sensible defaults. Never returns nil.

## IdentifierPatterns

```go
type IdentifierPatterns struct {
    Spec         string
    Test         string
    TestContract string
}
```

Defines regex patterns for extracting SPEC, TEST, and TEST-CONTRACT identifiers from documentation files.

## Validate

`Validate() error` checks configuration values and applies defaults where fields are missing or out of range. Returns nil on success. Mutates the Config in-place (sets defaults).
