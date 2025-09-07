package identity

import (
	"context"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// API Handler for
//

// Users
//
// GET /users
// GET /users/{id}
// POST /users
// DELETE /users/{id}
// GET /users/{id}/groups

// Groups
//
// GET /groups
// GET /groups/{id}
// POST /groups
// DELETE /groups/{id}
// GET /groups/{id}/members

type IdentityStoreApiHandler struct {
	store IdentityStore
	g     *gin.RouterGroup
}

func NewIdentityStoreHandler(store IdentityStore) *IdentityStoreApiHandler {
	return &IdentityStoreApiHandler{
		store: store,
		g:     nil,
	}
}

func (h *IdentityStoreApiHandler) SeedDefault() error {
	// Seed with some default users and groups
	users := []string{"alice", "bob", "carol", "dave"}
	for _, uid := range users {
		if err := h.store.AddUser(context.Background(), &User{ID: uid, Username: uid, Email: uid + "@example.com"}); err != nil {
			return err
		}
	}
	groups := []string{"admins", "users", "guests"}
	for _, gid := range groups {
		if _, err := h.store.AddGroup(context.Background(), &Group{ID: gid, Name: gid}); err != nil {
			return err
		}
	}
	h.store.AddGroupMember(context.Background(), "admins", "alice")
	h.store.AddGroupMember(context.Background(), "users", "bob")
	h.store.AddGroupMember(context.Background(), "users", "carol")
	h.store.AddGroupMember(context.Background(), "guests", "dave")
	return nil
}

func (h *IdentityStoreApiHandler) RegisterRoutes(rg *gin.Engine) {
	h.g = rg.Group("/api")
	// Users
	h.g.GET("/users", h.handleListUsers)
	h.g.POST("/users", h.handleCreateUser)
	h.g.GET("/users/:id", h.handleGetUser)
	h.g.DELETE("/users/:id", h.handleDeleteUser)
	h.g.GET("/users/:id/groups", h.handleGetUserGroups)
	// Groups
	h.g.GET("/groups", h.handleListGroups)
	h.g.POST("/groups", h.handleCreateGroup)
	h.g.GET("/groups/:id", h.handleGetGroup)
	h.g.DELETE("/groups/:id", h.handleDeleteGroup)
	h.g.GET("/groups/:id/members", h.handleGetGroupMembers)
	// Group Members
	h.g.POST("/groups/:id/members", h.handleAddGroupMember)
	h.g.DELETE("/groups/:id/members/:userId", h.handleRemoveGroupMember)
}

// Handlers (Users)

func (h *IdentityStoreApiHandler) handleListUsers(c *gin.Context) {
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	users, err := h.store.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}
func (h *IdentityStoreApiHandler) handleCreateUser(c *gin.Context) {
	// Use gin bind
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	if err := h.store.AddUser(c.Request.Context(), &user); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, user)
}
func (h *IdentityStoreApiHandler) handleGetUser(c *gin.Context) {
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	id := c.Param("id")
	user, err := h.store.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	c.JSON(200, user)
}
func (h *IdentityStoreApiHandler) handleDeleteUser(c *gin.Context) {
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	id := c.Param("id")
	user, err := h.store.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	if err := h.store.RemoveUser(c.Request.Context(), user.GetID()); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

func (h *IdentityStoreApiHandler) handleGetUserGroups(c *gin.Context) {
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	id := c.Param("id")
	user, err := h.store.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	groups := user.GetGroups()
	c.JSON(200, groups)
}

// Handlers (Groups)
func (h *IdentityStoreApiHandler) handleListGroups(c *gin.Context) {
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	groups, err := h.store.ListGroups(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, groups)
}
func (h *IdentityStoreApiHandler) handleCreateGroup(c *gin.Context) {
	// Use gin bind
	var group Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	createdGroup, err := h.store.AddGroup(c.Request.Context(), &group)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, createdGroup)
}
func (h *IdentityStoreApiHandler) handleGetGroup(c *gin.Context) {
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	id := c.Param("id")
	group, err := h.store.GetGroup(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if group == nil {
		c.JSON(404, gin.H{"error": "group not found"})
		return
	}
	c.JSON(200, group)
}

func (h *IdentityStoreApiHandler) handleDeleteGroup(c *gin.Context) {
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	id := c.Param("id")
	group, err := h.store.GetGroup(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if group == nil {
		c.JSON(404, gin.H{"error": "group not found"})
		return
	}
	if err := h.store.RemoveGroup(c.Request.Context(), group.ID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

func (h *IdentityStoreApiHandler) handleGetGroupMembers(c *gin.Context) {
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	id := c.Param("id")
	group, err := h.store.GetGroup(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if group == nil {
		c.JSON(404, gin.H{"error": "group not found"})
		return
	}
	members, err := h.store.ListGroupMembers(c.Request.Context(), group)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.Infof("Group %s has %d members", group.ID, len(members))
	c.JSON(200, members)
}

// Handlers (Group Members)

func (h *IdentityStoreApiHandler) handleAddGroupMember(c *gin.Context) {
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	id := c.Param("id")
	group, err := h.store.GetGroup(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if group == nil {
		c.JSON(404, gin.H{"error": "group not found"})
		return
	}
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := h.store.GetUser(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	if err := h.store.AddGroupMember(c.Request.Context(), group.ID, user.GetID()); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}

func (h *IdentityStoreApiHandler) handleRemoveGroupMember(c *gin.Context) {
	if h.store == nil {
		c.JSON(500, gin.H{"error": "identity store not configured"})
		return
	}
	id := c.Param("id")
	group, err := h.store.GetGroup(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if group == nil {
		c.JSON(404, gin.H{"error": "group not found"})
		return
	}
	userId := c.Param("userId")
	user, err := h.store.GetUser(c.Request.Context(), userId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	if err := h.store.RemoveGroupMember(c.Request.Context(), group.ID, user.GetID()); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(204)
}
