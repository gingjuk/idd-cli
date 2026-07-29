---
idd:
  version: "1.0"
  package: cmd/idd-cli
---

# Contracts: cmd/idd-cli

## Contract: CLI

**Guarantees:**

The CLI contract is the observable process boundary: accepted commands and
flags, filesystem scope, report destinations, stderr/stdout separation, and
exit status. Internal package types may change without affecting users as long
as these behaviors remain stable.

`run [project-root]` and `lint [project-root]` use the same concrete workflow:

1. Resolve the optional path as one absolute project workdir.
2. Load configuration from the explicit flag, `<root>/.idd.yaml`, or
   `<root>/config/.idd.yaml`.
3. Resolve documentation, source, and engine filesystem paths against that
   workdir without changing process cwd.
4. Merge identifiers and execute `Engine.Run`.
5. Write a finding-centered report.
6. Exit non-zero when the validation result is invalid.

Verbose diagnostics go to stderr so JSON stdout remains parseable.
The normal in-project gate is `run .`. Passing an absolute or relative path to
another project validates that tree without mixing it with source or
configuration from the caller's working directory.

`llm-markdown` is the agent repair format, Markdown is the human readout, and
JSON is the automation-facing result. JSON and LLM Markdown expose
path-specific split-document repair prompts in the detailed finding, while a
rule group remains path-neutral when it spans several files. An invalid graph
still writes its full report before exiting non-zero.

### Inputs and scope

The optional path selects the complete project root. Relative explicit
configuration paths resolve inside that root. Global flags select config,
format, output, verbosity, or config bypass. Invocation flags override loaded
output values without modifying configuration files. A relative report output
path remains relative to the invoking working directory.

The selected root is stored as invocation-local runtime configuration.
Collectors and Engine resolve every relative input through it; the command
never calls `os.Chdir`, and finding locations inside the root remain
project-relative.

### Outputs and errors

Successful validation writes exactly one selected report. Validation findings
produce a complete report and non-zero exit. Configuration, traversal,
collection, serialization, or output failures return command errors because a
complete validation result cannot be guaranteed.

Verbose messages are diagnostic only and remain on stderr. JSON stdout must
never contain banners, progress, or duplicate error prose.

### Side effects and semantic boundary

`run`, `lint`, `skills`, and `docs review-context` are read-only except for an
explicit report output. `generate skill` writes only its explicit output.
Document mutations are limited to `docs init` and `docs fix`.

A successful process status means enabled structure and traceability rules
passed. It does not certify the accuracy or sufficiency of human-authored
requirements and architecture.

## Document initialization

`docs init <package>...`:

- requires every requested project-relative package directory to exist;
- deduplicates repeated packages and preflights the whole batch before writing;
- creates or adopts exactly `design.md`, `contract.md`, `spec.md`, and
  `testing.md` under `docs/<package>/`;
- treats each nested source sub-package as an independent invocation and
  equally nested documentation set;
- writes only version/package identity because the exact basename is the sole
  role authority;
- never overwrites existing IDD or legacy metadata;
- refuses legacy marker or `related_files` metadata before writing anything;
- creates structural headings with stable `idd:scaffold` markers, never fake
  requirements or generated IDs;
- includes the initial deterministic incomplete-slot work list in its output.

The command may adopt an existing plain narrative or generic frontmatter, but
it refuses an already self-describing set because initialization is not an
update operation. All preflight checks complete before the first file write.
Successful output lists exactly the paths created or adopted.

## Document completion status

`docs status <docs-package-or-document>...` is read-only. It uses the same
shared role schema as generation and normal collection and reports schema
`idd.document_status.v1`, normalized first-occurrence targets, overall
`complete` or `incomplete` status, and a stable list of remaining slots. Each
work item contains file, line, role, slot, and reason. Directory input may
include multiple package document sets and reports missing sibling files as
incomplete work. Repeated or overlapping inputs do not duplicate a work item.
A split role name or IDD frontmatter on a non-canonical filename is rejected
rather than treated as another document role.

A split role name returns a directly executable agent prompt in the operational
error. It identifies source and canonical target, requires a complete
content-preserving merge, defers deletion until no-loss verification, forbids
another fragment, and supplies package-status and full-project validation
commands. The read-only status command never merges or removes files.

A generated slot is complete only when its scaffold marker is absent and the
bounded Markdown contains effective, non-placeholder authored content. Removing
a marker while leaving an empty heading does not make the status complete.
After a canonical record exists, missing required record fields appear as
their own stable slots; a half-authored record cannot make the status complete.

## Document repair

`docs fix <docs-package-or-document>...`:

- repairs minimal version/package identity while deriving role from the
  canonical filename;
- creates missing skeletons for a directory target;
- writes only the selected file for a file target;
- preserves the Markdown body byte-for-byte;
- therefore preserves scaffold markers and incomplete state;
- never reformats prose or invents semantic records;
- refuses malformed frontmatter rather than discarding unknown data.
- preflights every target, deduplicates overlapping write paths, and applies no
  writes when an expected input error occurs anywhere in the batch.

Repair neither accepts nor emits `idd.document` and never creates
`design-*`, `contract-*`, `spec-*`, or `testing-*` fragments. There is no
document line-count limit to repair around.

The file form has one possible write target. The directory form has four
structural targets and prepares every selected result before replacement.
An empty changed-path list is a successful idempotent repair.

