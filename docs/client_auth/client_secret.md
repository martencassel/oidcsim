---
## 🔑 What is the `client_secret` in OAuth/OIDC?
In OAuth 2.0 and OpenID Connect (OIDC), the `client_secret` is essentially the **application’s own password** — a shared secret known **only** to:
- The **client application** (e.g., your backend server)
- The **authorization server** (e.g., Google, Azure AD, Okta)

It’s issued when you register your app with the identity provider and is used to **prove that the app itself is legitimate**, not just the user.

---
## 🧠 Why is it needed?
Think of OAuth/OIDC as having **two layers of trust**:
1. **User authentication** – proving *who the user is* (via username/password, biometrics, etc.)
2. **Client authentication** – proving *which app is making the request*

The `client_secret` addresses the second layer. Without it, any malicious app could pretend to be your app and request tokens on behalf of users.

---
## ⚙️ How it works conceptually
Here’s the high-level flow when using `client_secret` authentication (e.g., in the **Authorization Code** or **Client Credentials** flow):

1. **App registration** – You register your app with the identity provider and receive:
   - `client_id` (public identifier)
   - `client_secret` (private key/password)

2. **Token request** – When your app exchanges an authorization code (or requests a token directly), it sends:
   - `client_id`
   - `client_secret`
   - Other required parameters (grant type, redirect URI, etc.)

3. **Server verification** – The authorization server checks:
   - Does this `client_id` exist?
   - Does the `client_secret` match what’s stored?

4. **Token issuance** – If valid, the server issues an access token (and possibly an ID token in OIDC).

---
## 📦 Common authentication methods using `client_secret`
OIDC/OAuth 2.0 supports several ways to send the secret:

| Method | How it’s sent | Notes |
|--------|--------------|-------|
| **HTTP Basic Auth** | `Authorization: Basic Base64(client_id:client_secret)` in the header | Simple, widely supported; must use HTTPS |
| **POST body** | As form parameters (`client_id`, `client_secret`) in the request body | Equivalent security to Basic Auth if over HTTPS |
| **Private Key JWT / mTLS** | More advanced, avoids static secrets | Used for higher security needs |

---
## 🚨 Security considerations
- **Never** embed `client_secret` in public clients (mobile apps, SPAs) — it can be extracted.
- Only use it in **confidential clients** (server-side apps) where it can be kept private.
- Always transmit over **HTTPS** to prevent interception.
- Rotate secrets periodically to reduce risk if compromised.

---
## Pseudo Code

function authenticateClientUsingSecret(request):
    # 0) MUST require TLS at token endpoint
    if not request.is_tls():
        return error_response(
            status=400,
            body=oauth_error("invalid_request", "TLS required")
        )
    # RFC 6749 requires TLS for token endpoint usage. (§1.6)

    # 1) Token endpoint expects application/x-www-form-urlencoded POST
    if request.method != "POST" or request.content_type != "application/x-www-form-urlencoded":
        return error_response(
            status=400,
            body=oauth_error("invalid_request", "POST with application/x-www-form-urlencoded required")
        )
    # (§3.2.1)

    # 2) Extract client credentials, enforcing “only one method” rule
    has_basic = request.headers.contains("Authorization")
                and scheme(request.headers["Authorization"]) == "Basic"
    has_post  = request.form.contains("client_id") or request.form.contains("client_secret")

    if has_basic and has_post:
        # Client MUST NOT use more than one method in a single request (§2.3.1)
        return error_response(
            status=400,
            body=oauth_error("invalid_request", "Multiple client authentication methods used")
        )

    client_id     = null
    client_secret = null
    used_method   = null

    if has_basic:
        # 3a) Parse Basic credentials
        # Per RFC 6749, client_id and client_secret are form-encoded BEFORE joining with ":" and base64-encoding
        # creds = Base64Decode( credentials_part_of_header )
        # raw = UTF8Decode(creds)
        username, password = split_on_first_colon(raw)  # split only on the first ":"
        client_id     = form_url_decode(username)
        client_secret = form_url_decode(password)
        used_method   = "client_secret_basic"
        # (§2.3.1)

    else if has_post:
        # 3b) Read from form body (client_secret_post)
        client_id     = request.form.get("client_id")
        client_secret = request.form.get("client_secret")
        used_method   = "client_secret_post"
        # (§2.3.1 permits credentials in the request body)

    else:
        # No client authentication provided (for this flow, we require it)
        # If Basic was not attempted, 400 is used for token endpoint errors (§5.2)
        return error_response(
            status=400,
            body=oauth_error("invalid_client", "Client authentication missing")
        )

    # 4) Look up client by client_id
    client = clients_repository.find_by_id(client_id)

    # Do NOT reveal whether client_id exists; keep timing consistent
    if client is null:
        return handleInvalidClient(used_method)

    # 5) Enforce the client's configured token_endpoint_auth_method (if your AS tracks per-client policy)
    # (Defined by Dynamic Client Registration metadata; values include client_secret_basic, client_secret_post, etc.)
    if not client.allows_auth_method(used_method):
        return handleInvalidClient(used_method)
    # (See RFC 7591 "token_endpoint_auth_method" metadata)

    # 6) Verify the secret using constant-time comparison against a stored hash
    # (e.g., Argon2/BCrypt/PBKDF2; never log secrets)
    if not crypto.verify_secret(client.hashed_secret, client_secret):
        return handleInvalidClient(used_method)

    # 7) Optional additional checks (typical for a token endpoint):

-----------

During client registration for `client_secret` authentication, you configure the following:

1. **Register your app with the authorization server**
   - You send a registration request (often via a web UI or API endpoint).
   - You provide metadata such as:
     - `client_name`
     - `redirect_uris`
     - `grant_types`
     - `token_endpoint_auth_method` (set to `client_secret_basic` or `client_secret_post`)

2. **The server generates and returns:**
   - `client_id` (public identifier)
   - `client_secret` (private, only shown once)

3. **Store the secret securely**
   - Your app must keep the `client_secret` safe (never expose it in public code or client-side apps).

4. **Use the secret for authentication**
   - When requesting tokens, your app authenticates using the configured method (Basic Auth or POST body).

**Example registration payload (API):**
```json
{
  "client_name": "My App",
  "redirect_uris": ["https://myapp.com/callback"],
  "grant_types": ["authorization_code"],
  "token_endpoint_auth_method": "client_secret_basic"
}
```

**Reference:**
See [OIDC Dynamic Client Registration](https://openid.net/specs/openid-connect-registration-1_0.html) for details.

Let me know if you want a sample curl or Go code for registration!
