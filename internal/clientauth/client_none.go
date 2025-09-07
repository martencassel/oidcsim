package clientauth

import (
	"context"

	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/store"
)

type ClientNone struct{}

func (a *ClientNone) Name() string { return "none" }

func (a *ClientNone) Authenticate(ctx context.Context, client store.Client, req dto.TokenRequest) (*store.Client, error) {
	// req.ClientID already populated by builder
	return &client, nil
}
