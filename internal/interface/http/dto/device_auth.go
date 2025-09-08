package dto

// Device Authorization (RFC 8628).

// DeviceAuthorizationRequest represents the parameters for a device authorization request.
// See: https://datatracker.ietf.org/doc/html/rfc8628#section-3.1
type DeviceAuthorizationRequest struct {
	ClientID string `form:"client_id" binding:"required"`
	Scope    string `form:"scope,omitempty"`
}

// DeviceAuthorizationResponse represents a successful response from the device authorization endpoint.
// See: https://datatracker.ietf.org/doc/html/rfc8628#section-3.2
type DeviceAuthorizationResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete,omitempty"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval,omitempty"`
}
