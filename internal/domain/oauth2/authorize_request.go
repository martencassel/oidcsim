package oauth2

import "errors"

// Encapsulate the semantic concept of an authorization request inside your business language.
//
// No HTTP tags, no direct coupling to transport.

type AuthorizeRequest struct {
	ResponseType        string
	ClientID            string
	RedirectURI         string
	Scope               []string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
	Nonce               string
	// Extra
	RequiredACR string
	MaxAge      int64
}

func NewAuthorizeRequest(
	responseType,
	clientID,
	redirectURI,
	scope,
	state,
	codeChallenge,
	codeChallengeMethod,
	nonce string,
) (AuthorizeRequest, error) {
	// Invariants ...
	if responseType == "" {
		return AuthorizeRequest{}, errors.New("missing response_type")
	}
	var scopeList []string
	if scope != "" {
		scopeList = append(scopeList, scope)
	}
	return AuthorizeRequest{
		ResponseType:        responseType,
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		Scope:               scopeList,
		State:               state,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		Nonce:               nonce,
	}, nil
}

func (r *AuthorizeRequest) RedirectURIWithParams(params map[string]string) string {
	// Build the redirect URI with query parameters
	redirectURI := r.RedirectURI
	if len(params) > 0 {
		// Append parameters to the redirect URI
		q := "?"
		for k, v := range params {
			if q != "?" {
				q += "&"
			}
			q += k + "=" + v
		}
		redirectURI += q
	}
	return redirectURI
}

// Helper methods for validation

func (r AuthorizeRequest) IsResponseTypeEmpty() bool {
	return r.ResponseType == ""
}

func (r AuthorizeRequest) IsClientIDEmpty() bool {
	return r.ClientID == ""
}

func (r AuthorizeRequest) IsRedirectURIEmpty() bool {
	return r.RedirectURI == ""
}

func (r AuthorizeRequest) IsScopeEmpty() bool {
	return len(r.Scope) == 0
}

func (r AuthorizeRequest) IsCodeChallengeProvided() bool {
	return r.CodeChallenge != ""
}

func (r AuthorizeRequest) IsCodeChallengeMethodProvided() bool {
	return r.CodeChallengeMethod != ""
}

func (r AuthorizeRequest) IsNonceProvided() bool {
	return r.Nonce != ""
}
