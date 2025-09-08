package delegation

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type DelegationID string
type ClientID string
type SubjectID string
type Audience string
type Scope string

type Delegation struct {
	ID        string
	UserID    string
	ClientID  string
	Scopes    []string
	CreatedAt time.Time
	RevokedAt *time.Time
	ExpiresAt *time.Time
}

// Future
// ExpiresAt   time.Time       // Optional expiry for time-limited consent
// RevokedAt   *time.Time      // If revoked, timestamp of revocation
// Claims      map[string]any  // Optional claims granted (e.g. email, profile)
// Remember    bool            // Whether user chose "remember this decision"
// PromptedAt  time.Time       // When the user was last shown a consent screen

func NewDelegation(userID, clientID string, scopes []string) (Delegation, error) {
	if userID == "" || clientID == "" {
		return Delegation{}, errors.New("missing user or client")
	}
	if len(scopes) == 0 {
		return Delegation{}, errors.New("no scopes granted")
	}
	return Delegation{
		ID:        uuid.New().String(),
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    scopes,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (d Delegation) IsRevoked() bool {
	return d.RevokedAt != nil
}

func (d Delegation) IsExpired(now time.Time) bool {
	return d.ExpiresAt != nil && now.After(*d.ExpiresAt)
}
