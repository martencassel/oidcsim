package dto

// ErrorResponse represents a standard error response in JSON format.
// See: https://openid.net/specs/openid-connect-core-1_0.html#AuthError
type ErrorResponse struct {
	Error            string `json:"error"`                       // e.g. "invalid_request"
	ErrorDescription string `json:"error_description,omitempty"` // Human-readable explanation
	ErrorURI         string `json:"error_uri,omitempty"`         // Link to documentation
}
