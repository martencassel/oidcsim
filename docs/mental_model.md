------------------------------------------------------------------------------------------------------------------------------------------------------------------

HTTP GET /authorization
   ↓
Adapter → AuthorizationRequestDTO
   ↓
ClientValidatorRegistry → validate client_id, redirect_uri, response_type
   ↓
UserSessionManager → check if user is logged in
   ├─ if not → redirect to login
   ↓
ConsentManager → check if user has granted requested scopes to this client
   ├─ if not → show consent UI
   ↓
ResponseTypeRegistry → pick handler for response_type (code, token, id_token, hybrid)
   ↓
ResponseTypeHandler → issue code/tokens, build redirect URI
   ↓
HTTP Redirect to client

------------------------------------------------------------------------------------------------------------------------------------------------------------------

HTTP POST /token
   ↓
Adapter.BuildDTO(r) → ClientAuthDTO
   ↓
AuthenticatorRegistry.Find(dto) → ClientAuthenticator
   ↓
Authenticate(dto) → Client
   ↓
GrantRegistry.Get(dto.GrantType) → GrantHandler
   ↓
GrantHandler.Handle(Client, dto) → calls TokenIssuer
   ↓
TokenIssuer pipeline:
    ClaimsMappingStrategy
    TTLPolicy
    TokenFormatStrategy
    TokenDecorators
    SignerFactory
    Observers
   ↓
HTTP JSON Response

------------------------------------------------------------------------------------------------------------------------------------------------------------------

Extensible: Add a new client auth method by implementing ClientAuthenticator and registering it — no changes to core logic.

Testable: Each authenticator can be tested with just a ClientAuthDTO.

Config‑driven: Admin can enable/disable auth methods per client or globally.

Separation of concerns: Adapter parses, authenticators verify, grant handlers issue tokens.

------------------------------------------------------------------------------------------------------------------------------------------------------------------

HTTP Request (/token)
   ↓
Adapter → ClientAuthDTO
   ↓
AuthenticatorRegistry → uses ClientStore
   ↓
GrantRegistry → GrantHandler
       ├─ authorization_code → uses AuthCodeStore, IdentityStore, ClientStore
       ├─ refresh_token → uses RefreshTokenStore, IdentityStore
       ├─ client_credentials → uses ClientStore
   ↓
ClaimsMappingStrategy → uses IdentityStore, ScopeStore
   ↓
TokenIssuer → uses KeyStore, AccessTokenStore, RefreshTokenStore
   ↓
Observers → uses AuditLogStore
   ↓
HTTP Response

------------------------------------------------------------------------------------------------------------------------------------------------------------------

type AuthorizationRequestDTO struct {
    ResponseType string   // "code", "token", "id_token", etc.
    ClientID     string
    RedirectURI  string
    Scope        []string
    State        string
    Nonce        string
    CodeChallenge       string
    CodeChallengeMethod string
}


---------

function handleCodeFlow(authReq, user):
    // Create authorization code
    code = generateRandomCode()
    authCodeStore.save(code, {
        client_id: authReq.ClientID,
        user_id: user.ID,
        scopes: authReq.Scope,
        redirect_uri: authReq.RedirectURI,
        code_challenge: authReq.CodeChallenge
    })

    // Build redirect URI
    redirect = authReq.RedirectURI + "?code=" + code + "&state=" + authReq.State
    return redirect


------------

func (h *AuthorizationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 1. Adapt HTTP → DTO
    authReq := h.adapter.BuildDTO(r)

    // 2. Validate client & request params
    if err := h.clientValidatorRegistry.Validate(authReq); err != nil {
        h.writeOAuthError(w, err, authReq.RedirectURI, authReq.State)
        return
    }

    // 3. Check user session
    user, err := h.userSessionManager.GetLoggedInUser(r)
    if err == session.ErrNotLoggedIn {
        // Redirect to login, preserving original request
        h.userSessionManager.RedirectToLogin(w, r)
        return
    }

    // 4. Check consent
    if !h.consentManager.HasConsent(user.ID, authReq.ClientID, authReq.Scope) {
        h.consentManager.ShowConsentUI(w, r, user, authReq)
        return
    }

    // 5. Pick response type handler
    respHandler := h.responseTypeRegistry.Get(authReq.ResponseType)
    if respHandler == nil {
        h.writeOAuthError(w, oauth.ErrUnsupportedResponseType, authReq.RedirectURI, authReq.State)
        return
    }

    // 6. Let the handler issue code/tokens and build redirect URI
    redirectURI, err := respHandler.Handle(authReq, user)
    if err != nil {
        h.writeOAuthError(w, err, authReq.RedirectURI, authReq.State)
        return
    }

    // 7. Redirect back to client
    http.Redirect(w, r, redirectURI, http.StatusFound)
}


---

[HTTP GET /authorization]
        │
        ▼
[Adapter] → AuthorizationRequestDTO
        │
        ▼
[ClientValidatorRegistry]
        │
        ▼
[UserSessionManager] ──┐
        │ logged in?    │ no
        ▼ yes           ▼
[ConsentManager]     Redirect to login
        │ consent?       ↑
        ▼ yes            │
[ResponseTypeRegistry]   │
        │
        ▼
[ResponseTypeHandler]
        │
        ▼
[HTTP Redirect to client]


--------------------------------------
