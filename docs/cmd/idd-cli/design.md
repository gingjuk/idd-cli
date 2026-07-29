---
idd:
  version: "1.0"
  package: cmd/idd-cli
---

# Design: cmd/idd-cli

## Component: IDDCLIModule

**Purpose:**

`IDDCLIModule` is the composition root for validation, document maintenance,
reporting, and embedded workflow commands. It turns Cobra arguments into calls
to concrete domain packages and translates their results into process output
and exit behavior.

The component owns command discovery, flag precedence, configuration
selection, collaborator construction, output routing, and user-facing error
boundaries. It does not parse Markdown records, build links, implement
validation rules, decide architectural meaning, or synthesize semantic
documentation. Keeping those responsibilities outside `main.go` prevents the
CLI surface from becoming a second implementation of the domain model.

### Design rationale

The command package deliberately uses concrete collaborators rather than
advertising extension interfaces that do not exist. This keeps the executable
easy to trace: every command path ends in a specific collector, engine,
reporter, embedded filesystem, or document-file operation.

The embedded Skill and the CLI ship together because format instructions and
validation behavior must evolve as one release. Embedding does not make the CLI
the semantic author; it gives the agent the correct authoring contract.

## Architecture

`IDDCLIModule` wires concrete collectors, the engine, and the reporter. Link
building and validation live inside `internal/engine`; no standalone Linker or
Validator interface exists.

Skill export, document maintenance, and validation are separate state and
failure domains that form one user workflow:

```text
generate skill ── embedded SkillsFS ── agent authoring instructions

docs init/status/fix ── scaffold work list + safe package-document structure

agent-authored Markdown + code annotations
        │
        ▼
run/lint ── DocCollector + CodeCollector ── Engine.Run ── Reporter
   ▲                                                        │
   └──────────────── llm-markdown repair loop ──────────────┘

final run . --format json ── project gate
```

The binary-embedded skill defines semantic authoring behavior.
`docs init/status/fix` never replace that behavior: they create visible work
slots, inspect completion, or normalize safe structure.
`run .` is the normal complete validity gate. When another directory is
supplied, `run` resolves it as one project root and executes configuration
discovery, document/source collection, and engine filesystem checks inside that
root. The caller's working directory is restored before command completion.

### Skill discovery and export

`skills` first inspects a requested or default local `skills/` directory.
When it does not exist or contains no Markdown skills, the command reads the
binary-embedded filesystem. It parses only frontmatter needed for listing and
does not reinterpret the Skill body.

`generate skill` accepts the supported `skill` target, reads
`skills/SKILL.md` from the embedded filesystem, and copies its bytes to stdout
or an explicit file. It does not create parent directories or install into an
agent product, because installation conventions are outside this repository.
The user is expected to install the exported file through the chosen agent's
normal mechanism.

### Document command path

`docs init <package>...` treats each package path as project-relative and
delegates to `InitDocumentPackages`. The collector layer performs path,
package-existence, legacy-metadata, and central-catalog preflight for the whole
batch before any write. It creates the four canonical basenames with
version/package-only frontmatter; a nested package is initialized independently
at its matching nested path. After creation it runs
`InspectDocumentCompletions` and serializes normalized targets, the changed-path
list, and the exact incomplete-slot work list.

`docs status <path>...` calls the same completion inspector without mutation.
It preserves normalized first-request target order, deduplicates overlapping
work, and globally sorts the result. JSON is the agent-oriented stable wire
shape; human output summarizes completion or lists each slot with its location
and reason. Reusing one inspector prevents generation, validation, and status
from drifting into different definitions of “filled”. The inspector advances
from collection and section slots to record-field slots as authors add records,
so every missing required field remains directly actionable.

When traversal sees a role-derived split filename, completion inspection stops
because a fragment cannot be assigned an independent completion state. Its
error carries the source, adjacent canonical target, content-preservation
rules, safe deletion order, and both verification commands. The normal
validation path carries the same source/target evidence into the reporter,
which emits a per-finding prompt in JSON and LLM Markdown. Group summaries stay
generic so a multi-package group does not select only its first file.

`docs fix <path>...` delegates file-vs-directory semantics to
`RepairDocumentTargets`. A file target can change only that named document; a
directory target owns version/package identity for all four canonical roles
and may create missing skeletons. All targets are planned before writes and
overlapping file plans are deduplicated. Role comes only from the exact
basename; the command never interprets a validation finding as permission to
rewrite prose or split a long document.

`run` and `lint` keep one positional project root. They build one graph from
documents and source in that same tree, preventing a caller in one worktree
from combining its source with documents selected from another. Focused
document inspection remains the responsibility of `docs status` and
`docs review-context --docs-path`.

`docs review-context <SPEC-ID>...` retains a distinct evidence-gathering scope
and bypasses the engine. `--docs-path` defaults to `.`, and the source root
remains the current working tree. Moving the path to a flag keeps every
positional value unambiguously a SPEC identifier.

The collectors run once, then the builder resolves up to ten unique SPEC owners
in first-request order. Every context gathers complete bounded record Markdown,
the named Contract, covering TEST records, and AST-bound implementation and
covering-TEST declaration excerpts. The reporter projects one context using
the original schema or several using a batch schema, without adding a finding,
score, or verdict.

