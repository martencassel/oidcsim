package grantvalidators

import (
	"context"

	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/errors"
	"github.com/martencassel/oidcsim/internal/store"
)

type RefreshTokenValidator struct{}

func (v *RefreshTokenValidator) GrantType() string { return "refresh_token" }
func (v *RefreshTokenValidator) IsGrantTypeAllowed(client store.Client) bool {
	return client.AllowsGrantType(v.GrantType())
}

func (v *RefreshTokenValidator) Validate(ctx context.Context, req dto.TokenRequest, client store.Client) error {
	if req.RefreshToken == "" {
		return errors.ErrInvalidRequest.WithDescription("missing refresh_token")
	}
	return nil
}
