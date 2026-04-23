---
related_files:
  spec: spec.md
  contract: contract.md
  design: design.md
  testing: testing.md
---

# Design (reporter)

**Status:** Done

## Architecture

The reporter module handles output formatting:

1. **Format Agnostic**: Reporter accepts format parameter and delegates formatting
2. **Output Abstraction**: Write method abstracts stdout vs file output
3. **Complete Reports**: Generate always creates full report structure with metadata

## Package Layout

```text
internal/reporter/
├── reporter.go       # Main reporter implementation
├── json.go           # JSON format output
└── markdown.go       # Markdown format output
```

## Function Composition

1. **Generate()** - Creates full report with metadata and results
2. **WriteJSON()** - Formats result as JSON
3. **WriteMarkdown()** - Formats result as Markdown
4. **Write()** - Handles output destination (stdout or file)

## Testability Hooks

- Reporter can write to bytes.Buffer for testing output
- Format functions are pure transformations
- No side effects in report generation

## Dependencies

- `internal/engine` - For ValidationResult type
- `github.com/fatih/color` - For colored terminal output (optional)
