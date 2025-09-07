## Bounded Context

domain/identitysources

## Core Question

“How do we verify a subject’s identity and retrieve authoritative identity attributes from trusted sources?”

## Core Goal

The goal here is to model the “truth about who the subject is” and “how we prove it” — without leaking protocol, persistence, or UI concerns into the domain.

## Responsibilities

### Authenticate a subject using a specific method (password, WebAuthn, federated login, etc.).

### Retrieve identity attributes (claims) from the source.

### Abstract over multiple sources (local DB, LDAP, external IdPs, social logins).

### Normalize identity data into a consistent internal model.

### Support account linking if a subject has multiple external identities.


## Boundaries

- Upstream callers
Authentication context (needs to verify a subject), possibly account management flows.

- Downstream dependencies:
Infrastructure implementations that talk to DBs, directories, or external IdPs.


## Why this works:

Domain is pure: only contracts and core types.
Application orchestrates and applies config.
Infrastructure handles the messy details of each source.
Adding a new identity source = implement IdentityProvider + register it.
