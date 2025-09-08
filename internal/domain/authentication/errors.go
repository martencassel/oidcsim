package authentication

// Domain Errors for authentication package
import "errors"

var (
	ErrNoAuthenticatorFound = errors.New("no authenticator found for the given method")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrPolicyViolation      = errors.New("authentication result does not satisfy the policy")
)
var ErrUnsupportedAuthMethod = errors.New("unsupported authentication method")
var ErrAuthenticationFailed = errors.New("authentication failed")
var ErrInvalidAuthRequest = errors.New("invalid authentication request")
var ErrEmptyAuthMethod = errors.New("authentication method is empty")
var ErrEmptyCredentials = errors.New("credentials are empty")
var ErrUnsupportedACR = errors.New("requested ACR is not supported by the authentication result")
