package infraoauth2

import (
	"context"

	"github.com/martencassel/oidcsim/internal/domain/oauth2"
)

type ClientRepo struct {
}

func (cl *ClientRepo) Get(ctx context.Context, clientID string) (*oauth2.Client, error) {
	return nil, nil
}
