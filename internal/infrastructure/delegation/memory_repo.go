package delegation

import (
	"context"
	"sync"

	delegationapp "github.com/martencassel/oidcsim/internal/application/delegation"
	"github.com/martencassel/oidcsim/internal/domain/delegation"
)

type MemoryRepo struct {
	mu   sync.RWMutex
	data map[string]delegation.Delegation // key: userID|clientID
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		data: make(map[string]delegation.Delegation),
	}
}

func (r *MemoryRepo) FindByUserAndClient(_ context.Context, userID, clientID string) (*delegation.Delegation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if d, ok := r.data[userID+"|"+clientID]; ok {
		return &d, nil
	}
	return nil, nil
}

func (r *MemoryRepo) Save(_ context.Context, d delegation.Delegation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[d.UserID+"|"+d.ClientID] = d
	return nil
}

func (r *MemoryRepo) Delete(_ context.Context, userID, clientID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, userID+"|"+clientID)
	return nil
}

// FindByID implements the Repository interface.
// Assumes the ID is in the format "userID|clientID".
func (r *MemoryRepo) FindByID(_ context.Context, id string) (*delegation.Delegation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if d, ok := r.data[id]; ok {
		return &d, nil
	}
	return nil, nil
}

// Ensure interface compliance
var _ delegationapp.Repository = (*MemoryRepo)(nil)
