package client

import (
	"context"
	"fmt"
	"sync"

	domain "github.com/martencassel/oidcsim/internal/domain/oauth2/client"
)

type inMemoryClientRepo struct {
	clients map[string]string
	mu      sync.RWMutex
}

func NewInMemoryClientRepo() *inMemoryClientRepo {
	return &inMemoryClientRepo{
		clients: make(map[string]string),
		mu:      sync.RWMutex{},
	}
}

func (r *inMemoryClientRepo) ListAll(ctx context.Context) ([]domain.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	clients := make([]domain.Client, 0, len(r.clients))
	for id, secret := range r.clients {
		clients = append(clients, domain.Client{ID: id, Secret: secret})
	}
	return clients, nil
}

func (r *inMemoryClientRepo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, id)
	return nil
}

// GetByID retrieves a client by its ID.
func (r *inMemoryClientRepo) GetByID(ctx context.Context, id string) (*domain.Client, error) {
	client, err := r.FindByID(id)
	if err != nil {
		return &domain.Client{}, err
	}
	if client == nil {
		return nil, fmt.Errorf("client with ID %s not found", id)
	}
	return client, nil
}

func (r *inMemoryClientRepo) FindByID(id string) (*domain.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	secret, exists := r.clients[id]
	if !exists {
		return nil, nil
	}
	return &domain.Client{ID: id, Secret: secret}, nil
}

func (r *inMemoryClientRepo) Save(client domain.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client.ID] = client.Secret
	return nil
}

var _ domain.Repository = (*inMemoryClientRepo)(nil)
