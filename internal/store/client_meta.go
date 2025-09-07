package store

import "github.com/martencassel/oidcsim/internal/dto"

// ClientMeta is for client registration and policy.
type ClientMeta struct {
	ID                 string                 // Client identifier
	AllowedAuthMethods []dto.ClientAuthMethod // Which auth methods this client can use
	SecretHash         string                 // Hashed secret for secret-based methods
	JWKSURI            string                 // JWKS URI for private_key_jwt
	JWKS               []byte                 // Static JWKS (optional)
	TLSAuthSubjectDN   string                 // Subject DN for mTLS
	TLSSANs            []string               // SANs for mTLS
	Enabled            bool                   // Is the client enabled
}
