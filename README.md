# IDD Link Validator

A CLI tool that validates bidirectional linkage consistency between IDD identifiers across documentation and source code.

## Quick Start

```bash
# Build
go build -o bin/idd-verify ./cmd/validator

# Run validation
./bin/idd-verify run .

# With config
./bin/idd-verify --config idd.yaml run .

# Or use Make
make build && make run
```

## Identifier Format

Format: `TYPE-MODULE-NUMBER` (e.g., `SPEC-BE-001`, `TEST-BE-001`)

| Prefix | Meaning |
|--------|---------|
| `SPEC-` | Functional specification |
| `CONTRACT-` | Interface/behavior contract |
| `TEST-` | Test case |
| `DESIGN-` | Architecture design decision |

## Code Annotations

```go
// @spec SPEC-BE-001
// @contract CONTRACT-BE-001
// @test TEST-BE-001
// @design DESIGN-BE-001
```

## Configuration

See `examples/idd.yaml` for a full configuration example with comments.

## IDD Documentation System

For IDD framework documentation, see `docs/IDD.md`.