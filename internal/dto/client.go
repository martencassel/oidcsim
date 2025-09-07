package dto

import "crypto/x509"

// ClientAuthDTO is the normalized, immutable representation of
// what the client presented for authentication.
type ClientAuthDTO struct {
	ClientID      string
	ClientSecret  string
	GrantType     string
	Scope         []string
	AuthMethod    string // basic, post, mtls, none
	TLSCert       *x509.Certificate
	RawAuthHeader string
	FormValues    map[string]string
}
