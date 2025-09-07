package domain

import (
	"context"
	"time"

	"github.com/martencassel/oidcsim/internal/security"
	"github.com/martencassel/oidcsim/internal/store"
)

// TokenIssuer issues access tokens in either opaque or JWT format.
type TokenIssuer interface {
	IssueAccessToken(ctx context.Context, client store.Client, user store.User, scope []string) (string, error)
}

// IDTokenIssuer always issues JWT ID tokens for OIDC.
type IDTokenIssuer interface {
	IssueIDToken(ctx context.Context, client store.Client, user store.User, nonce string, scope []string) (string, error)
}

type JWTTokenIssuer struct {
	Signer       security.JWTSigner
	Expiry       time.Duration
	ClaimsMapper ClaimsMapper
}

func (i *JWTTokenIssuer) IssueIDToken(ctx context.Context, client store.Client, user store.User, nonce string, scopes []string) (string, error) {
	claims := i.ClaimsMapper.MapIDTokenClaims(ctx, user, client, scopes)
	claims["nonce"] = nonce
	claims["exp"] = time.Now().Add(i.Expiry).Unix()
	return i.Signer.Sign(claims)
}
