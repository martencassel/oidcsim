package oidc

import (
	"context"

	user "github.com/martencassel/oidcsim/internal/domain/user"
)

// TokenValidator defines the interface for validating tokens.
type TokenValidator interface {
	ValidateAccessToken(ctxContext context.Context, token string) (sub string, scopes []string, err error)
}

// UserInfoProvider defines the interface for retrieving user information based on a token.
type UserInfoProvider interface {
	GetUserInfo(ctx context.Context, token string, scopes []string) (*user.User, error)
}
