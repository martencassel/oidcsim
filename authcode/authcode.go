package authcode

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

type Code struct {
	Value       string
	ClientID    string
	RedirectURI string
	Expiry      time.Time
}

type Store struct {
	codes map[string]Code
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewStore(ttl time.Duration) *Store {
	return &Store{
		codes: make(map[string]Code),
		ttl:   ttl,
	}
}

func (s *Store) Generate(clientID, redirectURI string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	code := base64.RawURLEncoding.EncodeToString(b)

	s.mu.Lock()
	defer s.mu.Unlock()
	s.codes[code] = Code{
		Value:       code,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Expiry:      time.Now().Add(s.ttl),
	}

	return code, nil
}

func (s *Store) Validate(code string) (*Code, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.codes[code]
	if !ok {
		return nil, errors.New("code not found")
	}
	if time.Now().After(c.Expiry) {
		return nil, errors.New("code expired")
	}

	return &c, nil
}

func (s *Store) Delete(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.codes, code)
}
