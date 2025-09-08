package dto

// For /logout endpoint (OIDC RP-Initiated Logout).

// LogoutRequest represents the parameters for an OIDC RP-Initiated Logout request.
// See: https://openid.net/specs/openid-connect-rpinitiated-1_0.html#RPLogout
type LogoutRequest struct {
	IDTokenHint           string `form:"id_token_hint,omitempty"`
	PostLogoutRedirectURI string `form:"post_logout_redirect_uri,omitempty"`
	State                 string `form:"state,omitempty"`
}

// LogoutResponse represents a successful response from the logout endpoint.
// See: https://openid.net/specs/openid-connect-rpinitiated-1_0.html#RPLogout
type LogoutResponse struct {
	RedirectURI string `json:"redirect_uri,omitempty"`
}

// For /check_session_frame endpoint (OIDC Session Management).

// SessionStatus represents the session status for the OIDC Session Management
// check_session_iframe endpoint.
// See: https://openid.net/specs/openid-connect-session-1_0.html#CheckSessionIframe
// /check_session_frame
type SessionStatus struct {
	SessionState string `json:"session_state"`
	Sub          string `json:"sub"`
}
