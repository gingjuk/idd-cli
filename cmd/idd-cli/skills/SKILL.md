---
name: intent-driven-development
description: Intent-Driven Development (IDD) workflow — from user intent to committed code through structured phases
license: MIT
compatibility: opencode
metadata:
  audience: agents
  workflow: development
  protected: true  # Modification requires explicit approval
---

# Intent-Driven Development (IDD)

> ⚠️ **Protected Skill** — Changes to this skill require explicit user approval. Do not modify without consulting the user first.

IDD is a framework-agnostic development workflow that transforms user intent into implemented, tested, and committed code through defined phases.

## Core Philosophy

**Four Pillars of IDD:**

1. **Intent is Truth** — Every deliverable originates from an approved `intent.md`. Nothing exists unless it is defined in intent and approved. Implementation without intent is guesswork.

2. **Specs are Code** — Specifications must be formalized, precise, and unambiguous. Like production code, specs undergo review, version control, and validation. Vague specs produce vague implementations.

3. **Tests Prove Intent** — Tests verify behavior against specs, not implementation details. A test proving "the function returns X for input Y" is valid; a test proving "the function calls dependency Z" is not.

4. **Nothing Guessed** — When anything is unclear, CLARIFY. Never guess, assume, or fill gaps with assumptions. Ambiguity resolved by the implementor (not the author) introduces interpretation risk. Ask, don't guess.

## Workflow

```text
User Intent
    ↓
Intent Capture (intent.md)
    ↓
Planning (plan.md)
    ↓
[Architecture] → Consulting Agent → design.md/spec.md/contract.md
    ↓
TDD Implementation → Implementation Agent + tdd-workflow
    ↓
Contract Validation (if applicable)
    ↓
Review → Review Agent / Consulting Agent
    ↓
Commit
```

**Phase Details:**

| Phase | Input | Output | When |
|-------|-------|--------|------|
| 1. Intent Capture | User request | `intent.md` | Always — before any code |
| 2. Planning | `intent.md` | `plan.md` | Always |
| 3. Architecture | `intent.md`, `plan.md` | `design.md`, `spec.md`, `contract.md`, `testing.md` | Complex features |
| 4. Implementation | Specs, contracts | Code with annotations | Always |
| 5. Contract Validation | `*_contract_test.go` | Validated interfaces | If `contract.md` exists |
| 6. Review | Changes | Approved code | Always |
| 7. Knowledge Capture | Implementation | Updated memory/docs | Always |

**TDD Workflow:**

1. Write test first (RED) — test should FAIL
2. Write minimal implementation (GREEN) — test should PASS
3. Refactor (IMPROVE) — verify coverage 80%+

**Troubleshoot:** check test isolation → verify mocks → fix implementation (not tests, unless tests are wrong).

## Document Dependency Chain

Documents are created in dependency order. Each document type builds upon or validates the previous:

```text
User Intent
     ↓
contract.md          ← User Intent defines WHAT needs to be built (interfaces, capabilities)
     ↓
design.md            ← Contract refines HOW to structure it (architecture, components)
     ↓
spec.md              ← Design breaks down HOW in detail (function signatures, behaviors)
     ↓
Implementation       ← Spec guides implementation (code annotated with @implement SPEC-*)
     ↓
testing.md           ← Each SPEC requires TEST coverage (annotated with @test TEST-*)

Contract Validation (separate from testing.md):
spec.md ──implements──→ contract.md
         └───implements──→ design.md

testing.md contains TEST identifiers (via `@test`) that validate SPECs

Contract test file (`*_contract_test.go`) uses `@test-contract SPEC` to validate SPEC satisfies contract
```

**Dependency Rules:**

| Document | Depends On | Provides To |
|---------|-----------|------------|
| `intent.md` | User request | Starting point for all planning |
| `contract.md` | `intent.md` | Interface signatures, capability definitions |
| `design.md` | `contract.md` | Architecture decisions, package layout |
| `spec.md` | `contract.md`, `design.md` | Detailed specifications, **implements** contract and design |
| `testing.md` | `spec.md` | TEST identifiers (via `@test`) that validate SPECs |
| Contract test (`*_contract_test.go`) | `spec.md`, `contract.md` | Uses `@test-contract` to validate SPEC satisfies contract |
| Code | `spec.md`, `design.md` | Implementation annotated with `@implement SPEC-*` |

**Note:** `spec.md` has two roles:

1. **Implements contract** — Each SPEC declares which contract interface it implements
2. **Implements design** — SPEC details (function signatures, behaviors) realize the architecture defined in design.md

