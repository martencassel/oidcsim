package oauth2app

import (
	"context"

	delegationapp "github.com/martencassel/oidcsim/internal/application/delegation"
	dom "github.com/martencassel/oidcsim/internal/domain/oauth2"
)

// uses FlowRegistry to get an AuthorizeFlow

type AuthorizeService struct {
	clients   dom.ClientRepo
	flows     FlowRegistry
	authCodes dom.AuthorizationCodeRepo
}

func NewAuthorizeService(delegationSvc delegationapp.DelegationService, flows FlowRegistry) *AuthorizeService {
	return &AuthorizeService{}
}

// HandleAuthorize handles the OAuth2 /authorize endpoint (GET)
func (s *AuthorizeService) HandleAuthorize(ctx context.Context, req dom.AuthorizeRequest, user dom.User) (string, error) {
	// Consent/delegation check happens here first...
	// Then resolve flow:
	flow := s.flows.Resolve(req.ResponseType)
	if err := flow.Validate(ctx, req); err != nil {
		return "", err
	}
	return flow.Handle(ctx, req, user.ID)
}
