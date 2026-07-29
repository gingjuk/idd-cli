---
idd:
  version: "1.0"
  package: internal/auth
---

# Design: internal/auth

## Component: AuthTypes

**Purpose:**

`AuthTypes` defines the request and response values exchanged at the
authentication boundary. The current package is deliberately a value-model
stub: it documents the shape and expected meaning of login data without
claiming that authentication behavior already exists.

### Responsibilities

`LoginRequest` groups the identifier and secret supplied by a future login
caller. `LoginResponse` groups the token-like value returned after a future
successful authentication operation. Keeping these values named gives future
handlers and services a stable vocabulary and a place to attach behavioral
SPECs when validation and authentication are implemented.

### Boundaries

The structs store data only. They do not validate email syntax, reject an empty
password, generate a token, parse a token, authenticate a user, manage a
session, authorize an action, or protect sensitive values automatically.
Those expectations are currently contract documentation and test evidence for
future implementations, not methods executed by the values themselves.

Callers remain responsible for the lifetime of plaintext passwords. Code that
uses `LoginRequest` must avoid logging, persistence, or accidental copying
beyond the request lifecycle. `LoginResponse.Token` does not prescribe JWT as
the permanent protocol even though the current contract tests use a
three-segment shape as representative evidence.

### Design rationale

Plain structs keep the stub dependency-free and easy to construct across
transports. Adding validation methods now would select policy before there is a
service boundary or error contract. The trade-off is that constraints are not
self-enforcing; future behavior must introduce an explicit component and
contract instead of relying on comments alone.

## Architecture

`AuthTypes` is a value-only module. It establishes request and response shapes
without choosing a credential store, token format implementation, service, or
transport.

The intended dependency direction is:

```text
transport or service
        │ constructs/consumes
        ▼
LoginRequest / LoginResponse
```

There is no reverse dependency from the value package to HTTP, persistence,
cryptography, configuration, or the validation tool. IDD annotations on the
types provide traceability but have no runtime effect.

## Package Layout

```text
internal/auth/
├── auth.go                  # LoginRequest and LoginResponse
└── auth_contract_test.go    # Documented value constraints
```

`auth.go` contains production values only. The contract test file records
expectations separately so readers do not confuse test validation with runtime
struct behavior.

## Function Composition

There is no call graph or initialization order because the package contains no
functions. Consumers construct and exchange the two value types directly.

A future authentication service should accept a `LoginRequest`, enforce input
and credential policy, and return a `LoginResponse` or typed error. That
service must receive its own Component, Contract, SPEC, and tests; it must not
be implied by extending this value-only description.

## Dependencies

The production package has no external or internal runtime dependencies.

Standard Go string values are the only field representation. This makes the
types portable but does not provide secret zeroization, format validation, or
domain-specific type safety.

## Testability Hooks

Plain structs are directly constructible in table-driven tests. Future behavior
should introduce its own interface and behavioral SPEC instead of adding
unstated semantics to these values.

The current contract tests exercise representative valid and invalid values as
documentation evidence. Their oracle is the documented constraint, not a
method on the structs. Tests must remain explicit about that distinction so a
green suite is not mistaken for production validation.

## Evolution constraints

- Adding validation behavior requires a defined error contract and a decision
  about which layer owns normalization.
- Replacing raw strings with domain types is a compatibility change for every
  caller and must preserve sensitive-data handling.
- Token generation or verification belongs to a separate capability with
  explicit cryptographic and lifecycle boundaries.
- Removing the stub types requires migrating all annotations, SPEC coverage,
  and human-readable contract context rather than deleting them as unused
  structs.