**Spec implements Design and Contract:**

Each SPEC in `spec.md` must state which design and contract it implements. Required fields must appear first, in this exact order: `Design`, `Contract`, `Requirement`, `Tests`. Optional fields must come after all required fields.

```markdown
## SPEC-AUTH-001: User Login

**Design:** `design.md` — implements architecture `AuthModule`

**Contract:** `contract.md` — implements interface `Authenticator`

**Requirement:** Login with email/password returning JWT token.

**Tests:** `TEST-AUTH-001`

**Implementation:** `internal/auth/auth.go`
```

**Required fields:**

| Field | Required | Description |
|-------|----------|-------------|
| Design | Yes | Which design component this SPEC implements |
| Contract | Yes | Which contract interface this SPEC implements |
| Requirement | Yes | What this SPEC describes |
| Tests | Yes | TEST identifiers that validate this SPEC |
| Implementation | No | Path to source code; must appear after required fields |

**Test validates SPEC (not contract directly):**

Each TEST in `testing.md` must state which SPEC it covers via **Spec Coverage:**

```markdown
## TEST-AUTH-001: User Login Test

**Spec Coverage:** `SPEC-AUTH-001`

**Purpose:** Verify login returns valid JWT for valid credentials.
```

**Note:** Contract validation is done via `@test-contract` in contract test files (`*_contract_test.go`), not in regular test files.

**Code implements SPEC (and implicitly implements Contract):**

```go
// @implement SPEC-AUTH-001
func Login(email, password string) (string, error) {
    // Implementation
}
```

**Contract Test File (`*_contract_test.go`):**

Contract tests explicitly validate that the implementation satisfies contract interfaces. Each test function is annotated with `@test-contract SPEC-<MODULE>-NNN` to declare which SPEC's contract is being verified.

```go
// Package auth provides contract tests for authentication.
//
// Spec: docs/internal/auth/spec.md
// Test: docs/internal/auth/testing.md
// Contract: docs/internal/auth/contract.md
package auth

// @test-contract SPEC-AUTH-001
func TestContractLogin(t *testing.T) {
    // Validate Login() satisfies Authenticator interface
}
```

## Document Hierarchy

All IDD documents stored in `docs/` directory, organized by package path:

```text
docs/
├── <pkg_path>/           # mirrors the package directory (e.g., internal/auth, pkg/pattern)
│   ├── spec.md           # SPEC-<MODULE>-NNN
│   ├── contract.md       # Interface signatures (references SPECs)
│   ├── testing.md        # TEST-<MODULE>-NNN
│   └── design.md         # Architecture decisions (references SPECs)
```

**Rules:**

- `<pkg_path>` mirrors the package directory relative to project root (e.g., `internal/auth` → `docs/internal/auth/`)
- `<MODULE>` in identifiers is derived from `<pkg_path>` using the abbreviation table (see **Module prefix** below)
- Document filenames are lowercase: `spec.md`, `contract.md`, `testing.md`, `design.md`
- `contract.md` and `design.md` do not have their own identifiers — they reference SPECs
- If spec.md exceeds **1500 lines**, split into `spec-<feat>.md`; main `spec.md` becomes an index

**Agent Responsibilities:**

Delegate to **specialized agents** for specific documents. Use the right agent for the right job:

| Document | Responsible Agent | Notes |
|----------|-------------------|-------|
| `intent.md` | Human (user) or planning agent | Captures user intent |
| `plan.md` | Planning agent | Creates and maintains plan |
| `contract.md` | Architecture/consulting agent | Defines interfaces from intent |
| `design.md` | Architecture/consulting agent | Designs system structure |
| `spec.md` | Deep implementation agent | Detailed specifications |
| `testing.md` | Deep implementation agent | Test cases |
| Code | Deep implementation agent | Annotated with `@implement` |
| Code review | Review agent | Quality, security check |

**Frontmatter:** Each IDD markdown file MUST include YAML frontmatter:

```yaml
---
markers:
  - id: SPEC-<MODULE>-NNN
    name: <description>

related_files:
  spec: docs/<pkg_path>/spec.md
  contract: docs/<pkg_path>/contract.md
  design: docs/<pkg_path>/design.md
  testing: docs/<pkg_path>/testing.md
---
```

`<pkg_path>` is the package directory path relative to the project root (e.g., `internal/auth`, `pkg/pattern`, `cmd/idd-cli`). All four entries are always required — use the full path even if the file does not exist yet.

## Document Templates

### Planning Templates

**`intent.md`** (in `<module>/.planning/`):

