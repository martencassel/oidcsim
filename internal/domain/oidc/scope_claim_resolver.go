package oidc

import (
	user "github.com/martencassel/oidcsim/internal/domain/user"
)

type ScopeClaimResolver interface {
	ResolveClaims(scopes []string, clientID string) []string
	MapClaims(user *user.User, allowedClaims []string) map[string]any
}

type DefaultScopeClaimResolver struct{}

func (r *DefaultScopeClaimResolver) ResolveClaims(scopes []string, clientID string) []string {
	return []string{}
}

func (r *DefaultScopeClaimResolver) MapClaims(user *user.User, allowed []string) map[string]any {
	claims := map[string]any{}
	for _, claim := range allowed {
		switch claim {
		case "preferred_username":
			claims["username"] = user.PrefferedUsername // renamed
		case "email_verified":
			claims["email_status"] = map[bool]string{true: "verified", false: "unverified"}[user.EmailVerified]
		case "email":
			claims["email"] = user.Email
		case "name":
			claims["name"] = user.Name
			// etc.
		}
	}
	return claims
}
