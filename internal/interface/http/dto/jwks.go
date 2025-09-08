package dto

// For /jwks endpoint (RFC 7517).

// JWKSResponse represents a JSON Web Key Set (JWKS).
// See: https://tools.ietf.org/html/rfc7517#section-5
type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a single JSON Web Key.
// See: https://tools.ietf.org/html/rfc7517
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use,omitempty"`
	Kid string `json:"kid,omitempty"`
	Alg string `json:"alg,omitempty"`
	N   string `json:"n,omitempty"`
	E   string `json:"e,omitempty"`
	Crv string `json:"crv,omitempty"`
	X   string `json:"x,omitempty"`
	Y   string `json:"y,omitempty"`
}
