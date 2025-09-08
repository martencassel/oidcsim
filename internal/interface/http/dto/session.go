package dto

type LogoutRequest struct {
	IDTokenHint           string `form:"id_token_hint,omitempty"`
	PostLogoutRedirectURI string `form:"post_logout_redirect_uri,omitempty"`
	State                 string `form:"state,omitempty"`
}

type LogoutResponse struct {
	RedirectURI string `json:"redirect_uri,omitempty"`
}

// /check_session_frame
type SessionStatus struct {
	SessionState string `json:"session_state"`
	Sub          string `json:"sub"`
}
