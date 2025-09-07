package authflows

import (
	"context"

	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/errors"
	"github.com/martencassel/oidcsim/internal/store"
)

type TokenIssuer interface {
	IssueAccessToken(ctx context.Context, client store.Client, user store.User, scope []string) (string, error)
}

type tokenIssuerImpl struct{}

func NewTokenIssuerImpl() *tokenIssuerImpl {
	return &tokenIssuerImpl{}
}

func (i *tokenIssuerImpl) IssueAccessToken(ctx context.Context, client store.Client, user store.User, scope []string) (string, error) {
	// In a real implementation, you would create a JWT or opaque token here.
	// For simplicity, we'll return a dummy token.
	return "dummy-access-token", nil
}

type TokenValidator struct{}

func (v *TokenValidator) ResponseType() string { return "token" }
func (v *TokenValidator) Validate(ctx context.Context, req dto.AuthorizeRequest, client store.Client) error {
	if !client.AllowsResponseType("token") {
		return errors.ErrUnauthorizedClient
	}
	return nil
}

type TokenHandler struct {
	TokenIssuer TokenIssuer
}

func (h *TokenHandler) ResponseType() string { return "token" }
func (h *TokenHandler) Handle(ctx context.Context, req dto.AuthorizeRequest, client store.Client, user store.User) (string, error) {
	// accessToken, err := h.TokenIssuer.IssueAccessToken(ctx, client, user, req.Scope)
	// if err != nil {
	// 	return "", err
	// }
	// log.Infof("Issued access token: %s", accessToken)
	return "", nil
	// return buildRedirect(req.RedirectURI, nil, map[string]string{
	// 	"access_token": accessToken,
	// 	"token_type":   "Bearer",
	// 	"expires_in":   "3600",
	// 	"state":        req.State,
	// })
}
