package oauth2app

import (
	"context"

	dom "github.com/martencassel/oidcsim/internal/domain/oauth2"
)

type AuthorizeHandler interface {
	ResponseType() string
	Handle(ctx context.Context, req dom.AuthorizeRequest, client dom.Client, user dom.User) (string, error)
}
