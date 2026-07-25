---
related_files:
  design: docs/architecture/design.md
---

# Design (architecture)

**Status:** Done

## Architecture

The embedded IDD skill is the authoring control plane; idd-cli is the
structural and traceability enforcement plane. The binary exports the same
skill version that describes its accepted document model.

idd-cli follows a graph-first validation architecture with two inputs:

1. `DocCollector` reads self-describing document sets or legacy Markdown metadata.
2. `CodeCollector` reads source annotations.
3. `Engine.buildGraph` creates the `LinkageGraph` and derives relationships.
4. Validation methods in `internal/engine/engine.go` check the graph and files.
5. `Reporter` converts results into finding-centered JSON or Markdown.

There is no separate Linker or Validator interface in the current
implementation.

```text
Embedded IDD skill ── generate skill ──→ Agent authoring workflow
                                             │
docs init/fix ── safe structure ─────────────┤
                                             ▼
                         Markdown records + code annotations
                                  │                  │
                                  ▼                  ▼
                            DocCollector       CodeCollector
                                  └─────────┬────────┘
                                            ▼
                                        Engine.Run
                                            │
                                            ▼
                                         Reporter
                                  llm-markdown │ JSON
                                               │
                           Agent repair loop ←──┘──→ CI gate
```

The final validity check runs as `idd-cli run .` from project root.
Documentation targets can narrow `DocCollector`, but `CodeCollector` still
scans the current project working tree.

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

1. **generate skill** - Exports the workflow embedded in the binary.
2. **docs init/fix** - Creates skeletons or normalizes safe identity.
3. **DocCollector.Collect()** - Parses self-describing sets first, then legacy Markdown.
4. **CodeCollector.Collect()** - Collects annotations from the project working tree.
5. **Engine.Run()** - Builds the graph and invokes enabled validation methods.
6. **Reporter.Write()** - Emits agent-, human-, or automation-oriented findings.

## Testability Hooks

- Collectors are tested with temporary package and documentation trees.
- IDD document tests cover minimal YAML identity and CommonMark records.
- Graph tests pre-populate identifiers and edges directly.
- Engine tests use temporary files and focused configuration patterns.
- Reporter tests write to temporary outputs and inspect structured findings.

## Dependencies

- `internal/collector` - For identifier collection
- `internal/engine` - For orchestration
- `internal/graph` - For graph structure
- `internal/model` - For data types
- `internal/reporter` - For output
- `internal/similarity` - For consistency checking
- `pkg/pattern` - For IDD patterns
- `pkg/walk` - For file traversal
- `gopkg.in/yaml.v3` - For config and document identity parsing
- `github.com/yuin/goldmark` - For CommonMark AST parsing and source positions
- `github.com/spf13/cobra` - For CLI
