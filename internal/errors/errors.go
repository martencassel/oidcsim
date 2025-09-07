// internal/errors/auth_errors.go
package errors

type AuthError string

func (e AuthError) Error() string {
	return string(e)
}

func (e AuthError) Description() string {
	switch e {
	case ErrInvalidRequest:
		return "The request is missing a required parameter, includes an unsupported parameter value, or is otherwise malformed."
	case ErrUnauthorizedClient:
		return "The client is not authorized to request an authorization code using this method."
	case ErrAccessDenied:
		return "The resource owner or authorization server denied the request."
	case ErrUnsupportedResponseType:
		return "The authorization server does not support obtaining an authorization code using this method."
	case ErrInvalidScope:
		return "The requested scope is invalid, unknown, or malformed."
	case ErrServerError:
		return "The authorization server encountered an unexpected condition that prevented it from fulfilling the request. (This error code is needed because a 500 Internal Server Error HTTP status code cannot be returned to the client via a HTTP redirect.)"
	case ErrTemporarilyUnavailable:
		return "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server. (This error code is needed because a 503 Service Unavailable HTTP status code cannot be returned to the client via a HTTP redirect.)"
	case ErrInvalidClient:
		return "Client authentication failed (e.g., unknown client, no client authentication included, or unsupported authentication method)."
	case ErrInvalidGrant:
		return "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client."
	case ErrUnsupportedGrantType:
		return "The authorization grant type is not supported by the authorization server."
	case ErrInvalidToken:
		return "The access token provided is expired, revoked, malformed, or invalid for other reasons."
	case ErrInsufficientScope:
		return "The request requires higher privileges than provided by the access token."
	case ErrAuthorizationPending:
		return "The authorization request is still pending as the end-user hasn't yet completed the user interaction steps."
	}
	return ""
}

func (e AuthError) WithDescription(desc string) error {
	return &AuthErrorWithDescription{Err: e, DescriptionText: desc}
}

type AuthErrorWithDescription struct {
	Err             error
	DescriptionText string
}

func (e *AuthErrorWithDescription) Error() string {
	return e.Err.Error() + ": " + e.DescriptionText
}

func (e *AuthErrorWithDescription) Unwrap() error {
	return e.Err
}

// ===== Authorization Endpoint Errors (RFC 6749 §4.1.2.1) =====
const (
	ErrInvalidRequest          = AuthError("invalid_request")
	ErrUnauthorizedClient      = AuthError("unauthorized_client")
	ErrAccessDenied            = AuthError("access_denied")
	ErrUnsupportedResponseType = AuthError("unsupported_response_type")
	ErrInvalidScope            = AuthError("invalid_scope")
	ErrServerError             = AuthError("server_error")
	ErrTemporarilyUnavailable  = AuthError("temporarily_unavailable")
)

// ===== Token Endpoint Errors (RFC 6749 §5.2) =====
const (
	ErrInvalidClient        = AuthError("invalid_client")
	ErrInvalidGrant         = AuthError("invalid_grant")
	ErrUnsupportedGrantType = AuthError("unsupported_grant_type")
)

// ===== Resource Server / Introspection Errors (RFC 6750 §3) =====
const (
	ErrInvalidToken      = AuthError("invalid_token")
	ErrInsufficientScope = AuthError("insufficient_scope")
)

// ===== Device Authorization Grant Errors (RFC 8628 §3.5) =====
const (
	ErrAuthorizationPending = AuthError("authorization_pending")
	ErrSlowDown             = AuthError("slow_down")
	ErrExpiredToken         = AuthError("expired_token")
)

// ===== OpenID Connect Specific Errors (OIDC Core §3.1.2.6) =====
const (
	ErrInteractionRequired      = AuthError("interaction_required")
	ErrLoginRequired            = AuthError("login_required")
	ErrAccountSelectionRequired = AuthError("account_selection_required")
	ErrConsentRequired          = AuthError("consent_required")
	ErrInvalidRequestURI        = AuthError("invalid_request_uri")
	ErrInvalidRequestObject     = AuthError("invalid_request_object")
	ErrRequestNotSupported      = AuthError("request_not_supported")
	ErrRequestURINotSupported   = AuthError("request_uri_not_supported")
	ErrRegistrationNotSupported = AuthError("registration_not_supported")
)
