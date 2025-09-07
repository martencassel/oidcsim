package granthandlers

import (
	"context"

	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/store"
)

type DeviceCodeHandler struct{}

func (h *DeviceCodeHandler) GrantType() string { return "urn:ietf:params:oauth:grant-type:device_code" }

// M2M token issuance
func (h *DeviceCodeHandler) Handle(ctx context.Context, req dto.TokenRequest, client store.Client) (dto.TokenResponse, error) {
	// 1. Validate client.
	//  - Allowed to use client credentials grant.
	//  - Determine allowed scopes.
	// 2. Issue access token (no refresh token).
	// 3. Return token response.
	return dto.TokenResponse{}, nil
}
