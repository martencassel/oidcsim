
The /authorize flow is a use case that coordinates multiple domain concepts:

- oauth2.AuthorizeRequest (OAuth2 domain)

- authentication.Context (Authentication domain)

- delegation.Delegation (Delegation domain)

That orchestration belongs in the application layer for the OAuth2 bounded context,
because it’s where you glue together domain objects and services to fulfill a protocol‑level use case.

Keeps domain pure The OAuth2 domain types (AuthorizeRequest, AuthorizationCode, etc.) stay free of HTTP, persistence, or cross‑domain orchestration. The application service (AuthorizeService) uses them, but doesn’t live in the domain package.
