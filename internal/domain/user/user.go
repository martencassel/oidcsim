package user

// User represents an authenticated user in the system.
type User struct {
	ID                string // maps to "sub" claim
	Name              string
	GivenName         string
	FamilyName        string
	PrefferedUsername string
	Email             string
	EmailVerified     bool
	UserName          string
	Password          string // In real life, passwords should be hashed!
}
