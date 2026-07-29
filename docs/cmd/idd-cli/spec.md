---
idd:
  version: "1.0"
  package: cmd/idd-cli
---

# Specifications: cmd/idd-cli

## SPEC-CMD_IDD_CLI-001: Validation execution

- **Design:** `IDDCLIModule`
- **Contract:** `CLI`

**Requirement:** Collect documentation and source identifiers, build and
validate the linkage graph, and report actionable findings through a
command-line workflow.

**Acceptance:**

`run` and `lint` resolve a target, load configuration, collect documentation
and source identifiers, merge them into one graph, execute the concrete engine
checks, and write a finding-centered report. Verbose diagnostics use stderr;
structured report output remains parseable on stdout.

A validation failure produces a report and a non-zero process status. A command
or I/O failure returns an error without manufacturing validation findings.

The positional path denotes one project root. Configuration discovery,
documentation collection, source collection, and engine filesystem checks must
all execute against that root, even when the command is launched from another
worktree. The command must resolve those paths through an invocation-local
workdir without changing process cwd. Validation must not mutate documents,
source, configuration, or embedded Skill content.

The command proves configured structural and traceability properties. It does
not infer missing requirements or treat a green result as proof that the
human-readable design is semantically complete.

## SPEC-CMD_IDD_CLI-002: Linkage graph

- **Design:** `IDDCLIModule`
- **Contract:** `LinkageGraph`

**Requirement:** Represent identifiers and directed traceability relationships
with deterministic lookup, verification, statistics, and snapshots.

**Acceptance:**

The graph is the common traceability model used by validation and reporting.
It retains every collected identifier while keeping edge direction explicit.

Document and code occurrences of one ID contribute distinct origin evidence to
one logical node. Node metadata must retain both descriptions, sources, and
TEST kinds needed by later checks. Adding a duplicate logical edge must not
erase source context or make lookup nondeterministic.

Typed inbound and outbound queries support validation without exposing callers
to index implementation. Verification marks whether a relationship has its
expected reverse evidence, while snapshots copy stable summaries for reports
without allowing mutation of the live graph.

## SPEC-CMD_IDD_CLI-003: Configuration precedence

- **Design:** `IDDCLIModule`
- **Contract:** `Config`

**Requirement:** Execute validation behavior from an explicit or discovered
IDD configuration with safe defaults.

**Acceptance:**

The command uses this search order:

1. the explicit `--config` path, resolved from the selected project root when
   relative;
2. `<project-root>/.idd.yaml`;
3. `<project-root>/config/.idd.yaml`;
4. in-memory defaults when no explicit path was required.

`--no-config` bypasses file loading. CLI output, format, and verbosity flags
override loaded values for the current invocation.

An explicitly requested configuration that cannot be read or validated is a
command error; the CLI must not silently fall back to defaults. When discovery
finds no configured file and no explicit path was required, in-memory defaults
provide a deterministic baseline.

This SPEC owns command-level selection and precedence. The parsing and
normalization behavior of an individual configuration file is owned by
`SPEC-CMD_IDD_CLI-006`.

## SPEC-CMD_IDD_CLI-004: Report boundaries

- **Design:** `IDDCLIModule`
- **Contract:** `Reporter`

**Requirement:** Convert validation results into stable JSON, Markdown, or
LLM-oriented findings without contaminating structured stdout.

**Acceptance:**

JSON output uses `idd.llm_report.v1` and groups repeated findings without
dropping their individual locations. Markdown remains available for people and
LLM-oriented Markdown for repair workflows.

Every document finding identifies its rule, severity, file and line when known,
affected identifier or field, and a repair hint. Graph statistics and optional
snapshots are derived from the same validation result.

Repeated findings may be summarized by rule and severity for navigation, but
the report must retain each concrete location. LLM Markdown emphasizes the
canonical owner and safe next action; it must not recommend `docs fix` for
semantic prose that only an author can repair.

For each `split-role` filename finding, JSON `suggested_fix` and LLM Markdown
`Fix` must name the exact split source and canonical target. The prompt must
require all unique still-valid semantic content to survive, prohibit
summary-only reduction and another role fragment, allow source removal only
after no-loss verification, and end with package `docs status` and project
`run` commands. A group that aggregates these findings must remain
path-neutral.

