---
idd:
  version: "1.0"
  package: internal/auth
---

# Design: internal/auth

## Component: AuthModule

**Purpose:**

`AuthModule` owns the decision that a supplied credential set represents a
known identity. It gives callers one authentication boundary without exposing
password-storage details or coupling authentication policy to HTTP, CLI, or
background-job transports.

### Responsibilities

The component coordinates credential lookup and secret verification, returns a
stable identity on success, and translates expected rejection cases into one
public invalid-credentials result. It also preserves cancellation and
dependency failures so callers can distinguish a rejected login from an
operation that could not be completed.

The component owns no long-lived session state. Each call receives all request
data explicitly and returns an immutable identity value or an error.

### Boundaries and collaboration

Callers own transport validation, rate limiting, request logging, and mapping
typed failures to protocol responses. A credential store owns account lookup;
a password verifier owns secret comparison. Neither dependency is exposed
through the public `Authenticator` contract.

Persistence adapters may retain password hashes, but plaintext credentials
exist only for the duration of a call and must not be logged or stored by
`AuthModule`.

The four package documents own their respective declarations. Source files do
not repeat document paths: idd-cli discovers implementation and test
declarations through `@implement`, `@test`, and `@test-contract`, then joins
those annotations to SPEC and TEST records by identifier. Design and Contract
relationships are resolved from the names declared in these Markdown records.

### Decisions and trade-offs

Unknown users and invalid secrets deliberately share one public error. This
reduces account-discovery risk at the cost of less precise caller diagnostics;
operations can still record a private reason without changing the contract.

Dependencies are injected behind narrow interfaces. That adds small
construction overhead but keeps verification deterministic, makes storage
replaceable, and prevents the authentication package from selecting a database
or cryptographic implementation.

## Architecture

The execution path is synchronous and request-scoped:

```text
transport
   │ validated request shape
   ▼
Authenticator.Authenticate
   ├── CredentialStore.Lookup
   ├── PasswordVerifier.Verify
   └── Identity or typed failure
   ▼
transport-specific response mapping
```

The component does not retry dependency failures because it lacks
transport-specific deadlines and idempotency policy. Cancellation flows from
the caller to the credential store and must stop work promptly.

## Package Layout

```text
internal/auth/
├── contract.go       # Authenticator and public value/error types
├── service.go        # orchestration implementation
└── service_test.go   # behavior-focused tests
```

Database adapters remain with persistence code and transports remain with their
delivery packages. Keeping those directions one-way prevents authentication
policy from importing infrastructure or protocol concerns.

## Function Composition

Construction supplies a credential store and password verifier. At runtime,
`Authenticate` validates the credential value, performs lookup, verifies the
secret, and constructs the returned identity. No global initialization or
mutable singleton participates in the call path.

The implementation annotation is attached to the concrete declaration whose
behavior realizes `SPEC-INTERNAL_AUTH-001`. It is traceability metadata for
idd-cli, not a replacement for the boundary and rationale recorded here.

## Dependencies

- `CredentialStore` resolves an account record without exposing persistence
  representation to callers.
- `PasswordVerifier` compares a supplied secret with a stored hash.
- `context.Context` carries cancellation and deadlines across the lookup
  boundary.

Both injected dependencies must be safe for the concurrency model used by the
calling service.

## Testability Hooks

Behavior tests use deterministic fake stores. Contract tests run the same
boundary cases against each production implementation.

The verifier is replaced with a deterministic fake so tests can cover accepted
and rejected secrets without expensive hashing. Failure-injection stores expose
cancellation and dependency-error paths. Tests assert public error identity and
returned values rather than internal call order.

## Evolution constraints

Token issuance, session persistence, authorization, and rate limiting require
separate components and SPECs. Adding them directly to `AuthModule` would blur
the credential-verification boundary described here.
