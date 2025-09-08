package delegation

import (
	"testing"

	"github.com/martencassel/oidcsim/internal/application/delegation"
	domain "github.com/martencassel/oidcsim/internal/domain/delegation"
	"github.com/stretchr/testify/assert"
)

func TestDelegationMemoryRepo(t *testing.T) {
	repo_ := NewMemoryRepo()
	var repo delegation.Repository
	repo = repo_
	repo.Save(nil, domain.Delegation{
		ID:       "delegation1",
		UserID:   "alice",
		ClientID: "client1",
		Scopes:   []string{"openid", "profile", "email"},
	})
	d, err := repo.FindByUserAndClient(nil, "alice", "client1")
	assert.NoError(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, "delegation1", d.ID)

	d2, err := repo.FindByID(nil, "alice|client1")
	assert.NoError(t, err)
	assert.NotNil(t, d2)
	assert.Equal(t, "delegation1", d2.ID)

	err = repo.Delete(nil, "alice", "client1")
	assert.NoError(t, err)

	d3, err := repo.FindByUserAndClient(nil, "alice", "client1")
	assert.NoError(t, err)
	assert.Nil(t, d3)

}
