package dto

// PAR (Pushed Authorization Requests) (RFC 9126).

type PushedAuthorizationRequest struct {
	ClientID string `form:"client_id" binding:"required"`
	Scope    string `form:"scope,omitempty"`
	// plus any other /authorize params
}

type PushedAuthorizationResponse struct {
	RequestURI string `json:"request_uri"`
	ExpiresIn  int    `json:"expires_in"`
}
