package security

import (
	"crypto/rand"
	"encoding/base64"
)

// generateRandomString returns a URL‑safe random string of n bytes.
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

type RandomStringGenerator interface {
	Generate(n int) (string, error)
}

type DefaultRandomStringGenerator struct{}

func (g DefaultRandomStringGenerator) Generate(n int) (string, error) {
	return GenerateRandomString(n)
}
