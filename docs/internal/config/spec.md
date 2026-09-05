---
idd:
  version: "1.1"
  package: internal/config
  namespace: INTERNAL_CONFIG
---

# Specifications: internal/config

## SPEC-INTERNAL_CONFIG-001: Root configuration ownership

- **Components:** `ConfigModule`
- **Contracts:** `ConfigurationSchema`

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

- **Components:** `ConfigModule`
- **Contracts:** `ConfigurationSchema`

**Requirement:**

Documentation configuration must carry file patterns, optional documentation
units with one or more project-relative source globs, and ignored paths without
treating any of those values as semantic document records.

**Acceptance:** The built-in-profile contract test observes the configured
documentation globs, unit mappings, and ignored paths through the
root configuration. Loading `examples/idd-config-example.yaml` must produce a
value deeply equal to `Default()`, so a missing or stale discovery value fails
the executable configuration example.

### Boundary

The package validates unit package paths and source globs, rejecting missing,
escaping, repeated, or malformed values. Empty document patterns are normalized
by `Validate`; units and ignore paths are not filled when omitted.

## SPEC-INTERNAL_CONFIG-003: Deprecated identifier pattern roles

- **Components:** `ConfigModule`
- **Contracts:** `ConfigurationSchema`
- **Status:** `deprecated`

**Requirement:**

Legacy `docs.identifier_patterns` input must remain recognizable only as
migration debt; it must not compete with the versioned IDD protocol as a second
authority for identifier syntax.

**Acceptance:** The public configuration model and built-in defaults expose no
identifier-pattern fields. Loading YAML that still contains the obsolete key
emits an actionable deprecation warning stating that the value is ignored and
should be removed. SPEC and TEST syntax remains consistent across document,
annotation, trace, and validation paths because `pkg/pattern` owns one grammar.

### Non-goals

This record preserves why the old configuration was removed. It does not make
legacy regex values effective or require current declarations to implement
the former API.

## SPEC-INTERNAL_CONFIG-004: Source discovery and annotation settings

- **Components:** `ConfigModule`
- **Contracts:** `ConfigurationSchema`

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

- **Components:** `ConfigModule`
- **Contracts:** `ConfigurationSchema`

**Requirement:**

Every independently configurable repository gate must have an explicit boolean
field, including graph consistency, coverage, document structure, doc/code
correspondence, annotation rules, and package document sets. Code-to-document
association is not a separate
file-header policy: it is the enabled correspondence rule joining matching
identifier evidence.

**Acceptance:** `TestDefault` observes an enabled representative gate, and
`TestDefaultMatchesExampleConfiguration` compares every validation switch in
the maintained YAML example with `Default()`. Adding, removing, or changing a
default gate without synchronizing the example makes that contract test fail.
Neither value contains the removed package-path-comment switch.

`require_pkg_doc_files` means every distinct directory containing a scanned,
non-ignored supported source file maps to the same relative path below
`docs/`, including nested sub-packages, and must contain the four canonical
role filenames. It is not a line-count or document-splitting policy.

### Loaded-file boundary

When YAML omits a boolean, loading leaves it false because files are decoded
into zero values rather than merged with `Default`. Callers wanting the full
built-in policy must choose `Default` or author the values explicitly.

## SPEC-INTERNAL_CONFIG-006: Deprecated consistency configuration alias

- **Components:** `ConfigModule`
- **Contracts:** `ConfigurationSchema`

**Requirement:**

Existing YAML containing `validation.consistency_check` must remain decodable
while the field is deprecated and behaviorally ignored, and its explicit
presence must produce an actionable deprecation warning.

**Acceptance:** Loading legacy `enabled` and `threshold` values succeeds, but
neither value changes validation behavior or receives default normalization.
Presence of the key produces exactly one diagnostic naming the source file and
`validation.consistency_check`, stating that the values are ignored, directing
the caller to remove the whole key, and naming
`docs review-context <SPEC-ID>...` as the semantic-review replacement. Presence
still warns when values are false, zero, or omitted.

### Edge cases

Unknown nested fields follow the package's existing YAML decoding behavior.
The alias must not be renamed into another active similarity rule.
Built-in defaults and loaded YAML without the alias produce no deprecation
diagnostic, and loading never rewrites the source file.

## SPEC-INTERNAL_CONFIG-007: Validated YAML loading

- **Components:** `ConfigModule`
- **Contracts:** `ConfigurationLoading`

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

Temporary-file tests cover successful deprecated-alias loading and invalid YAML;
separate cases cover missing paths and invalid annotation maps.

## SPEC-INTERNAL_CONFIG-008: Complete built-in profile

- **Components:** `ConfigModule`
- **Contracts:** `ConfigurationDefaults`

**Requirement:**

The default constructor must return a fresh non-nil profile that is immediately
usable by the CLI, with collection patterns, annotation roles, ignore paths,
deterministic validation gates, and no active lexical or semantic scoring
policy.

**Acceptance:** `TestDefault` observes version `1.0`, the complete supported
source-pattern set, an enabled validation gate, and a zero-valued deprecated
consistency alias. The example-configuration contract test then loads the
maintained YAML and requires deep equality with a fresh `Default()` value,
detecting stale collection, annotation, validation, or output defaults.

### Compatibility boundary

Default values influence every no-config invocation and are mirrored in the
repository example. Changes require an explicit compatibility decision and
example synchronization; partial YAML loading is not equivalent to calling
this constructor.

## SPEC-INTERNAL_CONFIG-009: In-place normalization and role-key validation

- **Components:** `ConfigModule`
- **Contracts:** `ConfigurationValidation`

**Requirement:**

Validation must mutate omitted core collection values to their documented
defaults, require exactly the three annotation-role keys, and return an error
for a missing or unknown role.

### Non-goals

Validation does not prove path existence, pattern syntax, output writability,
version compatibility, or semantic relationships between flags. A nil receiver
is not a supported input.

### Acceptance evidence

**Acceptance:**

Table-driven tests cover preservation of deprecated alias values and complete,
missing, partial, and extended annotation maps.

## SPEC-INTERNAL_CONFIG-010: Runtime workdir path resolution

- **Components:** `ConfigModule`
- **Contracts:** `RuntimeWorkdir`

**Requirement:**

One absolute, process-local workdir must resolve all validation filesystem
inputs without changing the process working directory.

**Acceptance:**

An absolute project root resolves relative document, source, configuration, and
engine paths into that tree. Paths reported from inside the root are converted
back to stable project-relative locations, while absolute paths outside the
root remain absolute. Configuration deprecation diagnostics loaded before the
workdir is installed are normalized by the same rule. Relative workdirs are
rejected and the runtime value is never decoded from or written to YAML.
