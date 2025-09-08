package oidc

import (
	"context"

	"github.com/martencassel/oidcsim/internal/domain/oauth2/client"
	user "github.com/martencassel/oidcsim/internal/domain/user"
)

type UserInfoService struct {
	userRepo           user.UserRepository
	clientRepo         client.ClientRepository
	scopeClaimResolver ScopeClaimResolver
}

func (s *UserInfoService) GetUserInfo(ctxContext context.Context, sub string, clientID string, scopes []string) (map[string]any, error) {
	u, _ := s.userRepo.FindByID(ctxContext, sub)
	client, _ := s.clientRepo.GetByID(ctxContext, clientID)
	allowedClaims := s.scopeClaimResolver.ResolveClaims(scopes, client.ID)
	return s.scopeClaimResolver.MapClaims(u, allowedClaims), nil
}
