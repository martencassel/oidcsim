Private Key JWT

The client signs a JWT with its private key and sends it to the token endpoint.
The server validates the signature using the registered public key.

----

Got it — let’s walk through **Private Key JWT** in OAuth 2.0 / OIDC conceptually, so you can see the moving parts and why it exists.

---

## 🧩 The core concept
Private Key JWT is a **client authentication method** where:
- The client proves its identity by **digitally signing** a short-lived JWT with its **private key**.
- The authorization server verifies that signature using the **public key** it already has on file for that client.

This replaces the need for a static `client_secret` and gives you stronger, asymmetric cryptographic proof of identity.

---

## 🔍 Why it’s used
- **No shared secret** to store on both sides — only the client holds the private key.
- **Short-lived, one-time assertions** — each JWT expires quickly and can’t be reused.
- **Replay protection** — unique IDs (`jti`) prevent attackers from reusing an intercepted token.
- **Better for high-security integrations** — especially in regulated industries (finance, healthcare).

---

## ⚙️ How it works step-by-step

1. **Key pair setup**
   - The client generates a **public/private key pair**.
   - The public key is registered with the authorization server during client onboarding.
   - The private key stays securely on the client’s server.

2. **JWT creation**
   - When the client needs to authenticate (e.g., at the token endpoint), it builds a JWT with claims like:
     - `iss` (issuer) → the client’s ID
     - `sub` (subject) → the client’s ID
     - `aud` (audience) → the token endpoint URL
     - `exp` (expiration) → a few minutes in the future
     - `jti` (JWT ID) → a unique random value
   - The client **signs** this JWT with its private key.

3. **Token request**
   - The client sends the signed JWT to the token endpoint with:
     - `client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer`
     - `client_assertion=<signed JWT>`
     - Plus the usual OAuth parameters (grant type, code, etc.).

4. **Server verification**
   - The authorization server looks up the client’s registered public key.
   - It verifies the JWT’s signature and checks the claims (audience matches, not expired, jti unused).

5. **Token issuance**
   - If everything checks out, the server issues the requested access token (and ID token in OIDC).

---

## 📦 Conceptual analogy
Think of it like a **wax-sealed letter**:
- Only you have the seal (private key).
- The recipient has a mold of your seal (public key) and can instantly tell if the letter really came from you.
- Each letter is dated and numbered so it can’t be reused later.

--- Setup step --------------

Alright — let’s walk through the **provisioning and configuration steps** for enabling **Private Key JWT** client authentication in OAuth 2.0 / OIDC, from a conceptual and practical standpoint.

---

## 🛠 1. Generate a key pair
You need an **asymmetric key pair** for signing and verification.
- **Private key** → stays securely on your client (server-side only).
- **Public key** → shared with the OIDC provider (authorization server).

**Typical ways to generate keys**:
- Using `openssl` (PEM format)
- Using `keytool` (Java keystore)
- Using a cloud KMS (AWS KMS, Azure Key Vault, GCP KMS)

---

## 🏷 2. Register the client with the OIDC provider
When you create or update your client application in the IdP’s admin console (Okta, Azure AD, Auth0, WSO2, etc.):
- Set **Client Authentication Method** to **Private Key JWT** (sometimes called `private_key_jwt`).
- Upload or register your **public key**:
  - As a raw PEM certificate
  - Or via a **JWKS (JSON Web Key Set)** endpoint your app hosts
- Ensure the **grant types** you need (e.g., Authorization Code, Client Credentials) are enabled.

---

## 📜 3. Configure JWT claim requirements
Your client will need to generate a JWT with specific claims when authenticating:

| Claim | Purpose |
|-------|---------|
| `iss` | Issuer — your client ID |
| `sub` | Subject — your client ID |
| `aud` | Audience — the token endpoint URL of the IdP |
| `exp` | Expiration — short-lived (e.g., +300 seconds) |
| `jti` | Unique ID — prevents replay attacks |

---

## 🔏 4. Implement signing logic in the client
When your app requests a token:
1. Build the JWT with the claims above.
2. Sign it with your **private key** using RS256 (or another supported algorithm).
3. Send it to the token endpoint with:
   - `client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer`
   - `client_assertion=<signed JWT>`
   - Plus your normal OAuth parameters (`grant_type`, `code`, etc.).

---

## ✅ 5. Authorization server verification
The OIDC provider will:
- Look up your registered public key.
- Verify the JWT signature.
- Validate claims (audience, expiry, jti uniqueness).
- If valid, issue the access token (and ID token if OIDC).

---

## 🔐 Security best practices
- **Rotate keys** periodically — update the public key in the IdP and deploy the new private key to your app.
- **Store private keys securely** — use a secure vault or HSM/KMS.
- **Short expiry** for assertions — 1–5 minutes is common.
- **Replay protection** — ensure `jti` is unique per request.

---

