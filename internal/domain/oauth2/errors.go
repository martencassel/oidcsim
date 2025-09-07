package oauth2

import "fmt"

// ErrInvalidRequest returns a domain-level error for invalid authorization requests.
func ErrInvalidRequest(desc string) error {
	return fmt.Errorf("invalid_request: %s", desc)
}

// ErrInvalidClient returns a domain-level error indicating client validation failed.
func ErrInvalidClient(desc string) error {
	return fmt.Errorf("invalid_client: %s", desc)
}

// ErrInvalidRedirectURI returns a domain-level error for redirect URI violations.
func ErrInvalidRedirectURI(desc string) error {
	return fmt.Errorf("invalid_redirect_uri: %s", desc)
}
