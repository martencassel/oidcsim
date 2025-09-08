package oauth2

import (
	"context"

	"github.com/martencassel/oidcsim/internal/domain/oauth2"
	"github.com/martencassel/oidcsim/internal/store"
)

type AuthorizeValidator interface {
	ResponseType() string
	Validate(ctx context.Context, req oauth2.AuthorizeRequest, client store.Client) error
}
