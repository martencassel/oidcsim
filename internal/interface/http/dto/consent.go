package dto

// ConsentView represents the data needed to render a consent page.
type ConsentView struct {
	ClientName string
	Scopes     []string
	UserID     string
	RequestID  string
}
