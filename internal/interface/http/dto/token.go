package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	GrantType string `form:"grant_type" json:"grant_type" binding:"required"` // e.g. "authorization_code", "refresh_token"

	// For authorization_code grant
	Code         string `form:"code" json:"code"`
	RedirectURI  string `form:"redirect_uri" json:"redirect_uri"`
	CodeVerifier string `form:"code_verifier" json:"code_verifier"`

	// For refresh_token grant
	RefreshToken string `form:"refresh_token" json:"refresh_token"`

	// Client authentication
	ClientID            string `form:"client_id" json:"client_id"`
	ClientSecret        string `form:"client_secret" json:"client_secret"`
	ClientAssertion     string `form:"client_assertion" json:"client_assertion"`
	ClientAssertionType string `form:"client_assertion_type" json:"client_assertion_type"`

	// Optional scope override
	Scope string `form:"scope" json:"scope"`
}

// Bind using go gin framework
func (tr *TokenRequest) Bind(c *gin.Context) error {
	return c.ShouldBind(tr)
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`            // The issued access token
	TokenType    string `json:"token_type"`              // Usually "Bearer"
	ExpiresIn    int    `json:"expires_in"`              // Lifetime in seconds
	RefreshToken string `json:"refresh_token,omitempty"` // Optional, if issued
	Scope        string `json:"scope,omitempty"`         // Space-delimited scopes
	IDToken      string `json:"id_token,omitempty"`      // Optional, for OIDC flows
}

// Bind using go gin framework
func (tr *TokenResponse) WriteResponse(c *gin.Context) error {
	c.JSON(http.StatusOK, tr)
	return nil
}
