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

**Acceptance:**

`run` and `lint` resolve a target, load configuration, collect documentation
and source identifiers, merge them into one graph, execute the concrete engine
checks, and write a finding-centered report. Verbose diagnostics use stderr;
structured report output remains parseable on stdout.

A validation failure produces a report and a non-zero process status. A command
or I/O failure returns an error without manufacturing validation findings.

Documentation and source scope are intentionally asymmetric: the positional
path selects documentation, while code collection starts at the project
working directory. The authoritative gate therefore uses `.` from the project
root. Validation must not mutate documents, source, configuration, or embedded
Skill content.

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

1. the explicit `--config` path;
2. `./.idd.yaml`;
3. `./config/.idd.yaml`;
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
annotations, version, and similarity threshold, clamps supported threshold
bounds, and rejects unknown or missing annotation categories that would make
document and source collection disagree.

This SPEC does not select which project file wins; command-level precedence is
owned by `SPEC-CMD_IDD_CLI-003`. It also does not make `.idd.yaml` a semantic
catalog: components, contracts, requirements, tests, and coverage remain in
their Markdown owners.

## SPEC-CMD_IDD_CLI-007: Semantic comparison

- **Design:** `IDDCLIModule`
- **Contract:** `TFIDF`

**Requirement:** Compare meaningful documentation and code descriptions with
TF-IDF while ignoring declaration-only function locators.

**Acceptance:**

TF-IDF comparison runs only when both documentation and code provide meaningful
descriptions. A string containing only `[function: Name]` is a declaration
locator, not semantic prose, and is excluded to avoid false warnings when one
behavioral SPEC annotates multiple declarations.

Tokenization normalizes case and punctuation, removes stop words and
single-character noise, and compares vectors with cosine similarity. The
configured threshold produces warnings rather than rewriting either side.

Similarity is a review signal, not semantic verification. A high score cannot
prove behavioral equivalence, and a low score can reflect legitimate
vocabulary differences; every warning must retain enough source context for a
human decision.

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

`docs init <package>` creates or adopts four role-owned, self-describing documents
only after package, traversal, overwrite, central-catalog, and legacy-metadata
preflight checks. New fill locations carry stable scaffold markers, and the
command returns their deterministic work list instead of presenting generated
guidance as completed documentation.

`docs status <path>` reads one role file or a document tree and returns schema
`idd.document_status.v1`, complete/incomplete status, and exact
file/line/role/slot/reason entries. A slot remains incomplete when its marker is
present, its required role structure is absent, or its bounded content is empty
or a known placeholder. Existing records also remain incomplete while any
role-schema required field is absent or placeholder-filled.

`docs fix <path>` normalizes only minimal identity metadata. A file target
modifies only that document; a directory target may create missing skeletons.
It preserves scaffold markers. Neither mutation mode reformats prose or infers
requirements, titles, contracts, designs, purposes, kinds, or coverage.

Generated skeleton guidance must ask authors for purpose, responsibilities,
boundaries, rationale, failure cases, acceptance evidence, scenarios,
fixtures, and oracles. The command still leaves those sections unauthored:
helpful prompts are not semantic completion, and validation must continue to
report `idd-document-incomplete`, missing declarations, or empty required
design sections.

Initialization may preserve an existing plain Markdown body or merge generic
frontmatter, but it refuses existing IDD or legacy semantic metadata before any
write. Repair preserves body bytes and uses atomic replacement. A directory
repair may create missing structural files; a file repair never writes a
sibling.

After documents, tests, code, and annotations form a coherent checkpoint, the
agent first uses `docs status`, then `run . --format llm-markdown` as its
repair loop. Repository tests and `run . --format json` form the final project
gate. Package-targeted
documentation validation is not presented as a package-only code check.

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