Report serialization errors or an unwritable destination are command failures.
Verbose diagnostics stay on stderr so stdout JSON remains a single parseable
document even when validation fails.

## SPEC-CMD_IDD_CLI-005: Identifier model

- **Design:** `IDDCLIModule`
- **Contract:** `Identifier`

**Requirement:** Preserve each identifier's type, origin, source, description,
links, and TEST kind through collection, graph construction, and reporting.

**Acceptance:**

Documentation and code collectors produce the same identifier model so later
stages do not need origin-specific traceability rules.

The model distinguishes a logical identifier from its occurrences. Origin,
source, line, description, links, and semantic TEST kind must survive set
merge and graph construction. Lookups that return one occurrence exist for
compatibility; correspondence validation must still be able to inspect every
origin.

The model does not own Markdown syntax, source-language parsing, validation
policy, or output formatting. Those stages exchange model values rather than
reaching into each other's internal representations.

## SPEC-CMD_IDD_CLI-006: Configuration loading

- **Design:** `IDDCLIModule`
- **Contract:** `Config`

**Requirement:** Load configuration by explicit flag and then project defaults,
while applying command-line output and verbosity overrides.

**Acceptance:**

Configuration discovery is deterministic, and invocation-only overrides do not
mutate the configuration file.

Loading decodes YAML into the concrete config value and validates annotation
key consistency. Validation supplies defaults for omitted patterns,
annotations, and version, preserves the deprecated consistency alias without
using it, emits an actionable warning when that YAML key is explicitly present,
and rejects unknown or missing annotation categories that would make document
and source collection disagree. Validation reports retain the warning as
structured data; evidence-only review-context output sends it to stderr so JSON
stdout remains parseable.

This SPEC does not select which project file wins; command-level precedence is
owned by `SPEC-CMD_IDD_CLI-003`. It also does not make `.idd.yaml` a semantic
catalog: components, contracts, requirements, tests, and coverage remain in
their Markdown owners.

## SPEC-CMD_IDD_CLI-007: Focused SPEC review contexts

- **Design:** `IDDCLIModule`
- **Contract:** `CLI`

**Requirement:** Assemble bounded, deterministic evidence bundles for one or
more SPECs so a human or LLM can judge document quality outside validation
with one shared repository scan.

**Acceptance:**

`docs review-context <SPEC-ID>...` accepts one to ten unique SPEC identifiers.
Repeated identifiers are deduplicated in first-request order. `--docs-path`
selects documentation input and defaults to `.`, while source remains rooted at
the current working tree. Documentation and source collection each run once.

Each identifier resolves exactly one canonical SPEC and returns its
Requirement, Acceptance, Design and Contract references, complete bounded
authored record Markdown including `Details`, the named Contract record, every
covering TEST record, and every attached
non-ignored source declaration annotated with that SPEC, plus declarations
annotated with each covering TEST ID. Declaration excerpts are bounded and
expose truncation. Output ordering is stable.

A one-SPEC JSON result retains schema `idd.spec_review_context.v1`. Multiple
SPECs use `idd.spec_review_context_batch.v1` and ordered independent contexts;
Markdown formats preserve the same separation and neutral review questions.
The command is read-only, does not invoke an LLM, does not emit a semantic score
or pass/fail verdict, and does not participate in `run` validity. An empty,
oversized, missing, malformed, or ambiguously owned request fails atomically as
an operational error.

## SPEC-CMD_IDD_CLI-008: Embedded workflow

- **Design:** `IDDCLIModule`
- **Contract:** `SkillsFS`

**Requirement:** Embed, list, and export the IDD authoring workflow so one
binary distributes instructions aligned with its validation behavior.

**Acceptance:**

The binary embeds the current IDD skill under `skills/*.md`.
`ListEmbeddedSkills` returns stable embedded paths and `ReadEmbeddedSkill`
returns one requested file. `skills` and `generate skill` use this fallback when
no external skill directory is available.

After an idd-cli upgrade, users regenerate the installed skill. The skill owns
semantic authoring and repair decisions; the binary owns structural document
operations and graph validation.

`skills` may list project-local Markdown Skills when present, falling back to
embedded files when no local Markdown is available. `generate skill` always
exports the binary-owned IDD workflow so its bytes are version-aligned with the
validator.

