package dto

// BackchannelAuthRequest represents the parameters for a CIBA backchannel authentication request.
// See: https://openid.net/specs/openid-client-initiated-backchannel-authentication-core-1_0.html#BackchannelAuthentication
type BackchannelAuthRequest struct {
	ClientID       string `form:"client_id" binding:"required"`
	Scope          string `form:"scope,omitempty"`
	LoginHint      string `form:"login_hint,omitempty"`
	BindingMessage string `form:"binding_message,omitempty"`
}

// BackchannelAuthResponse represents a successful response from the CIBA backchannel authentication endpoint.
// See: https://openid.net/specs/openid-client-initiated-backchannel-authentication-core-1_0.html#BackchannelAuthenticationResponse
type BackchannelAuthResponse struct {
	AuthReqID string `json:"auth_req_id"`
	ExpiresIn int    `json:"expires_in"`
	Interval  int    `json:"interval,omitempty"`
}
