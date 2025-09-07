package authn

import "errors"

var (
	ErrInsufficientACR = errors.New("insufficient ACR level")
	ErrMFARequired     = errors.New("MFA is required")
)

type Policy interface {
	Validate(result AuthResult, requestedACR string) error
}

func NewBasicPolicy(minACR int, requireMFA bool, sessionMaxAge int) Policy {
	return &basicPolicy{
		minACR:        minACR,
		requireMFA:    requireMFA,
		sessionMaxAge: sessionMaxAge,
	}
}

type basicPolicy struct {
	minACR        int
	requireMFA    bool
	sessionMaxAge int
}

func (p *basicPolicy) Validate(result AuthResult, requestedACR string) error {
	// Check ACR level
	if result.AssuranceLevel < p.minACR {
		return ErrInsufficientACR
	}
	// Check MFA requirement
	if p.requireMFA && !result.MFA {
		return ErrMFARequired
	}
	return nil
}

func (p *basicPolicy) SessionMaxAge() int {
	return p.sessionMaxAge
}
