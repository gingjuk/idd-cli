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

**Intent First** — Every implementation starts with explicit user intent, not assumptions. The workflow preserves intent through planning, architecture, implementation, and validation phases.

## Key Principle

Delegate to **specialized agents** for specific tasks. Use the right agent for the right job:
- Planning tasks → planning agent
- Architecture/consultation → architecture/consulting agent
- Code review → review agent
- Deep implementation → implementation agent

## Phases

### Phase 1: Intent Capture

Before ANY code is written, capture the user's intent in `<module>/.planning/intent.md`:

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

Use the `planning` skill to create this: `skill({ name: "planning" })`

### Phase 2: Planning

Create `<module>/.planning/plan.md` with:

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

### Phase 3: Architecture (when needed)

For complex features, consult the **architecture/consulting agent** to design architecture:

```
Invoke architecture/consulting agent when:
- New patterns or patterns unfamiliar to the codebase
- Multiple modules/systems involved
- Performance or security concerns
- Unfamiliar framework behavior
```

Architecture artifacts:
- `design.md` — System design, architecture decisions
- `spec.md` — Functionality specification
- `testing.md` — Testing strategy and approach
- `contract.md` — (15+ functions or shared lib) Signatures, types, error contracts

Location: `docs/<package_path>/`

### Phase 3.5: Documentation Linkage

Code and documentation are linked through numbered identifiers. This creates a bidirectional traceable relationship between specs, contracts, tests, and implementation.

**Identifier Format:** `<TYPE>-<MODULE>-<NUMBER>`

| Prefix | Meaning | Example |
|--------|---------|---------|
| `SPEC-` | Functionality specification | `SPEC-BE-001` |
| `CONTRACT-` | Interface/behavior contract | `CONTRACT-BE-001` |
| `TEST-` | Test case | `TEST-BE-001` |
| `DESIGN-` | Architecture design decision | `DESIGN-BE-001` |

**Pattern (Regex):** Used by IDD Link Validator for auto-detection
```yaml
identifier_patterns:
  spec: "SPEC-[A-Z]+-[0-9]+"
  contract: "CONTRACT-[A-Z]+-[0-9]+"
  test: "TEST-[A-Z]+-[0-9]+"
  design: "DESIGN-[A-Z]+-[0-9]+"

code_annotations:
  - "@spec"
  - "@contract"
  - "@test"
  - "@design"
```

**Spec File Size Rule:**
- If `spec.md` exceeds **1500 lines**, split into multiple files using the pattern `spec-<feat>.md`
- Each split file should focus on a specific feature or subdomain
- The main `spec.md` becomes an index that references all split files
- Example: `spec-auth.md`, `spec-trading.md`, `spec-portfolio.md`

**Spec Index Example (`spec.md`):**
```markdown
# Specification Index

## Modules
- [Authentication](spec-auth.md) — SPEC-BE-001 to SPEC-BE-010
- [Trading](spec-trading.md) — SPEC-BE-011 to SPEC-BE-030
- [Portfolio](spec-portfolio.md) — SPEC-BE-031 to SPEC-BE-050
```

**Module Prefixes:**

| Prefix | Module |
|--------|--------|
| `BE-` | api-server (backend) |
| `FE-` | web (frontend) |
| `E2E-` | End-to-end tests |

**Code Comment Format:**

```python
# =============================================================================
# @spec     SPEC-BE-001: Order execution must complete within 100ms
# @contract CONTRACT-BE-001: execute_order(order: Order) -> Execution
# @test     TEST-BE-001: Verify order execution meets SLA
# =============================================================================
def execute_order(self, order: Order) -> Execution:
    """Execute a market order. See CONTRACT-BE-001 for error handling."""
```

**Document Format (spec.md, contract.md, testing.md):**

```markdown
## SPEC-BE-001: Order Execution Performance

**Requirement:** Order execution must complete within 100ms for market orders.

**Implementation:** `api-server/app/services/backtest_engine/executor.py`

**Tests:** TEST-BE-001, TEST-BE-002
```

**Test Format:**

```python
# TEST-BE-001: Verify order execution meets SLA
def test_execute_order_performance():
    """SLA: 100ms p99 latency. See SPEC-BE-001."""
```

**Linkage Rules:**
- Each identifier must be unique within its type namespace
- Every spec/contract item should have at least one corresponding test
- Every code implementation should reference its spec and contract IDs
- Every test should reference the spec/contract it validates

### Document Storage Organization

**All IDD documents MUST be stored in the `docs/` directory, organized by module name.**

```
docs/
├── auth/
│   ├── spec.md           # SPEC-AUTH-001, SPEC-AUTH-002, ...
│   ├── contract.md      # CONTRACT-AUTH-001, ...
│   ├── testing.md       # TEST-AUTH-001, ...
│   └── design.md        # DESIGN-AUTH-001, ...
├── trading/
│   ├── spec.md           # SPEC-TRADING-001, ...
│   └── ...
└── ...
```

**Rules:**
- One subdirectory per module (use module name, lowercase or as appropriate)
- Each subdirectory contains the module's IDD documents (spec, contract, testing, design)
- Document filenames are lowercase (spec.md, contract.md, testing.md, etc.)
- Module prefix in identifiers must match the module name (e.g., `SPEC-AUTH-001` in `docs/auth/`)

