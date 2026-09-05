---
idd:
  version: "1.1"
  package: internal/auth
  namespace: INTERNAL_AUTH
---

# Contracts: internal/auth

## Contract: Authenticator

**Guarantees:**

`Authenticator` converts a complete credential attempt into either an
authenticated identity or a caller-visible failure. The contract describes
observable behavior shared by every implementation; storage layout, hash
selection, and internal diagnostic detail remain private.

```go
type Authenticator interface {
    Authenticate(ctx context.Context, credentials Credentials) (Identity, error)
}
```

### Inputs and ownership

`ctx` may carry a deadline or cancellation and must not be retained after the
call. `credentials` contains an identifier and plaintext secret owned by the
caller. Implementations may read those values for the duration of the call but
must not log, persist, or return the secret.

The identifier must be syntactically valid and the secret non-empty. Invalid
shape is rejected before dependency work.

### Outputs

On success, the returned `Identity` is complete, stable for the duration of the
caller operation, and independent of persistence-specific records. A non-zero
identity is never returned together with an error.

### Errors and side effects

Unknown identifiers and invalid secrets return the same
`ErrInvalidCredentials`. Cancellation remains discoverable through
`context.Canceled` or `context.DeadlineExceeded`. Credential-store failures are
preserved as operational errors and are not translated into credential
rejection.

Authentication may read credential data and record private operational
telemetry. It does not create sessions, issue tokens, mutate credentials, or
retry failed dependencies.

### Non-guarantees

The contract does not select a password-hash algorithm, persistence schema,
transport status code, session representation, or token format. It does not
promise that operational dependencies always succeed.

### Known limitations

Rate limiting, multi-factor challenges, session issuance, and authorization
are intentionally absent from this example boundary and require separate
Components, Contracts, and SPECs.

### Invariants and security

- Failure shape must not reveal whether an account exists.
- Plaintext secrets must not outlive the call or appear in logs and errors.
- The same input and dependency state must produce the same public outcome.
- Cancellation must stop downstream work as soon as the dependency permits.

### Compatibility commitments

New implementations may change storage or verification algorithms but must
preserve the public error categories, secret-ownership rules, and success
semantics. Adding multi-factor challenges or token issuance requires a new
contract rather than silently widening this one.

Contract evidence is declared by `TEST-INTERNAL_AUTH-002` in `testing.md`. Its
source declaration uses `@test-contract TEST-INTERNAL_AUTH-002`; the annotation
points to the TEST record, while that record names `Authenticator` through its
`Contracts` field. No source file needs to know the path of this contract
document.
