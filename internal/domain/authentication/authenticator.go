package authentication

import "context"

type Authenticator interface {
	Method() AuthMethod
	Authenticate(ctx context.Context, req AuthRequest) (AuthResult, error)
}
