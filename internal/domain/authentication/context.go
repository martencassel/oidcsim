package authn

import (
	"time"

	"github.com/martencassel/oidcsim/internal/domain/oauth2"
)

type Context struct {
	SubjectID string // user ID
	AuthTime  time.Time
	Claims    map[string]interface{}
}

// internal/domain/authentication/auth_context.go
func (c Context) IsValidFor(req oauth2.AuthorizeRequest) bool {
	if c.SubjectID == "" {
		return false
	}
	return true
}
