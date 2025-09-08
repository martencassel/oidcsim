package authentication

import (
	"context"
	"testing"

	"github.com/martencassel/oidcsim/internal/domain/authentication"
	"github.com/martencassel/oidcsim/internal/domain/oauth2"
	users "github.com/martencassel/oidcsim/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

type FakeUserRepo struct {
	users map[string]*users.User
}

type User struct {
	ID       string
	Username string
	Password string // In real life, passwords should be hashed!
}

func NewFakeUserRepo() *FakeUserRepo {
	return &FakeUserRepo{users: make(map[string]*users.User)}
}

func (r *FakeUserRepo) FindByID(ctx context.Context, id string) (*users.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func (r *FakeUserRepo) Authenticate(ctx context.Context, username, password string) (*users.User, error) {
	if user, exists := r.users[username]; exists && user.Password == password {
		return user, nil
	}
	return nil, nil
}

func (r *FakeUserRepo) Add(user *users.User) {
	r.users[user.UserName] = user
}

func (r *FakeUserRepo) FindByUsername(username string) (*users.User, bool) {
	user, exists := r.users[username]
	return user, exists
}

type FakeSessionStore struct {
	sessions map[string]map[string]interface{}
}

func NewFakeSessionStore() *FakeSessionStore {
	return &FakeSessionStore{sessions: make(map[string]map[string]interface{})}
}

func (s *FakeSessionStore) Create(sessionID string) {
	s.sessions[sessionID] = make(map[string]interface{})
}

func (s *FakeSessionStore) GetAuthFlow(sessionID string) []authentication.AuthStep {
	if sess, exists := s.sessions[sessionID]; exists {
		if methods, ok := sess["methods"].([]string); ok {
			var authSteps []authentication.AuthStep
			for _, m := range methods {
				authSteps = append(authSteps, authentication.AuthStep{Method: authentication.AuthMethod(m)})
			}
			return authSteps
		}
	}
	return nil
}

func (s *FakeSessionStore) Get(sessionID string) (authentication.AuthSession, bool) {
	sess, exists := s.sessions[sessionID]
	if !exists {
		return authentication.AuthSession{}, false
	}
	// Convert map[string]interface{} to authentication.AuthSession
	authSess := authentication.AuthSession{}
	if subjectID, ok := sess["subject_id"].(string); ok {
		authSess.SubjectID = subjectID
	}
	if methods, ok := sess["methods"].([]string); ok {
		for _, m := range methods {
			authSess.Methods = append(authSess.Methods, authentication.AuthMethod(m))
		}
	}
	return authSess, true
}

func (s *FakeSessionStore) GetAuthTime(sessionID string) int64 {
	if sess, exists := s.sessions[sessionID]; exists {
		if at, ok := sess["auth_time"].(int64); ok {
			return at
		}
	}
	return 0
}

func (s *FakeSessionStore) Set(sessionID, key string, value interface{}) {
	if sess, exists := s.sessions[sessionID]; exists {
		sess[key] = value
	}
}

func (s *FakeSessionStore) Delete(sessionID string) error {
	delete(s.sessions, sessionID)
	return nil
}

func (s *FakeSessionStore) GetCompletedSteps(sessionID string) []authentication.AuthMethod {
	if sess, exists := s.sessions[sessionID]; exists {
		if methods, ok := sess["completed_methods"].([]string); ok {
			var authMethods []authentication.AuthMethod
			for _, m := range methods {
				authMethods = append(authMethods, authentication.AuthMethod(m))
			}
			return authMethods
		}
	}
	return nil
}

func (s *FakeSessionStore) Clear(sessionID string) {
	delete(s.sessions, sessionID)
}

func (s *FakeSessionStore) SetSubject(sessionID, subjectID string) error {
	if sess, exists := s.sessions[sessionID]; exists {
		sess["subject_id"] = subjectID
		return nil
	}
	return nil
}

func (s *FakeSessionStore) GetSubjectID(sessionID string) string {
	if sess, exists := s.sessions[sessionID]; exists {
		if subjectID, ok := sess["subject_id"].(string); ok {
			return subjectID
		}
	}
	return ""
}

func (s *FakeSessionStore) IsAuthenticated(sessionID string) bool {
	if sess, exists := s.sessions[sessionID]; exists {
		if auth, ok := sess["authenticated"].(bool); ok {
			return auth
		}
	}
	return false
}

func (s *FakeSessionStore) MarkAuthenticated(sessionID string) error {
	if sess, exists := s.sessions[sessionID]; exists {
		sess["authenticated"] = true
		sess["auth_time"] = int64(1234567890) // mock auth time
		return nil
	}
	return nil
}

func (s *FakeSessionStore) Save(session authentication.AuthSession) error {
	if sess, exists := s.sessions[session.ID]; exists {
		var methods []string
		for _, m := range session.Methods {
			methods = append(methods, string(m))
		}
		sess["methods"] = methods
		return nil
	}
	return nil
}

func TestAuthService_LoginFlow(t *testing.T) {
	ctx := context.Background()
	sessionID := "sess-abc123"
	clientID := "client-xyz"

	// Setup: create a fake user repo and session store
	userRepo := NewFakeUserRepo()
	userRepo.Add(&users.User{ID: "user123", UserName: "alice", Password: "secret"})
	sessionStore := NewFakeSessionStore()
	authSvc := authentication.NewDefaultAuthService(sessionStore, userRepo)

	// Step 1: Initiate login flow
	spec, err := authSvc.Initiate(ctx, clientID)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(spec.Steps))
	assert.Equal(t, authentication.MethodPassword, spec.Steps[0].Method)

	// Step 2: Start password step
	ui, err := authSvc.StartStep(ctx, sessionID, authentication.MethodPassword)
	assert.NoError(t, err)
	assert.Contains(t, ui.Fields, "username")
	assert.Contains(t, ui.Fields, "password")

	// Step 3: Complete password step
	inputs := map[string]string{
		"username": "alice",
		"password": "secret",
	}
	done, err := authSvc.CompleteStep(ctx, sessionID, authentication.MethodPassword, inputs)
	assert.NoError(t, err)
	assert.False(t, done) // not done yet, OTP next

	// Step 4: Start OTP step
	ui, err = authSvc.StartStep(ctx, sessionID, authentication.MethodOTP)
	assert.NoError(t, err)
	assert.Contains(t, ui.Fields, "otp")

	// Step 5: Complete OTP step
	inputs = map[string]string{
		"otp": "123456",
	}
	done, err = authSvc.CompleteStep(ctx, sessionID, authentication.MethodOTP, inputs)
	assert.NoError(t, err)
	assert.True(t, done)

	// Step 6: Verify authentication context
	authCtx, ok, err := authSvc.Current(ctx, sessionID)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "user123", authCtx.SubjectID)
	assert.True(t, authCtx.IsValidFor(oauth2.AuthorizeRequest{}))
}
