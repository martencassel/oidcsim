package dto

// OIDC RP-Initiated Logout Token

type LogoutToken struct {
	Iss    string         `json:"iss"`
	Sub    string         `json:"sub,omitempty"`
	Sid    string         `json:"sid,omitempty"`
	Events map[string]any `json:"events"`
}
