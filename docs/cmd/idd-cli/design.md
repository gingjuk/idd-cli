---
idd:
  version: "1.0"
  package: cmd/idd-cli
  document: design
---

# Design: cmd/idd-cli

## Component: IDDCLIModule

`IDDCLIModule` is the composition root for validation, document maintenance,
reporting, and embedded workflow commands.

## Architecture

`IDDCLIModule` wires concrete collectors, the engine, and the reporter. Link
building and validation live inside `internal/engine`; no standalone Linker or
Validator interface exists.

Skill export, document maintenance, and validation are separate command paths
that form one user workflow:

```text
generate skill ── embedded SkillsFS ── agent authoring instructions

docs init/fix ── safe package-document structure

agent-authored Markdown + code annotations
        │
        ▼
run/lint ── DocCollector + CodeCollector ── Engine.Run ── Reporter
   ▲                                                        │
   └──────────────── llm-markdown repair loop ──────────────┘

final run . --format json ── project gate
```

The binary-embedded skill defines semantic authoring behavior. `docs init/fix`
never replace that behavior: they only create or normalize safe structure.
`run .` is the complete validity gate because source collection always starts
from the current project working directory.

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

## Function Composition

1. `generateSkill` exports the workflow paired with the current binary.
2. `InitDocuments` bootstraps a new package without semantic placeholders.
3. The skill authors documents, tests, code, and annotations.
4. `RepairDocuments` handles safe identity normalization when requested.
5. `DocCollector.Collect` gathers self-describing or legacy documentation IDs.
6. `CodeCollector.Collect` gathers source annotation IDs from the project.
7. `Engine.Run` builds and validates the linkage graph.
8. `Reporter.Write` emits LLM Markdown for repair or JSON for the final gate.

## Testability Hooks

- Collector tests use temporary package/document trees.
- Graph tests pre-populate nodes and edges directly.
- Engine tests use focused file patterns and concrete identifiers.
- Reporter tests inspect serialized findings.
- `cmd/idd-cli` remains thin; domain behavior is tested in internal packages.

## Dependencies

- `internal/collector` - For identifier collection
- `internal/engine` - For orchestration
- `internal/graph` - For graph structure
- `internal/model` - For data types
- `internal/reporter` - For output
- `internal/similarity` - For consistency checking
- `pkg/pattern` - For IDD patterns
- `pkg/walk` - For file traversal
- `gopkg.in/yaml.v3` - For configuration and minimal document identity parsing
- `github.com/yuin/goldmark` - For CommonMark record parsing with source positions
- `github.com/spf13/cobra` - For CLI