Batch preflight is not a cross-filesystem transaction guarantee. If an
unexpected I/O failure occurs while an already validated write plan is being
applied, earlier writes remain individually atomic but may already be visible.

Repair is not a formatter and cannot resolve semantic findings. Requirements,
contracts, design decisions, purposes, scenarios, and coverage remain authored
Markdown even when a thinner summary would be easier to synthesize.

## Concrete implementation boundary

Collectors are the concrete `DocCollector` and `CodeCollector` structs.
Link construction and validation rules are methods in `internal/engine`.
There is no production `Collector`, `Linker`, or `Rule` extension interface.

The command surface is the `CLI` boundary. Its concrete collaborators are
`LinkageGraph`, `Config`, `Engine`, `Identifier`, and `Reporter`.
Embedded skill access is provided by `SkillsFS`, `ListEmbeddedSkills`, and
`ReadEmbeddedSkill`; the `SkillInfo` type is the serialized skill metadata
boundary.

`generate skill` writes the exact workflow embedded in the binary. Users
regenerate their installed skill after upgrading idd-cli so semantic authoring
instructions stay aligned with validator behavior.

## Contract: Config

**Guarantees:**

`Config` defines validation patterns, paths, output, and deprecated
compatibility fields.
The CLI loads it from the explicit path or documented project defaults before
applying invocation-only flag overrides.

`.idd.yaml` is configuration only. Package Components, Contracts, SPECs, TESTs,
and coverage remain in their owning Markdown records.

Missing optional values receive deterministic defaults. Unknown or inconsistent
annotation keys return configuration errors before collection. Deprecated
consistency values are accepted and ignored behaviorally, but explicit YAML
presence produces a warning that names the key, source, cleanup action, and
review-context replacement. `run` and `lint` include it in their report;
`docs review-context` writes it to stderr. Loading and validation do not rewrite
the configuration file.

## Contract: Engine

**Guarantees:**

`Engine` consumes collected identifiers and returns one validation result from
its concrete graph-building and validation methods.

Structural collector findings may be added before `Run` and must survive graph
validation. The command also supplies `CodeCollector`'s normalized
Tree-sitter analyses for Go, TypeScript/TSX, JavaScript/JSX, C++, Java, and
Python. `Run` merges document and code origins, derives typed links, applies
enabled rules in one deterministic lifecycle, sorts findings, and optionally
attaches a graph snapshot. Invalid source syntax is reported and never handled
through regex fallback. Validation findings are result data; unexpected
execution failures are returned errors.

The engine does not read authoring intent or mutate evidence. Rules may verify
presence, shape, and traceability, but cannot declare a
thin design explanation semantically adequate.

## Contract: Identifier

**Guarantees:**

`Identifier` preserves type, origin, source location, narrative description,
legacy and typed traceability links, derived-node state, and TEST kind across
collection and reporting.

Multiple instances with one ID may coexist when documentation and code provide
independent evidence. Set operations must retain origins for correspondence
checks while offering deterministic lookup for graph construction. TEST kind
distinguishes behavioral and contract annotations without creating a second ID
namespace.

## Contract: LinkageGraph

**Guarantees:**

`LinkageGraph` stores identifiers and directed traceability relationships with
deterministic lookup, verification, statistics, and snapshots.

Adding an existing node returns the established node rather than erasing
metadata. Edges retain direction, type, source, and verification status.
Backlink and typed-edge queries expose graph semantics without requiring
callers to traverse internal slices. Snapshots are serialization values and do
not permit mutation of the live graph.

## Contract: Reporter

**Guarantees:**

`Reporter` renders the same validation result as JSON, Markdown, or
LLM-oriented Markdown without mixing diagnostics into structured stdout.
LLM-oriented Markdown supports the Skill repair loop; JSON supports CI and the
final project gate.

Every actionable finding retains rule, severity, problem, fix guidance,
location, related identifier, and field context when available. Repeated
findings may be grouped in summaries but cannot be dropped from detailed
output. Unsupported formats and unwritable destinations return errors rather
than falling back silently.

## Contract: SkillInfo

**Guarantees:**

`SkillInfo` is the stable serialized description returned when embedded
workflow skills are listed.

It exposes frontmatter metadata and the source path only. Listing does not
execute or rewrite a Skill, and an unreadable individual local file does not
change embedded content.

## Contract: SkillsFS

**Guarantees:**

`SkillsFS` exposes the embedded IDD workflow files used by skill listing,
reading, and generation commands. The embedded workflow is the canonical
companion to that binary.

Paths are relative to the embedded filesystem. Reading an unknown path returns
an error. Export copies the exact embedded bytes so an installed Skill can be
compared directly with the binary-owned version.

## SPEC review context

`docs review-context <SPEC-ID>...` emits deterministic evidence after one
documentation scan and one source scan. One SPEC uses schema
`idd.spec_review_context.v1`; multiple SPECs use
`idd.spec_review_context_batch.v1` with first-request ordering. `--docs-path`
selects documentation input without conflicting with positional SPEC IDs.

Each context includes the selected canonical SPEC, named Contract guarantees,
covering TEST records, and annotated declaration excerpts.
Declaration evidence includes both implementation annotations for the SPEC and
test annotations for its covering TEST IDs. It never emits a semantic approval,
warning, or score and never changes the validation result. Requests are
deduplicated, bounded to ten unique IDs, and fail atomically if any ID cannot
produce trustworthy evidence.
