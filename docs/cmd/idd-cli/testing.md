---
idd:
  version: "1.0"
  package: cmd/idd-cli
  document: testing
---

# Testing: cmd/idd-cli

## TEST-CMD_IDD_CLI-001: Core validation behavior

- **Kind:** `test`
- **Covers:** `SPEC-CMD_IDD_CLI-001`, `SPEC-CMD_IDD_CLI-002`,
  `SPEC-CMD_IDD_CLI-003`, `SPEC-CMD_IDD_CLI-004`,
  `SPEC-CMD_IDD_CLI-005`, `SPEC-CMD_IDD_CLI-006`,
  `SPEC-CMD_IDD_CLI-007`, `SPEC-CMD_IDD_CLI-008`,
  `SPEC-CMD_IDD_CLI-009`, `SPEC-CMD_IDD_CLI-010`

**Purpose:** Engine, collector, model, and graph tests verify the core IDD
validation workflow.

### Evidence and scenarios

This record proves deterministic identifier merging, directed link
construction, derived coverage, structural-finding retention, enabled rule
execution, result sorting, and invalid-result behavior. It also covers
self-describing and legacy document inputs, code annotation kinds, missing and
malformed evidence, seven-language syntax-tree binding, scaffold completion,
named Contract coverage, Component dependencies, lifecycle links, and report
construction from the final graph.

Positive cases establish a valid graph. Negative cases isolate each rule and
assert its owner, severity, location, identifier, and field context rather than
checking only that some error occurred.

### Fixtures and oracle

**Oracle:** Pure graph and model tests construct values directly. Collector and
file-dependent engine tests use temporary package trees with focused configs.
The oracle is the resulting node/link metadata and structured validation
finding, not internal call order.

### Exclusions

Cobra argument parsing and process exit are kept thin; repository-level binary
smoke checks provide composition evidence instead of duplicating domain cases
through subprocess-heavy unit tests.

## TEST-CMD_IDD_CLI-002: Module integration behavior

- **Kind:** `test`
- **Covers:** `SPEC-CMD_IDD_CLI-001`, `SPEC-CMD_IDD_CLI-002`,
  `SPEC-CMD_IDD_CLI-003`, `SPEC-CMD_IDD_CLI-004`,
  `SPEC-CMD_IDD_CLI-005`, `SPEC-CMD_IDD_CLI-006`,
  `SPEC-CMD_IDD_CLI-007`, `SPEC-CMD_IDD_CLI-008`,
  `SPEC-CMD_IDD_CLI-009`, `SPEC-CMD_IDD_CLI-010`

**Purpose:** Cross-package tests verify configuration, engine, reporting,
similarity, embedding, and command collaborators.

### Evidence and scenarios

Integration evidence covers configuration discovery/defaults, embedded Skill
listing and byte-exact export, collector/engine handoff, all report formats,
stdout/stderr isolation, document initialization/status/repair safety, and the
project-root validity gate. Fresh scaffold fixtures prove that `docs init`
returns the initial work list, `docs status` becomes complete only after
authored content replaces every marker and every required record field is
concrete, and `docs fix` preserves that state.

Failure scenarios include invalid config, malformed documents, migration debt,
unwritable outputs, invalid graphs that still emit reports, and document
operations that must abort before partial writes.

### Fixtures and oracle

**Oracle:** Tests combine temporary files, concrete packages, and serialized report
inspection. The final oracle is the built CLI validating this repository with
zero errors and warnings while an independently exported Skill matches the
embedded file.

### Exclusions

Agent-specific Skill installation and the semantic quality of authored prose
cannot be automated by this module. Review verifies those concerns using the
exported workflow and human-readable documents.

## TEST-CMD_IDD_CLI-003: Test-only Rule fixture success contract

- **Kind:** `contract`
- **Covers:** `SPEC-CMD_IDD_CLI-010`
- **Contracts:** `Engine`

**Purpose:**

Verify that the historical test-only `Rule` fixture exposes the documented
`Name` and `Validate` shape and can return no findings for a successful
implementation. This is evidence about a test fixture, not a production engine
extension point.