```markdown
# Intent: <feature name>

## User Request
[What the user asked for — verbatim or paraphrased]

## Success Criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Constraints
- [ ] Constraint 1
- [ ] Constraint 2

## Out of Scope
- Item 1
- Item 2
```

**`plan.md`** (in `<module>/.planning/`):

```markdown
# Plan: <feature name>

## Status
- [ ] Not started

## Phase 1: [Name]
| Task | Status | Notes |
|------|--------|-------|
| Task description | pending | |

## Risks
| Risk | Impact | Mitigation |
|------|--------|------------|
| Risk description | High/Med/Low | Mitigation approach |

## Next Action
[Immediate next step to take]
```

### IDD Document Templates

**`spec.md`** format:

````markdown
---
markers:
  - id: SPEC-<MODULE>-001
    name: <description>

related_files:
  spec: docs/<pkg_path>/spec.md
  contract: docs/<pkg_path>/contract.md
  design: docs/<pkg_path>/design.md
  testing: docs/<pkg_path>/testing.md
---

# Specification (<module>)

## SPEC-<MODULE>-001: <Title>

**Design:** implements architecture `<ComponentName>`

**Contract:** implements interface `<InterfaceName>`

**Requirement:** [What this spec describes]

**Tests:** `TEST-<MODULE>-001`

**Input Sources:** [API endpoints, user inputs, external dependencies]

**Implementation:** `<path to source code>`

**Implements:**

```go
func FunctionName(param1 string, param2 int) (Result, error)
```

**Parameters:**

- `param1`, Description of param1
- `param2`, Description of param2

**Returns:**

- `Result`, Description of return value
- `error`, Error description

````

**Required sections:** Design, Contract, Requirement, Tests
**Optional sections:** Input Sources, Implementation, Implements (with Parameters/Returns). Optional sections must come after all required sections.

**Interface/Class method naming:** When documenting an interface method or class member method, use `<InterfaceName>.<Method>` as the title (e.g., `Reporter.Write`). The Implements section signature should show the full method signature.

**`testing.md`** format:

```markdown
---
markers:
  - id: TEST-<MODULE>-001
    name: <test description>

related_files:
  spec: docs/<pkg_path>/spec.md
  contract: docs/<pkg_path>/contract.md
  design: docs/<pkg_path>/design.md
  testing: docs/<pkg_path>/testing.md
---

# Test Cases (<module>)

## TEST-<MODULE>-001: <Title>

**Purpose:** [What this test validates]

**Spec Coverage:** `SPEC-<MODULE>-001`
```

**`design.md`** format:

```markdown
---
related_files:
  spec: docs/<pkg_path>/spec.md
  contract: docs/<pkg_path>/contract.md
  design: docs/<pkg_path>/design.md
  testing: docs/<pkg_path>/testing.md
---

# Design (<module>)

## Architecture
[High-level system design]

## Package Layout
[Directory structure and package responsibilities]

## Function Composition
**Call graph:** `A → B → C`
**Initialization:** `InitA()` → `InitB()`

## Dependencies
**External:** [postgres, redis, kafka]
**Internal:** [pkg/auth, pkg/metrics]

## Testability Hooks
[Test strategy and contract test approach]
```

**Required sections:** Architecture, Package Layout, Function Composition, Dependencies, Testability Hooks

**`contract.md`** format:

````markdown
---
related_files:
  spec: docs/<pkg_path>/spec.md
  contract: docs/<pkg_path>/contract.md
  design: docs/<pkg_path>/design.md
  testing: docs/<pkg_path>/testing.md
---

# Contracts (<module>)

## Interface: <Name>

```go
type <Interface> interface {
    Method() error
}
```

**Implements:** `SPEC-<MODULE>-001`
````

**Required sections:** Interface definition with Implements field
**Note:** `contract.md` does not have its own identifiers — interfaces are referenced by SPECs via the **Contract:** field

### Code Annotation Templates

Where `<pkg_path>` is the package directory relative to project root (e.g., `internal/auth`, `pkg/pattern`). `<MODULE>` is the abbreviated identifier derived from `<pkg_path>` — see **Module prefix** rules below.

**Source file** (after package declaration):

```go
// Package mymodule provides <description>.
//
// Spec: docs/<pkg_path>/spec.md
// Contract: docs/<pkg_path>/contract.md
package mymodule

// @implement SPEC-<MODULE>-001
func PublicFunction() {
    // implementation
}
```

**Test file** (`<pkg>_test.go`):

