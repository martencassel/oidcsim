package dto

type ErrorResponse struct {
	Error            string `json:"error"`                       // e.g. "invalid_request"
	ErrorDescription string `json:"error_description,omitempty"` // Human-readable explanation
	ErrorURI         string `json:"error_uri,omitempty"`         // Link to documentation
}
