Great follow‑up, Mårten. For a **token service** (the whole set of OAuth2/OIDC endpoints), you typically want **multiple registries**, each responsible for a different axis of pluggability. The “client-auth registry” we discussed is just one. Below I’ll map the common endpoints and which registries they use, then show a clean way to wire them together.

---

## 1) Token service = multiple endpoints, multiple registries

Typical endpoints:

- **/oauth/token** (Token)
- **/oauth/introspect** (Introspection)
- **/oauth/revoke** (Revocation)
- **/oauth/device_authorization** (Device Authorization)
- **/oauth/par** (Pushed Authorization Requests)
- **/.well-known/jwks.json** (JWKS)
- **/userinfo** (OIDC UserInfo)

Each endpoint needs a different set of pluggable concerns. Think of “registries” as **dispatch maps** from a discriminant (method/format/grant_type) to a strategy implementation.

### Core registries you’ll likely need

1) **Client Authentication Registry** (we already designed it)
   - Key: **client authentication method** (basic, post, private_key_jwt, mtls, none)
   - Used by: token, introspection, revocation, PAR, device_authorization (public clients allowed)
   - Per‑endpoint policy: same registry, **different policy** knobs per endpoint.

2) **Grant Handler Registry**
   - Key: **grant_type** (`authorization_code`, `refresh_token`, `client_credentials`, `urn:ietf:params:oauth:grant-type:device_code`, `urn:ietf:params:oauth:grant-type:token-exchange`, etc.)
   - Used by: **/token** only.
   - Each handler validates grant-specific inputs, enforces scope/audience policy, and emits token issuance requests.

3) **Access Token Issuer Registry**
   - Key: **token format** or **issuer profile** (e.g., `jwt`, `opaque/db-reference`, `paseto`)
   - Used by: token (and sometimes device).
   - Selects signing key, algorithm, lifetimes, claims enrichment.

4) **ID Token Issuer / Claims Enricher Registry** (OIDC)
   - Key: **claims profile** (standard OIDC, enterprise custom claims, pairwise vs public subject, etc.)
   - Used by: token (when `openid` scope), userinfo.

5) **Access Token Verifier Registry**
   - Key: **token type/format** (JWT vs opaque reference), **proof-of-possession** binding (DPoP/mTLS), **issuer**
   - Used by: userinfo, introspection (if you accept Bearer to introspect), internal resource middlewares.

6) **Introspection Caller Authentication Registry**
   - Key: same as client auth methods, but **policy differs**: resource servers calling introspection often use basic/private_jwt/mtls or even Bearer PATs.
   - Used by: introspection endpoint (authenticate the *caller*, not the token subject).

7) **Revocation Caller Authentication Registry**
   - Often the same as introspection’s caller auth, but configurable separately.

8) **Request Object / PAR Validator Registry**
   - Key: **request object format** (`jwt` JAR vs plain), **signature alg**, **key source**
   - Used by: PAR (and Authorization Endpoint if you host it).

9) **DPoP Proof Validator Registry** (if you add DPoP)
   - Key: **JWS alg**, **nonce policy**, **clock skew profile**
   - Used by: token (to bind tokens), resource endpoints (to verify proof), introspection (optional).

10) **Scope / Audience Policy Registry**
   - Key: **client**, **grant_type**, **resource/audience**
   - Used by: token, device; optionally introspection to redact responses.

11) **Key Resolver Registry**
   - Key: **kid / issuer / client** → **JWK** (local KMS, JWKS URI, rotation)
   - Used by: token issuer, private_key_jwt auth, request object verification, token verification.

You don’t need all on day one; start with (1)–(4), add others as you enable features.

---

## 2) Which registries each endpoint uses

**/oauth/token**
- **ClientAuthRegistry**: authenticate the client (method and per‑client policy).
- **GrantRegistry**: pick handler by `grant_type`.
- **Scope/AudiencePolicy**: validate & compute granted scopes.
- **TokenIssuerRegistry**: mint access tokens (JWT or opaque).
- **IDTokenIssuer/Claims**: if OIDC (`openid` scope).
- **(Optional) DPoP Proof Validator**: if proof-of-possession.

**/oauth/introspect**
- **CallerAuthRegistry**: authenticate the **caller** (resource server). Can reuse ClientAuthRegistry with different endpoint policy.
- **AccessTokenVerifierRegistry**: parse/verify JWT or resolve opaque (DB).
- **IntrospectionPolicy**: mask/redact claims based on caller entitlement.

**/oauth/revoke**
- **CallerAuthRegistry**: authenticate the **caller**.
- **TokenLocator**: find and revoke (JWT blacklist or DB reference delete).
- **RevocationPolicy**: e.g., restrict who can revoke what.

