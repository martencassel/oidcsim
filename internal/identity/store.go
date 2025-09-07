package identity

import (
	"context"
	"sync"
)

type IdentityStore interface {
	// Users
	ListUsers(ctx context.Context) ([]UserIdentity, error)
	AddUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, idOrUsername string) (UserIdentity, error)
	RemoveUser(ctx context.Context, userID string) error
	GetUserGroups(ctx context.Context, userID string) ([]*Group, error)

	// Groups
	ListGroups(ctx context.Context) ([]*Group, error)
	AddGroup(ctx context.Context, group *Group) (*Group, error)
	GetGroup(ctx context.Context, groupID string) (*Group, error)
	RemoveGroup(ctx context.Context, groupID string) error

	// Members
	AddGroupMember(ctx context.Context, groupID string, userId string) error
	RemoveGroupMember(ctx context.Context, groupID string, userId string) error
	ListGroupMembers(ctx context.Context, group *Group) ([]*User, error)

	MergeUserClaims(ctx context.Context, user UserIdentity, sources ...IdentityStore) (UserIdentity, error)
}

type CoreIdentityStore struct {
	// Mutex
	mu         sync.RWMutex
	APIBaseURL string
	// Add fields for auth, http client, etc.
	users  map[string]*User
	groups map[string]*Group
}

func NewCoreIdentityStore(apiBaseURL string) *CoreIdentityStore {
	return &CoreIdentityStore{
		APIBaseURL: apiBaseURL,
		users:      make(map[string]*User),
		groups:     make(map[string]*Group),
		mu:         sync.RWMutex{},
	}
}

func (s *CoreIdentityStore) AddUser(ctx context.Context, user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
	return nil
}

func (s *CoreIdentityStore) GetUserGroups(ctx context.Context, userID string) ([]*Group, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[userID]
	if !exists || user == nil {
		return nil, nil // Or return an error indicating user does not exist
	}
	groups := make([]*Group, len(user.Groups))
	for i, g := range user.Groups {
		if group, exists := s.groups[g.ID]; exists {
			groups[i] = group
		}
	}
	return groups, nil
}

func (s *CoreIdentityStore) GetUser(ctx context.Context, idOrUsername string) (UserIdentity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Check by ID
	if user, exists := s.users[idOrUsername]; exists {
		return user, nil
	}
	// Check by Username
	for _, user := range s.users {
		if user.Username == idOrUsername {
			return user, nil
		}
	}
	return nil, nil // Not found
}

func (s *CoreIdentityStore) GetGroup(ctx context.Context, groupID string) (*Group, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if group, exists := s.groups[groupID]; exists {
		return group, nil
	}
	return nil, nil // Not found
}

func (s *CoreIdentityStore) MergeUserClaims(ctx context.Context, user UserIdentity, sources ...IdentityStore) (UserIdentity, error) {
	mergedClaims := make(map[string]interface{})
	// Start with the original user's claims
	for k, v := range user.GetClaims() {
		mergedClaims[k] = v
	}
	// Merge claims from each source
	for _, source := range sources {
		srcUser, err := source.GetUser(ctx, user.GetID())
		if err != nil || srcUser == nil {
			continue // Skip if user not found in this source
		}
		for k, v := range srcUser.GetClaims() {
			mergedClaims[k] = v // Overwrite with source claims
		}
	}
	// Create a new User with merged claims
	return &User{
		ID:       user.GetID(),
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Groups:   []Group{}, // Groups can be merged similarly if needed
		Claims:   mergedClaims,
	}, nil
}

func (s *CoreIdentityStore) AddGroup(ctx context.Context, group *Group) (*Group, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if group == nil {
		return nil, nil
	}
	// Check if group already exists
	if existingGroup, exists := s.groups[group.ID]; exists {
		return existingGroup, nil
	}
	if s.groups == nil {
		s.groups = make(map[string]*Group)
	}
	ok, exists := s.groups[group.ID]
	if exists && ok != nil {
		return ok, nil // Group already exists
	}
	// Append new group
	existingGroup := s.groups[group.ID]
	if existingGroup != nil {
		return existingGroup, nil // Group already exists
	}
	// Add new group
	s.groups[group.ID] = group
	return group, nil
}

func (s *CoreIdentityStore) AddGroupMember(ctx context.Context, groupId string, userId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if groupId == "" || userId == "" {
		return nil
	}
	user, userExists := s.users[userId]
	if !userExists || user == nil {
		return nil // Or return an error indicating user does not exist
	}

	// Ensure user exists
	if _, userExists := s.users[user.ID]; !userExists {
		return nil // Or return an error indicating user does not exist
	}
	// Ensure group exists
	group, groupExists := s.groups[groupId]
	if !groupExists || group == nil {
		return nil // Or return an error indicating group does not exist
	}
	// Check if user is already a member of the group
	for _, g := range user.Groups {
		if g.ID == group.ID {
			return nil // Already a member
		}
	}
	if user.Groups == nil {
		user.Groups = make([]Group, 0)
	}
	user.Groups = append(user.Groups, *group)
	return nil
}

func (s *CoreIdentityStore) RemoveGroupMember(ctx context.Context, groupId string, userId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if groupId == "" || userId == "" {
		return nil
	}
	// Ensure user exists
	user, userExists := s.users[userId]
	if !userExists || user == nil {
		return nil // Or return an error indicating user does not exist
	}
	// Ensure group exists
	group, groupExists := s.groups[groupId]
	if !groupExists || group == nil {
		return nil // Or return an error indicating group does not exist
	}
	// Remove user from group
	for i, member := range user.Groups {
		if member.ID == group.ID {
			user.Groups = append(user.Groups[:i], user.Groups[i+1:]...)
			return nil // Removed successfully
		}
	}
	return nil // User was not a member of the group
}

func (s *CoreIdentityStore) ListGroupMembers(ctx context.Context, group *Group) ([]*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if group == nil {
		return nil, nil
	}
	var members []*User
	for _, user := range s.users {
		for _, g := range user.Groups {
			if g.ID == group.ID {
				members = append(members, user)
				break
			}
		}
	}
	return members, nil
}

func (s *CoreIdentityStore) ListUsers(ctx context.Context) ([]UserIdentity, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var users []UserIdentity
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

func (s *CoreIdentityStore) ListGroups(ctx context.Context) ([]*Group, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var groups []*Group
	for _, group := range s.groups {
		groups = append(groups, group)
	}
	return groups, nil
}

// RemoveUser
func (s *CoreIdentityStore) RemoveUser(ctx context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, userExists := s.users[userID]; !userExists {
		return nil // Or return an error indicating user does not exist
	}
	delete(s.users, userID)
	return nil
}

// RemoveGroup
func (s *CoreIdentityStore) RemoveGroup(ctx context.Context, groupID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, groupExists := s.groups[groupID]; !groupExists {
		return nil // Or return an error indicating group does not exist
	}
	delete(s.groups, groupID)
	// Also remove this group from all users' group lists
	for _, user := range s.users {
		for i := 0; i < len(user.Groups); i++ {
			if user.Groups[i].ID == groupID {
				user.Groups = append(user.Groups[:i], user.Groups[i+1:]...)
				i-- // Adjust index after removal
			}
		}
	}
	return nil
}
