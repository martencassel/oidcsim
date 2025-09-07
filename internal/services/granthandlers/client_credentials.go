package granthandlers

import (
	"context"

	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/store"
)

type ClientCredentialsHandler struct{}

func (h *ClientCredentialsHandler) GrantType() string { return "client_credential" }

// M2M token issuance.
func (h *ClientCredentialsHandler) Handle(ctx context.Context, req dto.TokenRequest, client store.Client) (dto.TokenResponse, error) {
	return dto.TokenResponse{}, nil
}
