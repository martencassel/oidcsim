package granthandlers

import (
	"context"

	"github.com/martencassel/oidcsim/internal/dto"
)

type RefreshTokenHandler struct{}

func (h *RefreshTokenHandler) GrantType() string { return "refresh_token" }

// Exchange an auth code for tokens.

// Get a new access token withouth re-authenticating the user.
func (s *RefreshTokenHandler) Handle(ctx context.Context, req dto.TokenRequest) (dto.TokenResponse, error) {
	// 1. Look up refresh token in storage.
	// 2. Validate
	//    - Token exists, not expired, not revoked.
	//    - Belongs to the client.
	//    - Scope requested is within original scope.
	// 3. Rotate refresh token (if using rotation).
	// 4. Issue new access token (and refresh token if rotated).
	// 5. Return token response.
	return dto.TokenResponse{}, nil
}
