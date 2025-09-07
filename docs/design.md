## Adapter Layer
Extract what the client **presented** (headers, form fields, TLS peer cert) and normalize to an immutable payload for the domain layer.
  *No policy decisions here.*


## Policy and Client Registration
Handle the clients registration + policy:
- ID, AllowedAuthMethods, SecretHash, Enabled, JWKSURI / JWKS, TLSAuthSubjectDN, TLSSANs

**store.ClientStore**
Repository to fetch `ClientMeta` by ID.

