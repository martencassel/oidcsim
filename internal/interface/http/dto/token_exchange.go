package dto

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

type TokenExchangeResponse struct {
	AccessToken     string `json:"access_token"`
	IssuedTokenType string `json:"issued_token_type"`
	TokenType       string `json:"token_type"`
	ExpiresIn       int    `json:"expires_in"`
	Scope           string `json:"scope,omitempty"`
}
