package domain

import (
	"context"
	"strings"
	"time"

	"github.com/martencassel/oidcsim/internal/store"
)

type ClaimsMapper interface {
	MapIDTokenClaims(ctx context.Context, user store.User, client store.Client, scopes []string) map[string]interface{}
	MapAccessTokenClaims(ctx context.Context, user store.User, client store.Client, scopes []string) map[string]interface{}
}

type DefaultClaimsMapper struct{}

var defaultScopeClaims = map[string][]string{
	"openid": {"sub"},
	"profile": {"name", "family_name", "given_name", "middle_name", "nickname",
		"preferred_username", "profile", "picture", "website", "gender",
		"birthdate", "zoneinfo", "locale", "updated_at"},
	"email":   {"email", "email_verified"},
	"address": {"address"},
	"phone":   {"phone_number", "phone_number_verified"},
}

func (m *DefaultClaimsMapper) MapIDTokenClaims(ctx context.Context, user store.User, client store.Client, scopes []string) map[string]interface{} {
	claims := map[string]interface{}{
		"sub": user.ID,
		"aud": client.ID,
		"iat": time.Now().Unix(),
	}
	for _, scope := range scopes {
		for _, claimName := range defaultScopeClaims[scope] {
			if val := m.lookupUserClaim(user, claimName); val != nil {
				claims[claimName] = val
			}
		}
	}
	return claims
}
func (m *DefaultClaimsMapper) lookupUserClaim(user store.User, claimName string) interface{} {
	switch claimName {
	case "sub":
		return user.ID
	case "name":
		return user.FullName
	case "email":
		return user.Email
	case "email_verified":
		return user.EmailVerified
	case "roles":
		return user.Roles
	default:
		return nil
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (m *DefaultClaimsMapper) MapAccessTokenClaims(ctx context.Context, user store.User, client store.Client, scopes []string) map[string]interface{} {
	claims := map[string]interface{}{
		"sub":   user.ID,
		"aud":   client.ResourceServerID,
		"scope": strings.Join(scopes, " "),
		"iat":   time.Now().Unix(),
	}
	if contains(scopes, "roles") {
		claims["roles"] = user.Roles
	}
	return claims
}
