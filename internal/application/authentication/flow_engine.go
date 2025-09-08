package authentication

// import (
// 	"context"

// 	domain "github.com/martencassel/oidcsim/internal/domain/authentication"
// )

// type FlowEngine interface {
// 	// Given clientID, return the flow (e.g., pwd → totp)
// 	Plan(ctx context.Context, clientID string) (domain.FlowSpec, error)

// 	// Render/init a step: returns hints and transient state updates
// 	StartStep(ctx context.Context, step domain.StepSpec, state map[string]string) (map[string]string, error)

// 	// Complete a step with inputs; returns whether more steps remain and state updates
// 	CompleteStep(ctx context.Context, step domain.StepSpec, inputs map[string]string, state map[string]string) (done bool, updates map[string]string, err error)

// 	// Build final AuthContext from accumulated state
// 	BuildAuthContext(ctx context.Context, state map[string]string) (domain.Context, error)
// }