Export does not choose an agent installation directory, create missing parent
directories, or activate the Skill. Those actions remain with the user's agent
environment. An unknown generation target or embedded path returns an explicit
error.

## SPEC-CMD_IDD_CLI-009: Document commands

- **Design:** `IDDCLIModule`
- **Contract:** `CLI`

**Requirement:** Expose validation, skill generation, self-describing document
initialization, explicit completion status, and safe structural repair as one
phase-oriented Cobra workflow.

**Acceptance:**

`docs init <package>...` creates or adopts exactly `design.md`, `contract.md`,
`spec.md`, and `testing.md` for each unique package only after the whole batch
passes package, traversal, overwrite, central-catalog, and legacy-metadata
preflight checks. Each nested source sub-package maps to the matching nested
`docs/<package>/` path. Frontmatter contains only version/package identity; the
exact basename is the sole role authority. New fill locations carry stable
scaffold markers, and the command returns their deterministic work list instead
of presenting generated guidance as completed documentation.

`docs status <path>...` reads one or more role files or document trees and
returns schema `idd.document_status.v1`, normalized targets,
complete/incomplete status, and exact file/line/role/slot/reason entries.
Repeated or overlapping inputs do not duplicate work items. A slot remains
incomplete when its marker is present, its required role structure is absent,
or its bounded content is empty or a known placeholder. Existing records also
remain incomplete while any role-schema required field is absent or
placeholder-filled.

If the target file or tree contains a split role document, `docs status`
returns an operational error containing the same source/target,
content-preservation, deletion-order, and verification guidance. It does not
merge or delete the fragment.

`docs fix <path>...` normalizes only minimal identity metadata. A file target
modifies only that document; a directory target may create missing skeletons.
Repeated or overlapping targets produce one planned write per file. It
preserves scaffold markers. Neither mutation mode reformats prose or infers
requirements, titles, contracts, designs, purposes, kinds, or coverage.
Neither mode accepts or emits `idd.document`, creates split role files, or
imposes a document line-count limit.

Generated skeleton guidance must ask authors for purpose, responsibilities,
boundaries, rationale, failure cases, acceptance evidence, scenarios,
fixtures, and oracles. The command still leaves those sections unauthored:
helpful prompts are not semantic completion, and validation must continue to
report `idd-document-incomplete`, missing declarations, or empty required
design sections.

Initialization may preserve an existing plain Markdown body or merge generic
frontmatter, but it refuses existing IDD or legacy semantic metadata before any
write. Repair preserves body bytes and uses atomic replacement. Both mutators
preflight every batch member before applying writes. A directory repair may
create missing structural files; a file repair never writes a sibling.
Unexpected filesystem failure during application does not imply a
cross-filesystem transaction.

After documents, tests, code, and annotations form a coherent checkpoint, the
agent first uses `docs status`, then `run . --format llm-markdown` as its
repair loop. When semantic review is needed it requests one or more IDs through
`docs review-context <SPEC-ID>...` and reviews each context without treating
the batch as another check. Repository tests and `run . --format json` form the
final project gate. `run` and `lint` retain one complete project root rather
than an arbitrary list of package paths; package-focused document work uses
the dedicated `docs` commands.

## SPEC-CMD_IDD_CLI-010: Engine contract evidence

- **Design:** `IDDCLIModule`
- **Contract:** `Engine`

**Requirement:** Keep `Rule` as a test-only fixture that verifies concrete
engine and graph behavior without advertising a production rule extension
interface.

**Acceptance:**

Production validation remains implemented by concrete engine methods. The
fixture records expected behavior without creating a false public abstraction.

The test-only interface exists to express a narrow contract assertion in
`engine_contract_test.go`. It is not accepted by `Engine`, discovered at
runtime, configured by users, or promised as a plugin extension point.
Documentation and architecture diagrams must continue to name the concrete
engine ownership until a real extension lifecycle is designed and implemented.

## Concrete implementation notes

`SPEC-CMD_IDD_CLI-002`, `SPEC-CMD_IDD_CLI-005`, and
`SPEC-CMD_IDD_CLI-006` are realized by the concrete graph, model, and config
packages. `SPEC-CMD_IDD_CLI-010` describes the test-only `Rule` fixture; the
production engine exposes validation methods rather than a pluggable Rule
interface.
