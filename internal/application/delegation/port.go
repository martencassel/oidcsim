package delegationapp

import (
	"context"

	"github.com/martencassel/oidcsim/internal/domain/delegation"
)

type Repository interface {
	FindByUserAndClient(ctx context.Context, userID, clientID string) (*delegation.Delegation, error)
	Save(ctx context.Context, d delegation.Delegation) error
	Delete(ctx context.Context, userID, clientID string) error
}
