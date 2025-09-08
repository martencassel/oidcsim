package dto

// For /token endpoint with "urn:ietf:params:oauth:grant-type:token-exchange" grant type (RFC 8693).

// TokenExchangeRequest represents the parameters for a token exchange request.
// See: https://datatracker.ietf.org/doc/html/rfc8693#section-2.1
type TokenExchangeRequest struct {
	SubjectToken       string `form:"subject_token" binding:"required"`
	SubjectTokenType   string `form:"subject_token_type" binding:"required"`
	ActorToken         string `form:"actor_token,omitempty"`
	ActorTokenType     string `form:"actor_token_type,omitempty"`
	Resource           string `form:"resource,omitempty"`
	Audience           string `form:"audience,omitempty"`
	Scope              string `form:"scope,omitempty"`
	RequestedTokenType string `form:"requested_token_type,omitempty"`
}

// TokenExchangeResponse represents a successful response from the token exchange endpoint.
// See: https://datatracker.ietf.org/doc/html/rfc8693#section-2.2
type TokenExchangeResponse struct {
	AccessToken     string `json:"access_token"`
	IssuedTokenType string `json:"issued_token_type"`
	TokenType       string `json:"token_type"`
	ExpiresIn       int    `json:"expires_in"`
	Scope           string `json:"scope,omitempty"`
}
