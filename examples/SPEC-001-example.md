# SPEC-001 User Authentication

## Overview

This specification defines the user authentication system.

## Requirements

1. Users must be able to log in with email/password
2. Sessions must expire after 24 hours
3. Passwords must be hashed using bcrypt

## Implementation

See implementation in `internal/auth/`.

## Contracts

- CONTRACT-001 (Password Hashing)
- CONTRACT-002 (Session Management)

## Tests

- TEST-001 (Login flow)
- TEST-002 (Session expiration)

## Design

- DESIGN-001 (Auth flow diagram)