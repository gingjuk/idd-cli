---
idd:
  version: "1.0"
  package: internal/auth
  document: testing
---

# Testing: internal/auth

## TEST-INTERNAL_AUTH-001: Authentication behavior

- **Kind:** `test`
- **Covers:** `SPEC-INTERNAL_AUTH-001`

**Purpose:** Verify accepted and rejected credential outcomes.

**Oracle:** Each scenario's returned identity and public error category match
the expected row, failures return a zero identity, and captured telemetry
contains no plaintext secret.

### Evidence and scenarios

The behavior suite proves that valid credentials return the expected identity;
unknown identifiers and incorrect secrets both return
`ErrInvalidCredentials`; malformed values avoid dependency work; and
cancellation or injected store failures remain distinguishable.

It also verifies that no failure returns a partial identity and that repeated
calls do not share mutable authentication state.

### Fixtures, isolation, and oracle

A table-driven fake store supplies known, unknown, cancelled, and failed lookup
outcomes. A deterministic verifier accepts or rejects secrets without invoking
real password hashing. The oracle is the returned identity and error category,
plus captured telemetry checked for secret leakage.

### Exclusions

Transport response mapping, rate limiting, token issuance, and production
adapter conformance are proved in their owning packages or contracts.

### Source evidence

The behavior test declaration carries the TEST identifier that owns its oracle:

```go
// @test TEST-INTERNAL_AUTH-001
func TestAuthenticateBehavior(t *testing.T) {
    // Table-driven scenarios and assertions omitted from this example.
}
```

## TEST-INTERNAL_AUTH-002: Authenticator contract

- **Kind:** `contract`
- **Covers:** `SPEC-INTERNAL_AUTH-001`
- **Contracts:** `Authenticator`

**Purpose:** Verify every `Authenticator` implementation preserves the public
boundary.

**Oracle:** Every registered implementation produces the same contract-visible
identity and error categories for success, rejection, cancellation, and
dependency failure, without retaining or revealing the caller-owned secret.

### Evidence and scenarios

The contract suite runs the same success, rejection, cancellation, and
dependency-failure cases against every registered implementation. It proves
that implementations may differ internally without changing public error
categories, identity semantics, or secret ownership.

### Fixtures, isolation, and oracle

Each implementation factory receives isolated dependency fixtures. The shared
oracle compares only contract-visible values and errors, avoiding assertions
about internal call order or storage representation.

### Exclusions

Implementation-specific performance, migration, and persistence tests remain
with their adapters.

### Source evidence

The contract-test declaration points to this TEST record, not directly to a
SPEC identifier or Markdown path:

```go
// @test-contract TEST-INTERNAL_AUTH-002
func TestAuthenticatorContract(t *testing.T) {
    // Run the shared contract suite for every registered implementation.
}
```

idd-cli resolves `TEST-INTERNAL_AUTH-002` to this record, then follows
`Covers` to the SPEC and `Contracts` to `Authenticator`. That keeps source
association machine-checkable without making code own the document graph.

## Strategy

Behavior tests give fast, focused feedback for the primary implementation.
Contract tests protect substitutability across implementations. Integration
tests with production adapters verify persistence encoding separately, while
transport tests verify protocol mapping without duplicating authentication
semantics.

Sensitive values are synthetic and never written to snapshots or failure
messages. Test cases avoid timing-based assertions; any timing-attack analysis
belongs to a dedicated security evaluation.
