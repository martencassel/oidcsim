package dto

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthorizeRequest struct{}

var validate = validator.New()

func SplitAndTrim(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			if start < i {
				result = append(result, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}

func parseScope(scopeStr string) []string {
	var scopes []string
	for _, s := range SplitAndTrim(scopeStr, " ") {
		if s != "" {
			scopes = append(scopes, s)
		}
	}
	return scopes
}

func BuildAuthorizeRequest(r *http.Request) (AuthorizeRequest, error) {
	var req AuthorizeRequest
	if err := r.ParseForm(); err != nil {
		return req, err
	}
	formScopes := r.FormValue("scope")
	if formScopes != "" {
		//req.Scope = parseScope(formScopes)
	}
	// req.ResponseType = r.FormValue("response_type")
	// req.ClientID = r.FormValue("client_id")
	// req.RedirectURI = r.FormValue("redirect_uri")
	// req.Scope = parseScope(r.FormValue("scope"))
	// req.State = r.FormValue("state")
	// req.CodeChallenge = r.FormValue("code_challenge")
	// req.CodeChallengeMethod = r.FormValue("code_challenge_method")
	// req.Nonce = r.FormValue("nonce")

	return req, validate.Struct(req)
}

func BuildTokenRequest(r *http.Request) (TokenRequest, error) {
	var req TokenRequest
	if err := r.ParseForm(); err != nil {
		return req, err
	}
	req.GrantType = r.FormValue("grant_type")
	req.Code = r.FormValue("code")
	req.RedirectURI = r.FormValue("redirect_uri")
	req.ClientID = r.FormValue("client_id")
	req.ClientSecret = r.FormValue("client_secret")
	req.CodeVerifier = r.FormValue("code_verifier")
	req.RefreshToken = r.FormValue("refresh_token")
	req.Scope = r.FormValue("scope")
	req.DeviceCode = r.FormValue("device_code")
	return req, validate.Struct(req)
}

// Per-Flow DTO
//
// TokenRequest -> AuthCodeTokenRequest, RefreshTokenRequest
//
// PoP integration: Add fields for DPoP header, mTLS cert thumbprint, etc., to the relevant DTOs.