**Oracle:** Pass when the
historical test-only `Rule` fixture exposes the documented `Name` and `Validate` shape
and can return no findings for a successful implementation. This is evidence about a
test fixture, not a production engine extension point.

The fixture uses an in-memory graph and direct method assertions. It performs
no runtime registration, discovery, configuration, or engine injection.

## TEST-CMD_IDD_CLI-004: Test-only Rule fixture failure contract

- **Kind:** `contract`
- **Covers:** `SPEC-CMD_IDD_CLI-010`
- **Contracts:** `Engine`

**Purpose:**

Verify that the same test-only fixture can return a structured validation error
with its own rule name. The oracle is the returned finding count and rule
value, not interaction with production validation dispatch.

**Oracle:** Pass when the
same test-only fixture can return a structured validation error with its own rule name.
The oracle is the returned finding count and rule value, not interaction with production
validation dispatch.

This evidence deliberately excludes plugin lifecycle, ordering, error
aggregation, and compatibility guarantees because the production engine does
not consume the fixture interface.

## TEST-CMD_IDD_CLI-005: Stable command surface contract

- **Kind:** `contract`
- **Covers:** `SPEC-CMD_IDD_CLI-001`, `SPEC-CMD_IDD_CLI-002`,
  `SPEC-CMD_IDD_CLI-003`, `SPEC-CMD_IDD_CLI-004`,
  `SPEC-CMD_IDD_CLI-005`, `SPEC-CMD_IDD_CLI-006`,
  `SPEC-CMD_IDD_CLI-007`, `SPEC-CMD_IDD_CLI-009`,
  `SPEC-CMD_IDD_CLI-010`
- **Contracts:** `CLI`

**Purpose:** Verify the executable exposes the documented validation, Skill,
and document-management command hierarchy plus its stable persistent options.

**Oracle:** The Cobra command tree contains `run`, `lint`, `skills`,
`generate`, and `docs`; `docs` contains `init`, `fix`, and `status`; and the
root exposes configuration, output, format, verbosity, and no-config flags.

## TEST-CMD_IDD_CLI-006: Embedded Skill boundary contract

- **Kind:** `contract`
- **Covers:** `SPEC-CMD_IDD_CLI-008`, `SPEC-CMD_IDD_CLI-009`
- **Contracts:** `SkillInfo`, `SkillsFS`

**Purpose:** Verify the binary embeds exactly the paired IDD Skill and exposes
its frontmatter through the command package's stable Skill information shape.

**Oracle:** The embedded path resolves to `skills/SKILL.md`, its bytes contain
the IDD workflow, parsed metadata retains name, audience, workflow, and
protection state, and a missing embedded path returns an error.

## Strategy

The command package stays thin, so the two TEST records aggregate
evidence from the concrete internal packages instead of duplicating every Go
test name in this document.

## Core behavior

`TEST-CMD_IDD_CLI-001` covers identifier collection, typed graph construction,
derived coverage, validation rules, structured findings, document commands,
and safe failure behavior.

Unit tests use temporary package and documentation trees. Document tests
exercise real YAML identity parsing, CommonMark records, per-file write scope,
scaffold work-list stability, body preservation, and atomic repair boundaries.
Tree-sitter fixtures cover Go, TypeScript, TSX, JavaScript/JSX, C++, Java, and
Python declaration binding and parse failure. Engine tests construct focused
graphs and file fixtures.

## Integration behavior

`TEST-CMD_IDD_CLI-002` covers collaboration among configuration, collectors,
the engine, similarity scoring, reporters, and the embedded skill filesystem.

The repository-level acceptance check builds the actual `idd-cli` binary and
runs it against the repository. JSON output must parse successfully and contain
zero findings.

The acceptance gate also checks that structural document repair is idempotent
and that generated skeleton guidance covers semantic depth without claiming
the resulting empty scaffold is complete.

## Output isolation

Verbose-mode tests and smoke checks keep diagnostics on stderr so JSON stdout
can be piped directly to tools such as `jq`.
