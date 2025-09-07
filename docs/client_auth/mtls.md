Mutual TLS (mTLS)

Uses X.509 client certificates for strong authentication.
Often required for high-security environments and financial-grade APIs (FAPI).

---

Alright — let’s break down **Mutual TLS (mTLS)** in the context of OAuth 2.0 / OIDC so you can see the *conceptual model* and why it’s used in high‑security environments like **FAPI** (Financial‑grade API).

---

## 🧩 The core idea
In **regular TLS** (HTTPS), only the **server** presents a certificate — your browser or app verifies it to ensure you’re talking to the right server.

In **Mutual TLS**, **both sides** present and verify X.509 certificates during the TLS handshake:
- **Server certificate** → proves the server’s identity to the client.
- **Client certificate** → proves the client’s identity to the server.

This creates a **two‑way trust** before any OAuth or OIDC messages are even exchanged.

---

## 🔍 Why it’s used in OAuth/OIDC
- **Strong client authentication** — replaces static `client_secret` with cryptographic proof.
- **Sender‑constrained tokens** — access tokens can be bound to the client’s certificate, so even if stolen, they can’t be used without that certificate.
- **Meets regulatory requirements** — required in some open banking / FAPI profiles.

---

## ⚙️ Conceptual flow

1. **Provisioning**
   - The client generates an X.509 certificate (public/private key pair).
   - The public certificate is registered with the authorization server during client onboarding.

2. **TLS handshake with mutual authentication**
   - When the client connects to the token endpoint (or resource server), it presents its certificate.
   - The server verifies the certificate against its trust store and the registered client metadata.
   - The client also verifies the server’s certificate (normal HTTPS).

3. **OAuth/OIDC request**
   - Once the TLS session is established, the client sends its OAuth request (e.g., token request).
   - The server already knows which client it is, based on the certificate — no `client_secret` needed.

4. **Optional: Certificate‑bound access tokens**
   - The authorization server can issue tokens that are cryptographically bound to the client’s certificate.
   - Resource servers then check that the certificate used in the TLS session matches the one bound to the token.

---

## 📦 Two common mTLS client authentication methods (RFC 8705)
| Method | How it works | Typical use |
|--------|--------------|-------------|
| **PKI-based (`tls_client_auth`)** | Client uses a certificate issued by a trusted CA; server validates chain and subject | Enterprise, regulated industries |
| **Self-signed (`self_signed_tls_client_auth`)** | Client uses a self-signed cert; public key registered directly with server | Simpler setups, dev/test environments |

---

## 🔐 Security benefits
- Eliminates static secrets in public code or config.
- Prevents token replay — stolen tokens are useless without the matching certificate.
- Strong cryptographic identity at the network layer before application logic runs.

---

💡 **Analogy**:
Think of normal TLS as showing your ID to enter a secure building — only the guard (server) proves they’re legitimate.
Mutual TLS is like *both* you and the guard showing IDs to each other before you even step inside — and your visitor badge (access token) only works if you’re wearing the exact same ID card you used at the door.

---

If you want, I can also walk you through **how to provision and configure mTLS for an OAuth client** step‑by‑step, so you see exactly what’s needed on both the client and server sides. Would you like me to do that next?


Got it — let’s walk through the **end‑to‑end setup** for enabling **Mutual TLS (mTLS) client authentication** in OAuth 2.0 / OIDC from an **admin + client developer** perspective.

This is the **conceptual provisioning model** you’d follow in most IdPs (Okta, Auth0, Curity, Azure AD, etc.), even though the exact UI/CLI steps vary.

---

## 🛠 Step 1 — Prepare the client certificate & keys
**Who:** Client developer (or security team)

1. **Generate a key pair** (private + public key) and create an **X.509 client certificate**.
   - Can be **CA‑signed** (enterprise PKI) or **self‑signed** (for dev/test).
   - For CA‑signed, ensure the full certificate chain is available.
2. **Securely store** the private key on the client system (HSM, KMS, secure vault).
3. Export the **public certificate** (PEM or DER format) — this is what the IdP will trust.

---

## 🏷 Step 2 — Register the client in the Authorization Server
**Who:** Admin in the IdP

When creating or editing the OAuth/OIDC client:
- **Set the authentication method** to `tls_client_auth` or `self_signed_tls_client_auth` (per RFC 8705).
- **Upload the public certificate** or register a **JWKS endpoint** containing the public key.
- Configure **allowed grant types** (e.g., Authorization Code, Client Credentials).
- Configure **redirect URIs** (for browser-based flows).
- If using **certificate-bound access tokens**, enable sender‑constraining in the client settings.

---

## 🔐 Step 3 — Configure trust on the server side
**Who:** Admin in the IdP / API gateway

- Add the client’s public certificate (or CA chain) to the **trust store** used by the token endpoint and/or resource server.
- If using CA‑signed certs, ensure the CA root/intermediate certs are trusted.
- Configure TLS termination to **require client certificates** for the relevant endpoints.

---

## 📡 Step 4 — Client implementation
**Who:** Client developer

When making a token request:
1. Establish an **HTTPS connection with mutual TLS** to the token endpoint.
   - The client presents its certificate during the TLS handshake.
2. The server verifies the certificate before processing the OAuth request.
3. The OAuth request (e.g., `grant_type=client_credentials`) is sent **without a client_secret** — the TLS layer has already authenticated the client.

---

## ✅ Step 5 — Verification & testing
**Who:** Both admin + developer

- Test the handshake:
  - If the client certificate is missing or invalid, the TLS connection should fail before OAuth logic runs.
- Test token issuance:
  - Confirm the server issues tokens only when the correct certificate is presented.
- If using **certificate-bound tokens**, test that the resource server rejects tokens presented over a TLS session with a different certificate.

---

## 🔒 Security best practices
- **Rotate certificates** before expiry and update the IdP trust store.
- Use **short-lived certificates** for high-security environments.
- Store private keys in **HSM/KMS** — never in plain files.
- Monitor for failed TLS handshakes as potential intrusion attempts.

---

If you want, I can also give you a **visual diagram** showing where the mTLS handshake happens relative to the OAuth token request, so you can see exactly how the trust is established before any OAuth data is exchanged. That makes the flow much easier to picture. Would you like me to prepare that?