```go
// Package mymodule provides tests for <description>.
//
// Spec: docs/<pkg_path>/spec.md
// Test: docs/<pkg_path>/testing.md
package mymodule

// @test TEST-<MODULE>-001
func TestPublicFunction(t *testing.T) {
    // test implementation
}
```

**Contract test file** (`<pkg>_contract_test.go`):

Contract tests validate that the implementation satisfies the contract interface. Each test function must be annotated with `@test-contract SPEC-<MODULE>-NNN` — referencing the SPEC whose contract is being verified.

```go
// Package mymodule provides contract tests for <description>.
//
// Spec: docs/<pkg_path>/spec.md
// Test: docs/<pkg_path>/testing.md
// Contract: docs/<pkg_path>/contract.md
package mymodule

// @test-contract SPEC-<MODULE>-001
func TestContractPublicFunction(t *testing.T) {
    // Validate implementation matches contract
}
```

## Identifier Format

**Pattern:** `<TYPE>-<MODULE>-<NUMBER>` — Exactly 3 segments separated by hyphens

| Prefix | Meaning | Example |
| ------ | ------- | ------- |
| `SPEC-` | Functionality specification | `SPEC-INTERNAL_AUTH-001` |
| `TEST-` | Test case | `TEST-INTERNAL_AUTH-001` |

**Module prefix** derived from `<pkg_path>` by a deterministic rule — no lookup table required:

1. Split `<pkg_path>` by `/`.
2. Uppercase each component; replace hyphens (`-`) with underscores (`_`).
3. Join components with `_`.

**Examples:**

| `<pkg_path>`          | `<MODULE>`              |
|-----------------------|-------------------------|
| `internal/auth`       | `INTERNAL_AUTH`         |
| `internal/engine`     | `INTERNAL_ENGINE`       |
| `internal/collector`  | `INTERNAL_COLLECTOR`    |
| `internal/config`     | `INTERNAL_CONFIG`       |
| `pkg/pattern`         | `PKG_PATTERN`           |
| `pkg/walk`            | `PKG_WALK`              |
| `cmd/idd-cli`         | `CMD_IDD_CLI`           |
| `auth`                | `AUTH`                  |

**Invalid (will be flagged by idd-cli):**

- `SPEC-BE-007-007` — Triple-segment identifier
- `PATTERN-*`, `WALK-*` — Section identifiers, not IDD identifiers

## Linkage Rules

**Bidirectional Links:**

IDD documents have bidirectional links that must remain consistent:

| Link Pair | Forward Direction | Backward Direction |
|-----------|------------------|-------------------|
| Contract ↔ SPEC | SPEC: `**Contract:**` | contract.md: `**Implements:**` |
| SPEC ↔ TEST | SPEC: `**Tests:**` | testing.md: `**Spec Coverage:**` |

**Complete Tracking Chain:**

```text
contract.md (Authenticator interface)
    ↑ Implements
    |
SPEC-AUTH-001 (in spec.md)
    ├── **Design:** design.md — implements architecture AuthModule
    ├── **Contract:** contract.md — implements interface Authenticator
    └── **Tests:** TEST-AUTH-001
            ↑ Spec Coverage
            |
TEST-AUTH-001 (in testing.md)

CODE: @implement SPEC-AUTH-001    →  implements Authenticator interface
TEST: @test TEST-AUTH-001         →  validates SPEC-AUTH-001
CONTRACT TEST: @test-contract SPEC-AUTH-001  →  validates SPEC satisfies contract
```

**Linkage Verification Rules:**

*Contract ↔ SPEC:*

- Every SPEC with `**Contract:**` must have a matching `**Implements:**` in contract.md
- Every contract interface with `**Implements:**` must have a matching SPEC that declares it
- The interface name in SPEC's `implements interface <Name>` must exist in contract.md

*SPEC ↔ TEST:*

- Every SPEC's `**Tests:**` field must have corresponding `**Spec Coverage:**` in testing.md
- Every TEST's `**Spec Coverage:**` must reference an existing SPEC
- The TEST identifier in SPEC's `**Tests:**` must match the TEST's own identifier

**Frontmatter ↔ Body Sync:**

1. Every frontmatter marker must have a corresponding heading
2. Every body reference must be in frontmatter markers
3. Bidirectional consistency between `markers[].id` and headings

**Annotation Placement:**

- `@implement`, `@test`, `@test-contract` are ONLY allowed on function/type declarations
- NOT permitted inside function bodies, closures, or inline expressions

## Testing Requirements

### Coverage

Minimum coverage: 80%

