// jwksutil/jwks.go
package jwksutil

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
)

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid,omitempty"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func base64url(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func LoadRSAPublicKeyFromPEMFile(path string) (*rsa.PublicKey, error) {
	pemData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid PEM block type")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}

func ConvertToJWK(pub *rsa.PublicKey, kid string) JWK {
	nBytes := pub.N.Bytes()
	eBytes := big.NewInt(int64(pub.E)).Bytes()

	return JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: kid,
		N:   base64url(nBytes),
		E:   base64url(eBytes),
	}
}

func GenerateJWKS(pub *rsa.PublicKey, kid string) ([]byte, error) {
	jwk := ConvertToJWK(pub, kid)
	jwks := JWKS{Keys: []JWK{jwk}}
	return json.MarshalIndent(jwks, "", "  ")
}
