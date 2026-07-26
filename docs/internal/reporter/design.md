---
idd:
  version: "1.0"
  package: internal/reporter
---

# Design: internal/reporter

## Component: ReportRenderer

**Purpose:**

`ReportRenderer` turns the engine's structured result into output for three
different readers: automation consumes finding-centered JSON, humans consume a
linkage summary in Markdown, and repair agents consume action-oriented LLM
Markdown. It also creates the complete internal report envelope that supplies
tool, version, timestamp, configuration summary, and validation result.

The component does not run validation or repair documents. It must preserve the
engine's severity and evidence while adding presentation metadata and practical
rule guidance.

### Responsibilities

- retain the chosen format and shared configuration reference;
- build a complete report envelope with execution metadata;
- select stdout or a caller-specified file destination;
- render human Markdown with findings, statistics, and optional graph details;
- project results into schema `idd.llm_report.v1`;
- enrich raw validation errors with titles, explanations, structured
  locations, expected/actual values, fix hints, and related identifiers;
- rank common rules and group repeated findings deterministically; and
- render the same finding model as agent-readable Markdown.

### Boundaries and side effects

`Generate` reads configuration and current time but performs no output.
`Write` owns output side effects. A non-empty path is created or truncated
before the format switch; an unsupported format can therefore leave an empty
file before returning its error. Empty output and `-` write directly to process
stdout.

The reporter does not create parent directories, write atomically, recover a
partially written file, close stdout, or read `Config.Output.File` inside
`Write`; the CLI passes the final destination explicitly.

### Finding enrichment boundary

Rule metadata is a curated presentation map, not validation policy. Known rules
receive specific explanations and repair hints; unknown rules receive
humanized titles and generic guidance. Primary identifiers are inferred from
structured link evidence first and then message, code, and source text.

Scaffold findings direct agents to `docs status` and never suggest that
`docs fix` can author content. Source-parse findings explain the deliberate
absence of regex fallback. Named Contract coverage and Component dependency
cycle findings point to `Contracts` and `Depends on`, preserving the canonical
relationship owner.

Split-role filename findings are the one path-sensitive repair specialization.
The structured raw finding supplies the split source, canonical basename, and
`split-role` code. `suggestedFix` joins the basename to the source directory and
builds an executable prompt: read both documents, structurally create a missing
target if necessary, merge all still-valid semantics without shortening them
to a summary, verify no loss, remove the split source, and rerun package and
project checks. The reporter does not inspect or mutate either file.

Related identifiers come from the raw finding text and optional graph edges.
This is context for repair, not proof that every related node caused the
finding. The reporter must not alter validity or convert warnings into errors.

### Decisions and trade-offs

JSON intentionally emits the finding-centered schema rather than the complete
internal `Report`. This makes each issue useful without global joins but omits
some run metadata from the public JSON surface. Human Markdown retains the
traditional stats and graph view.

Errors are emitted before warnings in finding order. Top rules are ranked by
frequency, then name, and limited to five. Groups are ranked by count, severity
string, rule, and first occurrence; files and identifiers within groups are
sorted. Group guidance comes from rule metadata rather than the first
path-specific finding, so an aggregate cannot imply that repairing one file
repairs every member.

## Architecture

```text
ValidationResult + Config
          |
          v
       Generate
          |
          v
        Report
          |
          +--> human Markdown
          |
          +--> buildLLMReport --> JSON
                             └--> LLM Markdown
```

Graph inclusion is decided upstream when the result is built. The reporter
renders a graph only when the snapshot is present and non-empty.

## Package Layout

`reporter.go` contains the renderer, rule presentation metadata, enrichment
helpers, grouping, and all three writers. These operations share one output
schema and stay together so JSON and LLM Markdown cannot silently diverge.
Serializable data types remain in `internal/model`.

## Function Composition

`New` chooses JSON only when format is empty; original casing is retained until
`Write` performs case-insensitive dispatch. `Generate` copies report metadata
and the validation result into a complete envelope.

The LLM projection builds errors then warnings, parses locations, derives any
split-document prompt, enriches each finding, derives top rules and groups, and
is used by both JSON and LLM Markdown. Human Markdown formats the complete
report directly.

## Dependencies

The package uses `internal/config` and `internal/model`, plus standard JSON,
I/O, filesystem, regex, sorting, string, numeric, and time support. It has no
template engine or third-party renderer.

## Testability Hooks

Formatting helpers accept `io.Writer` or build pure model values, so tests can
use strings and temporary files. A fixed sample report exercises location
parsing, identifiers, graph relationships, grouping, severity, and both
LLM-oriented formats.

Tests do not currently cover file-creation failure, short writes, atomicity,
unsupported-format truncation, nil config/report/result inputs, every fallback
rule, Windows-style locations, more than five top rules, or deterministic group
ties. Split-document prompt tests cover project-relative paths; they do not
exercise filenames containing Markdown delimiter characters.
