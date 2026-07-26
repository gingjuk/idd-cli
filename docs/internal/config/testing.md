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

Prove threshold normalization and the complete annotation-role invariant.
Table-driven cases distinguish accepted maps from missing, partial, unknown,
empty, non-`@`-prefixed, whitespace-bearing, and duplicate-value maps, while
post-validation bounds are the numerical oracle.

**Oracle:** Valid maps return nil and retain distinct configured tokens; every
invalid key or value case returns a configuration error. Threshold cases finish
inside `[0, 1]`, using `0.3` for non-positive input and `1.0` as the upper cap.

## TEST-INTERNAL_CONFIG-002: Built-in profile

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_CONFIG-001`, `SPEC-INTERNAL_CONFIG-002`, `SPEC-INTERNAL_CONFIG-003`, `SPEC-INTERNAL_CONFIG-005`, `SPEC-INTERNAL_CONFIG-006`, `SPEC-INTERNAL_CONFIG-008`
- **Contracts:** `ConfigurationDefaults`, `cmd/idd-cli#Config`

**Purpose:**

Prove the built-in profile is immediately usable and the maintained example is
an exact executable representation of it, including seven-language source
patterns, ignore paths, all validation switches, opt-in consistency with the
`0.3` threshold, annotations, and output defaults.

**Oracle:** The representative field assertions pass and loading
`examples/idd-config-example.yaml` produces a value deeply equal to a fresh
`Default()` result. Any omitted, reordered semantically different, or stale
default value fails the contract test.

## TEST-INTERNAL_CONFIG-003: YAML load success and failures

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_CONFIG-007`
- **Contracts:** `ConfigurationLoading`, `cmd/idd-cli#Config`

**Purpose:**

Prove successful decoding and validation from an isolated temporary YAML file,
plus errors for a missing path and malformed YAML. The loaded threshold and
enabled flag are the success oracle; non-nil errors are the failure oracle.

**Oracle:** The test passes only when its assertions confirm successful
decoding and validation from an isolated temporary YAML file, plus errors for a missing
path and malformed YAML. The loaded threshold and enabled flag are the success oracle;
non-nil errors are the failure oracle.

## Strategy

Tests invoke public behavior without mocks. Temporary directories isolate file
loading, and table-driven in-memory configurations isolate normalization.
Assertions target observable values and error presence rather than YAML library
internals.

Important exclusions are unknown YAML fields, omitted validation booleans in
partial files, regex/glob validity, output path writability, and mutation
aliasing. Those boundaries are documented explicitly and need dedicated cases
before the contract is strengthened.
