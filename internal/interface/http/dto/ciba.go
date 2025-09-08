package dto

type BackchannelAuthRequest struct {
	ClientID       string `form:"client_id" binding:"required"`
	Scope          string `form:"scope,omitempty"`
	LoginHint      string `form:"login_hint,omitempty"`
	BindingMessage string `form:"binding_message,omitempty"`
}

type BackchannelAuthResponse struct {
	AuthReqID string `json:"auth_req_id"`
	ExpiresIn int    `json:"expires_in"`
	Interval  int    `json:"interval,omitempty"`
}
