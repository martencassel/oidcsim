package dto

import (
	"net/url"
	"testing"
)

func TestBuildAuthorizeRequest_Valid(t *testing.T) {
	form := url.Values{}
	form.Set("response_type", "code")
	form.Set("client_id", "my-client")
	form.Set("redirect_uri", "https://client.example/cb")
	form.Set("scope", "openid profile")
	form.Set("state", "xyz")
	form.Set("code_challenge", "abc123")
	form.Set("code_challenge_method", "S256")
	form.Set("nonce", "n-0S6_WzA2Mj")
	//req := httptest.NewRequest(http.MethodGet, "/authorize?"+form.Encode(), nil)
	// ar, err := BuildAuthorizeRequest(req)
	// assert.NoError(t, err)
	// assert.Equal(t, "code", ar.ResponseType)
	// assert.Equal(t, "my-client", ar.ClientID)
	// assert.Equal(t, "https://client.example/cb", ar.RedirectURI)
	// assert.Equal(t, "openid profile", ar.Scope)
	// assert.Equal(t, "xyz", ar.State)
	// assert.Equal(t, "abc123", ar.CodeChallenge)
	// assert.Equal(t, "S256", ar.CodeChallengeMethod)
	// assert.Equal(t, "n-0S6_WzA2Mj", ar.Nonce)
}

func TestBuildTokenRequest_ValidAuthCode(t *testing.T) {
}
