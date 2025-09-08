package dto

import "github.com/gin-gonic/gin"

type AuthorizeRequest struct {
	ResponseType        string `form:"response_type" query:"response_type"`
	ClientID            string `form:"client_id" query:"client_id"`
	RedirectURI         string `form:"redirect_uri" query:"redirect_uri"`
	Scope               string `form:"scope" query:"scope"`
	State               string `form:"state" query:"state"`
	ResponseMode        string `form:"response_mode" query:"response_mode"`
	Nonce               string `form:"nonce" query:"nonce"`
	Display             string `form:"display" query:"display"`
	Prompt              string `form:"prompt" query:"prompt"`
	MaxAge              string `form:"max_age" query:"max_age"`
	CodeChallenge       string `form:"code_challenge" query:"code_challenge"`
	CodeChallengeMethod string `form:"code_challenge_method" query:"code_challenge_method"`
}

// Bind using go gin framework
func (ar *AuthorizeRequest) Bind(c *gin.Context) error {
	ar.ClientID = c.Query("client_id")
	ar.ResponseType = c.Query("response_type")
	ar.RedirectURI = c.Query("redirect_uri")
	ar.Scope = c.Query("scope")
	ar.State = c.Query("state")
	ar.ResponseMode = c.Query("response_mode")
	ar.Nonce = c.Query("nonce")
	ar.Display = c.Query("display")
	ar.Prompt = c.Query("prompt")
	ar.MaxAge = c.Query("max_age")
	ar.CodeChallenge = c.Query("code_challenge")
	ar.CodeChallengeMethod = c.Query("code_challenge_method")
	return nil
}

type AuthorizeResponse struct {
	RedirectURI string
}

type AuthorizeErrorRedirect struct {
	RedirectURI      string
	State            string
	Error            string
	ErrorDescription string
}
