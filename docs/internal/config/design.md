---
related_files:
  spec: docs/internal/config/spec.md
  contract: docs/internal/config/contract.md
  design: docs/internal/config/design.md
  testing: docs/internal/config/testing.md
---

# Design (config)

**Status:** Done

## Architecture

The config module provides a single source of truth for all idd-cli settings. It is loaded once at startup and passed through the call stack via dependency injection.

```text
internal/config/
└── config.go   # Config struct, Load(), Default(), Validate()
```

## Package Layout

- **Config** — root struct aggregating all sub-configs
- **DocsConfig** — documentation scan settings (patterns, identifier_patterns, ignore_paths)
- **CodeConfig** — source code scan settings (patterns, annotations, ignore_paths)
- **ValidationConfig** — all validation rule toggles
- **ConsistencyCheck** — similarity threshold for doc/code description comparison
- **OutputConfig** — report output settings (file, format, verbosity)

## Function Composition

`Load(path)` → `os.ReadFile` → `yaml.Unmarshal` → `Validate()`

`Default()` returns a hard-coded Config with all fields populated. Used when no config file is found.

## Testability Hooks

- `Default()` is a pure function with no side effects; easy to call in tests
- `Load()` depends only on `os.ReadFile`; test with temp files

## Dependencies

- `gopkg.in/yaml.v3` — YAML parsing