| Type | What | When |
| ---- | --- | ---- |
| Unit tests | Individual functions, utilities | Always |
| Integration tests | API endpoints, database operations | Always |
| E2E tests | Critical user flows (Playwright) | Critical paths |

**Contract tests:** For every `contract.md`, there MUST exist `*_contract_test.go` in the same directory as source.

## Anti-Patterns

- **Skip planning** — "Just start coding" leads to scope creep and rework
- **Write code then retrofit docs** — Intent and plan must come first
- **Skip tests** — TDD is mandatory for all new features
- **Skip coverage verification** — 80%+ is the minimum
- **Ignore blockers** — Document and escalate
- **Delete .planning/** when done — Keep for continuity
- **Annotate inside function bodies** — Annotations only on function/type declarations
- **Triple-segment identifiers** — Use `SPEC-BE-007` not `SPEC-BE-007-007`
- **Mismatched frontmatter-body** — Breaks doc-link-consistency

## IDD CLI Tool

idd-cli validates doc-link-consistency between IDD identifiers across documentation and source code.

**Installation:**

```bash
go build -o idd-cli ./cmd/idd-cli
```

**Usage:**

```bash
idd-cli run --config .idd.yaml
```

**Configuration (`.idd.yaml`):**

```yaml
version: "1.0"

docs:
  patterns:
    - "docs/**/*.md"
  identifier_patterns:
    spec: "SPEC-[A-Z0-9_]+-[0-9]+"
    test: "TEST-[A-Z0-9_]+-[0-9]+"
    test_contract: "TEST-[A-Z0-9_]+-[0-9]+"
  ignore_paths: []          # glob patterns for doc paths to skip

code:
  patterns:
    - "**/*.go"
  annotations:
    - "@implement"
    - "@test"
    - "@test-contract"
  ignore_paths: []          # glob patterns for source files to skip

validation:
  # Link consistency between SPEC **Tests:** and TEST **Spec Coverage:**
  require_doc_link_consistency: true
  # Disallow identifiers with no connections in the graph
  allow_orphans: false
  # Every SPEC must have at least one TEST in its **Tests:** field
  require_spec_test_coverage: true
  # Every CONTRACT identifier must have a @test-contract annotation
  require_contract_test_coverage: true
  # design.md files must contain all required sections
  require_design_sections: true
  # Every identifier must appear in both docs and code annotations
  require_doc_code_correspondence: true
  # Every public function/type must have an @implement annotation
  require_public_func_annotation: true
  # Every .go file must have a package doc comment with Spec/Contract paths
  require_package_doc_comment: true
  # Every doc file must have a related_files field in frontmatter
  require_related_files: true
  # Every test function must have an @test or @test-contract annotation
  require_test_annotation: true
  # @implement/@test/@test-contract must include an identifier
  require_annotation_identifier: true
  # Multiple annotations of the same type must be comma-separated on one line
  require_annotation_on_same_line: true
  # Semantic similarity check between doc describe and code comments
  consistency_check:
    enabled: true
    threshold: 0.3          # 0.0–1.0; below this triggers a warning

output:
  file: "idd-report.json"
  include_graph: false
  verbose: false
```

**Validation Rules:**

1. **Completeness** — Every SPEC must have at least one TEST in its **Tests:** field
2. **Bidirectional (CRITICAL)** — SPEC→TEST link must match TEST→SPEC backlink
3. **Orphan Detection** — No identifiers with zero connections
4. **Identifier Format** — Must match `<TYPE>-<MODULE>-<NUMBER>` (3 segments)
5. **Heading Format** — Section headings MUST be `## <Identifier>: <description>`
6. **Frontmatter-Body Sync** — All identifiers in body must be in frontmatter, and vice versa
7. **Annotation Placement** — `@implement`, `@test`, `@test-contract` must be on function/type declarations only
8. **Package Doc Comment** — Every .go file must have a package doc comment after `package` declaration
9. **Duplicate Heading Identifier (CRITICAL)** — H2 section headings MUST NOT have duplicate identifiers within the same file. If a file contains multiple sections (e.g., `## Test` and `## Contract Test`), each section must use its own sequential numbering range (e.g., `TEST-INTERNAL_ENGINE-001` to `TEST-INTERNAL_ENGINE-007` for Test, and `TEST-INTERNAL_ENGINE-013` to `TEST-INTERNAL_ENGINE-019` for Contract Test). idd-cli will detect and report this error. Never reuse identifiers across sections.

**Output:** JSON report with `valid`, `errors`, `warnings`, `stats`, and `graph` snapshot.
