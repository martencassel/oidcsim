
---

## 🧩 The core idea
Instead of proving the client’s identity with a **shared secret** (like `client_secret`),
the client proves who it is by sending a **signed JSON Web Token (JWT)** — called a **client assertion** — to the authorization server.

This is often referred to as **Private Key JWT client authentication**.
- The JWT is signed with the client’s **private key**.
- The authorization server verifies it using the **public key** you registered earlier.

---

## 🔍 Why use it?
- **No shared secret** to store or transmit — reduces risk of leaks.
- **Short-lived & one-time use** — each assertion has an expiration and unique ID, making replay attacks harder.
- **Cryptographic proof** — only the holder of the private key can produce a valid signature.

---

## ⚙️ How it works conceptually

1. **Key setup**
   - You generate a public/private key pair.
   - You give the **public key** to the authorization server during client registration.
   - You keep the **private key** safe on your server.

2. **When requesting a token**
   - Your app creates a JWT with specific claims:
     - `iss` (issuer) → your client ID
     - `sub` (subject) → your client ID
     - `aud` (audience) → the token endpoint URL
     - `exp` (expiration) → very short lifetime (e.g., 1–5 minutes)
     - `jti` (JWT ID) → unique identifier to prevent reuse
   - You sign this JWT with your **private key**.

3. **Send to the token endpoint**
   - Along with the normal OAuth parameters, you send:
     - `client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer`
     - `client_assertion=<your signed JWT>`

4. **Server verification**
   - The authorization server checks the signature using your public key.
   - It validates the claims (audience matches, not expired, jti not reused).
   - If valid, it issues the access token (and ID token in OIDC).

---

## 📦 Where it’s used
- **Authorization Code flow** (server-side apps)
- **Client Credentials flow** (machine-to-machine)
- High-security integrations (e.g., financial APIs, healthcare systems)

---

## 🚨 Security benefits over `client_secret`
| Client Secret | Client Assertion JWT |
|---------------|----------------------|
| Static secret stored on both sides | No static secret — uses asymmetric keys |
| Sent on every request | Only sends a short-lived signed JWT |
| If leaked, attacker can impersonate indefinitely | If leaked, expires quickly and can be revoked by key rotation |

---

**How it works:**
- During registration (or update), your client sends a JSON object to the auth server’s registration endpoint.
- Include your public key as a JWK (or a JWKS containing one or more JWKs) in
  the `jwks` field, or provide a URL in `jwks_uri` where the server can fetch your JWKS.

**Example registration payload:**
```json
{
  "client_id": "your-client-id",
  "client_name": "My App",
  "jwks": {
    "keys": [
      {
        "kty": "RSA",
        "kid": "my-key-id",
        "use": "sig",
        "alg": "RS256",
        "n": "<base64url modulus>",
        "e": "<base64url exponent>"
      }
    ]
  }
}
```
Or, using `jwks_uri`:
```json
{
  "client_id": "your-client-id",
  "client_name": "My App",
  "jwks_uri": "https://your-app.com/.well-known/jwks.json"
}
```

**How to update:**
- Send a PATCH or PUT request to the client management endpoint with the new `jwks` or `jwks_uri`.

**Reference:**
See [OIDC Dynamic Client Registration spec](https://openid.net/specs/openid-connect-registration-1_0.html) for details.

