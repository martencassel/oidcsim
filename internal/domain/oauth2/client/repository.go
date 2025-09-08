package client

import "context"

type Repository interface {

	GetByID(ctx context.Context, clientID string) (*Client, error)
	ListAll(ctx context.Context) ([]Client, error)
}
