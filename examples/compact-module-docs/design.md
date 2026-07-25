---
idd:
  version: "1.0"
  package: internal/auth
  document: design
---

# Design: internal/auth

## Component: AuthModule

`AuthModule` owns credential verification while keeping transport and
persistence details outside the module boundary.

## Architecture

`AuthModule` separates credential verification from transport concerns. Callers
depend on the `Authenticator` boundary and do not access password storage.

## Package Layout

The package owns authentication types and the `Authenticator` implementation.
HTTP handlers and persistence adapters remain in their respective packages.

## Function Composition

The handler validates request shape, calls `Authenticator.Authenticate`, and
maps the returned identity or typed failure to a transport response.

## Dependencies

The module depends on a credential store and a password verifier supplied
through private implementation fields.

## Testability Hooks

Behavior tests use deterministic fake stores. Contract tests run the same
boundary cases against each production implementation.
