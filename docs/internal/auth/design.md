---
idd:
  version: "1.0"
  package: internal/auth
  document: design
---

# Design: internal/auth

## Component: AuthTypes

`AuthTypes` defines the request and response values exchanged at the
authentication boundary.

## Architecture

`AuthTypes` is a value-only module. It establishes request and response shapes
without choosing a credential store, token format implementation, service, or
transport.

## Package Layout

```text
internal/auth/
├── auth.go                  # LoginRequest and LoginResponse
└── auth_contract_test.go    # Documented value constraints
```

## Function Composition

There is no call graph or initialization order because the package contains no
functions. Consumers construct and exchange the two value types directly.

## Dependencies

The production package has no external or internal runtime dependencies.

## Testability Hooks

Plain structs are directly constructible in table-driven tests. Future behavior
should introduce its own interface and behavioral SPEC instead of adding
unstated semantics to these values.
