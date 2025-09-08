package user

import "context"

// UserRepository defines how to retrieve user data from a persistence layer or identity source.
type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
}
