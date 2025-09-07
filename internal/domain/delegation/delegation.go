package delegation

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Delegation struct {
	ID        string
	UserID    string
	ClientID  string
	Scopes    []string
	CreatedAt time.Time
}

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
