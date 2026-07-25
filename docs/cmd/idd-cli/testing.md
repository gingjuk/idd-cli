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

## TEST-CMD_IDD_CLI-002: Module integration behavior

- **Kind:** `test`
- **Covers:** `SPEC-CMD_IDD_CLI-001`, `SPEC-CMD_IDD_CLI-002`,
  `SPEC-CMD_IDD_CLI-003`, `SPEC-CMD_IDD_CLI-004`,
  `SPEC-CMD_IDD_CLI-005`, `SPEC-CMD_IDD_CLI-006`,
  `SPEC-CMD_IDD_CLI-007`, `SPEC-CMD_IDD_CLI-008`,
  `SPEC-CMD_IDD_CLI-009`, `SPEC-CMD_IDD_CLI-010`

**Purpose:** Cross-package tests verify configuration, engine, reporting,
similarity, embedding, and command collaborators.

## Strategy

The command package stays thin, so the two TEST records aggregate
evidence from the concrete internal packages instead of duplicating every Go
test name in this document.

## Core behavior

`TEST-CMD_IDD_CLI-001` covers identifier collection, graph construction,
bidirectional coverage, validation rules, structured findings, document
commands, and safe failure behavior.

Unit tests use temporary package and documentation trees. Document tests
exercise real YAML identity parsing, CommonMark records, per-file write scope,
body preservation, and atomic repair boundaries; engine tests construct focused
graphs and file fixtures.

## Integration behavior

`TEST-CMD_IDD_CLI-002` covers collaboration among configuration, collectors,
the engine, similarity scoring, reporters, and the embedded skill filesystem.

The repository-level acceptance check builds the actual `idd-cli` binary and
runs it against the repository. JSON output must parse successfully and contain
zero findings.

## Output isolation

Verbose-mode tests and smoke checks keep diagnostics on stderr so JSON stdout
can be piped directly to tools such as `jq`.
