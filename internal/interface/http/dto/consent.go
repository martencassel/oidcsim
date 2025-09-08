package dto

type ConsentView struct {
	ClientName string
	Scopes     []string
	UserID     string
	RequestID  string
}
