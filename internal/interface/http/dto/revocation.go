package dto

// /revoke endpoint (RFC 7009).

// RevocationRequest represents the parameters for an OAuth2 token revocation request.
// See: https://datatracker.ietf.org/doc/html/rfc7009#section-2.1
type RevocationRequest struct {
	Token         string `form:"token" binding:"required"`
	TokenTypeHint string `form:"token_type_hint"` // "access_token" or "refresh_token"
}
