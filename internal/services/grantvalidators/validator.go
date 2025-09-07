package grantvalidators

import (
	"context"

	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/store"
)

type GrantValidator interface {
	GrantType() string
	IsGrantTypeAllowed(client store.Client) bool
	Validate(ctx context.Context, req dto.TokenRequest, client store.Client) error
}
