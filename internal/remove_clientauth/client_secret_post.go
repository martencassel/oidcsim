package clientauth

// import (
// 	"context"

// 	"github.com/martencassel/oidcsim/internal/dto"
// 	"github.com/martencassel/oidcsim/internal/errors"
// 	"github.com/martencassel/oidcsim/internal/store"
// )

// type ClientSecretPost struct{}

// func (a *ClientSecretPost) Name() string { return "client_secret_post" }

// func (a *ClientSecretPost) Authenticate(ctx context.Context, client store.Client, req dto.TokenRequest) (*store.Client, error) {
// 	// req.ClientID and req.ClientSecret already populated by builder
// 	if !checkSecret(req.ClientID, req.ClientSecret) {
// 		return nil, errors.ErrInvalidClient
// 	}
// 	return &client, nil
// }
// func checkSecret(clientID, clientSecret string) bool {
// 	return clientID == "client" && clientSecret == "secret"
// }
