---

## 1) Domain objects (what each one is and why it exists)

**Adapter / DTO**
- **`adapter.BuildDTO(*http.Request) -> dto.ClientAuthDTO`**
  *Purpose*: Extract what the client **presented** (headers, form fields, TLS peer cert) and normalize to an immutable payload for the domain layer.
  *No policy decisions here.*

- **`dto.ClientAuthDTO`**
  Carries:
  - `Method` (what appears to be used: `basic`, `post`, `private_key_jwt`, `tls_client_auth`, `none`)
  - `ClientID` (if present explicitly; may be blank for JWT-only until inspected)
  - `ClientSecret`, `ClientAssertion`, TLS cert data, `TokenEndpointURI`, etc.

---

**Policy & Client Registration**
- **`store.ClientMeta`**
  *Purpose*: The client’s registration + policy:
  - `ID`
  - `AllowedAuthMethods []ClientAuthMethod` (policy per client)
  - `SecretHash` (hash only for secret methods)
  - `JWKSURI`/`JWKS` (for `private_key_jwt`)
  - `TLSAuthSubjectDN`, `TLSSANs` (for mTLS)
  - `Enabled`

- **`store.ClientStore`**
  *Purpose*: Repository to fetch `ClientMeta` by ID.

---

**Strategies (auth methods)**
- **`strategy.Authenticator` interface**
  *Purpose*: Validate **one method**. No cross-method policy.
  ```go
  type Authenticator interface {
      Method() dto.ClientAuthMethod
      Authenticate(ctx context.Context, d *dto.ClientAuthDTO, meta *store.ClientMeta)
          (*authctx.ClientPrincipal, *errors.OAuth2Error)
  }
  ```

  Examples:
  - `basic.Authenticator` (client_secret_basic)
  - `post.Authenticator` (client_secret_post)
  - `privatejwt.Authenticator` (private_key_jwt)
  - `mtls.Authenticator` (tls_client_auth/self_signed)
  - `none.Authenticator` (public clients)

  Optional helper:
  - **`ClientIDInspector`** (implemented by `private_key_jwt`)
    Allows the registry to **peek client_id** from unverified assertion **only to lookup metadata**. True validation still happens later.

---

**Registry (orchestrator)**
- **`registry.Registry`**
  *Purpose*: Central coordinator:
  1) Determine/resolve `client_id` for metadata lookup
  2) Load `ClientMeta`
  3) Enforce per-client **AllowedAuthMethods**
  4) Select and call the **Strategy**
  5) Return a **ClientPrincipal** or an OAuth2-aligned error

---

**Identity & Context**
- **`authctx.ClientPrincipal`**
  *Purpose*: Authenticated client identity (`ID`, `Method`, `AuthenticatedAt`, optional `KeyID`/`CertSubject`/`AMR`).
- Helpers to attach/extract principal from `context.Context`.

---

**Errors**
- **`errors.OAuth2Error`** with codes like `invalid_client`, `unauthorized_client`, `invalid_request`, `unsupported_client_authentication`.
  *Purpose*: Proper error contract for token endpoint responses.

---

**Crypto helpers (infrastructure)**
- **`crypto.SecretVerifier`** (bcrypt/argon2id wrappers)
- **`crypto.JWKResolver`** (JWKS URI/static with caching)
- **`crypto.JWTValidator`** (JWS signature; `aud`/`iss`/`sub`/`exp`/`nbf`/`iat`)

> These are **technical services** used by strategies; they don’t own OAuth policy.

---

## 2) Static relationships (compile-time)

```
+-----------------------+         +----------------------+
|      adapter          |         |         dto          |
|  BuildDTO(*http.Request) -----> | ClientAuthDTO        |
+-----------------------+         +----------------------+
                                           |
                                           v
                                  +----------------------+
                                  |      registry        |
                                  |  (orchestrator)      |
                                  +----+-----------+-----+
                                       |           |
                                       |           v
                                       |    +----------------------+
                                       |    |   store.ClientStore  |
                                       |    | -> ClientMeta        |
                                       |    +----------------------+
                                       |
                                       v
                            +--------------------------+
                            | strategy.Authenticator   |
                            |  (method-specific)       |
                            +-----------+--------------+
                                        |
                                        v
                              +-----------------------+
                              |       crypto          |
                              | (secret/jwk/jwt/mtls) |
                              +-----------------------+
```

---

## 3) Runtime sequence (happy path + where each decision happens)

**Scenario: `private_key_jwt`**

1. **Adapter** parses HTTP → `ClientAuthDTO{ Method=private_key_jwt, ClientAssertion, TokenEndpointURI, ClientID? }`
2. **Registry** is called with DTO.
3. **Registry** selects the **strategy** by `DTO.Method`.
4. **Registry** resolves `client_id`:
   - If DTO has `ClientID`, use it.
   - If blank and strategy implements `ClientIDInspector`, call `InspectClientID` to extract `iss/sub` **unverified** for **lookup only**.
   - If still blank → `invalid_request`.
