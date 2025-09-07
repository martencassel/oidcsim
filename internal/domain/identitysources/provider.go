package identitysources

import "context"

// IdentityProvider is the core interface for any identity source.
type IdentityProvider interface {
	AuthenticatePassword(ctx context.Context, username, password string) (SubjectID, error)
	AuthenticateExternal(ctx context.Context, assertion interface{}) (SubjectID, error)
	GetClaims(ctx context.Context, subjectID SubjectID) (map[string]interface{}, error)
}

// ClaimsProvider is a narrower interface for fetching claims only.
type ClaimsProvider interface {
	GetClaims(ctx context.Context, subjectID SubjectID) (map[string]interface{}, error)
}
