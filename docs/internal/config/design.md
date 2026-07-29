---
idd:
  version: "1.0"
  package: internal/config
---

# Design: internal/config

## Component: ConfigModule

**Purpose:**

`ConfigModule` is the typed boundary between authored YAML and the collectors,
engine, and reporter. It owns the configuration schema, the built-in default
profile, file loading, and a small normalization pass. Downstream packages
receive one `*Config` and should not reinterpret YAML independently.

The package is intentionally not a configuration-merging framework. Loading a
file starts from Go zero values, unmarshals the file, and then applies only the
defaults implemented by `Validate`; it does not overlay the document on the
complete value returned by `Default`.

### Responsibilities

- model document patterns, identifier patterns, source patterns, annotations,
  deterministic validation flags, deprecated compatibility settings, and
  output settings;
- return the full built-in profile used when the CLI chooses defaults;
- read and decode a specific YAML file;
- normalize selected omitted collection values;
- reject missing or unknown annotation-map keys and ambiguous or unusable
  prefix values; and
- wrap read, parse, and validation failures with stage context.

The component does not search for configuration files, expand environment
variables, merge multiple files, validate glob or regular-expression syntax,
or decide which CLI flag wins. Cobra command setup owns search order and
overrides; collectors and the engine consume the values.

### State, mutation, and failure boundaries

`Default` allocates fresh structs, slices, and maps. `Load` allocates and returns
a new configuration. `Validate` mutates its receiver in place: it fills empty
version and pattern lists, installs all seven-language source globs and the
complete annotation map when their values are empty. Deprecated consistency
values remain untouched and have no runtime effect.

Missing YAML fields that are not explicitly normalized stay at Go zero values.
In particular, `Load` does not automatically enable every validation boolean,
populate identifier regexes, or install default ignore paths. This distinction
between “use `Default`” and “load a partial file” is observable and must remain
clear to callers.

### Decisions and trade-offs

Plain structs with YAML tags keep the configuration easy to inspect and
serialize. In-place validation avoids a second copy, but callers must not
assume their input remains unchanged. Unknown top-level YAML fields are not
rejected because `yaml.Unmarshal` is used without strict known-field mode.

Annotation keys are validated as a closed set because collectors and
identifier patterns depend on the three semantic roles. Values remain
configurable, but validation requires distinct `@`-prefixed tokens without
whitespace. This prevents two roles from competing for the same syntax-tree
comment and lets findings quote the actual configured spelling.

## Architecture

```text
CLI config selection
       |
       +--> no file: Default()
       |
       +--> chosen path: ReadFile -> YAML Unmarshal -> Validate
                                                |
                                                v
                                   shared *config.Config
                           / collector / engine / reporter
```

Configuration is read before validation execution. The package has no global
mutable configuration, cache, watcher, reload lifecycle, or concurrency.

## Package Layout

`config.go` contains schema and lifecycle operations together. This keeps the
default profile adjacent to the fields it must populate and makes the
repository's config-example synchronization rule reviewable. CLI path discovery
does not belong here because it depends on command flags and working-directory
policy.

## Function Composition

`Load` performs read, decode, and `Validate` in that order, returning no partial
configuration on failure. `Validate` first supplies collection defaults, then
ensures annotation keys are exactly `spec`, `test`, and `test_contract`,
and validates their lexical values. It deliberately leaves the deprecated
consistency fields unchanged; no similarity threshold is normalized or used.
`Default` does not call `Validate`; it constructs the intended profile
directly.

## Dependencies

The package uses `os` for file reads and `gopkg.in/yaml.v3` for decoding.
Filesystem search, glob expansion, regex compilation, and report output remain
outside the package.

## Testability Hooks

`Default` and `Validate` are deterministic. Load tests use temporary YAML files
to exercise successful decoding and parse failures, while missing paths test
read errors without mocks. Table-driven annotation-key cases expose
normalization and rejection behavior, while deprecated-threshold cases prove
that compatibility values are preserved without affecting validation.

The current tests compare the complete example profile with `Default` but do
not cover unknown YAML fields, partial-file zero-value booleans, regex
validity, or aliasing of returned slices and maps.