5. **Registry** loads `ClientMeta` from **ClientStore** by `client_id`; ensures `Enabled`.
6. **Registry Policy**: check `AllowedAuthMethods` contains `DTO.Method`. If not → `unauthorized_client`.
7. **Registry** calls the chosen **strategy** with `(ctx, DTO, ClientMeta)`.
8. **Strategy** performs **method-specific verification**:
   - `private_key_jwt`: resolve JWKS, verify signature/claims; ensure `iss == sub == client_id`; ensure `aud` includes `TokenEndpointURI`.
   - On success return `ClientPrincipal`.
9. **Registry** returns principal to middleware; middleware attaches to `context.Context` and passes to token handler.

**Failure points** and mapped errors:
- Bad/missing request parts → `invalid_request` (adapter or strategy)
- Unknown/disabled client → `invalid_client` (registry)
- Method not allowed for that client → `unauthorized_client` (registry)
- Bad credentials → `invalid_client` (strategy)
- Server doesn’t support presented method → `unsupported_client_authentication` (registry)

---

## 4) How the **registry** connects and what it guarantees

**Inputs:** `ctx`, `dto.ClientAuthDTO`
**Collaborators:**
- `store.ClientStore` (to get `ClientMeta`)
- `strategy.Authenticator` (chosen by `DTO.Method`)
- Optional `ClientIDInspector` (to derive lookup key for JWT)

**Guarantees and responsibilities:**
- **Selects the strategy** purely by the method **presented** in the DTO (not by guesswork).
- **Resolves client_id** before policy enforcement; for JWT it may temporarily inspect unverified claims for **lookup only**.
- **Fetches `ClientMeta`**, validates `Enabled`.
- **Enforces per-client policy** via `AllowedAuthMethods`.
- **Delegates credential verification** to the **strategy**.
- **Maps outcomes** to RFC-compliant OAuth2 errors.
- **Never logs secrets**; **never** exposes reason details that would leak verification specifics beyond OAuth2 semantics.

**Non-responsibilities (intentionally):**
- Parsing HTTP (adapter does that).
- Method-specific credential checks (strategy does that).
- Cryptographic verification (crypto helpers do that).

---

## 5) Optional extension points that the registry can use

- **`ClientIDInspector`**: enable client_id resolution for `private_key_jwt` without trusting claims.
- **`MethodSelector`**: if multiple methods are presented (e.g., `Basic` + `client_assertion`), decide which to honor. Many deployments choose **strictness**: reject if multiple are presented.
- **`PolicyDecider`**: if you want policies beyond “allowed methods” (e.g., time-of-day, IP range, or per-env), add a hook after retrieving `ClientMeta`.

---

## 6) Example: registry glue (short and opinionated)

```go
type ClientIDInspector interface {
    InspectClientID(d *dto.ClientAuthDTO) (string, bool)
}

type Registry struct {
    store  store.ClientStore
    method map[dto.ClientAuthMethod]strategy.Authenticator
}

func NewRegistry(cs store.ClientStore, authenticators ...strategy.Authenticator) *Registry {
    m := make(map[dto.ClientAuthMethod]strategy.Authenticator, len(authenticators))
    for _, a := range authenticators {
        m[a.Method()] = a
    }
    return &Registry{store: cs, method: m}
}

func (r *Registry) Authenticate(ctx context.Context, d *dto.ClientAuthDTO) (*authctx.ClientPrincipal, *errors.OAuth2Error) {
    // 1) Choose strategy by presented method
    a, ok := r.method[d.Method]
    if !ok {
        return nil, errors.Unsupported("auth method not supported by server")
    }

    // 2) Resolve client_id for metadata lookup
    clientID := d.ClientID
    if clientID == "" {
        if insp, ok := a.(ClientIDInspector); ok {
            if id, ok := insp.InspectClientID(d); ok {
                clientID = id // unverified, lookup only
            }
        }
        if clientID == "" {
            return nil, errors.InvalidRequest("missing client_id")
        }
    }

    // 3) Load client metadata
    meta, err := r.store.GetByID(clientID)
    if err != nil || meta == nil || !meta.Enabled {
        return nil, errors.InvalidClient("unknown or disabled client")
    }

    // 4) Enforce per-client policy
    if !containsMethod(meta.AllowedAuthMethods, d.Method) {
        return nil, errors.UnauthorizedClient("client not allowed to use presented auth method")
    }

    // 5) Delegate method-specific verification
    pr, e := a.Authenticate(ctx, d, meta)
    if e != nil {
        return nil, e
    }
    if pr == nil {
        return nil, errors.InvalidClient("authentication failed")
    }
    if pr.AuthenticatedAt.IsZero() {
        pr.AuthenticatedAt = time.Now()
    }
    return pr, nil
}

func containsMethod(list []dto.ClientAuthMethod, m dto.ClientAuthMethod) bool {
    for _, v := range list {
        if v == m {
            return true
        }
    }
    return false
}
```

---

## 7) “When” the registry does what (timeline)

