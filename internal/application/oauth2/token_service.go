package oauth2

import (
	"context"

	"github.com/martencassel/oidcsim/internal/application/delegation"
	"github.com/martencassel/oidcsim/internal/domain/oauth2/client"
)

type TokenService struct {
	codeRepo      CodeRepository
	clientRepo    client.Repository
	delegationSvc delegation.DelegationService
	jwtIssuer     JWTIssuer
	refreshRepo   RefreshTokenRepository
}

func NewTokenService(codeRepo CodeRepository, clientRepo client.Repository, delegationSvc delegation.DelegationService, jwtIssuer JWTIssuer, refreshRepo RefreshTokenRepository) *TokenService {
	return &TokenService{
		codeRepo:      codeRepo,
		clientRepo:    clientRepo,
		delegationSvc: delegationSvc,
		jwtIssuer:     jwtIssuer,
		refreshRepo:   refreshRepo,
	}
}

func (s *TokenService) ExchangeCodeForTokens(ctx context.Context, req TokenRequest) (*TokenResponse, error) {
	// 1. Validate code
	// 2. Validate client
	// 3. Validate redirect URI and PKCE
	// 4. Validate delegation
	// 5. Issue access + refresh tokens
	return &TokenResponse{}, nil
}
