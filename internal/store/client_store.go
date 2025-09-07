package store

import (
	"context"
	"database/sql"
	"fmt"
)

type Client struct {
	ID               string
	Secret           string
	ResourceServerID string // for access tokens
	Name             string
	RedirectURIs     []string
	AuthMethod       string   // e.g. "client_secret_basic"
	Grants           []string // allowed grant types
	Scopes           []string // allowed scopes

	Public bool
}

func (c Client) AllowsResponseType(responseType string) bool {
	switch responseType {
	case "code":
		return c.AllowsGrantType("authorization_code")
	case "token":
		return c.AllowsGrantType("implicit")
	case "id_token":
		return c.AllowsGrantType("implicit")
	case "code token":
		return c.AllowsGrantType("authorization_code") && c.AllowsGrantType("implicit")
	case "code id_token":
		return c.AllowsGrantType("authorization_code") && c.AllowsGrantType("implicit")
	case "id_token token":
		return c.AllowsGrantType("implicit")
	case "code id_token token":
		return c.AllowsGrantType("authorization_code") && c.AllowsGrantType("implicit")
	default:
		return false
	}
}

func (c Client) IsRedirectURIMatching(uri string) bool {
	for _, r := range c.RedirectURIs {
		if r == uri {
			return true
		}
	}
	return false
}

func (c Client) AllowsGrantType(grant string) bool {
	for _, g := range c.Grants {
		if g == grant {
			return true
		}
	}
	return false
}

type ClientStore interface {
	GetByID(ctx context.Context, id string) (Client, error)
	Save(ctx context.Context, client Client) error
	List(ctx context.Context) ([]Client, error)
}

type InMemoryClientStore struct {
	clients map[string]Client
}

func NewInMemoryClientStore() *InMemoryClientStore {
	return &InMemoryClientStore{clients: make(map[string]Client)}
}

func (s *InMemoryClientStore) GetByID(ctx context.Context, id string) (Client, error) {
	c, ok := s.clients[id]
	if !ok {
		return Client{}, fmt.Errorf("client not found")
	}
	return c, nil
}

func (s *InMemoryClientStore) List(ctx context.Context) ([]Client, error) {
	clients := make([]Client, 0, len(s.clients))
	for _, c := range s.clients {
		clients = append(clients, c)
	}
	return clients, nil
}

func (s *InMemoryClientStore) Save(ctx context.Context, client Client) error {
	s.clients[client.ID] = client
	return nil
}

type SQLClientStore struct {
	db *sql.DB
}

func (s *SQLClientStore) GetByID(ctx context.Context, id string) (Client, error) {
	// SELECT ... FROM clients WHERE id = ?
	return Client{}, nil
}

func (s *SQLClientStore) List(ctx context.Context) ([]Client, error) {
	// SELECT ... FROM clients
	return nil, nil
}
