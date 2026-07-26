---
idd:
  version: "1.0"
  package: internal/config
---

# Contracts: internal/config

## Contract: ConfigurationSchema

**Guarantees:**

`Config` is the complete in-memory input to collection, validation, and
reporting. Its nested values have separate ownership:

- `DocsConfig` selects documentation globs, identifier regex strings, and
  ignored document paths;
- `CodeConfig` selects source globs, semantic annotation prefixes, and ignored
  source paths;
- `ValidationConfig` enables individual validation policies and embeds the
  advisory consistency settings;
- `ConsistencyCheck` provides an enable flag and numeric threshold; and
- `OutputConfig` provides default output location, graph inclusion, and
  verbosity.

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
disabled-by-default lexical consistency hint with a `0.3` threshold, and
non-verbose output without graph inclusion.

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

A consistency threshold less than or equal to zero becomes `0.3`; a value
greater than one becomes `1.0`; values in `(0, 1]` are preserved. An explicit
zero therefore does not disable checking—the `Enabled` field does.

The method does not validate identifier regex or glob syntax, version support,
output paths, ignore patterns, or relationships between boolean flags. It does
not fill omitted identifier patterns, ignore paths, output values, or
validation booleans in a partially loaded YAML file.
