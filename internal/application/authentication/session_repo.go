package authentication

import (
	"context"
	"time"

	domain "github.com/martencassel/oidcsim/internal/domain/authentication"
)

type SessionRepo interface {
	GetAuthContext(ctx context.Context, sid string) (domain.Context, bool, error)
	SetAuthContext(ctx context.Context, sid string, ctxVal domain.Context, ttl time.Duration) error
	ClearAuthContext(ctx context.Context, sid string) error

	// flow state between steps (opaque to repo)
	GetFlowState(ctx context.Context, sid string) (map[string]string, bool, error)
	SetFlowState(ctx context.Context, sid string, state map[string]string, ttl time.Duration) error
	ClearFlowState(ctx context.Context, sid string) error
}
