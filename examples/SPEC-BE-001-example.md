# SPEC-BE-001 User Authentication

## Overview

This specification defines the user authentication system.

## Requirements

1. Users must be able to log in with email/password
2. Sessions must expire after 24 hours
3. Passwords must be hashed using bcrypt

## Implementation

See implementation in `internal/auth/`.

## Contracts

- CONTRACT-BE-001 (Password Hashing)
- CONTRACT-BE-002 (Session Management)

## Tests

- TEST-BE-001 (Login flow)
- TEST-BE-002 (Session expiration)

## Design

- DESIGN-BE-001 (Auth flow diagram)