1. **Immediately** after adapter produces DTO → registry is invoked.
2. **Before** any cryptographic verification → registry resolves `client_id` & loads `ClientMeta` (JWT may require unverified inspection to get the lookup key).
3. **Before** calling strategy → registry enforces **allowed methods** policy.
4. **During** strategy execution → cryptographic checks happen (secret compare, JWT signature/claims, mTLS subject/SAN checks).
5. **After** success → registry timestamps `ClientPrincipal` and hands it off via context to endpoint handler.

---

## 8) Test seams (where to unit test each concern)

- **Adapter tests**: Given various HTTP requests, assert `ClientAuthDTO` correctness (method detection, parsing).
- **Strategy tests** (no HTTP): Provide DTO + fake `ClientMeta` and stub crypto to test method verification.
- **Registry tests**: Mock `ClientStore`, plug in a fake strategy, verify:
  - Unknown client → `invalid_client`
  - Disallowed method → `unauthorized_client`
  - Strategy failure → `invalid_client` (or appropriate)
  - `ClientIDInspector` path for JWT
- **Crypto tests**: Unit test secret verification and JWT checks with known vectors.

---

## 9) Observability and safety

- **Do not log** `client_secret` or raw JWTs. Redact if needed.
- **Audit events**: Registry is a good choke point to emit: (client_id, method, result, elapsed) without secrets.
- **Rate limits**: Apply per client_id (and per IP as backstop) around the registry or strategy failure path.
- **mTLS + proxies**: If you trust `XFCC`, lock it down to trusted hops and validate formatting.

---

## 10) TL;DR on “how the registry connects”

- The **registry** receives the **DTO** (from adapter), picks the **strategy** (based on `Method`), gets **ClientMeta** from **ClientStore**, enforces **per-client policy**, then invokes the **strategy** which uses **crypto** to verify the presented credentials, and finally returns a **ClientPrincipal** (or an OAuth2 error).
- Each piece has a clear responsibility, and the registry is the **only place** that ties policy and method selection together.

---

If you want, I can generate a **PlantUML sequence diagram** or a **Mermaid diagram** you can paste into docs, and/or wire this into a minimal `main.go` showing middleware + token handler.

Do you want **strict rejection** when multiple auth presentations appear, or should we add a **MethodSelector** to prefer, say, `private_key_jwt > mtls > basic > post > none`?


The **registry** is an **application-layer orchestrator**—a coordination point that ties together the other domain objects without owning their internal logic. Think of it as the **policy enforcement and dispatch hub** for client authentication.

---

### **What kind of idea is this?**
It’s a **Domain Service** in **DDD terms** (Domain-Driven Design) or an **Application Service** in Clean Architecture. Its job is to:

- **Centralize orchestration**: Instead of scattering logic across HTTP handlers or strategies, the registry ensures a single place for:
  - Selecting the correct **Authenticator Strategy** (based on the method presented in the request).
  - Enforcing **per-client policy** (AllowedAuthMethods).
  - Handling **error mapping** to OAuth2 semantics.
- **Keep strategies pure**: Each strategy only knows how to validate its method, not how to pick itself or enforce global policy.
- **Decouple HTTP from domain**: The registry works on DTOs, not raw `http.Request`.

---

### **How it connects to other types**
Here’s the conceptual wiring:

```
HTTP Request
   │
   ▼
Adapter (HTTP → DTO)
   │
   ▼
Registry.Authenticate(ctx, DTO)
   │
   ├─> ClientStore.GetByID(client_id) → ClientMeta
   │
   ├─> Policy check: AllowedAuthMethods
   │
   └─> Strategy.Authenticate(ctx, DTO, ClientMeta)
           │
           └─> Crypto helpers (secret hash, JWT, JWKS, mTLS cert)
   │
   ▼
ClientPrincipal (identity) or OAuth2Error
```

---

### **Why not just call strategies directly?**
Because:
- You need **per-client policy enforcement** (AllowedAuthMethods).
- You need **method selection logic** (and possibly fallback or strictness rules).
- You need **consistent error handling** (RFC6749-compliant).
- You want to keep **HTTP parsing** and **crypto details** out of your core orchestration logic.

---

### **Analogy**
Think of the registry as the **dispatcher in a payment system**:
- The adapter says: “This looks like a credit card payment.”
- The registry says: “Is this merchant allowed to use credit cards? Yes → call the credit card processor.”
- The processor (strategy) says: “I’ll validate the card number and CVV.”
- The registry then returns: “Payment accepted” or “Declined.”

---

✅ This pattern is common in **pluggable authentication systems** (e.g., OAuth2 servers, SSO gateways, payment gateways) because it gives you:
- **Extensibility**: Add new strategies without touching the registry logic.
- **Policy control**: Different clients, different allowed methods.
- **Testability**: Strategies and registry can be tested in isolation.

---

Would you like me to **draw a diagram** (Mermaid or PlantUML) showing the registry in the middle and all its collaborators (adapter, DTO, strategies, store, crypto)? Or should I **write a short Go snippet** showing how the registry is wired and called in `main.go`?

-----------------------------------