**/oauth/device_authorization**
- **ClientAuthRegistry**: likely allows `none` for public clients.
- **DeviceFlowIssuer**: create device_code/user_code + polling interval.
- **ScopePolicy**.

**/oauth/par**
- **ClientAuthRegistry**
- **RequestObjectValidator**: validate signed `request`/`request_uri`.

**/.well-known/jwks.json**
- **KeyResolverRegistry**: surface the keys the TokenIssuerRegistry uses.

**/userinfo**
- **AccessTokenVerifierRegistry** (and binding checks)
- **ClaimsEnricher** (subset of ID token claims, per spec and scopes)

---

## 3) Wiring shape: one “Service Kit” with endpoint-specific policy

Create an **application layer** that bundles registries and endpoint policy. Each endpoint picks which registries to use and with what policy.

```go
// Axis registries
type ServiceKit struct {
    ClientAuth        *authn.Registry            // methods: basic/post/jwt/mtls/none
    GrantHandlers     *grant.Registry            // grant_type -> handler
    TokenIssuers      *token.Registry            // "jwt" / "opaque" / profiles
    IDTokenIssuers    *idtoken.Registry          // OIDC id token profiles
    TokenVerifiers    *verify.Registry           // format/binding -> verifier
    RequestValidators *requestobj.Registry       // JAR/JWT request object
    ScopePolicy       *policy.ScopeRegistry      // rules for scopes/audience
    KeyResolver       *keys.Registry             // kid/issuer/client -> key
    // Optional
    DPoPValidators    *dpop.Registry
    IntrospectionAuth *authn.Registry            // could reuse with different policy
    RevocationAuth    *authn.Registry
}

// Per-endpoint policy knobs
type EndpointPolicy struct {
    // which client auth methods are allowed at each endpoint (per client will further restrict)
    AllowedClientMethods map[string][]dto.ClientAuthMethod // endpoint->allowed
    // defaults, lifetimes, alg preferences, max scopes per grant_type, etc.
}
```

**Invoking an endpoint** (example: /token):

```go
func (s *TokenService) HandleToken(w http.ResponseWriter, r *http.Request) {
    // 1) HTTP → DTO
    dto, err := adapter.BuildDTO(r, s.AdapterCfg)
    if err != nil { writeOAuthError(w, errors.InvalidRequest(err.Error())); return }

    // 2) Endpoint-level method allowlist: quick rejection before store fetch (optional)
    if !s.EndpointPolicy.AllowAt("token", dto.Method) {
        writeOAuthError(w, errors.Unsupported("auth method not supported at token endpoint"))
        return
    }

    // 3) Client authentication (per-client policy enforced inside registry)
    client, e := s.Kit.ClientAuth.Authenticate(r.Context(), dto)
    if e != nil { writeOAuthError(w, e); return }

    // 4) Grant handling
    gt := r.PostFormValue("grant_type")
    gh := s.Kit.GrantHandlers.Get(gt)
    if gh == nil {
        writeOAuthError(w, errors.InvalidRequest("unsupported grant_type"))
        return
    }

    // 5) Apply scope/audience policy
    reqScopes := parseScopes(r.PostFormValue("scope"))
    granted, perr := s.Kit.ScopePolicy.Evaluate(client, gt, reqScopes, r.Context())
    if perr != nil { writeOAuthError(w, perr); return }

    // 6) Issue tokens
    issueReq := grant.IssueRequest{Client: client, GrantType: gt, Scopes: granted, Request: r}
    at, idt, terr := gh.Execute(r.Context(), issueReq, s.Kit.TokenIssuers, s.Kit.IDTokenIssuers)
    if terr != nil { writeOAuthError(w, terr); return }

    // 7) Respond
    writeTokenResponse(w, at, idt)
}
```

> Note: The **client auth registry** enforces **per‑client policy** via `ClientMeta.AllowedAuthMethods`. The **endpoint policy** is a coarse global allowlist that lets you differ across endpoints (e.g., allow `none` only on **device_authorization**).

---

## 4) Making per‑endpoint policy per‑client

You can keep your `ClientMeta` shape and add an optional **endpoint scoping**:

```go
type ClientMeta struct {
    ID string
    Enabled bool
    // Global allowed methods (fallback)
    AllowedAuthMethods []dto.ClientAuthMethod
    // Optional: per-endpoint overrides
    AllowedByEndpoint map[string][]dto.ClientAuthMethod // e.g., "token": ["private_key_jwt"], "introspect": ["mtls"]
}
```

Then in `registry.Authenticate` (or a policy helper used by it):

```go
func isAllowed(meta *store.ClientMeta, endpoint string, method dto.ClientAuthMethod) bool {
    if v, ok := meta.AllowedByEndpoint[endpoint]; ok {
        return contains(v, method)
    }
    return contains(meta.AllowedAuthMethods, method)
}
```

This gives you **both** per‑client and per‑endpoint policy, which is very common:

