---
idd:
  version: "1.0"
  package: internal/collector
  document: design
---

# Design: internal/collector

## Component: CollectorModule

`CollectorModule` translates human-readable documentation and source
annotations into the common identifier model used by the validation engine.

## Architecture

`CollectorModule` contains two concrete collectors; there is no shared
Collector interface.

```text
DocCollector
    ├── idd_document.go    CommonMark records, validation, and derived links
    ├── document_files.go  non-overwriting init and per-file repair
    ├── doc_collector.go   self-describing/legacy mode selection
    └── frontmatter.go     legacy Markdown parsing

CodeCollector
    └── code_collector.go  source annotations and declaration context
```

`DocCollector.Collect` performs two passes. It first discovers directories
where one of the four fixed files has an `idd` frontmatter block, then parses
each four-file set as a unit. Markdown elsewhere continues through the legacy
frontmatter parser.

## Package Layout

```text
internal/collector/
├── idd_document.go
├── document_files.go
├── code_collector.go
├── doc_collector.go
└── frontmatter.go
```

## Function Composition

1. `DocCollector.Collect` discovers and sorts inputs.
2. `collectIDDDocumentSet` parses identity frontmatter and role-owned Markdown
   records.
3. `addIDDDocumentIdentifiers` derives SPEC backlinks from TEST coverage.
4. `validateIDDDocumentMarkdown` checks body references without redefining IDs.
5. `collectFile` handles a legacy Markdown package.
6. `CodeCollector.Collect` independently collects source annotations.

Initialization and repair are separate from collection:
`InitDocuments → MarshalIDDDocument → writeExclusiveFile`, while
`RepairDocuments → normalizeIDDDocument → replaceFileAtomically`.

## Dependencies

- `internal/config` supplies ignore patterns and collection settings.
- `internal/model` supplies identifiers and validation errors.
- `pkg/pattern` validates identifier syntax.
- `gopkg.in/yaml.v3` parses and emits identity and legacy frontmatter.
- `github.com/yuin/goldmark` parses CommonMark records with source positions.

## Testability Hooks

- Pure IDD document parsing and marshaling accept byte slices.
- Identity diagnostics use YAML node locations; semantic diagnostics use
  Markdown AST source positions.
- File operations are tested in temporary package trees.
- Self-describing and legacy modes share the same `IdentifierSet` output
  contract.
