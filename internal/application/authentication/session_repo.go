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

type FlowEngine interface {
	// Given clientID, return the flow (e.g., pwd → totp)
	Plan(ctx context.Context, clientID string) (domain.FlowSpec, error)

	// Render/init a step: returns hints and transient state updates
	StartStep(ctx context.Context, step domain.StepSpec, state map[string]string) (map[string]string, error)

	// Complete a step with inputs; returns whether more steps remain and state updates
	CompleteStep(ctx context.Context, step domain.StepSpec, inputs map[string]string, state map[string]string) (done bool, updates map[string]string, err error)

	// Build final AuthContext from accumulated state
	BuildAuthContext(ctx context.Context, state map[string]string) (domain.Context, error)
}
