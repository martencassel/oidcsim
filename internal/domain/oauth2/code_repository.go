package oauth2

import "context"

type CodeStore interface {
	Save(ctx context.Context, code *AuthorizationCode) error
	Get(ctx context.Context, code string) (*AuthorizationCode, error)
	Delete(ctx context.Context, code string) error
}
