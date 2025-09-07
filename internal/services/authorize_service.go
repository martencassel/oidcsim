package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/martencassel/oidcsim/internal/dto"
	"github.com/martencassel/oidcsim/internal/registry"
	"github.com/martencassel/oidcsim/internal/services/authflows"
)

type AuthorizeService interface {
	HandleAuthorize(ctx context.Context, req dto.AuthorizeRequest) (redirectURL string, err error)
}
type authorizeServiceImpl struct {
	registry *registry.Registry[authflows.AuthorizeFlow]
}

func NewAuthorizeServiceImpl(registry *registry.Registry[authflows.AuthorizeFlow]) AuthorizeService {
	return &authorizeServiceImpl{
		registry: registry,
	}
}

func generateAuthCode() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err) // or handle gracefully
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func (s *authorizeServiceImpl) isValidClient(clientID string, redirectURI string) bool {
	return true
}

func (s *authorizeServiceImpl) HandleAuthorize(ctx context.Context, req dto.AuthorizeRequest) (string, error) {
	// // 0. Fail fast if response_type is missing.
	// if req.ResponseType == "" {
	// 	return "", errors.ErrInvalidRequest.WithDescription("missing response_type")
	// }

	// // 1. Load and validate the client.
	// ok := s.isValidClient(req.ClientID, req.RedirectURI)
	// if !ok {
	// 	return "", errors.ErrInvalidClient
	// }

	// // 2. Resolve the flow
	// flow, err := s.registry.Get(req.ResponseType)
	// if err != nil {
	// 	return "", errors.ErrUnsupportedResponseType
	// }

	// // 3. Validate the request
	// err = flow.Validator.Validate(ctx, req, store.Client{ID: req.ClientID})
	// if err != nil {
	// 	return "", err
	// }

	// // 4. Authenticate the user (omitted here, assume user is authenticated)
	// user := store.User{ID: "user123"} // Placeholder for authenticated user

	// // 5. Handle the authorization request
	// redirectURL, err := flow.Handler.Handle(ctx, req, store.Client{ID: req.ClientID}, user)
	// if err != nil {
	// 	return "", err
	// }
	// return redirectURL, nil
	return "", nil
}
