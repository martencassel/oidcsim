package dto

type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required,oneof=authorization_code refresh_token client_credentials"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri" binding:"required,url"`
	ClientID     string `form:"client_id" binding:"required"`
	ClientSecret string `form:"client_secret" validate:"omitempty"`
	CodeVerifier string `form:"code_verifier"`
	RefreshToken string `form:"refresh_token"`
	Scope        string `form:"scope"`
	DeviceCode   string `form:"device_code"`
}
