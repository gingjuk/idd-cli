# idd-cli

A CLI tool that validates bidirectional linkage consistency between IDD (Intent-Driven Development) identifiers across documentation and source code. Designed to help AI agents maintain traceable, well-documented code through structured specification, contract, test, and design documents.

## Purpose

idd-cli is an **Intent-Driven Development (IDD)** tool that enables:

- **Traceable development**: Link specs → contracts → tests → code with bidirectional references
- **AI agent assistance**: Provides structured documentation conventions that AI agents can understand and follow
- **Consistency validation**: Ensures all identifiers are properly linked and documented
- **Self-documenting code**: Code annotations (`@spec`, `@contract`, `@test`, `@design`) create traceable links to formal documentation

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

For IDD framework documentation, see `docs/idd-intro.md`.

## Installing Skills for AI Agents

idd-cli includes embedded skills that define the IDD workflow. Copy the skills directory to your project's `.claude/skills/` folder to enable IDD-aware agent behavior.

### Claude Code / OpenCode

```bash
# Copy skills to your project
mkdir -p .claude/skills
cp -r cmd/validator/skills/* .claude/skills/

# Or create a symlink (recommended for development)
ln -s $(pwd)/cmd/validator/skills .claude/skills/idd
```

### Codex ( Anthropic )

Add to your agent's system prompt or project instructions:

```
This project uses the IDD (Intent-Driven Development) framework.
Load skills from: cmd/validator/skills/SKILL.md
```

### Verifying Installation

```bash
# List installed skills
ls -la .claude/skills/

# Verify skills are accessible
idd-verify skills
```

## Project Structure

```
idd-link-validator/
├── cmd/validator/       # CLI entry point
├── internal/
│   ├── collector/        # Doc/code identifier collection
│   ├── engine/          # Validation engine
│   ├── graph/           # Linkage graph
│   └── model/           # Data models
├── docs/                # IDD documentation
│   └── backend/        # Backend module docs
├── examples/            # Example configurations
└── skills/             # IDD skill definitions
```
