package dto

// OIDC RP-Initiated Logout Token

// See: https://openid.net/specs/openid-connect-rpinitiated-1_0.html#LogoutToken
type LogoutToken struct {
	Iss    string         `json:"iss"`
	Sub    string         `json:"sub,omitempty"`
	Sid    string         `json:"sid,omitempty"`
	Events map[string]any `json:"events"`
}
