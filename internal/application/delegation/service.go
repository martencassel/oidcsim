package delegation

import (
	"context"
	"fmt"
	"time"

	delegation "github.com/martencassel/oidcsim/internal/domain/delegation"
)

type DelegationService interface {
	// EnsureConsent ensures that a consent/delegation exists for the given user and client with the requested scopes.
	EnsureConsent(ctx context.Context, userID string, clientID string, scopes []string) (*delegation.ConsentResult, error)

	// GetDelegation retrieves an existing delegation by its ID.
	GetDelegation(ctx context.Context, delegationID string) (delegation.Delegation, error)

	// RevokeDelegation revokes an existing delegation by its ID.
	RevokeDelegation(ctx context.Context, delegationID string) error

	// ValidateDelegationForRefresh checks if the delegation is valid for use in a refresh token flow.
	ValidateDelegationForRefresh(ctx context.Context, delegationID string) error
}

type delegationServiceImpl struct {
	repo Repository
}

func NewDelegationService(repo Repository) DelegationService {
	return &delegationServiceImpl{
		repo: repo,
	}
}

// EnsureConsent checks whether the user has already granted consent to the client for the requested scopes.
//
// CURRENT BEHAVIOR:
// - Consent is auto-approved for all clients and scopes.
// - A new Delegation is created and persisted unconditionally.
//
// FUTURE EXTENSIONS:
// - Lookup existing Delegation and check if it covers requested scopes.
// - Apply consent policy (e.g. trusted clients, sensitive scopes, prompt=none).
// - Redirect to consent UI if required.
//
// This method is called during the /authorize flow after authentication is confirmed.
func (s *delegationServiceImpl) EnsureConsent(ctx context.Context, userID string, clientID string, scopes []string) (*delegation.ConsentResult, error) {
	// Always auto-approve for now
	d, err := delegation.NewDelegation(userID, clientID, scopes)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, d); err != nil {
		return nil, err
	}
	return &delegation.ConsentResult{
		Decision:     delegation.ConsentStatusGranted,
		DelegationId: d.ID,
	}, nil
}

// GetDelegation retrieves a Delegation by its ID.
//
// Used for:
// - Inspecting previously granted consents.
// - Displaying connected apps in a user dashboard.
// - Validating delegation status during token issuance.
//
// Returns an error if the delegation is not found or if the repository fails.
func (s *delegationServiceImpl) GetDelegation(ctx context.Context, delegationID string) (delegation.Delegation, error) {
	d, err := s.repo.FindByID(ctx, delegationID)
	if err != nil {
		return delegation.Delegation{}, err
	}
	if d == nil {
		return delegation.Delegation{}, fmt.Errorf("delegation not found: %s", delegationID)
	}
	return *d, nil
}

// RevokeDelegation deletes a Delegation by its ID.
//
// Use cases:
// - User-initiated revocation (e.g. "disconnect this app" from a dashboard).
// - Admin-triggered revocation (e.g. suspicious activity, scope abuse).
// - Expiry or rotation of Delegations (e.g. time-based invalidation).
//
// This method removes the delegation from the repository, effectively revoking the client's access.

func (s *delegationServiceImpl) RevokeDelegation(ctx context.Context, delegationID string) error {
	delegation, err := s.repo.FindByID(ctx, delegationID)
	if err != nil {
		return err
	}
	if delegation == nil {
		return fmt.Errorf("delegation not found: %s", delegationID)
	}
	return s.repo.Delete(ctx, delegation.ClientID, delegation.ID)
}

// ValidateDelegationForRefresh ensures that a refresh token is still backed by a valid Delegation.
//
// Use cases:
// - Enforces ongoing user consent when clients use refresh tokens.
// - Prevents token renewal if the Delegation has been revoked or expired.
// - Supports dynamic consent revocation (e.g. user disconnects app, admin disables access).
//
// This method is typically called during the refresh token grant flow.
// If the Delegation is missing, revoked, or expired, the refresh request should be denied.
func (s *delegationServiceImpl) ValidateDelegationForRefresh(ctx context.Context, delegationID string) error {
	d, err := s.repo.FindByID(ctx, delegationID)
	if err != nil {
		return err
	}
	if d == nil {
		return fmt.Errorf("delegation not found: %s", delegationID)
	}
	if d.IsRevoked() {
		return fmt.Errorf("delegation revoked")
	}
	if d.IsExpired(time.Now()) {
		return fmt.Errorf("delegation expired")
	}
	// 4. Valid
	return nil
}
