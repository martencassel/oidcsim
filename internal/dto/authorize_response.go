package dto

type AuthorizeResponse struct {
	RedirectURI string
	Code        string
	State       string
	Error       string
	ErrorDesc   string
}
