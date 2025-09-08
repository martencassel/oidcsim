package dto

// PAR (Pushed Authorization Requests) (RFC 9126).

// PushedAuthorizationRequest represents the parameters for a pushed authorization request.
// See: https://datatracker.ietf.org/doc/html/rfc9126#section-2.1
type PushedAuthorizationRequest struct {
	ClientID string `form:"client_id" binding:"required"`
	Scope    string `form:"scope,omitempty"`
	// plus any other /authorize params
}

// PushedAuthorizationResponse represents a successful response from the pushed authorization request endpoint.
// See: https://datatracker.ietf.org/doc/html/rfc9126#section-2.2
type PushedAuthorizationResponse struct {
	RequestURI string `json:"request_uri"`
	ExpiresIn  int    `json:"expires_in"`
}
