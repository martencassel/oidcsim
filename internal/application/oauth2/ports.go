package oauth2

import (
	"context"
	"time"

	"github.com/martencassel/oidcsim/internal/domain/oauth2"
	oauth2client "github.com/martencassel/oidcsim/internal/domain/oauth2/client"
)

type AuthorizationCodeRepo interface {
	Save(ctx context.Context, code oauth2.AuthorizationCode) error
	Get(ctx context.Context, codeValue string) (*oauth2.AuthorizationCode, error)
	Delete(ctx context.Context, codeValue string) error
}

type ClientRepo interface {
	Get(ctx context.Context, clientID string) (*oauth2client.Client, error)
}

type Clock interface {
	Now() time.Time
}

type NonceGenerator interface {
	Generate() (string, error)
}

// defines CodeRepository
type CodeRepository interface {
	// methods for code repository
}

// defines JWTIssuer
type JWTIssuer interface {
	// methods for JWT issuer
}

// defines RefreshTokenRepository
type RefreshTokenRepository interface {
	// methods for refresh token repository
}
