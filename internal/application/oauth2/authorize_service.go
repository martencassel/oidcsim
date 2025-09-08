package oauth2

import (
	"context"

	"github.com/martencassel/oidcsim/internal/domain/authorization"
	dom "github.com/martencassel/oidcsim/internal/domain/oauth2"
	oauth2client "github.com/martencassel/oidcsim/internal/domain/oauth2/client"
)

// uses FlowRegistry to get an AuthorizeFlow

type AuthorizationService interface {
	HandleAuthorize(ctx context.Context, req dom.AuthorizeRequest, user dom.User) (string, error)
}

type AuthorizeServiceImpl struct {
	clients   oauth2client.Repository
	flows     FlowRegistry
	authCodes dom.AuthorizationCodeRepo
}

func NewAuthorizeService(delegationSvc authorization.DelegationService, flows FlowRegistry) *AuthorizeServiceImpl {
	return &AuthorizeServiceImpl{}
}

func (s *AuthorizeServiceImpl) HandleAuthorize(ctx context.Context, req dom.AuthorizeRequest, user dom.User) (string, error) {
	flow := s.flows.Resolve(req.ResponseType)
	if err := flow.Validate(ctx, req); err != nil {
		return "", err
	}
	return flow.Handle(ctx, req, user.ID)
}
