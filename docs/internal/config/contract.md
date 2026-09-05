---
idd:
  version: "1.1"
  package: internal/config
  namespace: INTERNAL_CONFIG
---

# Contracts: internal/config

## Contract: ConfigurationSchema

**Guarantees:**

`Config` is the complete in-memory input to collection, validation, and
reporting. Its nested values have separate ownership:

- `DocsConfig` selects documentation globs, optional documentation-unit source
  mappings, and ignored document paths;
- `CodeConfig` selects source globs, semantic annotation prefixes, and ignored
  source paths;
- `ValidationConfig` enables deterministic validation policies and retains the
  deprecated `consistency_check` decode shape;
- `ConsistencyCheck` is a compatibility-only value whose fields are ignored;
  and
- `OutputConfig` provides default output location, graph inclusion, and
  verbosity.

`Config.DeprecationWarnings()` returns a copy of load-time warnings derived
from explicitly present obsolete YAML paths. Warnings identify
`validation.consistency_check` or `docs.identifier_patterns`, their source,
ignored behavior, and removal action. Programmatic defaults have no warnings.

Slices and maps are caller-visible mutable values. The package does not clone
them after construction or loading. Consumers may read them concurrently only
if no caller mutates the configuration.

## Contract: ConfigurationLoading

**Guarantees:**

```go
func Load(path string) (*Config, error)
```

The function reads exactly the supplied path, decodes YAML into a zero-valued
`Config`, calls `Validate`, and returns the resulting pointer. It does not
search fallback paths or merge with `Default`.

Read failures are wrapped as `failed to read config file`; YAML failures as
`failed to parse config`; validation failures as `invalid config`. On any
failure the returned configuration is nil. Unknown YAML fields are currently
ignored by the decoder.

## Contract: ConfigurationDefaults

**Guarantees:**

```go
func Default() *Config
```

Each call returns a non-nil, independently allocated built-in profile. It
contains the version, documentation and source patterns, three annotation
roles, default ignore paths, enabled deterministic validation gates, a
zero-valued deprecated consistency alias, and non-verbose output without graph
inclusion.

The complete values are an operational compatibility surface and must stay
synchronized with `examples/idd-config-example.yaml` whenever the function
changes.

## Contract: ConfigurationValidation

**Guarantees:**

```go
func (c *Config) Validate() error
```

Validation mutates a non-nil receiver. It supplies version `1.0`, default
document patterns, the complete supported source-pattern set, and the complete
annotation map when the respective values are empty. Source defaults cover Go,
TypeScript/TSX, JavaScript/JSX, C++, Java, and Python extensions. The
annotation map must contain exactly `spec`, `test`, and `test_contract`; any
missing or additional key returns an error. Each value must be a trimmed,
non-empty, whitespace-free token beginning with `@`, and values must be unique
case-insensitively so one source comment cannot map to two semantic roles.

`validation.consistency_check` remains decodable so existing configuration
files do not fail during migration. Both `enabled` and `threshold` are ignored:
validation does not normalize them and the engine does not consume them.
Explicit YAML presence is retained as a deprecation warning even when the
decoded values are zero. Semantic review is requested explicitly through
`docs review-context`.

The method validates documentation-unit source globs but does not validate
general scan glob syntax, version support, output paths, ignore patterns, or
relationships between boolean flags. It does not fill omitted units, ignore paths, output values, or
validation booleans in a partially loaded YAML file.

## Contract: RuntimeWorkdir

**Guarantees:**

`SetWorkdir` accepts an absolute project root or an empty value that disables
rooted resolution. `ResolvePath` joins relative runtime paths to that root
without calling `os.Chdir`. `DisplayPath` converts paths inside the root back to
project-relative form for stable findings and leaves outside paths absolute.

The workdir is in-memory invocation state, not part of the YAML schema.
Collectors and the validation engine share the same configured value, so
configuration discovery, filesystem reads, and reported locations cannot mix
two project trees.
