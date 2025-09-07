package flows

import (
	"context"
	"errors"
	"time"

	"github.com/martencassel/oidcsim/internal/domain/oauth2"
	"github.com/martencassel/oidcsim/internal/infrastructure/security"
)

type AuthorizationCodeRepo interface {
	Save(ctx context.Context, code oauth2.AuthorizationCode) error
}

type ClientRepo interface {
	Get(ctx context.Context, clientID string) (oauth2.Client, error)
}

type Clock interface {
	Now() int64
	NowTime() time.Time
	After(d time.Duration) <-chan time.Time
}

type CodeFlow struct {
	codes   AuthorizationCodeRepo
	clients ClientRepo
	clock   Clock
}

func NewCodeFlow(codes AuthorizationCodeRepo, clients ClientRepo, clock Clock) *CodeFlow {
	return &CodeFlow{
		codes:   codes,
		clients: clients,
		clock:   clock,
	}
}

func (f *CodeFlow) Validate(ctx context.Context, req oauth2.AuthorizeRequest) error {
	if req.ResponseType != "code" {
		return errors.New("unsupported response_type")
	}
	client, err := f.clients.Get(ctx, req.ClientID)
	if err != nil {
		return err
	}
	if !client.AllowsRedirect(req.RedirectURI) {
		return errors.New("invalid redirect_uri")
	}
	return nil
}

func (f *CodeFlow) Handle(ctx context.Context, req oauth2.AuthorizeRequest, userID string) (string, error) {
	rsg := security.DefaultRandomStringGenerator{}
	params := oauth2.AuthorizationCodeParams{
		Gen:                 rsg,
		SubjectID:           userID,
		ClientID:            req.ClientID,
		RedirectURI:         req.RedirectURI,
		Scope:               parseScopes(req.Scope),
		State:               req.State,
		SessionID:           "",
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
		Nonce:               req.Nonce,
		AuthTime:            f.clock.NowTime(),
		TTL:                 600, // 10 minutes
	}
	code, err := oauth2.NewAuthorizationCodeFromParams(params)
	if err != nil {
		return "", err
	}
	return req.RedirectURIWithParams(map[string]string{
		"code":  code.GetCode(),
		"state": req.State,
		"nonce": req.Nonce,
	}), nil
}

func parseScopes(scopes []string) string {
	if len(scopes) == 0 {
		return "default"
	}
	result := ""
	for i, s := range scopes {
		if i > 0 {
			result += " "
		}
		result += s
	}
	return result
}
