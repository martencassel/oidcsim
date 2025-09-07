🔹 How to integrate this into your own IdP
In your own implementation (like the one we’ve been building):

DelegationService.EnsureConsent is where you enforce this policy.

If policy says “auto‑approve” (first‑party, remembered, basic scopes) → create the Delegation silently and return ConsentGranted.

If policy says “require consent” → return ConsentRequired so /authorize can redirect to /consent.

If prompt=none and no consent exists → return InteractionRequired so you can send an OAuth2 error back to the client.
