package infraoauth2

import "context"

type JWTTokenService struct{}

func (ts *JWTTokenService) MintAccessToken(ctx context.Context, subjectID, clientID string, scopes []string) (string, error) {
	return "", nil
}

func (ts *JWTTokenService) MintIDToken(ctx context.Context, subjectID, clientID string, nonce string) (string, error) {
	return "", nil
}
