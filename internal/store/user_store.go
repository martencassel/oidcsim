package store

import (
	"fmt"
	"time"
)

type User struct {
	ID            string
	Username      string
	FullName      string
	EmailVerified bool

	Roles     []string
	Email     string
	Password  string // In a real implementation, passwords should be hashed
	SessionID string
	Groups    []string
	AuthTime  time.Time
}

type UserStore interface {
	GetByID(id string) (User, error)
	GetByUsername(username string) (User, error)
	Save(user User) error
	List() ([]User, error)
}

type InMemoryUserStore struct {
	users map[string]User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{users: make(map[string]User)}
}

func (s *InMemoryUserStore) GetByID(id string) (User, error) {
	u, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *InMemoryUserStore) GetByUsername(username string) (User, error) {
	for _, u := range s.users {
		if u.Username == username {
			return u, nil
		}
	}
	return User{}, fmt.Errorf("user not found")
}

func (s *InMemoryUserStore) List() ([]User, error) {
	users := make([]User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users, nil
}

func (s *InMemoryUserStore) Save(user User) error {
	s.users[user.ID] = user
	return nil
}
