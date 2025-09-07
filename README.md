
-> /authorize: "I need a logged-in subject before I can issue an authorization code"

-> Device Code Flow: "I need to authenticate the user before making the device code as approved"

-> Session Management / Web UI: "The web login handler or UI layer that starts the authentication process"




[Grant Flow: /authorize]
    → calls Authentication Service (upstream)
        → Authentication Service calls Identity Sources (downstream)
        → Authentication Service applies MFA policy from Configuration (downstream)
    ← returns AuthResult (SubjectID, ACR, AMR, auth_time, claims)
    → Grant Flow passes SubjectID to Delegation Service (downstream consumer)
    → Delegation approves → Token Issuance builds tokens



Delegation is the core domain of authorization:

- Your grant flows become application services that ask the delegation domain for a decision.

Given this incoming request, using this grant type, can i issue the correct tokens to this client
and so, how ?

- Your token service becomes a consumer of delegation decisions, not the place where
  those decisions are made.
- You can evolve consent logic, delegation expiry rules, or scope policies withouth touching
  token issuance or identity integration.


## Grant Flow Context
Is the flow engine of your AS

When a request hits /authorize or /token, the Grant Flows context is responsible for:

1. Identity the grant type / flow
2. Validating the request against the rules for that flow
3. Interacting with other bounded contexts
- Delegation / Consent: Is the client allowed to get these scopes for this subject ?
- Identity Sources: Who is the subject ? What authentication performed ?
- Token Issuance: Mint me an access token (and maybe an ID token) with these claims.
4. Producing the correct output for the flow
- Access token, refresh token, ID token (if OIDC)
- Error response if validation fails

--------------------
1. Grant Flow orchestrates
1.1 Validates the incoming request
1.2 Calls the Delegation Service

2. Delegation Service
2.1 Yes -> returns the DelegationID and the effective rights (scopes, audiences, expiry, constraints)
2.2 No -> returns an error (no consent, expired, revoked, insufficient scope)

3. Token Service consumes that decision:
- Uses the DelegationID and righs to build the token payload:
    - sub = SubjectID
    - aud = allowed audiences
    - scp = granted scopes
    - exp = min(Delegation expiry, token policy TTL)
    - Any claim mappings from Identity Sources
- Signs and returns the token

------------------------------
Delegation: Owns "can they?"
- Business rules, consent records, policy enforcment
Token Service: How do we represent that in a token ?
Format, signing, cryptography, claim embedding

[Authorize Handler]
    ↓ parse request
    ↓ check session
    ├──> [Authentication Service] → [Identity Sources Domain]
    ↓ have SubjectID
    ├──> [Delegation Service] → check consent
    ↓ if needed, render consent UI
    ↓ on approval, record consent
    ↓ issue authorization code

[Token Handler]
    ↓ parse request
    ├──> [Grant Flow Service]
           ├──> validate grant
           ├──> [Delegation Service] → confirm rights
           ├──> [Token Service] → build & sign tokens
    ↓ return token response








1. Client hits /authorize -> handlers/authorize.go
- Parses request into dto.AuthorizeRequest.
- Validates client via clientauth methods.
- Delegates to services/authorize_service.go and authflow/code.go

2. Authorization Code Issued -> stores in store/code_store.go

3. Client hits /token -> handlers/token.go
- Parses dto.TokenRequest
- Authenticates the client (clientauth)
- Chooses correct grant handler from services/granthandlers/
- Validates via services/grantvalidators/
- Issues tokens via services/token_service.go and security/jwt_signer.go

4. Tokens returned -> dto.TokenResponse

Next Steps

- Add middleware for logging, metrics, and request tracing

# oidcsim
