A Delegation is the relationship between:

1. Principal: the actor granting authority (often a User, but could be a Client in M2M flows)

2. Delegate: the actor receiving authority (usually a Client)

3. Scope(s): the permissions granted

4. Constraints: expiry, audience, resource restrictions, PKCE requirements etc

5. Grant Source: how the delegation was established (auth code, device code, client credentials, etc)

type Delegation struct {
    ID              string
    PrincipalID     string // user_id or service account id
    DelegateClientID string
    Scopes          []string
    Audience        []string
    IssuedAt        time.Time
    ExpiresAt       time.Time
    GrantType       string
    ConsentRecordID *string // link to stored consent if applicable
    Metadata        map[string]string
}


GET /authorize

1. After the user consents, you create a Delegation object representing that consent.
2. The AuthorizationCode then points to this delegation.

POST /token

1. When exchanging an AuthorizationCode or RefreshToken, you look up the delegation to know what scopes, audience, and expiry rules to apply.

2. The AccessToken is essentially a signed/serialized view of the delegation.

User (Principal) ──┐
                   ├── Delegation ──> AccessToken(s)
Client (Delegate) ─┘                  RefreshToken(s)

AuthorizationCode → references a Delegation

RefreshToken → references a Delegation

AccessToken → references a Delegation

---------------------
-
💡 Why Model It Explicitly
Auditability: You can answer “Who delegated what to whom, and when?”

Revocation: You can revoke a delegation and instantly invalidate all tokens tied to it.

Reusability: Multiple tokens (access + refresh) can share the same delegation record.

Extensibility: Easy to add constraints like “valid only from IP range X” or “valid only for resource server Y”.


          ┌───────────────────┐
          │       User         │  (Principal / Resource Owner)
          │  user_id           │
          │  claims[...]       │
          └─────────┬─────────┘
                    │ 1
                    │
                    │ N
          ┌─────────▼─────────┐
          │   Delegation       │  (Core contract: who, to whom, for what)
          │  delegation_id     │
          │  principal_id      │ → User.user_id
          │  client_id         │ → Client.client_id
          │  scopes[...]       │
          │  audience[...]     │
          │  issued_at         │
          │  expires_at        │
          │  grant_type        │
          │  consent_record_id │ → ConsentRecord.id (optional)
          └───────┬───────────┘
                  │ 1
                  │
                  │ N
   ┌──────────────▼──────────────┐
   │ AuthorizationCode           │
   │ code                        │
   │ redirect_uri                 │
   │ pkce_challenge               │
   │ expires_at                   │
   │ delegation_id → Delegation   │
   └──────────────┬──────────────┘
                  │
                  │ N
   ┌──────────────▼──────────────┐
   │ AccessToken                  │
   │ token                        │
   │ expires_at                   │
   │ delegation_id → Delegation   │
   └──────────────┬──────────────┘
                  │
                  │ 0..N
   ┌──────────────▼──────────────┐
   │ RefreshToken                 │
   │ token                        │
   │ expires_at                   │
   │ delegation_id → Delegation   │
   └─────────────────────────────┘

   ┌───────────────────┐
   │      Client        │  (Delegate)
   │  client_id         │
   │  redirect_uris[]   │
   │  allowed_grants[]  │
   │  scopes[]          │
   └───────────────────┘


/authorize endpoint
Parse request → AuthorizationRequest (transient object).

Validate Client and requested scopes.

Authenticate User.

Check or create ConsentRecord.

Create Delegation (User ↔ Client ↔ Scopes).

Create AuthorizationCode linked to that Delegation.

Redirect back to client with code.


/token endpoint (authorization_code grant)
Parse request → TokenRequest.

Validate Client authentication.

Look up AuthorizationCode → get Delegation.

Validate Delegation (expiry, scopes, audience).

Issue AccessToken (and RefreshToken) linked to Delegation.

Return TokenResponse.
