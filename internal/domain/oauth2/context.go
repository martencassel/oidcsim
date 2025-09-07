package oauth2

import (
	"time"
)

type Context struct {
	SubjectID string
	ACR       string
	AMR       []string
	AuthTime  time.Time
	Claims    map[string]interface{}
}

func (c Context) IsZero() bool { return c.SubjectID == "" }

func (c Context) IsValidFor(req AuthorizeRequest) bool {
	if c.IsZero() {
		return false
	}
	if req.RequiredACR != "" && c.ACR != req.RequiredACR {
		return false
	}
	if req.MaxAge > 0 && time.Since(c.AuthTime) > time.Duration(req.MaxAge)*time.Second {
		return false
	}
	return true
}
