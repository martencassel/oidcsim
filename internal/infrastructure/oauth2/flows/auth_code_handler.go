package flows

import (
	"context"

	appdelegation "github.com/martencassel/oidcsim/internal/application/delegation"
	autzdomain "github.com/martencassel/oidcsim/internal/domain/authorization"
	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/store"
	log "github.com/sirupsen/logrus"
)

type AuthCodeHandler struct {
	svc *appdelegation.DelegationService
}

func (h *AuthCodeHandler) GrantType() string { return "authorization_code" }

// Exchange an auth code for tokens.

func (h *AuthCodeHandler) Handle(ctx context.Context, req dto.TokenRequest, client store.Client) (dto.TokenResponse, error) {
	// Map DTO → domain input
	in := autzdomain.ExchangeInput{
		// CodeID:       autzdomain.AuthorizationCodeID(req.Code),
		// Client:       autzdomain.ClientID(client.ID),
		ClientSecret: req.ClientSecret, // if confidential
		RedirectURI:  req.RedirectURI,
		PKCEVerifier: req.CodeVerifier,
		//	Audience:     autzdomain.Audience("audience"),
	}
	log.Infof("Auth code exchange request: %+v", in)
	// Call domain service
	// res, err := h.svc.ExchangeCodeForTokens(ctx, in)
	// if err != nil {
	// 	return dto.TokenResponse{}, err
	// }
	// Map domain output → DTO
	return dto.TokenResponse{
		// AccessToken:  res.Token.AccessToken,
		// RefreshToken: res.Token.RefreshToken,
		// ExpiresIn:    res.Token.ExpiresIn,
		Scope:     "scope",
		TokenType: "Bearer",
	}, nil
}
