package services

// import (
// 	"context"
// 	"strings"

// 	"github.com/martencassel/oidcsim/internal/clientauth"
// 	"github.com/martencassel/oidcsim/internal/dto"
// 	"github.com/martencassel/oidcsim/internal/errors"
// 	"github.com/martencassel/oidcsim/internal/registry"
// 	"github.com/martencassel/oidcsim/internal/services/grantflows"
// 	"github.com/martencassel/oidcsim/internal/services/granthandlers"
// 	"github.com/martencassel/oidcsim/internal/services/grantvalidators"
// 	"github.com/martencassel/oidcsim/internal/store"
// )

// type TokenServiceDeps struct {
// 	AuthRegistry      *registry.Registry[clientauth.Authenticator]
// 	GrantValidatorReg *registry.Registry[grantvalidators.GrantValidator]
// 	GrantHandlerReg   *registry.Registry[granthandlers.GrantHandler]
// 	GrantFlowReg      *registry.Registry[grantflows.GrantFlow]
// 	ClientStore       store.ClientStore
// }

// type TokenService interface {
// 	IssueToken(ctx context.Context, req dto.TokenRequest) (dto.TokenResponse, error)
// }

// type tokenServiceImpl struct {
// 	authRegistry           *registry.Registry[clientauth.Authenticator]
// 	grantValidatorRegistry *registry.Registry[grantvalidators.GrantValidator]
// 	grantHandlerRegistry   *registry.Registry[granthandlers.GrantHandler]
// 	grantFlowRegistry      *registry.Registry[grantflows.GrantFlow]
// 	clientStore            store.ClientStore
// }

// func NewTokenServiceImpl(deps TokenServiceDeps) TokenService {
// 	return &tokenServiceImpl{
// 		authRegistry:           deps.AuthRegistry,
// 		grantValidatorRegistry: deps.GrantValidatorReg,
// 		grantHandlerRegistry:   deps.GrantHandlerReg,
// 		grantFlowRegistry:      deps.GrantFlowReg,
// 		clientStore:            deps.ClientStore,
// 	}
// }

// func (s *tokenServiceImpl) authenticateClient(ctx context.Context, req *dto.TokenRequest) (*store.Client, error) {
// 	client, err := s.clientStore.GetByID(ctx, req.ClientID)
// 	if err != nil {
// 		return nil, errors.ErrInvalidClient
// 	}
// 	authenticator, err := s.authRegistry.Get(client.AuthMethod)
// 	if err != nil {
// 		return nil, errors.ErrUnauthorizedClient
// 	}
// 	return authenticator.Authenticate(ctx, client, *req)
// }

// // IssueToken processes a token request and issues a token response.
// func (s *tokenServiceImpl) IssueToken(ctx context.Context, req dto.TokenRequest) (dto.TokenResponse, error) {
// 	// 0. Fail fast if grant_type is missing
// 	if strings.TrimSpace(req.GrantType) == "" {
// 		return dto.TokenResponse{}, errors.ErrInvalidRequest.WithDescription("missing grant_type")
// 	}

// 	// 1. Authenticate client
// 	client, err := s.authenticateClient(ctx, &req)
// 	if err != nil {
// 		return dto.TokenResponse{}, err
// 	}
// 	if !client.AllowsGrantType(req.GrantType) {
// 		return dto.TokenResponse{}, errors.ErrUnauthorizedClient
// 	}

// 	// 2. Get flow and validate
// 	flow, err := s.grantFlowRegistry.Get(req.GrantType)
// 	if err != nil {
// 		return dto.TokenResponse{}, errors.ErrUnsupportedGrantType
// 	}
// 	if err := flow.Validator.Validate(ctx, req, *client); err != nil {
// 		return dto.TokenResponse{}, err
// 	}

// 	// 3. Handle it
// 	return flow.Handler.Handle(ctx, req, *client)
// }