**Frontmatter:** Each IDD markdown file MUST include YAML frontmatter listing all identifiers it contains with brief descriptions:

```yaml
---
markers:
  - id: SPEC-AUTH-001
    name: User login with email/password
  - id: SPEC-AUTH-002
    name: Password hashing interface
---
```

This enables the IDD Link Validator to quickly locate and parse identifiers without reading full file content.

**Validator Config:** Ensure `idd.yaml` patterns cover `docs/**/*.md` to auto-detect all module documents.

### Phase 4: TDD Implementation

**TDD Workflow:**
1. Write test first (RED) — test should FAIL
2. Write minimal implementation (GREEN) — test should PASS
3. Refactor (IMPROVE) — verify coverage 80%+

```
delegate_task(
  category="<implementation>",
  load_skills=["tdd-workflow"],
  prompt="Implement <feature> following TDD:
1. Write failing test (RED)
2. Write minimal implementation (GREEN)
3. Refactor (IMPROVE)
Target: 80%+ coverage"
)
```

**Troubleshoot:** check test isolation → verify mocks → fix implementation (not tests, unless tests are wrong).

### Phase 5: Contract Validation (when applicable)

For modules with `contract.md`:

```
Run contract_test to validate module interfaces match contracts.
If contract_test fails, fix implementation — not contracts.
```

### Phase 6: Review

```
delegate_task(
  category="<review>",
  prompt="Review <changes> for quality, security, maintainability"
)
```

Or invoke specific agents:
- Architecture/consulting agent — Architecture, high-level review
- Explore agent — Code pattern consistency

### Phase 7: Knowledge Capture

- Personal debugging notes → auto memory
- Team/project knowledge → existing docs structure
- If task produces relevant docs/comments, don't duplicate elsewhere

### Phase 8: Commit

**Commit format:** `<type>: <description>`

Types: `feat`, `fix`, `refactor`, `docs`, `test`, `chore`, `perf`, `ci`

## Required Documents Summary

| Document | Location | Triggered By | Contents |
|----------|----------|--------------|----------|
| `intent.md` | `<module>/.planning/` | Planning | User intent import |
| `plan.md` | `<module>/.planning/` | Planning | Development state tracking, phase breakdown, risks |
| `design.md` | `docs/<package_path>/` | Architecture Design | System design, architecture decisions |
| `spec.md` | `docs/<package_path>/` | Architecture Design | Functionality specification (split if >1500 lines → `spec-<feat>.md`) |
| `testing.md` | `docs/<package_path>/` | Architecture Design | Testing strategy and approach |
| `contract.md` | `docs/<package_path>/` | Architecture Design (15+ funcs or shared lib) | Signatures, types, error contracts |

## Testing Requirements

**Minimum coverage: 80%**

| Type | What | When |
|------|------|------|
| Unit tests | Individual functions, utilities, components | Always |
| Integration tests | API endpoints, database operations | Always |
| E2E tests | Critical user flows (Playwright) | Critical paths |

## Workflow Summary

```
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

## Anti-Patterns

- **Skip planning** — "Just start coding" leads to scope creep and rework
- **Write code then retrofit docs** — Intent and plan must come first
- **Skip tests** — TDD is mandatory for all new features
- **Skip coverage verification** — 80%+ is the minimum
- **Ignore blockers** — Document and escalate
- **Delete .planning/** when done — Keep for continuity and future reference

## IDD Link Validator Tool

The IDD Link Validator (`idd-cli`) is a CLI tool that validates bidirectional linkage consistency between IDD identifiers across documentation and source code.

### Installation

```bash
go build -o idd-cli ./cmd/validator
```

### Usage

```bash
idd-cli run --config idd.yaml
```

### Configuration (`idd.yaml`)

The validator uses patterns defined in this skill. Default configuration:

```yaml
version: "1.0"

docs:
  patterns:
    - "docs/**/*.md"
  identifier_patterns:
    spec: "SPEC-[A-Z]+-[0-9]+"
    contract: "CONTRACT-[A-Z]+-[0-9]+"
    test: "TEST-[A-Z]+-[0-9]+"
    design: "DESIGN-[A-Z]+-[0-9]+"

code:
  patterns:
    - "**/*.go"
  annotations:
    - "@spec"
    - "@contract"
    - "@test"
    - "@design"

validation:
  require_bidirectional: true
  allow_orphans: false
  require_spec_test_coverage: true

output:
  file: "idd-report.json"
  include_graph: true
  verbose: true
```

### Validation Rules

1. **Completeness** — Every SPEC must have at least one TEST link (and vice versa)
2. **Bidirectional** — If SPEC→TEST exists, TEST→SPEC backlink must also exist
3. **Orphan Detection** — No identifiers with zero connections
4. **Consistency** — Cross-reference chain consistency (CODE→SPEC→TEST)

### Output

JSON report with:
- `valid` — Boolean pass/fail status
- `errors` — Failed validation rules
- `warnings` — Consistency issues
- `stats` — Identifier counts by type
- `graph` — Linkage graph snapshot
