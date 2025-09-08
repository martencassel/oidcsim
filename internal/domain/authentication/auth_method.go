package authentication

import (
	"time"

	"github.com/martencassel/oidcsim/internal/domain/oauth2"
)

type AuthMethod string

const (
	MethodPassword AuthMethod = "password"
	MethodOTP      AuthMethod = "otp"
	MethodWebAuthn AuthMethod = "webauthn"
)

type AuthStep struct {
	Method AuthMethod
}

type AuthFlowSpec struct {
	Steps []AuthStep
}

type StepUI struct {
	Prompt string
	Fields []string
}

type AuthContext struct {
	SubjectID string
	AuthTime  time.Time
}

func (a AuthContext) IsValidFor(req oauth2.AuthorizeRequest) bool {
	if a.SubjectID == "" {
		return false
	}
	return true
}
