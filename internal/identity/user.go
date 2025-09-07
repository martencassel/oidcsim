package identity

type UserIdentity interface {
	GetID() string
	GetUsername() string
	GetEmail() string
	GetGroups() []GroupIdentity
	GetClaims() map[string]interface{}
}

type GroupIdentity interface {
	GetID() string
	GetName() string
	GetClaims() map[string]interface{}
}

// Concrete implementation of GroupIdentity
type Group struct {
	ID     string                 `json:"id"`
	Name   string                 `json:"name"`
	Claims map[string]interface{} `json:"claims"`
}

func (g *Group) GetID() string                     { return g.ID }
func (g *Group) GetName() string                   { return g.Name }
func (g *Group) GetClaims() map[string]interface{} { return g.Claims }

// Concrete implementation of UserIdentity

type User struct {
	ID       string                 `json:"id"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Groups   []Group                `json:"groups"`
	Claims   map[string]interface{} `json:"claims"`
}

func (u *User) GetID() string       { return u.ID }
func (u *User) GetUsername() string { return u.Username }
func (u *User) GetEmail() string    { return u.Email }
func (u *User) GetGroups() []GroupIdentity {
	groups := make([]GroupIdentity, len(u.Groups))
	for i := range u.Groups {
		groups[i] = &u.Groups[i]
	}
	return groups
}
func (u *User) GetClaims() map[string]interface{} { return u.Claims }
