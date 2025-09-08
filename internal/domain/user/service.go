package user

import "context"

// UserInfoService defines the domain-level contract for retrieving OIDC claims.
type UserInfoService interface {
	GetUserInfo(ctx context.Context, sub string, scopes []string) (*User, error)
}
