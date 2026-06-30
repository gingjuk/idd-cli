---
related_files:
  spec: docs/internal/reporter/spec.md
  contract: docs/internal/reporter/contract.md
  design: docs/internal/reporter/design.md
  testing: docs/internal/reporter/testing.md
---

# Design (reporter)

**Status:** Done

## Architecture

The reporter module handles output formatting:

1. **Format Agnostic**: Reporter accepts format parameter and delegates formatting
2. **Output Abstraction**: Write method abstracts stdout vs file output
3. **Complete Reports**: Generate always creates full report structure with metadata
4. **LLM Findings**: JSON and LLM Markdown derive self-contained findings from the complete validation result
5. **Finding Groups**: Repeated rule failures are summarized by severity and rule before detailed findings

## Package Layout

```text
internal/reporter/
└── reporter.go       # Main reporter implementation and output formatters
```

## Function Composition

1. **Generate()** - Creates full report with metadata and results
2. **Write()** - Handles output destination (stdout or file)
3. **writeMarkdown()** - Formats result as Markdown
4. **buildLLMReport()** - Converts complete validation reports into finding-centered LLM reports
5. **writeLLMMarkdown()** - Formats LLM findings as Markdown

## Testability Hooks

- Reporter can write to bytes.Buffer for testing output
- Format functions are pure transformations
- No side effects in report generation
- LLM report construction is testable without invoking the CLI

## Dependencies

- `internal/model` - For report, validation result, graph snapshot, and LLM report types
