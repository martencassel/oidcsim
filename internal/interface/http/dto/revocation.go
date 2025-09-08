package dto

// /revoke endpoint (RFC 7009).

type RevocationRequest struct {
	Token         string `form:"token" binding:"required"`
	TokenTypeHint string `form:"token_type_hint"` // "access_token" or "refresh_token"
}
