# access token

It is not just a "proof of user identity" - it's a credential that represents:

**Who** the token is about (sub - the subject, usually a user or service account)
**Who** it was issued to (client_id - the OAuth2 client application)
**What** it can be used for (scope, aud, resource)
**When** it's valid (exp, nbf, iat)
**Where** it came from (iss - the issuer)

From a domain perspective, the token is the binding between:

A **resource owner** (user)
A **client application** (identified by client_id)
A set of permissions (scopes)

If you remove client_id from that binding, you lose the part of the contractual context:

"this user consented to this client having these rights"

# Why the domain needs client_id at /userinfo

When you /userinfo use case runs, the domain logic often needs to know:

- Which client is making the call
- What the client is allowed to see (per-client claim release rules)
- Whether the token was actually issued to that client

Withouth client_id in the token (or retrivable from a token store), you have to guess
or trust the caller's self-declared identity - which breaks the security model.

# How it fits into the DDD model

- Domain concept: A Delegation or Consent is always between a Subject and a Client
- Access token: A value object that encapsulates that delegation
- client_id: A core attribute of that value object - it's part of the identity of the delegation.

"This token was issued to Client X for User Y with Scopes Z."

{
  "sub": "user123",
  "client_id": "my-client",
  "scope": "openid profile email",
  "exp": 1699999999
}

- Auth0, Okta, Keycloak, ForgeRock, Curity
client_id

- Google Identity Platform, Microsoft Entra ID
azp ("authorized party")

- Sometimes set to the client_id if the token is for that client's own API
aud ("audience")


# Facts

1. A token is meaningless withouth knowing which client it was issued to.
2. Many downstream decisions (claim filtering, consent checks, auditing) depend on that.
3. The client_id is part of the contract between the Authorization Server and the Resource Server.


{
  "sub": "user123",
  "azp": "my-client",
  "scope": "openid profile email"
}


# Introspection

If the IdP issues opaque tokens (random strings), the client_id isn’t in the token itself. Instead:
The IdP stores a record in its token database keyed by the token value.
That record contains client_id, sub, scope, expiry, etc.
Resource servers call the token introspection endpoint (RFC 7662) to retrieve those details.

{
  "active": true,
  "client_id": "my-client",
  "sub": "user123",
  "scope": "openid profile email"
}

JWT: Self‑contained, no lookup needed — resource server can read client_id directly after verifying the signature.

Opaque: More control for the IdP (can revoke instantly, rotate keys without breaking tokens), but requires introspection to get client_id.

For OIDC‑compliant IdPs, the most common pattern is:

JWT access tokens with a client_id claim (or azp if following Google/Microsoft conventions).

Opaque tokens with client_id returned via introspection.
