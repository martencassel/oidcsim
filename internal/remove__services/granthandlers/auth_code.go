package granthandlers

// internal/services/granthandlers/auth_code.go

// import (
// 	"context"

// 	"github.com/martencassel/oidcsim/internal/application/delegation"
// 	authz "github.com/martencassel/oidcsim/internal/domain/authorization"
// 	"github.com/martencassel/oidcsim/internal/dto"
// 	"github.com/martencassel/oidcsim/internal/store"
// )

// type AuthCodeHandler struct {
// 	svc *delegation.DelegationService
// }

// func (h *AuthCodeHandler) GrantType() string { return "authorization_code" }

// // Exchange an auth code for tokens.

// func (h *AuthCodeHandler) Handle(ctx context.Context, req dto.TokenRequest, client store.Client) (dto.TokenResponse, error) {
// 	// Map DTO → domain input
// 	in := authz.ExchangeInput{
// 		// CodeID:       authz.AuthorizationCodeID(req.Code),
// 		// Client:       authz.ClientID(client.ID),
// 		ClientSecret: req.ClientSecret, // if confidential
// 		RedirectURI:  req.RedirectURI,
// 		PKCEVerifier: req.CodeVerifier,
// 		//Audience:     authz.Audience("audience"),
// 	}
// 	// Call domain service
// 	res, err := h.svc.ExchangeCodeForTokens(ctx, in)
// 	if err != nil {
// 		return dto.TokenResponse{}, err
// 	}
// 	// Map domain output → DTO
// 	return dto.TokenResponse{
// 		AccessToken:  res.Token.AccessToken,
// 		RefreshToken: res.Token.RefreshToken,
// 		ExpiresIn:    res.Token.ExpiresIn,
// 		Scope:        "scope",
// 		TokenType:    "Bearer",
// 	}, nil
// }
