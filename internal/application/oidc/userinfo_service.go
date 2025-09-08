	package oidc

import (
	"context"
	"fmt"

	"github.com/martencassel/oidcsim/internal/interface/http/dto"
)

// Orchestrates token validation + domain call

type UserInfoAppService struct {
	tokenValidator TokenValidator
	userInfoSvc    UserInfoProvider
}

func NewUserInfoAppService(tv TokenValidator, uip UserInfoProvider) *UserInfoAppService {
	return &UserInfoAppService{
		tokenValidator: tv,
		userInfoSvc:    uip,
	}
}

func (a *UserInfoAppService) HandleUserInfo(ctx context.Context, accessToken string) (*dto.UserInfoResponse, error) {
	sub, scopes, err := a.tokenValidator.ValidateAccessToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	user, err := a.userInfoSvc.GetUserInfo(ctx, sub, scopes)
	if err != nil {
		return nil, err
	}
	return &dto.UserInfoResponse{
		Sub:           user.ID,
		Name:          user.Name,
		Email:         user.Email,
		GivenName:     user.GivenName,
		FamilyName:    user.FamilyName,
		EmailVerified: user.EmailVerified,
	}, nil
}