Mutation commands report structural changes as JSON by default or as a
human-readable path list for other formats. An empty change list is a
successful idempotent result, not an error.

### Validation command path

`run` and `lint` share the same handler. The handler:

1. resolves the positional directory as a project root, defaulting to `.`;
2. enters that root and loads configuration in documented precedence order
   unless `--no-config` is set;
3. applies invocation-only format, output, and verbosity flags;
4. collects documentation and parses supported source files from that same
   root into normalized Tree-sitter declarations and attached annotations;
5. adds structural and source-parse findings and supplies source analyses to
   the engine before graph validation;
6. builds a report from the completed result and writes it once;
7. restores the caller's working directory; and
8. returns a validation error after report emission when the result is invalid.

Collection or I/O failures stop the path because the report would be based on
incomplete evidence. Validation failures do not prevent report generation,
because the report is the artifact required for repair.

### Output and process boundary

JSON stdout must remain parseable. Verbose progress and diagnostics therefore
use stderr. `-o` selects a file explicitly; otherwise the reporter uses stdout.
The three formats are projections of the same result rather than independent
validation runs:

- JSON is stable for CI and programmatic repair tooling;
- Markdown is a concise human readout;
- LLM Markdown emphasizes ownership, location, problem, and repair guidance
  for validation, or neutral review questions for a review-context bundle.

A zero exit proves that enabled structural and traceability checks passed. It
does not prove that the human-authored design or requirement is semantically
sufficient; that review remains in the paired Skill workflow.

## State and mutation boundaries

Validation is read-only with respect to project sources and documents. It may
write only the requested report destination. Skill generation writes only the
requested output file.

Document initialization and repair are the only commands that mutate package
documentation; status is read-only. Their domain implementation performs all
safety checks and atomic replacements, so Cobra handlers never partially
manage a document set.
There is no daemon, database, global cache, or background process; command
state ends with the process.

## Package Layout

```text
idd-cli/
├── cmd/idd-cli/       # CLI entry point
├── internal/
│   ├── collector/        # Doc/code identifier collection
│   ├── engine/           # Validation engine
│   ├── graph/            # Linkage graph
│   ├── model/            # Data models
│   └── reporter/         # Output formatting
└── pkg/
    ├── pattern/          # IDD regex patterns
    └── walk/             # File system traversal
```

`main.go` owns Cobra composition and serialization of command results.
`embed.go` owns the embedded filesystem boundary. Domain packages stay
independent of Cobra so they can be tested and reused without a process-level
fixture.

## Function Composition

1. `rootCmd.Execute` dispatches a Cobra command and prints only command-level
   errors to stderr.
2. `listSkills → loadSkills/loadEmbeddedSkills → parseSkillFrontmatter`
   produces stable `SkillInfo` values without modifying Skill content.
3. `generateSkill → ReadEmbeddedSkill` copies the paired workflow.
4. `initPackageDocs → InitDocuments → writeDocChanges` creates a new structural
   set and reports exact changed paths plus its incomplete work list.
5. `statusPackageDocs → InspectDocumentCompletion → writeDocumentStatus`
   reports remaining slots without mutation.
6. `repairPackageDocs → RepairDocuments → writeDocChanges` normalizes derived
   identity while preserving narrative bodies.
7. `run → DocCollector.Collect + CodeCollector.CollectWithErrors →
   Engine.SetSourceAnalyses → Engine.Run →
   BuildReport → Reporter.Write` performs the full validation lifecycle.

Semantic authoring sits deliberately between the document and validation
commands: the Skill edits design, contract, SPEC, TEST, source, and test
evidence; the CLI then checks the resulting graph.

## Testability Hooks

- CLI handlers are small enough that domain behavior is exercised in the
  internal package that owns it rather than through fragile process mocks.
- Collector tests use temporary package/document trees for real path,
  preflight, preservation, and idempotence behavior.
- Graph tests pre-populate nodes and edges directly, while engine tests use
  focused configurations and repository fixtures only for file-dependent
  rules.
- Reporter tests inspect JSON and Markdown serialization independently from
  command dispatch.
- Build-and-run acceptance checks exercise the real embedded Skill, command
  help, report output, and repository-wide validity gate.

## Dependencies

- `internal/collector` - For identifier collection
- `internal/engine` - For orchestration
- `internal/graph` - For graph structure
- `internal/model` - For data types
- `internal/reporter` - For output
- `internal/collector` - For bounded single-SPEC review evidence
- `pkg/pattern` - For IDD patterns
- `pkg/walk` - For file traversal
- `gopkg.in/yaml.v3` - For configuration and minimal document identity parsing
- `github.com/yuin/goldmark` - For CommonMark record parsing with source positions
- `github.com/tree-sitter/go-tree-sitter` and pinned official grammars - For
  declaration-aware binding in the seven supported languages
- `github.com/spf13/cobra` - For CLI

## Trade-offs and extension boundary

Keeping `cmd/idd-cli` thin means some end-to-end behavior is distributed across
collector, engine, and reporter tests rather than command-package unit tests.
The repository-level binary self-check supplies the missing composition
evidence.

New output formats belong in the reporter and new document semantics belong in
the collector/engine contract. New Cobra commands should compose those
capabilities without duplicating their rules. A future plugin interface would
require a deliberate public lifecycle and compatibility contract; this design
does not imply one today.
