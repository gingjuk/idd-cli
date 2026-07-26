---
idd:
  version: "1.0"
  package: internal/config
---

# Specifications: internal/config

## SPEC-INTERNAL_CONFIG-001: Root configuration ownership

- **Design:** `ConfigModule`
- **Contract:** `ConfigurationSchema`

**Requirement:**

The configuration root must group document discovery, source discovery,
validation policy, and output policy into one value that can be passed through
the CLI without downstream YAML parsing. Output settings are part of this root
schema rather than a separate mutable global.

### Implementation boundary

The schema stores configuration; it does not execute patterns, validations, or
report writes. Nested slices and maps remain caller-owned mutable values and
are not deep-copied.

### Acceptance evidence

**Acceptance:**

Default and load tests observe representative nested values through the root
object, and repository execution demonstrates that collectors, engine, and
reporter consume the same structure.

## SPEC-INTERNAL_CONFIG-002: Documentation discovery settings

- **Design:** `ConfigModule`
- **Contract:** `ConfigurationSchema`

**Requirement:**

Documentation configuration must carry file patterns, role-specific identifier
pattern strings, and ignored path patterns without treating any of those values
as semantic document records.

**Acceptance:** Documentation
configuration must carry file patterns, role-specific identifier pattern strings, and
ignored path patterns without treating any of those values as semantic document records.

### Boundary

The package does not compile or execute the patterns. Empty document patterns
are normalized by `Validate`; identifier patterns and ignore paths are not
filled when omitted from a loaded file.

## SPEC-INTERNAL_CONFIG-003: Identifier pattern roles

- **Design:** `ConfigModule`
- **Contract:** `ConfigurationSchema`

**Requirement:**

Identifier pattern settings must expose distinct strings for SPEC, TEST, and
contract-TEST recognition so collectors can preserve annotation kind while
using a common TEST identifier namespace.

**Acceptance:** Identifier pattern
settings must expose distinct strings for SPEC, TEST, and contract-TEST recognition so
collectors can preserve annotation kind while using a common TEST identifier namespace.

### Non-goals

This specification does not require regex compilation or guarantee that an
authored pattern agrees with the built-in lexical grammar. Invalid regex text
is detected only by the consuming collector or engine path.

## SPEC-INTERNAL_CONFIG-004: Source discovery and annotation settings

- **Design:** `ConfigModule`
- **Contract:** `ConfigurationSchema`

**Requirement:**

Source configuration must carry code globs, ignored paths, and a semantic map
from `spec`, `test`, and `test_contract` roles to their source annotation
prefixes. When patterns are omitted, validation must install the complete
Tree-sitter-supported extension set for Go, TypeScript/TSX,
JavaScript/JSX, C++, Java, and Python.

**Acceptance:** Default and normalization tests observe every supported source
glob and verify that caller-authored non-empty patterns are preserved. Map
validation accepts exactly `spec`, `test`, and `test_contract` with distinct
valid tokens and rejects missing, unknown, empty, malformed, or duplicate
roles before collection.

### Invariants

The role-key set is closed and validated. Prefix values remain configurable,
but each must be a distinct, trimmed, whitespace-free `@` token. This package
does not parse source annotations itself; it prevents configurations whose
lexical roles would be empty or ambiguous.

## SPEC-INTERNAL_CONFIG-005: Validation policy switches

- **Design:** `ConfigModule`
- **Contract:** `ConfigurationSchema`

**Requirement:**

Every independently configurable repository gate must have an explicit boolean
field, including graph consistency, coverage, document structure, doc/code
correspondence, annotation rules, package document sets, and advisory
description consistency. Code-to-document association is not a separate
file-header policy: it is the enabled correspondence rule joining matching
identifier evidence.

**Acceptance:** Every independently
configurable repository gate must have an explicit boolean field, including graph
consistency, coverage, document structure, doc/code correspondence, annotation rules,
package document sets, and advisory description consistency. Defaults and the
example configuration contain no package-path-comment switch.

`require_pkg_doc_files` means every distinct directory containing a scanned,
non-ignored supported source file maps to the same relative path below
`docs/`, including nested sub-packages, and must contain the four canonical
role filenames. It is not a line-count or document-splitting policy.

### Loaded-file boundary

When YAML omits a boolean, loading leaves it false because files are decoded
into zero values rather than merged with `Default`. Callers wanting the full
built-in policy must choose `Default` or author the values explicitly.

## SPEC-INTERNAL_CONFIG-006: Consistency warning policy

- **Design:** `ConfigModule`
- **Contract:** `ConfigurationSchema`

**Requirement:**

Description-consistency configuration must separate whether scoring runs from
the numeric warning threshold, and validation must normalize that threshold to
the supported interval.

**Acceptance:** Description-consistency
configuration must separate whether scoring runs from the numeric warning threshold, and
validation must normalize that threshold to the supported interval.

### Edge cases

Thresholds at or below zero become `0.3`; thresholds above one become one.
Disabling the check requires `enabled: false`, not a zero threshold.

## SPEC-INTERNAL_CONFIG-007: Validated YAML loading

- **Design:** `ConfigModule`
- **Contract:** `ConfigurationLoading`

**Requirement:**

Loading must read the requested YAML file, decode it, apply configuration
validation, and return either a complete pointer or a stage-qualified error
with no partial value.

### Failure cases

Missing or unreadable files, malformed YAML, and invalid annotation-role keys
must remain distinguishable through wrapped error text. Unknown YAML fields are
currently accepted and ignored.

### Acceptance evidence

**Acceptance:**

Temporary-file tests cover successful threshold loading and invalid YAML;
separate cases cover missing paths and invalid annotation maps.

## SPEC-INTERNAL_CONFIG-008: Complete built-in profile

- **Design:** `ConfigModule`
- **Contract:** `ConfigurationDefaults`

**Requirement:**

The default constructor must return a fresh non-nil profile that is immediately
usable by the CLI, with collection patterns, annotation roles, ignore paths,
deterministic validation gates, and an explicitly disabled lexical consistency
hint whose threshold remains populated for opt-in use.

**Acceptance:** The default constructor
must return a fresh non-nil profile that is immediately usable by the CLI, with
collection patterns, annotation roles, ignore paths, deterministic validation gates,
and an explicitly disabled lexical consistency hint whose threshold remains populated
for opt-in use.

### Compatibility boundary

Default values influence every no-config invocation and are mirrored in the
repository example. Changes require an explicit compatibility decision and
example synchronization; partial YAML loading is not equivalent to calling
this constructor.

## SPEC-INTERNAL_CONFIG-009: In-place normalization and role-key validation

- **Design:** `ConfigModule`
- **Contract:** `ConfigurationValidation`

**Requirement:**

Validation must mutate omitted core collection values to their documented
defaults, require exactly the three annotation-role keys, normalize the
similarity threshold, and return an error for a missing or unknown role.

### Non-goals

Validation does not prove path existence, pattern syntax, output writability,
version compatibility, or semantic relationships between flags. A nil receiver
is not a supported input.

### Acceptance evidence

**Acceptance:**

Table-driven tests cover in-range, below-range, and above-range thresholds and
complete, missing, partial, and extended annotation maps.
