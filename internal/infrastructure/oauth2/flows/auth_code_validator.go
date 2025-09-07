package infraoauth2flows

import (
	"context"

	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/errors"
	"github.com/martencassel/oidcsim/internal/store"
)

type AuthCodeValidator struct{}

func (v *AuthCodeValidator) GrantType() string { return "authorization_code" }
func (v *AuthCodeValidator) IsGrantTypeAllowed(client store.Client) bool {
	return client.AllowsGrantType(v.GrantType())
}
func (v *AuthCodeValidator) Validate(ctx context.Context, req dto.TokenRequest, client store.Client) error {
	if req.Code == "" || req.RedirectURI == "" {
		return errors.ErrInvalidRequest.WithDescription("missing code or redirect_uri")
	}
	return nil
}
