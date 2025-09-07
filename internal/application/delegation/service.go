package delegationapp

/*
EnsureConsent:
  1. Look up existing delegation
  2. Apply policy
  3. Create & save if needed

GetDelegation:
  1. Look up by ID
  2. Return or error

RevokeDelegation:
  1. Look up by ID
  2. Delete if exists

*/

import (
	"context"

	dom "github.com/martencassel/oidcsim/internal/domain/delegation"
)

type ConsentResult struct {
	Status       ConsentStatus
	DelegationID string
}

type DelegationService interface {
	// EnsureConsent ensures that a consent/delegation exists for the given user and client with the requested scopes.
	EnsureConsent(ctx context.Context, userID string, clientID string, scopes []string) (ConsentStatus, error)

	// GetDelegation retrieves an existing delegation by its ID.
	GetDelegation(ctx context.Context, delegationID string) (dom.Delegation, error)

	// RevokeDelegation revokes an existing delegation by its ID.
	RevokeDelegation(ctx context.Context, delegationID string) error
}

type ClientMeta struct {
}

type ClientRepo interface {
	GetMeta(ctx context.Context, clientID string) (ClientMeta, error)
}

type delegationServiceImpl struct {
	repo    Repository
	clients ClientRepo
}

// type ConsentResult struct {
// 	Decision   ConsentDecision
// 	Delegation *delegation.Delegation
// }

// func (s *delegationServiceImpl) EnsureConsent(ctx context.Context, userID string, clientID string, scopes []string) (string, error) {
// 	meta, err := s.clients.GetMeta(ctx, clientID)
// 	if err != nil {
// 		return ConsentResult{}, err
// 	}
// 	return "", nil
// }

// func (s *delegationServiceImpl) GetDelegation(ctx context.Context, delegationID string) (dom.Delegation, error) {
// 	return dom.Delegation{}, nil
// }

// func (s *delegationServiceImpl) RevokeDelegation(ctx context.Context, delegationID string) error {
// 	return nil
// }