- At `/token`, a client may use `private_key_jwt`.
- At `/introspect`, the same client must use `mtls`.

---

## 5) Example: registries per endpoint

```go
type TokenService struct {
    Kit           *ServiceKit
    AdapterCfg    adapter.Config
    EndpointPolicy EndpointPolicy
}
type IntrospectionService struct {
    Kit           *ServiceKit
    AdapterCfg    adapter.Config
    EndpointPolicy EndpointPolicy
}
type RevocationService struct { /* similar */ }
type DeviceService struct     { /* similar, allows 'none' */ }
type PARService struct        { /* similar with RequestValidators */ }
type UserInfoService struct   { /* uses TokenVerifiers + ClaimsEnricher */ }
```

You can supply **different auth registries** per endpoint if you want **completely different method lineups** (e.g., `RevocationAuth`), or **reuse one** with a different endpoint allowlist.

---

## 6) Optional registries you’ll grow into

- **Token Exchange (RFC 8693) Registry**
  - Key by **subject_token_type** / **actor_token_type**; different validators for JWT, SAML, opaque.

- **Audience Resolver Registry**
  - Key by **resource indicator** → **audience/claims** mapping.

- **Consent/Approval Registry**
  - Key by **client** / **scope** / **user** (server-side approval repository).

- **Throttle/Anti‑abuse Registry**
  - Keyed by **client_id** / **grant_type** / **endpoint**, to enforce rate limits.

---

## 7) Minimal types for the new registries (sketches)

**Grant handlers**
```go
package grant

type Handler interface {
    Type() string // "authorization_code", "client_credentials", ...
    Execute(ctx context.Context, req IssueRequest, issuers *token.Registry, idIssuers *idtoken.Registry) (AccessToken, *IDToken, *errors.OAuth2Error)
}

type Registry struct{ m map[string]Handler }
func (r *Registry) Get(t string) Handler { return r.m[t] }
```

**Token issuer**
```go
package token

type Issuer interface {
    Profile() string // "jwt", "opaque"
    Issue(ctx context.Context, in IssueInput) (AccessToken, *errors.OAuth2Error)
}

type Registry struct{ m map[string]Issuer }
func (r *Registry) Get(p string) Issuer { return r.m[p] }
```

**Token verifier**
```go
package verify

type Verifier interface {
    Kind() string // "jwt", "opaque"
    Verify(ctx context.Context, token string) (Claims, *errors.OAuth2Error)
}

type Registry struct{ m map[string]Verifier }
func (r *Registry) ForToken(token string) Verifier {
    // sniff/parse to choose "jwt" vs "opaque"
}
```

This mirrors the **client-auth registry** you already have: pick by discriminant, run a strategy.

---

## 8) Practical defaults & policies

- **Endpoint allowlist**:
  - `/token`: allow `basic`, `private_key_jwt`, `mtls`; **disable `post` by default**.
  - `/introspect`: allow `mtls` and `private_key_jwt`; disallow `none`.
  - `/device_authorization`: allow `none` for registered public clients.
  - `/par`: allow `private_key_jwt`, `mtls`.

- **Grant handlers in v1**: `authorization_code`, `refresh_token`, `client_credentials`, `device_code`.

- **Token issuer in v1**: `jwt` access tokens (easier to verify everywhere). Add `opaque` later if you need revocation-by-reference.

- **Verifier in v1**: JWT verifier with your signing keys; add an opaque verifier when you add DB-backed opaque tokens.

---

## 9) How this stays testable

- Each registry is a tiny map + interface: easy to mock and unit-test.
- Endpoints are just **pipelines** that pull from registries and apply policy checks in order:
  1) Adapter → DTO
  2) Endpoint allowlist (fast fail)
  3) ClientAuthRegistry
  4) GrantRegistry
  5) ScopePolicy
  6) Issuer(s) → response

Write table-driven tests per endpoint with fake registries/strategies.

---

## 10) TL;DR

- There isn’t only one “registry.” For a production token service, use **multiple registries**, each handling a pluggable dimension: **client auth**, **grant type**, **token issuing**, **verification**, **claims enrichment**, **request object validation**, etc.
- **Per‑endpoint policy** determines *which* of those are allowed/used for that endpoint.
- **Per‑client policy** (in `ClientMeta`) further restricts allowed methods for each client.
- The handler simply composes: **Adapter → [Endpoint allowlist] → ClientAuthRegistry → GrantRegistry → ScopePolicy → TokenIssuer(s)**.

---

If you want, I can:
- Generate a **Mermaid diagram** of the whole pipeline per endpoint.
- Provide concrete **Go skeletons** for `grant.Registry`, `token.Registry`, and a `TokenService` that wires them.
- Show how to express **per‑endpoint** and **per‑client** policies via a YAML or JSON config and load them into the registries.
