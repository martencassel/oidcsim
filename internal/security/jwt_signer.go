package security

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
)

type JWTSigner interface {
	// Sign takes a set of claims and returns a signed JWT string.
	Sign(claims map[string]interface{}) (string, error)

	// KeyID returns the kid (Key ID) for the current signing key.
	KeyID() string
}

type RS256Signer struct {
	privateKey *rsa.PrivateKey
	keyID      string
}

func NewRS256Signer(priv *rsa.PrivateKey, kid string) *RS256Signer {
	return &RS256Signer{privateKey: priv, keyID: kid}
}

func (s *RS256Signer) Sign(claims map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claims))
	token.Header["kid"] = s.keyID
	return token.SignedString(s.privateKey)
}

func (s *RS256Signer) KeyID() string {
	return s.keyID
}
