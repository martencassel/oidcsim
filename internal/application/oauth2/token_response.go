package oauth2

type TokenResponse struct {
	AccessToken  string   // The issued access token
	TokenType    string   // Usually "Bearer"
	ExpiresIn    int      // Lifetime in seconds
	RefreshToken string   // Optional, if issued
	Scope        []string // Scopes granted
	IDToken      string   // Optional, for OIDC flows
}

func (r *TokenResponse) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"access_token": r.AccessToken,
		"token_type":   r.TokenType,
		"expires_in":   r.ExpiresIn,
	}
	if r.RefreshToken != "" {
		m["refresh_token"] = r.RefreshToken
	}
	if len(r.Scope) > 0 {
		m["scope"] = r.Scope
	}
	if r.IDToken != "" {
		m["id_token"] = r.IDToken
	}
	return m
}

func NewTokenResponse(accessToken, tokenType string, expiresIn int, refreshToken string, scope []string, idToken string) *TokenResponse {
	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		ExpiresIn:    expiresIn,
		RefreshToken: refreshToken,
		Scope:        scope,
		IDToken:      idToken,
	}
}
