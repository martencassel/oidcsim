package oauth2

type TokenRequest struct {
	GrantType string // authorization_code, refresh_token, etc.

	// For authorization_code grant
	Code         string
	RedirectURI  string
	ClientID     string
	CodeVerifier string // for PKCE

	// For refresh_token grant
	RefreshToken string

	// Optional: client authentication
	ClientSecret     string
	ClientAssertion  string // JWT assertion for client authentication
	ClientRepository string // e.g., "client_secret_basic", "client_secret_post", "private_key_jwt"

	// Optional: requested scopes (space-separated)
	Scope []string
}
