package client

import "context"

type ClientRepository interface {
	GetByID(ctx context.Context, clientID string) (*Client, error)
	ListAll(ctx context.Context) ([]Client, error)
}
