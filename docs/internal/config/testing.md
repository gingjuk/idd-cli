---
idd:
  version: "1.0"
  package: internal/config
---

# Testing: internal/config

## TEST-INTERNAL_CONFIG-001: Validation normalization and annotation roles

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_CONFIG-004`, `SPEC-INTERNAL_CONFIG-006`, `SPEC-INTERNAL_CONFIG-009`
- **Contracts:** `ConfigurationSchema`, `ConfigurationValidation`, `cmd/idd-cli#Config`

**Purpose:**

Prove preservation of the ignored deprecated consistency alias, detection of
its explicit YAML presence, and the complete annotation-role invariant.
Table-driven cases distinguish accepted maps from
missing, partial, unknown, empty, non-`@`-prefixed, whitespace-bearing, and
duplicate-value maps.

**Oracle:** Valid maps return nil, retain distinct configured tokens, and leave
legacy consistency values unchanged; loaded YAML with the alias returns one
actionable warning even for zero values, YAML without it returns none, and
every invalid key or value case returns a configuration error.

## TEST-INTERNAL_CONFIG-002: Built-in profile

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_CONFIG-001`, `SPEC-INTERNAL_CONFIG-002`, `SPEC-INTERNAL_CONFIG-003`, `SPEC-INTERNAL_CONFIG-005`, `SPEC-INTERNAL_CONFIG-006`, `SPEC-INTERNAL_CONFIG-008`
- **Contracts:** `ConfigurationDefaults`, `cmd/idd-cli#Config`

**Purpose:**

Prove the built-in profile is immediately usable and the maintained example is
an exact executable representation of it, including seven-language source
patterns, ignore paths, all deterministic validation switches, a zero-valued
deprecated consistency alias, annotations, and output defaults.

**Oracle:** The representative field assertions pass and loading
`examples/idd-config-example.yaml` produces a value deeply equal to a fresh
`Default()` result. Any omitted, reordered semantically different, or stale
default value fails the contract test.

## TEST-INTERNAL_CONFIG-003: YAML load success and failures

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_CONFIG-006`, `SPEC-INTERNAL_CONFIG-007`
- **Contracts:** `ConfigurationLoading`, `cmd/idd-cli#Config`

**Purpose:**

Prove successful decoding and validation from an isolated temporary YAML file,
including source-aware deprecation diagnostics, plus errors for a missing path
and malformed YAML. Preserved deprecated alias values and its warning are the
success oracle; non-nil errors are the failure oracle.

**Oracle:** Loading the valid fixture returns the authored deprecated values
unchanged and one source-aware removal warning. A fixture without the key has
no warning, while zero or null alias mappings still warn. Missing paths and
malformed YAML return non-nil stage-qualified errors and no configuration.

## Strategy

Tests invoke public behavior without mocks. Temporary directories isolate file
loading, and table-driven in-memory configurations isolate normalization.
Assertions target observable values and error presence rather than YAML library
internals.

Important exclusions are unknown YAML fields, omitted validation booleans in
partial files, regex/glob validity, output path writability, and mutation
aliasing. Those boundaries are documented explicitly and need dedicated cases
before the contract is strengthened.
