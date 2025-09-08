package oauth2

import (
	"context"

	"github.com/martencassel/oidcsim/internal/domain/oauth2"
)

type AuthorizeFlow interface {
	Validate(ctx context.Context, req oauth2.AuthorizeRequest) error
	Handle(ctx context.Context, req oauth2.AuthorizeRequest, userID string) (redirectURI string, err error)
}
