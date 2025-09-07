package identity

import (
	"context"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func NewUser(id string) *User {
	return &User{
		ID:       id,
		Username: id,
		Email:    id + "@example.com",
	}
}

func TestIdentity(t *testing.T) {
	cis := NewCoreIdentityStore("http://example.com/api")

	users := []string{"user1", "user2", "user3", "user4", "user5", "user6"}
	for _, uid := range users {
		err := cis.AddUser(context.Background(), NewUser(uid))
		assert.NoError(t, err)
	}
	for _, uid := range users {
		u, err := cis.GetUser(context.Background(), uid)
		assert.NoError(t, err)
		assert.NotNil(t, u)
		assert.Equal(t, uid, u.GetID())
	}
	groups := []string{"group1", "group2", "group3"}
	for _, gid := range groups {
		g := &Group{
			ID:   gid,
			Name: gid,
		}
		addedGroup, err := cis.AddGroup(context.Background(), g)
		assert.NoError(t, err)
		assert.NotNil(t, addedGroup)
		assert.Equal(t, gid, addedGroup.ID)
	}
	for _, gid := range groups {
		g, err := cis.GetGroup(context.Background(), gid)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		assert.Equal(t, gid, g.ID)
	}

	// Members (group, user) list
	members := [][2]string{
		{"group1", "user1"},
		{"group1", "user2"},
		{"group2", "user3"},
		{"group2", "user4"},
		{"group3", "user5"},
		{"group3", "user6"},
	}
	for _, m := range members {
		g, err := cis.GetGroup(context.Background(), m[0])
		assert.NoError(t, err)
		assert.NotNil(t, g)
		u, err := cis.GetUser(context.Background(), m[1])
		assert.NoError(t, err)
		assert.NotNil(t, u)
		err = cis.AddGroupMember(context.Background(), g.ID, u.GetID())
		assert.NoError(t, err)

	}
	// Print all members per group
	for _, gid := range groups {
		g, err := cis.GetGroup(context.Background(), gid)
		assert.NoError(t, err)
		assert.NotNil(t, g)
		members, err := cis.ListGroupMembers(context.Background(), g)
		assert.NoError(t, err)
		assert.NotNil(t, members)
		t.Logf("Members of group %s:", gid)
		for _, member := range members {
			t.Logf(" - %s", member.GetID())
		}
	}
}
