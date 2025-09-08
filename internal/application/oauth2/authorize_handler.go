package oauth2

import (
	"context"

	"github.com/martencassel/oidcsim/internal/domain/oauth2"
	oauth2dom "github.com/martencassel/oidcsim/internal/domain/oauth2"
	dom "github.com/martencassel/oidcsim/internal/domain/oauth2/client"
)

type AuthorizeHandler interface {
	ResponseType() string
	Handle(ctx context.Context, req oauth2dom.AuthorizeRequest, client dom.Client, user oauth2.User) (string, error)
}
