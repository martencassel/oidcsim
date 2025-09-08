package oidc

var standardScopeClaims = map[string][]string{
	"openid":  {"sub"},
	"profile": {"name", "family_name", "given_name", "preferred_username"},
	"email":   {"email", "email_verified"},
}

func ClaimsForScope(scope string) []string {
	return standardScopeClaims[scope]
}
