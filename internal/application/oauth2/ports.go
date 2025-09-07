package oauth2app

import (
	"context"
	"time"

	"github.com/martencassel/oidcsim/internal/domain/oauth2"
)

type AuthorizationCodeRepo interface {
	Save(ctx context.Context, code oauth2.AuthorizationCode) error
	Get(ctx context.Context, codeValue string) (*oauth2.AuthorizationCode, error)
	Delete(ctx context.Context, codeValue string) error
}

type TokenService interface {
	MintAccessToken(ctx context.Context, subjectID, clientID string, scopes []string) (string, error)
	MintIDToken(ctx context.Context, subjectID, clientID string, nonce string) (string, error)
}

type ClientRepo interface {
	Get(ctx context.Context, clientID string) (*oauth2.Client, error)
}

type Clock interface {
	Now() time.Time
}

type NonceGenerator interface {
	Generate() (string, error)
}
