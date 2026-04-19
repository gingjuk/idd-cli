---
markers:
  - id: DESIGN-AUTH-001
    name: Authentication Architecture
---

# Design Decisions (auth)

## DESIGN-AUTH-001: Authentication Architecture

**Date:** 2026-04-18
**Status:** Accepted

### Context

idd-cli needs a simple authentication module to support future features like authenticated API access and user-specific configurations.

### Decision

The authentication module uses a minimal type-based design:

1. **LoginRequest/LoginResponse pattern** — Simple request-response types
2. **Token-based authentication** — JWT or session tokens
3. **Type-only implementation** — No actual auth logic yet (stub module)

### Architecture

```text
internal/auth/
├── auth.go          # Type definitions with @spec, @contract annotations
└── (future)
    ├── service.go    # Auth service implementation
    ├── middleware.go # Auth middleware
    └── token.go      # Token generation/validation
```

### Consequences

**Positive:**

- Clear contract defined for future implementation
- Type definitions provide clear API surface
- Minimal overhead for initial implementation

**Negative:**

- No actual authentication functionality yet
- Authorization module (`SPEC-AUTH-002`) not yet implemented

**Related:** `SPEC-AUTH-001`, `SPEC-AUTH-002`
