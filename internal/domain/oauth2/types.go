package oauth2

import (
	"time"

	infrasec "github.com/martencassel/oidcsim/internal/infrastructure/security"
)

type AuthorizationCodeParams struct {
	Gen                 infrasec.RandomStringGenerator
	ClientID            string
	RedirectURI         string
	Scope               string
	State               string
	SubjectID           string
	SessionID           string
	CodeChallenge       string
	CodeChallengeMethod string
	Nonce               string
	AuthTime            time.Time // optional: zero => time.Now()
	TTL                 time.Duration
}

// domain/oauth2/code_factory.go
func NewAuthorizationCodeFromParams(p AuthorizationCodeParams) (AuthorizationCode, error) {
	if p.Gen == nil {
		p.Gen = infrasec.DefaultRandomStringGenerator{} // provide default
	}
	if p.AuthTime.IsZero() {
		p.AuthTime = time.Now()
	}
	if p.TTL == 0 {
		p.TTL = 5 * time.Minute
	}
	// generate code, populate fields, return
	return NewAuthorizationCode(p.Gen, p.ClientID, p.RedirectURI, p.Scope, p.State, p.SubjectID, p.SessionID, p.CodeChallenge, p.CodeChallengeMethod, p.Nonce, p.AuthTime, p.TTL)
}
