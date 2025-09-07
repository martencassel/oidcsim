package clientauth

import (
	"context"

	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/store"
)

type ClientPrivateJWT struct{}

func (a *ClientPrivateJWT) Name() string { return "private_key_jwt" }

func (a *ClientPrivateJWT) Authenticate(ctx context.Context, client store.Client, req dto.TokenRequest) (*store.Client, error) {
	return &client, nil
}
