package oauth2

import (
	"context"
	"time"

	authz "github.com/martencassel/oidcsim/internal/domain/authorization"
	infrasec "github.com/martencassel/oidcsim/internal/infrastructure/security"
)

type AuthorizationCodeID string

type AuthorizationCode struct {
	// Core Identifiers
	ID           AuthorizationCodeID
	DelegationID authz.DelegationID // Link to consent/delegation record
	ClientID     authz.ClientID

	// Protocol parameters
	RedirectURI         string
	Scope               string // space-delimited scopes granted
	State               string // opaque value from client
	Nonce               string // OIDC nonce
	PKCEChallenge       string // PKCE code challenge
	PKCEChallengeMethod string // S256 or plain

	// Lifecycle
	ExpiresAt         time.Time
	UsedAt            time.Time
	RevokedAt         time.Time
	InvalidatedReason string // if set, the code was invalidated for this reason

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c AuthorizationCode) IsExpired(at time.Time) bool {
	return at.After(c.ExpiresAt)
}

func (c AuthorizationCode) IsUsed() bool {
	return !c.UsedAt.IsZero()
}

func (c AuthorizationCode) GetCode() string {
	return string(c.ID)
}

type AuthorizationCodeRepo interface {
	Issue(ctx context.Context, code AuthorizationCode) error
	GetForRedemption(ctx context.Context, id AuthorizationCodeID) (AuthorizationCode, error)
	MarkUsed(ctx context.Context, id AuthorizationCodeID, at time.Time) error
}

func NewAuthorizationCode(gen infrasec.RandomStringGenerator, clientID, redirectURI, scope, state, subjectID, sessionID, codeChallenge, codeChallengeMethod, nonce string, authTime time.Time, ttl time.Duration) (AuthorizationCode, error) {
	if gen == nil {
		gen = infrasec.DefaultRandomStringGenerator{}
	}
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	if authTime.IsZero() {
		authTime = time.Now()
	}
	r, err := gen.Generate(32)
	if err != nil {
		return AuthorizationCode{}, err
	}
	code := AuthorizationCode{
		ID:                  AuthorizationCodeID(r),
		DelegationID:        authz.DelegationID(""), // to be set when saving
		ClientID:            authz.ClientID(clientID),
		RedirectURI:         redirectURI,
		Scope:               scope,
		State:               state,
		Nonce:               nonce,
		PKCEChallenge:       codeChallenge,
		PKCEChallengeMethod: codeChallengeMethod,
		ExpiresAt:           time.Now().Add(ttl),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	return code, nil
}

// Revoke(id AuthorizationCode, reason string) error
// GetByCode(ctx context.Context, code string) (AuthorizationCode, error)
// ListByDelegation(ctx context.Context, delegationID authz.DelegationID) ([]AuthorizationCode, error)
// ListByClient(ctx context.Context, clientID authz.ClientID) ([]AuthorizationCode, error)
// Redeem(ctx context.Context, code string, at time.Time) (AuthorizationCode, error)
