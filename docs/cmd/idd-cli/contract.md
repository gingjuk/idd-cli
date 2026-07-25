---
idd:
  version: "1.0"
  package: cmd/idd-cli
  document: contract
---

# Contracts: cmd/idd-cli

## Contract: CLI

`run [path]` and `lint [path]` use the same concrete workflow:

1. Load configuration from the explicit flag, `./.idd.yaml`, or
   `./config/.idd.yaml`.
2. Collect documentation identifiers from the target path.
3. Collect code annotations from the project working directory.
4. Merge identifiers and execute `Engine.Run`.
5. Write a finding-centered report.
6. Exit non-zero when the validation result is invalid.

Verbose diagnostics go to stderr so JSON stdout remains parseable.
The complete project gate is `run .`: a narrower documentation target does not
narrow source collection from the current working directory.

`llm-markdown` is the agent repair format, Markdown is the human readout, and
JSON is the automation-facing result. An invalid graph still writes its full
report before exiting non-zero.

## Document initialization

`docs init <package>`:

- requires an existing project-relative package directory;
- creates or adopts four self-describing Markdown files under
  `docs/<package>/`;
- never overwrites existing IDD or legacy metadata;
- refuses legacy marker or `related_files` metadata before writing anything;
- creates structural headings without fake requirements or generated IDs.

## Document repair

`docs fix <docs-package-or-document>`:

- repairs the minimal version, package, and document identity;
- creates missing skeletons for a directory target;
- writes only the selected file for a file target;
- preserves the Markdown body byte-for-byte;
- never reformats prose or invents semantic records;
- refuses malformed frontmatter rather than discarding unknown data.

## Concrete implementation boundary

Collectors are the concrete `DocCollector` and `CodeCollector` structs.
Link construction and validation rules are methods in `internal/engine`.
There is no production `Collector`, `Linker`, or `Rule` extension interface.

The command surface is the `CLI` boundary. Its concrete collaborators are
`LinkageGraph`, `Config`, `Engine`, `Identifier`, `Reporter`, and `TFIDF`.
Embedded skill access is provided by `SkillsFS`, `ListEmbeddedSkills`, and
`ReadEmbeddedSkill`; the `SkillInfo` type is the serialized skill metadata
boundary.

`generate skill` writes the exact workflow embedded in the binary. Users
regenerate their installed skill after upgrading idd-cli so semantic authoring
instructions stay aligned with validator behavior.

## Contract: Config

`Config` defines validation patterns, paths, output, and consistency thresholds.
The CLI loads it from the explicit path or documented project defaults before
applying invocation-only flag overrides.

`.idd.yaml` is configuration only. Package Components, Contracts, SPECs, TESTs,
and coverage remain in their owning Markdown records.

## Contract: Engine

`Engine` consumes collected identifiers and returns one validation result from
its concrete graph-building and validation methods.

## Contract: Identifier

`Identifier` preserves type, origin, source location, narrative description,
traceability links, and TEST kind across collection and reporting.

## Contract: LinkageGraph

`LinkageGraph` stores identifiers and directed traceability relationships with
deterministic lookup, verification, statistics, and snapshots.

## Contract: Reporter

`Reporter` renders the same validation result as JSON, Markdown, or
LLM-oriented Markdown without mixing diagnostics into structured stdout.
LLM-oriented Markdown supports the Skill repair loop; JSON supports CI and the
final project gate.

## Contract: SkillInfo

`SkillInfo` is the stable serialized description returned when embedded
workflow skills are listed.

## Contract: SkillsFS

`SkillsFS` exposes the embedded IDD workflow files used by skill listing,
reading, and generation commands. The embedded workflow is the canonical
companion to that binary.

## Contract: TFIDF

`TFIDF` compares meaningful descriptions while excluding declaration-only
function locators from semantic consistency scoring.

**Related Specs:** `SPEC-CMD_IDD_CLI-001`,
`SPEC-CMD_IDD_CLI-004`, `SPEC-CMD_IDD_CLI-009`
