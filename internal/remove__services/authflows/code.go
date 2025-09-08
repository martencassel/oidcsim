package authflows

// import (
// 	"context"
// 	"net/url"
// 	"time"

// 	"github.com/martencassel/oidcsim/internal/dto"
// 	"github.com/martencassel/oidcsim/internal/errors"
// 	infrasec "github.com/martencassel/oidcsim/internal/infrastructure/security"
// 	"github.com/martencassel/oidcsim/internal/store"
// )

// type CodeValidator struct{}

// func (v *CodeValidator) ResponseType() string { return "code" }

// func (v *CodeValidator) Validate(ctx context.Context, req dto.AuthorizeRequest, client store.Client) error {
// 	// if req.RedirectURI == "" {
// 	// 	return errors.ErrInvalidRequest.WithDescription("missing redirect_uri")
// 	// }
// 	// if !client.IsRedirectURIMatching(req.RedirectURI) {
// 	// 	return errors.ErrInvalidRequest.WithDescription("redirect_uri mismatch")
// 	// }
// 	if !client.AllowsResponseType("code") {
// 		return errors.ErrUnauthorizedClient
// 	}
// 	// if len(req.Scope) == 0 {
// 	// 	return errors.ErrInvalidRequest.WithDescription("missing scope")
// 	// }
// 	return nil
// }

// func generateAuthCode(req dto.AuthorizeRequest, client store.Client, user store.User) *store.AuthorizationCode {
// 	code, _ := infrasec.GenerateRandomString(32)
// 	return &store.AuthorizationCode{
// 		Code:     code,
// 		ClientID: client.ID,
// 		// RedirectURI: req.RedirectURI,
// 		// Scope:       req.Scope,
// 		// State:       req.State,
// 		UserID: user.ID,
// 	}
// }

// type CodeHandler struct {
// 	TTL time.Duration
// 	//CodeStore store.CodeStore
// }

// func (h *CodeHandler) ResponseType() string { return "code" }

// func (h *CodeHandler) Handle(ctx context.Context, req dto.AuthorizeRequest, client store.Client, user store.User) (string, error) {
// 	// // 1. Build code
// 	// code, err := oauth2.BuildAuthorizationCode(req, client, user, h.TTL)
// 	// if err != nil {
// 	// 	return "", err
// 	// }
// 	// // 2. Save code
// 	// if err := h.CodeStore.Save(ctx, code); err != nil {
// 	// 	return "", err
// 	// }
// 	// // 3. Build redirect URI
// 	// redirectURL, err := buildRedirect(req.RedirectURI, map[string]string{
// 	// 	"code":  code.Code,
// 	// 	"state": req.State,
// 	// }, nil)
// 	// if err != nil {
// 	// 	return "", err
// 	// }
// 	// return redirectURL, nil
// 	return "", nil
// }

// func buildRedirect(base string, queryParams, fragmentParams map[string]string) (string, error) {
// 	u, err := url.Parse(base)
// 	if err != nil {
// 		return "", err
// 	}
// 	if len(queryParams) > 0 {
// 		q := u.Query()
// 		for k, v := range queryParams {
// 			if v != "" {
// 				q.Set(k, v)
// 			}
// 		}
// 		u.RawQuery = q.Encode()
// 	}
// 	if len(fragmentParams) > 0 {
// 		f := url.Values{}
// 		for k, v := range fragmentParams {
// 			if v != "" {
// 				f.Set(k, v)
// 			}
// 		}
// 		u.Fragment = f.Encode()
// 	}
// 	return u.String(), nil
// }

// // buildRedirectURI appends the given params to the redirect URI and returns the full string.
// func buildRedirectURI(base string, params map[string]string) (string, error) {
// 	u, err := url.Parse(base)
// 	if err != nil {
// 		return "", err
// 	}
// 	q := u.Query()
// 	for k, v := range params {
// 		if v != "" {
// 			q.Set(k, v)
// 		}
// 	}
// 	u.RawQuery = q.Encode()
// 	return u.String(), nil
// }
