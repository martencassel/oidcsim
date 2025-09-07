package store

import (
	"context"
	"fmt"
	"time"
)

type Nonce string

type CodeChallenge string

type CodeChallengeMethod string

const (
	CodeChallengePlain CodeChallengeMethod = "plain"
	CodeChallengeS256  CodeChallengeMethod = "S256"
)

// Snapshot of all important facts from the /authorize request (plus some server context) you need at /token
type AuthorizationCode struct {
	Code        string
	ClientID    string   // Extract value from AuthorizationRequest parameter
	RedirectURI string   // Extract value from AuthorizationRequest parameter
	Scope       []string // Extract value from AuthorizationRequest parameter
	State       string   // Echoed back immediately in /authorize response; not needed at /token.

	UserID    string // foreign key to UserStore
	SessionID string // foreign key to SessionStore
	AuthTime  time.Time

	CodeChallenge       CodeChallenge
	CodeChallengeMethod CodeChallengeMethod
	Nonce               Nonce

	ExpiresAt time.Time

	Used bool
}


type InMemoryCodeStore struct {
	codes map[string]*AuthorizationCode
}

func NewInMemoryCodeStore() *InMemoryCodeStore {
	return &InMemoryCodeStore{codes: make(map[string]*AuthorizationCode)}
}

func (s *InMemoryCodeStore) Save(ctx context.Context, code *AuthorizationCode) error {
	s.codes[code.Code] = code
	return nil
}

func (s *InMemoryCodeStore) Get(ctx context.Context, code string) (*AuthorizationCode, error) {
	c, ok := s.codes[code]
	if !ok {
		return nil, fmt.Errorf("code not found")
	}
	return c, nil
}

func (s *InMemoryCodeStore) Delete(ctx context.Context, code string) error {
	delete(s.codes, code)
	return nil
}


