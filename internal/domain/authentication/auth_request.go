package authentication

import "time"

// SubjectID is a stable identifier for an authenticated subject.
type SubjectID string

// AuthMethod is the mechanism used to authenticate.
type AuthMethod string

// AuthResult is the outcome of an successful authentication.
type AuthResult struct {
	Subject         SubjectID
	AssuranceLevel  int      // assurance level
	AuthMethodsUsed []string // methods used
	AuthTime        time.Time
	MFA             bool                   // whether MFA was used
	Claims          map[string]interface{} // identity attributes
}

// AuthRequest is the input to the authentication process.
type AuthRequest struct {
	Method                  AuthMethod
	CallbackData            map[string]string // For post-IdP callback handling
	Credentials             map[string]string
	RequestedAssuranceLevel string
	SessionID               string
}
