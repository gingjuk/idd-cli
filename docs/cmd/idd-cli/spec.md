---
idd:
  version: "1.0"
  package: cmd/idd-cli
  document: spec
---

# Specifications: cmd/idd-cli

## SPEC-CMD_IDD_CLI-001: Validation execution

- **Design:** `IDDCLIModule`
- **Contract:** `CLI`

**Requirement:** Collect documentation and source identifiers, build and
validate the linkage graph, and report actionable findings through a
command-line workflow.

`run` and `lint` resolve a target, load configuration, collect documentation
and source identifiers, merge them into one graph, execute the concrete engine
checks, and write a finding-centered report. Verbose diagnostics use stderr;
structured report output remains parseable on stdout.

A validation failure produces a report and a non-zero process status. A command
or I/O failure returns an error without manufacturing validation findings.

## SPEC-CMD_IDD_CLI-002: Linkage graph

- **Design:** `IDDCLIModule`
- **Contract:** `LinkageGraph`

**Requirement:** Represent identifiers and directed traceability relationships
with deterministic lookup, verification, statistics, and snapshots.

The graph is the common traceability model used by validation and reporting.
It retains every collected identifier while keeping edge direction explicit.

## SPEC-CMD_IDD_CLI-003: Configuration precedence

- **Design:** `IDDCLIModule`
- **Contract:** `Config`

**Requirement:** Execute validation behavior from an explicit or discovered
IDD configuration with safe defaults.

The command uses this search order:

1. the explicit `--config` path;
2. `./.idd.yaml`;
3. `./config/.idd.yaml`;
4. in-memory defaults when no explicit path was required.

`--no-config` bypasses file loading. CLI output, format, and verbosity flags
override loaded values for the current invocation.

## SPEC-CMD_IDD_CLI-004: Report boundaries

- **Design:** `IDDCLIModule`
- **Contract:** `Reporter`

**Requirement:** Convert validation results into stable JSON, Markdown, or
LLM-oriented findings without contaminating structured stdout.

JSON output uses `idd.llm_report.v1` and groups repeated findings without
dropping their individual locations. Markdown remains available for people and
LLM-oriented Markdown for repair workflows.

Every document finding identifies its rule, severity, file and line when known,
affected identifier or field, and a repair hint. Graph statistics and optional
snapshots are derived from the same validation result.

## SPEC-CMD_IDD_CLI-005: Identifier model

- **Design:** `IDDCLIModule`
- **Contract:** `Identifier`

**Requirement:** Preserve each identifier's type, origin, source, description,
links, and TEST kind through collection, graph construction, and reporting.

Documentation and code collectors produce the same identifier model so later
stages do not need origin-specific traceability rules.

## SPEC-CMD_IDD_CLI-006: Configuration loading

- **Design:** `IDDCLIModule`
- **Contract:** `Config`

**Requirement:** Load configuration by explicit flag and then project defaults,
while applying command-line output and verbosity overrides.

Configuration discovery is deterministic, and invocation-only overrides do not
mutate the configuration file.

## SPEC-CMD_IDD_CLI-007: Semantic comparison

- **Design:** `IDDCLIModule`
- **Contract:** `TFIDF`

**Requirement:** Compare meaningful documentation and code descriptions with
TF-IDF while ignoring declaration-only function locators.

TF-IDF comparison runs only when both documentation and code provide meaningful
descriptions. A string containing only `[function: Name]` is a declaration
locator, not semantic prose, and is excluded to avoid false warnings when one
behavioral SPEC annotates multiple declarations.

## SPEC-CMD_IDD_CLI-008: Embedded workflow

- **Design:** `IDDCLIModule`
- **Contract:** `SkillsFS`

**Requirement:** Embed, list, and export the IDD authoring workflow so one
binary distributes instructions aligned with its validation behavior.

The binary embeds the current IDD skill under `skills/*.md`.
`ListEmbeddedSkills` returns stable embedded paths and `ReadEmbeddedSkill`
returns one requested file. `skills` and `generate skill` use this fallback when
no external skill directory is available.

After an idd-cli upgrade, users regenerate the installed skill. The skill owns
semantic authoring and repair decisions; the binary owns structural document
operations and graph validation.

## SPEC-CMD_IDD_CLI-009: Document commands

- **Design:** `IDDCLIModule`
- **Contract:** `SkillInfo`

**Requirement:** Expose validation, skill generation, and safe self-describing
document initialization and repair as one phase-oriented Cobra workflow.

`docs init <package>` creates or adopts four compact, self-describing documents
only after package, traversal, overwrite, central-catalog, and legacy-metadata
preflight checks.

`docs fix <path>` normalizes only minimal identity metadata. A file target
modifies only that document; a directory target may create missing skeletons.
Neither mode reformats prose or infers requirements, titles, contracts,
designs, purposes, kinds, or coverage.

After documents, tests, code, and annotations form a coherent checkpoint, the
agent uses `run . --format llm-markdown` as its repair loop. Repository tests
and `run . --format json` form the final project gate. Package-targeted
documentation validation is not presented as a package-only code check.

## SPEC-CMD_IDD_CLI-010: Engine contract evidence

- **Design:** `IDDCLIModule`
- **Contract:** `Engine`

**Requirement:** Keep `Rule` as a test-only fixture that verifies concrete
engine and graph behavior without advertising a production rule extension
interface.

Production validation remains implemented by concrete engine methods. The
fixture records expected behavior without creating a false public abstraction.

## Concrete implementation notes

`SPEC-CMD_IDD_CLI-002`, `SPEC-CMD_IDD_CLI-005`, and
`SPEC-CMD_IDD_CLI-006` are realized by the concrete graph, model, and config
packages. `SPEC-CMD_IDD_CLI-010` describes the test-only `Rule` fixture; the
production engine exposes validation methods rather than a pluggable Rule
interface.
