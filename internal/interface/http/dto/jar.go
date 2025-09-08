package dto

// For /authorize endpoint (OAuth2 + OIDC).

// JWTRequestObject represents the parameters for a JWT Request Object in an authorization request.
// See: https://openid.net/specs/openid-connect-core-1_0.html#JWTRequests
type JWTRequestObject struct {
	Request    string `form:"request" json:"request"`
	RequestURI string `form:"request_uri" json:"request_uri"`
}
